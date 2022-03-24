"""web_assets

This is a simple implementation for "web_assets" that copies srcs -> outs
copy_to_bin cannot be used because some downstream rules need to be aware of "DeclarationInfo"
"""

# These are the functions used in the copy_file macro. We can't use a macro inside of a rule implementation so we grab these
load("@bazel_skylib//rules/private:copy_file_private.bzl", "copy_bash", "copy_cmd")  # buildifier: disable=bzl-visibility
load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo")

def _web_assets_impl(ctx):
    all_dst = []
    for src in ctx.files.srcs:
        if not src.is_source:
            fail("A source file must be specified in web_assets rule, %s is not a source file." % src.path)
        dst = ctx.actions.declare_file(src.basename, sibling = src)
        if ctx.attr.is_windows:
            copy_cmd(ctx, src, dst)
        else:
            copy_bash(ctx, src, dst)
        all_dst.append(dst)

    return [
        DeclarationInfo(
            declarations = depset([]),
            transitive_declarations = depset(ctx.attr.deps),
        ),
        DefaultInfo(
            files = depset(all_dst),
        ),
    ]

_web_assets = rule(
    implementation = _web_assets_impl,
    attrs = {
        "is_windows": attr.bool(mandatory = True, doc = "Automatically set by macro"),
        "srcs": attr.label_list(allow_files = True),
        "deps": attr.label_list(),
    },
)

def web_assets(name, srcs, deps = [], **kwargs):
    _web_assets(
        name = name,
        is_windows = select({
            "@bazel_tools//src/conditions:host_windows": True,
            "//conditions:default": False,
        }),
        srcs = srcs,
        deps = deps,
        **kwargs
    )
