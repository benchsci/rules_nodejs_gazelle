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

## Directives

Gazelle can be configured with *directives*, which are written as top-level
comments in build files.

Directive comments have the form `# gazelle:key value`.
[More information here](https://github.com/bazelbuild/bazel-gazelle#directives)

Directives apply in the directory where they are set *and* in subdirectories.
This means, for example, if you set `# gazelle:prefix` in the build file
in your project's root directory, it affects your whole project. If you
set it in a subdirectory, it only affects rules in that subtree.

The following directives are recognized by this plugin:


<table>
<thead>
  <tr>
    <th>Directive</th>
    <th>Default value</th>
  </tr>
</thead>
<tbody>

  <tr>
    <td><code># gazelle:js_extension</code></td>
    <td><code>enabled</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Controls whether the JS extension is enabled or not. Sub-packages inherit this value. Can be either "enabled" or "disabled".</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_lookup_types true|false</code></td>
    <td><code>false</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Causes Gazelle to try and find a matching @npm//types dependency for each @npm dependency</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_package_file package.json</code></td>
    <td><code>none</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Instructs Gazelle to use a package.json file to lookup imports from dependencies and devDependencies</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_import_alias some_folder other</code></td>
    <td><code>none</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Specifies partial string substitutions applied to imports before resolving them. Eg. <code># gazelle:js_import_alias foo bar</code> means that <code>import "foo/module"</code> will resolve to the package <code>bar/module</code>. This directive can be used several times.</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_visibility label</code></td>
    <td><code>none</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">By default, internal packages are only visible to its siblings. This directive adds a label internal packages should be visible to additionally. This directive can be used several times, adding a list of labels.</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_root</code></td>
    <td><code>workspace root</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Specifies the current package (folder) as a JS root. Imports for JS and TS consider this folder the root level for relative and absolute imports. This is used on monorepos with multiple Python projects that don't share the top-level of the workspace as the root.</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_aggregate_modules true|false</code></td>
    <td><code>false</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Generate 1 js_library, or ts_project rule per package when a <code>index.ts</code> or <code>index.js</code> file is found, rather than 1 per file. The js_root pkg cannot be a module</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_aggregate_web_assets true|false</code></td>
    <td><code>false</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Causes Gazelle to generate 1 web_assets rule, rather than 1 per file</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_aggregate_all_assets true|false</code></td>
    <td><code>false</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Generates a <code>web_assets</code> rule in the configured <code>web_root</code> that refers to all of the <code>web_assets</code> rules in child packages using</p></td>
  </tr>

  <tr>
    <td><code># gazelle:js_web_asset .json,.css,.scss</code></td>
    <td><code>none</code></td>
  </tr>
  <tr>
    <td colspan="2"><p dir="auto">Files with a matching suffix will have <code>web_assets</code> rules created for them</p></td>
  </tr>

</tbody>
