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
)

func TestDirnameRel(t *testing.T) {
	input := "foo/bar/baz"
	expected := "baz"
	if result := PkgName(input); result != expected {
		t.Logf("expected %s, got %s", expected, result)
		t.FailNow()
	}
}

func TestDirnameAbs(t *testing.T) {
	input := "/abc/def/ghi"
	expected := "ghi"
	if result := PkgName(input); result != expected {
		t.Logf("expected %s, got %s", expected, result)
		t.FailNow()
	}
}
func TestDirnameTrailing(t *testing.T) {
	input := "abc/def/ghi/"
	expected := "ghi"
	if result := PkgName(input); result != expected {
		t.Logf("expected %s, got %s", expected, result)
		t.FailNow()
	}
}

func TestDirnameSingle(t *testing.T) {
	input := "abc"
	expected := "abc"
	if result := PkgName(input); result != expected {
		t.Logf("expected %s, got %s", expected, result)
		t.FailNow()
	}
}

func TestDirnameSingleAbs(t *testing.T) {
	input := "/abc"
	expected := "abc"
	if result := PkgName(input); result != expected {
		t.Logf("expected %s, got %s", expected, result)
		t.FailNow()
	}
}

func TestDirnameSingleTrailing(t *testing.T) {
	input := "/abc/"
	expected := "abc"
	if result := PkgName(input); result != expected {
		t.Logf("expected %s, got %s", expected, result)
		t.FailNow()
	}
}
