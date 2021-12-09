# JS rules for Bazel
Ecosia specific JS Bazel rules to be used with the NodeJS rules

## Setup

```py
http_archive(
    name = "benchsci_bazel_rules_nodejs_contrib",
    TODO
)
```

## Build file generation

Build file generation is provided as a plugin for [gazelle](https://github.com/bazelbuild/bazel-gazelle) and is still WIP 

To setup the gazlle plugin follow the installation instructions provided by the repository and additionally add the following to your root level `BUILD`:

```py
load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle", "gazelle_binary")

# gazelle:exclude node_modules

gazelle(
    name = "gazelle",
    gazelle = ":gazelle_js",
)

gazelle_binary(
    name = "gazelle_js",
    languages = DEFAULT_LANGUAGES + [
        "@benchsci_bazel_rules_nodejs_contrib//gazelle/js:go_default_library",
    ],
)
```
