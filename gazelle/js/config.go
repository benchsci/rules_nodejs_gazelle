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
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// Configs is an extension of map[string]*Config. It provides finding methods
// on top of the mapping.
type JsConfigs map[string]*JsConfig

// ParentForPackage returns the parent Config for the given Bazel package.
func (c *JsConfigs) ParentForPackage(pkg string) *JsConfig {
	dir := filepath.Dir(pkg)
	if dir == "." {
		dir = ""
	}
	parent := (map[string]*JsConfig)(*c)[dir]
	return parent
}

// JsConfig contains configuration values related to js/ts.
//
// This type is public because other languages need to generate rules based
// on JS, so this configuration may be relevant to them.
type JsConfig struct {
	Enabled         bool
	PackageFile     string
	NpmDependencies struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	LookupTypes        bool
	ImportAliases      map[string]string
	ImportAliasPattern *regexp.Regexp
	Visibility         Visibility
	AggregateModules   bool
	AggregateWebAssets bool
	AggregateAllAssets bool
	AggregatedAssets   map[string]bool
	Fix                bool
	JSRoot             string
	WebAssetSuffixes   map[string]bool
	Quiet              bool
	Verbose            bool
}

func NewJsConfig() *JsConfig {
	return &JsConfig{
		Enabled:     true,
		PackageFile: "package.json",
		NpmDependencies: struct {
			Dependencies    map[string]string "json:\"dependencies\""
			DevDependencies map[string]string "json:\"devDependencies\""
		}{
			Dependencies:    make(map[string]string),
			DevDependencies: make(map[string]string),
		},
		ImportAliases:      make(map[string]string),
		ImportAliasPattern: regexp.MustCompile("$^"),
		Visibility: Visibility{
			Labels: []string{},
		},
		AggregateModules:   false,
		AggregateWebAssets: false,
		AggregateAllAssets: false,
		AggregatedAssets:   make(map[string]bool),
		Fix:                false,
		JSRoot:             "/",
		WebAssetSuffixes:   make(map[string]bool),
		Quiet:              false,
		Verbose:            false,
	}
}

// NewChild creates a new child JsConfig. It inherits desired values from the
// current JsConfig and sets itself as the parent to the child.
func (parent *JsConfig) NewChild() *JsConfig {

	child := NewJsConfig()

	child.Enabled = parent.Enabled
	child.PackageFile = parent.PackageFile
	child.NpmDependencies = parent.NpmDependencies // This is treated immutably
	child.ImportAliases = make(map[string]string)  // copy map
	for k, v := range parent.ImportAliases {
		child.ImportAliases[k] = v
	}
	child.ImportAliasPattern = parent.ImportAliasPattern // Regenerated on change to ImportAliases
	child.Visibility = Visibility{
		Labels: make([]string, len(parent.Visibility.Labels)), // copy slice
	}
	for i := range parent.Visibility.Labels {
		child.Visibility.Labels[i] = parent.Visibility.Labels[i]
	}
	child.AggregateModules = parent.AggregateModules
	child.AggregateWebAssets = parent.AggregateWebAssets
	child.AggregateAllAssets = parent.AggregateAllAssets
	child.AggregatedAssets = parent.AggregatedAssets // Reinitialized on change to JSRoot
	child.JSRoot = parent.JSRoot
	child.WebAssetSuffixes = make(map[string]bool) // copy map
	for k, v := range parent.WebAssetSuffixes {
		child.WebAssetSuffixes[k] = v
	}
	child.Quiet = parent.Quiet
	child.Verbose = parent.Verbose

	return child
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

// RegisterFlags registers command-line flags used by the extension. This
// method is called once with the root configuration when Gazelle
// starts. RegisterFlags may set an initial values in Config.Exts. When flags
// are set, they should modify these values.
func (lang *JS) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {}

// CheckFlags validates the configuration after command line flags are parsed.
// This is called once with the root configuration when Gazelle starts.
// CheckFlags may set default values in flags or make implied changes.
func (lang *JS) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	return nil
}

// KnownDirectives returns a list of directive keys that this Configurer can
// interpret. Gazelle prints errors for directives that are not recoginized by
// any Configurer.
func (*JS) KnownDirectives() []string {
	return []string{
		"js_extension",
		"js_root",
		"js_lookup_types",
		"js_fix",
		"js_package_file",
		"js_import_alias",
		"js_visibility",
		"js_aggregate_modules",
		"js_aggregate_web_assets",
		"js_aggregate_all_assets",
		"js_web_asset",
		"js_quiet",
		"js_verbose",
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

	// Create the root config.
	if _, exists := c.Exts[languageName]; !exists {
		rootConfig := NewJsConfig()
		rootConfig.JSRoot = "."
		rootConfig.AggregatedAssets = make(map[string]bool)
		c.Exts[languageName] = JsConfigs{
			"": rootConfig,
		}
	}

	jsConfigs := c.Exts[languageName].(JsConfigs)

	jsConfig, exists := jsConfigs[rel]
	if !exists {
		parent := jsConfigs.ParentForPackage(rel)
		jsConfig = parent.NewChild()
		jsConfigs[rel] = jsConfig
	}

	// Read directives from existing file
	if f != nil {

		for _, directive := range f.Directives {

			switch directive.Key {

			case "js_extension":
				switch directive.Value {
				case "enabled":
					jsConfig.Enabled = true
				case "disabled":
					jsConfig.Enabled = false
				default:
					log.Fatalf(Err("failed to read directive %s: %s, only \"enabled\", and \"disabled\" are valid", directive.Key, directive.Value))
				}

			case "js_lookup_types":
				jsConfig.LookupTypes = readBoolDirective(directive)

			case "js_fix":
				jsConfig.Fix = readBoolDirective(directive)

			case "js_package_file":
				jsConfig.PackageFile = directive.Value

				data, err := ioutil.ReadFile(path.Join(c.RepoRoot, f.Pkg, jsConfig.PackageFile))
				if err != nil {
					log.Fatalf(Err("failed to open %s: %v", directive.Value, err))
				}

				// Clear any existing dependencies
				for k := range jsConfig.NpmDependencies.Dependencies {
					delete(jsConfig.NpmDependencies.Dependencies, k)
				}
				for k := range jsConfig.NpmDependencies.DevDependencies {
					delete(jsConfig.NpmDependencies.DevDependencies, k)
				}

				// Read dependencies from file
				if err := json.Unmarshal(data, &jsConfig.NpmDependencies); err != nil {
					log.Fatalf(Err("failed to parse %s: %v", directive.Value, err))
				}

			case "js_import_alias":
				vals := strings.SplitN(directive.Value, " ", 2)
				jsConfig.ImportAliases[vals[0]] = vals[1]

				// Regenerate ImportAliasPattern
				keyPatterns := make([]string, 0, len(jsConfig.ImportAliases))
				for k := range jsConfig.ImportAliases {
					keyPatterns = append(keyPatterns, fmt.Sprintf("(^%s)", regexp.QuoteMeta(k)))
				}

				var err error
				if jsConfig.ImportAliasPattern, err = regexp.Compile(strings.Join(keyPatterns, "|")); err != nil {
					log.Fatalf(Err("failed to parse %s: %v", directive.Value, err))
				}

			case "js_visibility":
				jsConfig.Visibility.Set(directive.Value)

			case "js_root":
				jSRoot, err := filepath.Rel(".", f.Pkg)
				if err != nil {
					log.Fatalf(Err("failed to read directive %s: %v", directive.Key, err))
				} else {
					jsConfig.JSRoot = jSRoot
					jsConfig.AggregatedAssets = make(map[string]bool)
				}

			case "js_aggregate_modules":
				jsConfig.AggregateModules = readBoolDirective(directive)

			case "js_aggregate_web_assets":
				jsConfig.AggregateWebAssets = readBoolDirective(directive)

			case "js_aggregate_all_assets":
				jsConfig.AggregateAllAssets = readBoolDirective(directive)

			case "js_web_asset":
				vals := strings.SplitN(directive.Value, " ", 2)
				suffixes := vals[0]
				status := false
				if len(vals) > 1 {
					val, err := strconv.ParseBool(directive.Value)
					if err != nil {
						log.Fatalf(Err("failed to read directive %s: %v", directive.Key, err))
					}
					status = val
				}
				for _, suffix := range strings.Split(suffixes, ",") {
					jsConfig.WebAssetSuffixes[suffix] = status
				}

			case "js_quiet":
				jsConfig.Quiet = readBoolDirective(directive)
				if jsConfig.Quiet {
					jsConfig.Verbose = false
				}

			case "js_verbose":
				jsConfig.Verbose = readBoolDirective(directive)
				if jsConfig.Verbose {
					jsConfig.Quiet = false
				}
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

func readBoolDirective(directive rule.Directive) bool {
	if directive.Value == "" {
		return true
	} else {
		val, err := strconv.ParseBool(directive.Value)
		if err != nil {
			log.Fatalf(Err("failed to read directive %s: %v", directive.Key, err))
		}
		return val
	}
}
