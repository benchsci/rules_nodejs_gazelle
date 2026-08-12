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
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
)

func TestPattern(t *testing.T) {
	t.Logf("%v", indexFilePattern)
	if !indexFilePattern.MatchString("index.jsx") {
		t.FailNow()
	}
}

func TestUnmapKind(t *testing.T) {
	std := map[string]config.MappedKind{
		"jest_test":  {FromKind: "jest_test", KindName: "vitest_test"},
		"ts_project": {FromKind: "ts_project", KindName: "mid"},
		"mid":        {FromKind: "mid", KindName: "outer"},
		"web_assets": {FromKind: "web_assets", KindName: "web_assets"},
	}
	cycle := map[string]config.MappedKind{
		"x": {FromKind: "x", KindName: "y"},
		"y": {FromKind: "y", KindName: "x"},
	}
	for _, tc := range []struct {
		desc     string
		kindMap  map[string]config.MappedKind
		in, want string
	}{
		{"single step reverse", std, "vitest_test", "jest_test"},
		{"chained reverse to builtin", std, "outer", "ts_project"},
		{"partial chain", std, "mid", "ts_project"},
		{"self map is identity", std, "web_assets", "web_assets"},
		{"unmapped builtin unchanged", std, "jest_test", "jest_test"},
		{"completely unmapped", std, "js_library", "js_library"},
		{"cycle terminates", cycle, "x", "y"},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			c := &config.Config{KindMap: tc.kindMap}
			if got := unmapKind(c, tc.in); got != tc.want {
				t.Errorf("unmapKind(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
