load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

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

go_library(
    name = "tests_lib",
    testonly = True,
    srcs = ["test_runner.go"],
    importpath = "github.com/benchsci/rules_nodejs_gazelle/tests",
    visibility = ["//visibility:private"],
    deps = [
        "@bazel_gazelle//testtools:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel:go_default_library",
    ],
)

go_binary(
    name = "test_runner",
    testonly = True,
    data = ["//:test_gazelle_bin"],
    embed = [":tests_lib"],
    visibility = ["//visibility:private"],
)

[
    sh_test(
        name = "%s_test" % t,
        srcs = [":test_runner"],
        args = [t],
        data = glob(["%s/**" % t]),
    )
    for t in [
        "collect_all",
        "collect_all_nested",
        "collect_all_test_shards",
        "collect_asset_modules",
        "collect_asset_singletons",
        "default_npm_label",
        "disabled",
        "disjoint_module",
        "dynamic_import",
        "fix",
        "import_alias",
        "jest_mock",
        "jsx_conversion",
        "lookup_types",
        "module_self_import",
        "react_example",
        "simple_barrel",
        "simple_library",
        "simple_npm_library",
        "ts_conversion",
        "visibility",
        "web_assets_module",
        "monorepo",
    ]
]
