# Javascript/typescipt Gazelle Plugin
Bazel [Gazelle](https://github.com/bazelbuild/bazel-gazelle) rule
that generates BUILD file content for javascript/typescript code.

## Setup

First, you'll need to add Gazelle to your `WORKSPACE` file.
Follow the instructions at https://github.com/bazelbuild/bazel-gazelle#running-gazelle-with-bazel

Next, we need to fetch the third-party Go libraries that the python extension
depends on.

Add this to your `WORKSPACE`:

```starlark
http_archive(
    name = "com_github_benchsci_rules_nodejs_gazelle",
    sha256 = "1493c2d10628a7a59934dfd56b862051c4ad89d95cb8f465673695c0b6c4ba71",
    strip_prefix = "rules_nodejs_gazelle-fcf758ea027a266bd5a9c7ab6440f9f086422ab2",
    urls = [
        "https://github.com/benchsci/rules_nodejs_gazelle/archive/fcf758ea027a266bd5a9c7ab6440f9f086422ab2.tar.gz",
    ],
)
```
Add the following preferably to your root BUILD file.
 
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
