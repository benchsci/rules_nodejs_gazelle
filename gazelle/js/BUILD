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
    name = "go_default_library",
    srcs = [
        "colors.go",
        "config.go",
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

go_test(
    name = "go_default_test",
    srcs = [
        "parse_test.go",
        "pkgname_test.go",
    ],
    embed = [":go_default_library"],
    importpath = "benchsci.com/gazelle/v2",
    rundir = ".",
    deps = [
        "@bazel_gazelle//label:go_default_library",
        "@bazel_gazelle//resolve:go_default_library",
        "@bazel_gazelle//testtools:go_default_library",
    ],
)

go_test(
    name = "gazelle_test",
    srcs = ["js_test.go"],
    data = [
        "//:gazelle_bin",
    ] + glob(["testdata/**"]),
    deps = [
        "@bazel_gazelle//testtools:go_default_library",
        "@com_github_emirpasic_gods//lists/singlylinkedlist",
        "@com_github_ghodss_yaml//:yaml",
        "@io_bazel_rules_go//go/tools/bazel:go_default_library",
    ],
)
