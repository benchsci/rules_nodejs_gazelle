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
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var quotePattern = regexp.MustCompile(`([/][/].*)|(?:[/][*](?:\n|.)*?[*][/])`)

func ParseJS(data []byte) ([]string, int, error) {

	lastCommentMatchIndex := 0
	codeBlocks := make([][]int, 0)
	for _, match := range quotePattern.FindAllIndex(data, -1) {
		codeBlocks = append(codeBlocks, []int{lastCommentMatchIndex, match[0]})
		lastCommentMatchIndex = match[1]
	}
	codeBlocks = append(codeBlocks, []int{lastCommentMatchIndex, len(data)})

	imports := make([]string, 0)
	jestTestCount := 0

	for _, block := range codeBlocks {
		blockImports, blockTestCount, err := parseCodeBlock(data[block[0]:block[1]])
		if err != nil {
			return nil, 0, err
		}
		imports = append(imports, blockImports...)
		jestTestCount += blockTestCount
	}
	sort.Strings(imports)

	return imports, jestTestCount, nil
}

const (
	IMPORT         = 1
	REQUIRE        = 2
	EXPORT         = 3
	JEST_MOCK      = 4
	DYNAMIC_IMPORT = 5
)

var jsImportPattern = compileJsImportPattern()

func compileJsImportPattern() *regexp.Regexp {
	stringLiteralPattern := `'(?:[^)\n]+|")*'|"(?:[^)\n]+|')*"`
	importPattern := `^import\s(?:(?:.|\n)+?from )??(?P<import>` + stringLiteralPattern + `)`
	requirePattern := `^\s*?(?:const .+ = )?require\((?P<require>` + stringLiteralPattern + `)\)`
	exportPattern := `^export\s(?:(?:.|\n)+?from )??(?P<export>` + stringLiteralPattern + `)`
	jestMockPattern := `^\s*?(?:const .+ = )?jest.mock\((?P<jestMock>` + stringLiteralPattern + `),`
	dynamicImportPattern := `^.*?import\((?P<dynamicImport>` + stringLiteralPattern + `)\)`
	return regexp.MustCompile(`(?m)` + strings.Join([]string{importPattern, requirePattern, exportPattern, jestMockPattern, dynamicImportPattern}, "|"))
}

var jestTestPattern = regexp.MustCompile(`(?m)^\s*it\(`)

func parseCodeBlock(data []byte) ([]string, int, error) {

	imports := make([]string, 0)
	for _, match := range jsImportPattern.FindAllSubmatch(data, -1) {
		switch {
		case match[IMPORT] != nil:
			unquoted, err := unquoteImportString(match[IMPORT])
			if err != nil {
				return nil, 0, fmt.Errorf("unquoting string literal %s from js, %v", match[IMPORT], err)
			}
			imports = append(imports, unquoted)

		case match[REQUIRE] != nil:
			unquoted, err := unquoteImportString(match[REQUIRE])
			if err != nil {
				return nil, 0, fmt.Errorf("unquoting string literal %s from js, %v", match[REQUIRE], err)
			}
			imports = append(imports, unquoted)

		case match[EXPORT] != nil:
			unquoted, err := unquoteImportString(match[EXPORT])
			if err != nil {
				return nil, 0, fmt.Errorf("unquoting string literal %s from js, %v", match[EXPORT], err)
			}
			imports = append(imports, unquoted)

		case match[JEST_MOCK] != nil:
			unquoted, err := unquoteImportString(match[JEST_MOCK])
			if err != nil {
				return nil, 0, fmt.Errorf("unquoting string literal %s from js, %v", match[JEST_MOCK], err)
			}
			imports = append(imports, unquoted)

		case match[DYNAMIC_IMPORT] != nil:
			unquoted, err := unquoteImportString(match[DYNAMIC_IMPORT])
			if err != nil {
				return nil, 0, fmt.Errorf("unquoting string literal %s from js, %v", match[DYNAMIC_IMPORT], err)
			}
			imports = append(imports, unquoted)

		default:
			// Comment matched. Nothing to extract.
		}
	}
	sort.Strings(imports)

	jestTestCount := len(jestTestPattern.FindAll(data, -1))

	return imports, jestTestCount, nil
}

// unquoteImportString takes a string that has a complex quoting around it
// and returns a string without the complex quoting.
func unquoteImportString(quoted []byte) (string, error) {
	// Adjust quotes so that Unquote is happy. We need a double quoted string
	// without unescaped double quote characters inside.
	noQuotes := bytes.Split(quoted[1:len(quoted)-1], []byte{'"'})
	if len(noQuotes) != 1 {
		for i := 0; i < len(noQuotes)-1; i++ {
			if len(noQuotes[i]) == 0 || noQuotes[i][len(noQuotes[i])-1] != '\\' {
				noQuotes[i] = append(noQuotes[i], '\\')
			}
		}
		quoted = append([]byte{'"'}, bytes.Join(noQuotes, []byte{'"'})...)
		quoted = append(quoted, '"')
	}
	if quoted[0] == '\'' {
		quoted[0] = '"'
		quoted[len(quoted)-1] = '"'
	}

	result, err := strconv.Unquote(string(quoted))
	if err != nil {
		return "", fmt.Errorf("unquoting string literal %s from js: %v", quoted, err)
	}
	return result, err
}
