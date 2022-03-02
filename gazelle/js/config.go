// Copyright 2019 The Bazel Authors. All rights reserved.
// Modifications copyright (C) 2021 BenchSci Analytics Inc.
// Modifications copyright (C) 2018 Ecosia GmbH

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package js

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// JsConfig contains configuration values related to js/ts.
//
// This type is public because other languages need to generate rules based
// on JS, so this configuration may be relevant to them.
type JsConfig struct {
	PackageFile     string
	NpmDependencies struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	ImportAliases      Aliases
	ImportAliasPattern *regexp.Regexp
	Visibility         Visibility
	Ignores            Patterns
	AggregateModules   bool
	NoAggregateLike    Patterns
	AggregateWebAssets bool
	AggregateAllAssets bool
	Fix                bool
	WebRoot            string
}

func NewJsConfig() *JsConfig {
	return &JsConfig{
		PackageFile: "package.json",
		NpmDependencies: struct {
			Dependencies    map[string]string "json:\"dependencies\""
			DevDependencies map[string]string "json:\"devDependencies\""
		}{
			Dependencies:    make(map[string]string),
			DevDependencies: make(map[string]string),
		},
		ImportAliases: Aliases{
			aliases: make(map[string]string),
		},
		ImportAliasPattern: nil,
		Visibility: Visibility{
			Labels: []string{},
		},
		AggregateModules:   false,
		AggregateWebAssets: false,
		AggregateAllAssets: false,
		WebRoot:            "",
	}
}

var aliasPattern = regexp.MustCompile(`([^:\s]+):([^:\s]+)`)

type Aliases struct {
	aliases map[string]string
}

func (a *Aliases) String() string {
	return fmt.Sprintf("%v", a.aliases)
}

func (a *Aliases) Set(value string) error {
	match := aliasPattern.FindStringSubmatch(value)
	if len(match) == 0 {
		return fmt.Errorf("invalid key:value pair")
	}
	// assign key:value
	a.aliases[match[1]] = match[2]

	return nil
}

type Visibility struct {
	Labels []string
}

func (v *Visibility) String() string {
	return fmt.Sprintf("%v", v.Labels)
}

func (v *Visibility) Set(value string) error {
	v.Labels = append(v.Labels, value)
	return nil
}

type Patterns struct {
	Patterns []regexp.Regexp
}

func (i *Patterns) String() string {
	return fmt.Sprintf("%v", i.Patterns)
}

func (i *Patterns) Set(value string) error {
	r, err := regexp.Compile(value)
	if err != nil {
		return err
	}
	i.Patterns = append(i.Patterns, *r)
	return nil
}

// RegisterFlags registers command-line flags used by the extension. This
// method is called once with the root configuration when Gazelle
// starts. RegisterFlags may set an initial values in Config.Exts. When flags
// are set, they should modify these values.
func (lang *JS) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	c.Exts[lang.Name()] = &lang.Config
	fs.BoolVar(&lang.Config.Fix, "fix", false, "Fix deprecated rules (same as \"gazelle fix\")")
	fs.Var(&lang.Config.Ignores, "ignore", "ignore directories matching this regex")
	fs.StringVar(&lang.Config.PackageFile, "package_file", "package.json", "path to the package.json used to resolve npm deps")
	fs.Var(&lang.Config.ImportAliases, "import_alias", "a key:value re-mapping an import path prefix")
	fs.Var(&lang.Config.Visibility, "visibility", "a default visibility to add to each rule")
	fs.StringVar(&lang.Config.WebRoot, "web_root", "./", "relative path to the web root of the modules")
	fs.BoolVar(&lang.Config.AggregateModules, "aggregate_modules", false, "aggregate pkg/index.js rules")
	fs.Var(&lang.Config.NoAggregateLike, "no_aggregate", "don't aggregate pkg/index.js rules matching this regex")
	fs.BoolVar(&lang.Config.AggregateWebAssets, "aggregate_web_assets", false, "aggregate aggregate_web_assets rules")
	fs.BoolVar(&lang.Config.AggregateAllAssets, "aggregate_all_assets", false, "aggregate web_assets rules at the configured web_root (does not imply -aggregate_web_assets)")
}

// CheckFlags validates the configuration after command line flags are parsed.
// This is called once with the root configuration when Gazelle starts.
// CheckFlags may set default values in flags or make implied changes.
func (lang *JS) CheckFlags(fs *flag.FlagSet, c *config.Config) error {

	data, err := ioutil.ReadFile(path.Join(c.RepoRoot, lang.Config.PackageFile))
	if err != nil {
		log.Fatalf(Err("failed to open %s: %v", lang.Config.PackageFile, err))
	}
	if err := json.Unmarshal(data, &lang.Config.NpmDependencies); err != nil {
		log.Fatalf(Err("failed to parse %s: %v", lang.Config.PackageFile, err))
	}

	keyPatterns := make([]string, 0, len(lang.Config.ImportAliases.aliases))
	for k := range lang.Config.ImportAliases.aliases {
		keyPatterns = append(keyPatterns, fmt.Sprintf("(^%s)", regexp.QuoteMeta(k)))
	}

	lang.Config.ImportAliasPattern, err = regexp.Compile(strings.Join(keyPatterns, "|"))
	if err != nil {
		return err
	}

	// for some reason "gazelle fix isn't working, so this flag will do instead"
	c.ShouldFix = c.ShouldFix || lang.Config.Fix

	return nil
}

// KnownDirectives returns a list of directive keys that this Configurer can
// interpret. Gazelle prints errors for directives that are not recoginized by
// any Configurer.
func (*JS) KnownDirectives() []string {
	return []string{
		"types",
	}
}

// Configure modifies the configuration using directives and other information
// extracted from a build file. Configure is called in each directory.
//
// c is the configuration for the current directory. It starts out as a copy
// of the configuration for the parent directory.
//
// rel is the slash-separated relative path from the repository root to
// the current directory. It is "" for the root directory itself.
//
// f is the build file for the current directory or nil if there is no
// existing build file.
func (*JS) Configure(c *config.Config, rel string, f *rule.File) {

	// Read directives from existing file
	if f != nil {
		for _, directive := range f.Directives {
			// directive.Key = directive.Value
			if directive.Key == "ts-auto-types" {
				val, _ := strconv.ParseBool(directive.Value)
				c.Exts["ts-auto-types"] = val
			}
		}
	}
}

var tsDefsExtensions = []string{
	".d.ts",
	".d.tsx",
}

var jsTestExtensions = []string{
	".test.js",
	".test.jsx",
}

var tsTestExtensions = []string{
	".test.ts",
	".test.tsx",
}

var tsExtensions = []string{
	".ts",
	".tsx",
}

var jsExtensions = []string{
	".js",
	".jsx",
}

var tsDefsExtensionsPattern *regexp.Regexp
var jsTestExtensionsPattern *regexp.Regexp
var tsTestExtensionsPattern *regexp.Regexp
var tsExtensionsPattern *regexp.Regexp
var jsExtensionsPattern *regexp.Regexp

func init() { tsDefsExtensionsPattern = extensionPattern(tsDefsExtensions) }
func init() { tsTestExtensionsPattern = extensionPattern(tsTestExtensions) }
func init() { jsTestExtensionsPattern = extensionPattern(jsTestExtensions) }
func init() { tsExtensionsPattern = extensionPattern(tsExtensions) }
func init() { jsExtensionsPattern = extensionPattern(jsExtensions) }

func extensionPattern(extensions []string) *regexp.Regexp {
	escaped := make([]string, len(extensions))
	for i := range extensions {
		escaped[i] = fmt.Sprintf("(%s$)", regexp.QuoteMeta(extensions[i]))
	}
	return regexp.MustCompile(strings.Join(escaped, "|"))
}

var indexFilePattern *regexp.Regexp
var trimExtPattern *regexp.Regexp

func init() {
	escaped := make([]string, len(tsExtensions)+len(jsExtensions))
	for i, ext := range append(tsExtensions, jsExtensions...) {
		escaped[i] = regexp.QuoteMeta(ext)
	}
	indexFilePattern = regexp.MustCompile(
		fmt.Sprintf(`(index)(%s)$`,
			strings.Join(escaped, "|"),
		),
	)
	trimExtPattern = regexp.MustCompile(
		fmt.Sprintf(`(\S+)(%s)$`,
			strings.Join(escaped, "|"),
		),
	)
}

func trimExt(baseName string) string {
	matches := trimExtPattern.FindStringSubmatch(baseName)
	if len(matches) > 0 {
		return matches[1]
	}
	return baseName
}

func isModuleFile(baseName string) bool {
	return indexFilePattern.MatchString(baseName)
}
