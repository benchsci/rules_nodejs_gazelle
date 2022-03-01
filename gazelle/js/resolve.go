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
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// BUILTINS list taken from https://github.com/sindresorhus/builtin-modules/blob/master/builtin-modules.json
var BUILTINS = map[string]bool{
	"assert":         true,
	"async_hooks":    true,
	"buffer":         true,
	"child_process":  true,
	"cluster":        true,
	"console":        true,
	"constants":      true,
	"crypto":         true,
	"dgram":          true,
	"dns":            true,
	"domain":         true,
	"events":         true,
	"fs":             true,
	"http":           true,
	"http2":          true,
	"https":          true,
	"inspector":      true,
	"module":         true,
	"net":            true,
	"os":             true,
	"path":           true,
	"perf_hooks":     true,
	"process":        true,
	"punycode":       true,
	"querystring":    true,
	"readline":       true,
	"repl":           true,
	"stream":         true,
	"string_decoder": true,
	"timers":         true,
	"tls":            true,
	"trace_events":   true,
	"tty":            true,
	"url":            true,
	"util":           true,
	"v8":             true,
	"vm":             true,
	"wasi":           true,
	"worker_threads": true,
	"zlib":           true,
}

// maps resolve.Resolver -> *JS
// Resolver is an interface that language extensions can implement to resolve
// dependencies in rules they generate.

// Name returns the name of the language. This should be a prefix of the
// kinds of rules generated by the language, e.g., "go" for the Go extension
// since it generates "go_library" rules.
func (*JS) Name() string {
	return "js"
}

// Imports returns a list of ImportSpecs that can be used to import the rule
// r. This is used to populate RuleIndex.
//
// If nil is returned, the rule will not be indexed. If any non-nil slice is
// returned, including an empty slice, the rule will be indexed.
func (lang *JS) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {

	srcs := r.AttrStrings("srcs")

	module := false
	// look for index.js and mark this rule as a module rule
	for _, src := range srcs {
		if isModuleFile(src) {
			module = true
			break
		}
	}

	n := len(srcs)
	if module {
		n += 1
	}

	importSpecs := make([]resolve.ImportSpec, n)

	// index each source file
	for i, src := range srcs {
		filePath := path.Join(f.Pkg, src)
		importSpecs[i] = resolve.ImportSpec{
			Lang: lang.Name(),
			Imp:  strings.ToLower(filePath),
		}
	}

	// modules can be resolved via the directory containing them
	if module {
		importSpecs[n-1] = resolve.ImportSpec{
			Lang: lang.Name(),
			Imp:  strings.ToLower(f.Pkg),
		}
	}

	return importSpecs
}

// Embeds returns a list of labels of rules that the given rule embeds. If
// a rule is embedded by another importable rule of the same language, only
// the embedding rule will be indexed. The embedding rule will inherit
// the imports of the embedded rule.
func (*JS) Embeds(r *rule.Rule, from label.Label) []label.Label {
	return nil
}

// https://www.typescriptlang.org/docs/handbook/module-resolution.html#classic
// Resolve translates imported libraries for a given rule into Bazel
// dependencies. Information about imported libraries is returned for each
// rule generated by language.GenerateRules in
// language.GenerateResult.Imports. Resolve generates a "deps" attribute (or
// the appropriate language-specific equivalent) for each import according to
// language-specific rules and heuristics.
func (lang *JS) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, _imports interface{}, from label.Label) {

	packageJSON := "//:package"
	packageResolveResult := lang.tryResolve("package.json", c, ix, from)
	if packageResolveResult.err != nil {
		log.Print(Err("%v", packageResolveResult.err))
		return
	}
	if packageResolveResult.selfImport {
		// ignore self imports
		return
	}
	if packageResolveResult.label != label.NoLabel {
		// add discovered label
		lbl := packageResolveResult.label
		packageJSON = lbl.Abs(from.Repo, from.Pkg).String()
	}

	imports := _imports.(*imports)
	depSet := make(map[string]bool)
	dataSet := make(map[string]bool)

	for name := range imports.set {

		// is it a package.json import?
		if name == "package" || name == "package.json" {
			depSet[packageJSON] = true
			continue
		}

		// is it an npm dependency?
		if lang.isNpmDependency(name) {
			s := strings.Split(name, "/")
			name = s[0]
			if strings.HasPrefix(name, "@") {
				name += "/" + s[1]
			}
			depSet["@npm//"+name] = true

			ifc := c.Exts["ts-auto-types"]
			autoTypes := ifc.(bool)
			if autoTypes && r.Kind() == "ts_project" {
				// does it have a corresponding @types/[...] declaration?
				if lang.isNpmDependency("@types/" + name) {
					depSet["@npm//@types/"+name] = true
				}
			}

			continue
		}

		// is it a builtin?
		if _, ok := BUILTINS[name]; ok {
			// Built in module -> ignore
			continue
		}

		// fix aliases
		match := lang.Config.ImportAliasPattern.FindStringSubmatch(name)
		if len(match) > 0 {
			prefix := match[0]
			alias := lang.Config.ImportAliases.aliases[prefix]
			name = alias + strings.TrimPrefix(name, prefix)
		}

		lang.resolveWalkParents(name, depSet, dataSet, c, ix, rc, r, from)

	}

	deps := []string{}
	for dep := range depSet {
		deps = append(deps, dep)
	}
	if len(deps) > 0 {
		r.SetAttr("deps", deps)
	}

	data := []string{}
	for d := range dataSet {
		data = append(data, d)
	}
	if len(data) > 0 {
		r.SetAttr("data", data)
	}
}

func (lang *JS) resolveWalkParents(name string, depSet map[string]bool, dataSet map[string]bool, c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, from label.Label) {

	parents := ""
	tries := []string{}

	for {

		if name == "package" {
			name = "package.json"
		}

		localDir := path.Join(from.Pkg, parents)
		target := path.Join(localDir, name)

		// add supported extensions to target name to get a filePath
		for _, ext := range append(append([]string{""}, tsExtensions...), jsExtensions...) {

			filePath := target + ext
			tries = append(tries, filePath)

			// try to find a rule providing the filePath
			resolveResult := lang.tryResolve(filePath, c, ix, from)
			if resolveResult.err != nil {
				log.Print(Err("%v", resolveResult.err))
				return
			}
			if resolveResult.selfImport {
				// ignore self imports
				return
			}
			if resolveResult.label != label.NoLabel {
				// add discovered label
				lbl := resolveResult.label
				dep := lbl.Rel(from.Repo, from.Pkg).String()
				depSet[dep] = true
				return
			}
			if resolveResult.fileName != "" {
				// add discovered file
				pkgName := path.Dir(target)
				data := fmt.Sprintf("//%s:%s", pkgName, resolveResult.fileName)
				dataSet[data] = true
				return
			}

		}

		// don't look higher than web root for files
		if localDir == lang.Config.WebRoot {
			// unable to resolve import
			log.Print(Err("[%s] import %v not found", from.Abs(from.Repo, from.Pkg).String(), name))
			log.Print(Warn("tried @npm//%s", name))
			for _, try := range tries {
				log.Print(Warn("tried %s", try))
			}
			return
		}

		// continue to search one directory higher
		parents += "../"
	}

}

//  https://nodejs.org/api/modules.html#modules_all_together
func (lang *JS) isNpmDependency(imp string) bool {

	// These prefixes cannot be NPM dependencies
	var prefixes = []string{".", "/", "../", "~/", "@/", "~~/"}
	if hasPrefix(prefixes, imp) {
		return false
	}

	// Assume all @ imports are npm dependencies
	if strings.HasPrefix(imp, "@") {
		return true
	}

	// Grab the first part of the import (ie "foo/bar" -> "foo")
	packageRoot := imp
	for i := range imp {
		if imp[i] == '/' {
			packageRoot = imp[:i]
			break
		}
	}

	// Is the package root found in package.json ?
	if _, ok := lang.Config.NpmDependencies.Dependencies[packageRoot]; ok {
		return true
	}

	if _, ok := lang.Config.NpmDependencies.DevDependencies[packageRoot]; ok {
		return true
	}

	return false
}

func hasPrefix(suffixes []string, x string) bool {
	for _, suffix := range suffixes {
		if strings.HasPrefix(x, suffix) {
			return true
		}
	}
	return false
}

type resolveResult struct {
	label      label.Label
	selfImport bool
	fileName   string
	err        error
}

func (lang *JS) tryResolve(target string, c *config.Config, ix *resolve.RuleIndex, from label.Label) resolveResult {

	importSpec := resolve.ImportSpec{
		Lang: lang.Name(),
		Imp:  strings.ToLower(target),
	}

	matches := ix.FindRulesByImportWithConfig(c, importSpec, lang.Name())

	// too many matches
	if len(matches) > 1 {
		return resolveResult{
			label:      label.NoLabel,
			selfImport: false,
			fileName:   "",
			err:        fmt.Errorf("multiple rules (%s and %s) provide %s", matches[0].Label, matches[1].Label, target),
		}
	}

	// no matches
	if len(matches) == 0 {

		// no rule is found for this file
		// it could be a regular file w/o a target
		if fileInfo, err := os.Stat(path.Join(c.RepoRoot, target)); err == nil && !fileInfo.IsDir() {
			// found a file matching the target
			return resolveResult{
				label:      label.NoLabel,
				selfImport: false,
				fileName:   fileInfo.Name(),
				err:        nil,
			}

		}
		return resolveResult{
			label:      label.NoLabel,
			selfImport: false,
			fileName:   "",
			err:        nil,
		}
	}

	if matches[0].IsSelfImport(from) {
		return resolveResult{
			label:      label.NoLabel,
			selfImport: true,
			fileName:   "",
			err:        nil,
		}
	}

	return resolveResult{
		label:      matches[0].Label,
		selfImport: false,
		fileName:   "",
		err:        nil,
	}

}
