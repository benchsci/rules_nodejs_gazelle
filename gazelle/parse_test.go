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
	"reflect"
	"testing"
)

func TestParseJS(t *testing.T) {
	for _, tc := range []struct {
		desc, name, js string
		want           []string
	}{
		{
			desc: "empty",
			name: "empty^file.js",
			js:   "",
			want: []string{},
		},
		{
			desc: "import single quote",
			name: "single.js",
			js:   `export * from 'date-fns';`,
			want: []string{"date-fns"},
		},
		{
			desc: "import single quote",
			name: "single.js",
			js:   `import dateFns from 'date-fns';`,
			want: []string{"date-fns"},
		},
		{
			desc: "import double quote",
			name: "double.sass",
			js:   `import dateFns from "date-fns";`,
			want: []string{"date-fns"},
		}, {
			desc: "import two",
			name: "two.sass",
			js: `import {format} from 'date-fns'
import Puppy from '@/components/Puppy';`,
			want: []string{"@/components/Puppy", "date-fns"},
		}, {
			desc: "import depth",
			name: "deep.sass",
			js:   `import package from "from/internal/package";`,
			want: []string{"from/internal/package"},
		}, {
			desc: "import multiline",
			name: "multiline.js",
			js: `import {format} from 'date-fns'
import {
	CONST1,
	CONST2,
	CONST3,
} from '~/constants';`,
			want: []string{"date-fns", "~/constants"},
		},
		{
			desc: "simple require",
			name: "require.js",
			js:   `const a = require("date-fns");`,
			want: []string{"date-fns"},
		},
		{
			desc: "ignores incorrect imports",
			name: "incorrect.js",
			js:   `@import "~mapbox.js/dist/mapbox.css";`,
			want: []string{},
		},
		{
			desc: "ignores commented out imports",
			name: "comment.js",
			js: `
    // takes ?inline out of the aliased import path, only if it's set
    // e.g. ~/path/to/file.svg?inline -> ~/path/to/file.svg
    '^~/(.+\\.svg)(\\?inline)?$': '<rootDir>$1',
// const a = require("date-fns");
// import {format} from 'date-fns';
`,
			want: []string{},
		},
		{
			desc: "full import",
			name: "comment.js",
			js: `import "mypolyfill";
import "mypolyfill2";`,
			want: []string{"mypolyfill", "mypolyfill2"},
		},
		{
			desc: "full require",
			name: "full_require.js",
			js:   `require("mypolyfill2");`,
			want: []string{"mypolyfill2"},
		},
		{
			desc: "imports and full imports",
			name: "mixed_imports.js",
			js: `import Vuex, { Store } from 'vuex';
import { createLocalVue, shallowMount } from '@vue/test-utils';

import '~/plugins/intersection-observer-polyfill';
import '~/plugins/intersect-directive';
import ClaimsSection from './claims-section';
`,
			want: []string{"./claims-section", "@vue/test-utils", "vuex", "~/plugins/intersect-directive", "~/plugins/intersection-observer-polyfill"},
		},
		{
			desc: "dynamic require",
			name: "dynamic_require.js",
			js: `
if (process.ENV.SHOULD_IMPORT) {
    // const old = require('oldmapbox.js');
    const leaflet = require('mapbox.js');
}
`,
			want: []string{"mapbox.js"},
		},
		{
			desc: "dynamic import",
			name: "dynamic_import.js",
			js: ` () => import('dynamic_module.js');
			const foo = import('dynamic_module2.js')`,
			want: []string{"dynamic_module.js", "dynamic_module2.js"},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {

			imports, err := ParseJS([]byte(tc.js))
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			if !reflect.DeepEqual(imports, tc.want) {
				t.Errorf("Inequalith.\ngot  %#v;\nwant %#v", imports, tc.want)
			}
		})
	}
}
