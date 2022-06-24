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

go_test(
    name = "gazelle_aggregate_asset_modules",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/aggregate_asset_modules/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_aggregate_asset_singletons",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/aggregate_asset_singletons/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_disabled",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/disabled/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_fix",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/fix/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_import_alias",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/import_alias/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_jest_mock",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/jest_mock/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_lookup_types",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/lookup_types/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_module_self_import",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/module_self_import/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_simple_library",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/simple_library/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_simple_module",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/simple_module/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_simple_npm_library",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/simple_npm_library/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_visibility",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/visibility/**"]),
    deps = TEST_DEPS,
)

go_test(
    name = "gazelle_web_assets_module",
    srcs = ["gazelle_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/web_assets_module/**"]),
    deps = TEST_DEPS,
)
