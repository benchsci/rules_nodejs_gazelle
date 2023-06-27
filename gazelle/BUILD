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

# gazelle:exclude gazelle_test.go
# gazelle:exclude generate_test.go
# gazelle:exclude parse_test.go
# gazelle:exclude pkgname_test.go

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
    importpath = "github.com/benchsci/rules_nodejs_gazelle/gazelle",
    visibility = ["//visibility:public"],
    deps = [
        "@bazel_gazelle//config:go_default_library",
        "@bazel_gazelle//label:go_default_library",
        "@bazel_gazelle//language:go_default_library",
        "@bazel_gazelle//repo:go_default_library",
        "@bazel_gazelle//resolve:go_default_library",
        "@bazel_gazelle//rule:go_default_library",
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
    srcs = ["generate_test.go"],
    embed = [":gazelle"],
    importpath = "benchsci.com/gazelle/v2",
    rundir = ".",
    deps = TEST_DEPS,
)

go_test(
    name = "parse_test",
    srcs = ["parse_test.go"],
    embed = [":gazelle"],
    importpath = "benchsci.com/gazelle/v2",
    rundir = ".",
    deps = TEST_DEPS,
)

go_test(
    name = "pkgname_test",
    srcs = ["pkgname_test.go"],
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
            "//:test_gazelle_bin",
        ] + glob(["testdata/{}/**".format(t)]),
        deps = TEST_DEPS,
    )
    for t in [
        "aggregate_asset_modules",
        "aggregate_asset_singletons",
        "disabled",
        "disjoint_module",
        "dynamic_import",
        "fix",
        "folder_rules",
        "import_alias",
        "jest_mock",
        "jsx_conversion",
        "lookup_types",
        "module_self_import",
        "npm_label",
        "react_example",
        "simple_library",
        "simple_module",
        "simple_npm_library",
        "ts_conversion",
        "visibility",
        "web_assets_module",
    ]
]
