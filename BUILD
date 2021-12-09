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

load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle", "gazelle_binary")

gazelle(
    name = "gazelle",
    gazelle = ":gazelle_bin",
    prefix = "github.com/benchsci/bazel_rules_nodejs_contrib",
    visibility = ["//visibility:public"],
)

gazelle_binary(
    name = "gazelle_bin",
    languages = DEFAULT_LANGUAGES + [
        "//gazelle/js:go_default_library",
    ],
    visibility = ["//visibility:public"],
)
