"""web_assets

This is a simple implementation for "web_assets" that just directly maps srcs -> outs
filegroup cannot be used because some downstream rules need to be aware of "DeclarationInfo"
"""

load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo")

def _web_assets_impl(ctx):
    return [
        DeclarationInfo(
            declarations = depset([]),
            transitive_declarations = depset([]),
        ),
        DefaultInfo(
            files = depset(ctx.files.srcs),
        ),
    ]

web_assets = rule(
    implementation = _web_assets_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = True),
    },
)
