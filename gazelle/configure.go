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
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/buildtools/labels"
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
	ImportAliases      []struct{ From, To string }
	ImportAliasPattern *regexp.Regexp
	Visibility         Visibility
	CollectBarrels     bool
	CollectWebAssets   bool
	CollectAllAssets   bool
	CollectedAssets    map[string]bool
	CollectAll         bool
	CollectAllRoot     string
	CollectAllSources  map[string]bool
	Fix                bool
	JSRoot             string
	WebAssetSuffixes   map[string]bool
	Quiet              bool
	Verbose            bool
	DefaultNpmLabel    string
	JestConfig         string
	JestTestsPerShard  int
	JestSize           string
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
		LookupTypes:        true,
		ImportAliases:      []struct{ From, To string }{},
		ImportAliasPattern: regexp.MustCompile("$^"),
		Visibility: Visibility{
			Labels: []string{},
		},
		CollectBarrels:    false,
		CollectWebAssets:  false,
		CollectAllAssets:  false,
		CollectedAssets:   make(map[string]bool),
		CollectAll:        false,
		CollectAllRoot:    "",
		CollectAllSources: make(map[string]bool),
		Fix:               false,
		JSRoot:            "/",
		WebAssetSuffixes:  make(map[string]bool),
		Quiet:             false,
		Verbose:           false,
		DefaultNpmLabel:   "//:node_modules/",
		JestTestsPerShard: -1,
		JestConfig:        "",
	}
}

// NewChild creates a new child JsConfig. It inherits desired values from the
// current JsConfig and sets itself as the parent to the child.
func (parent *JsConfig) NewChild() *JsConfig {

	child := NewJsConfig()

	child.Enabled = parent.Enabled

	child.PackageFile = parent.PackageFile

	// copy maps
	child.NpmDependencies = struct {
		Dependencies    map[string]string "json:\"dependencies\""
		DevDependencies map[string]string "json:\"devDependencies\""
	}{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}
	for k, v := range parent.NpmDependencies.Dependencies {
		child.NpmDependencies.Dependencies[k] = v
	}
	for k, v := range parent.NpmDependencies.DevDependencies {
		child.NpmDependencies.DevDependencies[k] = v
	}

	child.LookupTypes = parent.LookupTypes
	child.ImportAliases = parent.ImportAliases
	child.ImportAliases = make([]struct{ From, To string }, len(parent.ImportAliases)) // copy slice
	for i := range parent.ImportAliases {
		child.ImportAliases[i] = parent.ImportAliases[i]
	}
	child.ImportAliasPattern = parent.ImportAliasPattern // Regenerated on change to ImportAliases

	child.Visibility = Visibility{
		Labels: make([]string, len(parent.Visibility.Labels)), // copy slice
	}
	for i := range parent.Visibility.Labels {
		child.Visibility.Labels[i] = parent.Visibility.Labels[i]
	}
	child.CollectBarrels = parent.CollectBarrels
	child.CollectWebAssets = parent.CollectWebAssets
	child.CollectAllAssets = parent.CollectAllAssets
	child.CollectedAssets = parent.CollectedAssets // Reinitialized on change to JSRoot

	child.CollectAll = parent.CollectAll
	child.CollectAllRoot = parent.CollectAllRoot
	child.CollectAllSources = parent.CollectAllSources // Copy reference, reinitialized on change to CollectAll

	child.JestTestsPerShard = parent.JestTestsPerShard
	child.JestSize = parent.JestSize
	child.JestConfig = parent.JestConfig

	child.JSRoot = parent.JSRoot
	child.WebAssetSuffixes = make(map[string]bool) // copy map
	for k, v := range parent.WebAssetSuffixes {
		child.WebAssetSuffixes[k] = v
	}
	child.Quiet = parent.Quiet
	child.Verbose = parent.Verbose
	child.DefaultNpmLabel = parent.DefaultNpmLabel

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

func newJsConfigsWithRootConfig() JsConfigs {
	rootConfig := NewJsConfig()
	rootConfig.JSRoot = "."
	rootConfig.CollectedAssets = make(map[string]bool)
	return JsConfigs{
		"": rootConfig,
	}
}

// RegisterFlags registers command-line flags used by the extension. This
// method is called once with the root configuration when Gazelle
// starts. RegisterFlags may set an initial values in Config.Exts. When flags
// are set, they should modify these values.
func (lang *JS) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	c.Exts[languageName] = newJsConfigsWithRootConfig()
}

// CheckFlags validates the configuration after command line flags are parsed.
// This is called once with the root configuration when Gazelle starts.
// CheckFlags may set default values in flags or make implied changes.
func (lang *JS) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	return nil
}

// KnownDirectives returns a list of directive keys that this Configurer can
// interpret. Gazelle prints errors for directives that are not recognized by
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
		"js_collect_barrels",
		"js_aggregate_modules",
		"js_collect_web_assets",
		"js_aggregate_web_assets",
		"js_collect_all_assets",
		"js_aggregate_all_assets",
		"js_collect_all",
		"js_jest_test_per_shard",
		"js_jest_size",
		"js_jest_config",
		"js_web_asset",
		"js_quiet",
		"js_verbose",
		"js_default_npm_label",
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
		c.Exts[languageName] = newJsConfigsWithRootConfig()
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
				values := strings.Split(directive.Value, " ")
				if len(values) != 2 {
					log.Fatalf(Err("failed to read directive %s: %s, expected 2 values", directive.Key, directive.Value))
				}
				jsConfig.PackageFile = values[0]
				npmLabel := values[1]
				if strings.HasPrefix(npmLabel, ":") {
					npmLabel = labels.ParseRelative(npmLabel, f.Pkg).Format()
				}
				if !strings.HasSuffix(npmLabel, ":") && !strings.HasSuffix(npmLabel, "/") {
					npmLabel += "/"
				}

				data, err := os.ReadFile(path.Join(c.RepoRoot, f.Pkg, jsConfig.PackageFile))
				if err != nil {
					log.Fatalf(Err("failed to open %s: %v", directive.Value, err))
				}

				// Read dependencies from file
				newDeps := struct {
					Dependencies    map[string]string "json:\"dependencies\""
					DevDependencies map[string]string "json:\"devDependencies\""
				}{
					Dependencies:    make(map[string]string),
					DevDependencies: make(map[string]string),
				}
				if err := json.Unmarshal(data, &newDeps); err != nil {
					log.Fatalf(Err("failed to parse %s: %v", directive.Value, err))
				}

				// Store npmLabel in dependencies
				for k, _ := range newDeps.Dependencies {
					jsConfig.NpmDependencies.Dependencies[k] = npmLabel
				}
				for k, _ := range newDeps.DevDependencies {
					jsConfig.NpmDependencies.DevDependencies[k] = npmLabel
				}

			case "js_import_alias":
				vals := strings.SplitN(directive.Value, " ", 2)
				jsConfig.ImportAliases = append(jsConfig.ImportAliases, struct{ From, To string }{From: vals[0], To: strings.TrimSpace(vals[1])})

				// Regenerate ImportAliasPattern
				keyPatterns := make([]string, 0, len(jsConfig.ImportAliases))
				for _, alias := range jsConfig.ImportAliases {
					keyPatterns = append(keyPatterns, fmt.Sprintf("(^%s)", regexp.QuoteMeta(alias.From)))
				}

				var err error
				if jsConfig.ImportAliasPattern, err = regexp.Compile(strings.Join(keyPatterns, "|")); err != nil {
					log.Fatalf(Err("failed to parse %s: %v", directive.Value, err))
				}

			case "js_visibility":
				jsConfig.Visibility.Set(directive.Value)
			case "js_default_npm_label":
				jsConfig.DefaultNpmLabel = directive.Value
				if !strings.HasSuffix(jsConfig.DefaultNpmLabel, ":") && !strings.HasSuffix(jsConfig.DefaultNpmLabel, "/") {
					jsConfig.DefaultNpmLabel += "/"
				}

			case "js_root":
				jSRoot, err := filepath.Rel(".", f.Pkg)
				if err != nil {
					log.Fatalf(Err("failed to read directive %s: %v", directive.Key, err))
				} else {
					jsConfig.JSRoot = jSRoot
					jsConfig.CollectedAssets = make(map[string]bool)
				}

			case "js_collect_barrels":
				jsConfig.CollectBarrels = readBoolDirective(directive)

			case "js_aggregate_modules":
				jsConfig.CollectBarrels = readBoolDirective(directive)

			case "js_collect_web_assets":
				jsConfig.CollectWebAssets = readBoolDirective(directive)

			case "js_aggregate_web_assets":
				jsConfig.CollectWebAssets = readBoolDirective(directive)

			case "js_collect_all_assets":
				jsConfig.CollectAllAssets = readBoolDirective(directive)

			case "js_aggregate_all_assets":
				jsConfig.CollectAllAssets = readBoolDirective(directive)

			case "js_collect_all":
				collectRoot, err := filepath.Rel(".", f.Pkg)
				if err != nil {
					log.Fatalf(Err("failed to read directive %s: %v", directive.Key, err))
				} else {
					jsConfig.CollectAllRoot = collectRoot
					jsConfig.CollectAll = true
					jsConfig.CollectAllSources = make(map[string]bool)
				}

			case "js_jest_config":
				jsConfig.JestConfig = labels.ParseRelative(directive.Value, f.Pkg).Format()

			case "js_jest_test_per_shard":
				jsConfig.JestTestsPerShard = readIntDirective(directive)

			case "js_jest_size":
				jsConfig.JestSize = directive.Value

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

var jsTestExtensionsPattern *regexp.Regexp
var tsTestExtensionsPattern *regexp.Regexp
var tsExtensionsPattern *regexp.Regexp
var jsExtensionsPattern *regexp.Regexp

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
var reactFilePattern *regexp.Regexp

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
	reactFilePattern = regexp.MustCompile(`\.(jsx|tsx)$`)
}

func trimExt(baseName string) string {
	matches := trimExtPattern.FindStringSubmatch(baseName)
	if len(matches) > 0 {
		return matches[1]
	}
	return baseName
}

func isBarrelFile(baseName string) bool {
	return indexFilePattern.MatchString(baseName) && !isReactFile(baseName)
}

func isReactFile(baseName string) bool {
	return reactFilePattern.MatchString(baseName)
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

func readIntDirective(directive rule.Directive) int {
	if directive.Value == "" {
		return -1
	} else {
		val, err := strconv.ParseInt(directive.Value, 10, 32)
		if err != nil {
			log.Fatalf(Err("failed to read directive %s: %v", directive.Key, err))
		}
		return int(val)
	}
}
