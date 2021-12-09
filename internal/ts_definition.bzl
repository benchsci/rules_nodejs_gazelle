"""ts_definition

This is a simple implementation for "ts_definition" that just directly maps srcs -> definitions
while passing through transitive dependencies
"""

load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo", "ExternalNpmPackageInfo", "JSModuleInfo")

def _ts_definition_impl(ctx):
    return [
        DeclarationInfo(
            declarations = depset(ctx.attr.srcs),
            transitive_declarations = depset(ctx.attr.deps),
        ),
        DefaultInfo(
            files = depset(ctx.files.srcs),
        ),
    ]

ts_definition = rule(
    implementation = _ts_definition_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = True),
        "deps": attr.label_list(),
    },
)
