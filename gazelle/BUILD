# Copyright 2019 The Bazel Authors. All rights reserved.
# Modifications copyright (C) 2021 BenchSci Analytics Inc.
# Modifications copyright (C) 2018 Ecosia GmbH

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "gazelle",
    srcs = [
        "colors.go",
        "configure.go",
        "generate.go",
        "kinds.go",
        "lang.go",
        "parse.go",
        "pkgname.go",
        "resolve.go",
    ],
    importpath = "benchsci.com/gazelle/v2",
    visibility = ["//visibility:public"],
    deps = [
        "@bazel_gazelle//config:go_default_library",
        "@bazel_gazelle//label:go_default_library",
        "@bazel_gazelle//language:go_default_library",
        "@bazel_gazelle//repo:go_default_library",
        "@bazel_gazelle//resolve:go_default_library",
        "@bazel_gazelle//rule:go_default_library",
        "@com_github_bazelbuild_buildtools//build:go_default_library",
    ],
)

TEST_DEPS = [
    "@bazel_gazelle//label:go_default_library",
    "@bazel_gazelle//resolve:go_default_library",
    "@bazel_gazelle//testtools:go_default_library",
    "@com_github_emirpasic_gods//lists/singlylinkedlist",
    "@com_github_ghodss_yaml//:yaml",
    "@io_bazel_rules_go//go/tools/bazel:go_default_library",
]

go_test(
    name = "generate_test",
    srcs = [
        "generate_test.go",
    ],
    embed = [":gazelle"],
    importpath = "benchsci.com/gazelle/v2",
    rundir = ".",
    deps = TEST_DEPS,
)

go_test(
    name = "parse_test",
    srcs = [
        "parse_test.go",
    ],
    embed = [":gazelle"],
    importpath = "benchsci.com/gazelle/v2",
    rundir = ".",
    deps = TEST_DEPS,
)

go_test(
    name = "pkgname_test",
    srcs = [
        "pkgname_test.go",
    ],
    embed = [":gazelle"],
    importpath = "benchsci.com/gazelle/v2",
    rundir = ".",
    deps = TEST_DEPS,
)

[
    go_test(
        name = "gazelle_{}".format(t),
        srcs = ["gazelle_test.go"],
        data = [
            "//:gazelle_bin",
        ] + glob(["testdata/{}/**".format(t)]),
        deps = TEST_DEPS,
    )
    for t in [
        "fix",
        "aggregate_asset_modules",
        "disabled",
        "aggregate_asset_singletons",
        "web_assets_module",
        "npm_label",
        "visibility",
        "ts_conversion",
        "simple_module",
        "simple_npm_library",
        "simple_library",
        "module_self_import",
        "lookup_types",
        "jsx_conversion",
        "jest_mock",
        "import_alias",
        "next_app",
    ]
]
