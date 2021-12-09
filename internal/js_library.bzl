"""js_library

This is a minimal js_library implementation that directly maps srcs -> outs,
while also passing through transitive dependencies
"""

load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo", "ExternalNpmPackageInfo", "JSModuleInfo")

def _inputs(ctx):
    # Also include files from npm fine grained deps as inputs.
    # These deps are identified by the ExternalNpmPackageInfo provider.
    inputs_depsets = []
    for d in ctx.attr.deps:
        if ExternalNpmPackageInfo in d:
            inputs_depsets.append(d[ExternalNpmPackageInfo].sources)
        if JSModuleInfo in d:
            inputs_depsets.append(d[JSModuleInfo].sources)
        if DeclarationInfo in d:
            inputs_depsets.append(d[DeclarationInfo].declarations)
    return depset(ctx.files.deps, transitive = inputs_depsets)

def _js_library_impl(ctx):
    return [
        JSModuleInfo(
            direct_sources = depset(ctx.attr.srcs + ctx.attr.data),
            sources = _inputs(ctx),
        ),
        DefaultInfo(
            files = depset(ctx.files.srcs),
        ),
    ]

js_library = rule(
    implementation = _js_library_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = True),
        "deps": attr.label_list(),
        "data": attr.label_list(allow_files = True),
    },
)
