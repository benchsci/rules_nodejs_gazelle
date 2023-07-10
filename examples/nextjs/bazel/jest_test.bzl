load("@aspect_rules_jest//jest:defs.bzl", _jest_test = "jest_test")
load("//bazel:ts_project.bzl", "ts_project")

def jest_test(name = "", srcs = [], deps=[], data=[], **kwargs):
    """Provides defaults for jest_test"""

    node_modules = kwargs.pop("node_modules", "//:node_modules")
    tags = kwargs.pop("tags", ["jest"])

    ts_project(
        name = "%s_js" % name,
        srcs = srcs,
        deps = deps,
    )

    data.append(":%s_js" % name)

    _jest_test(
        name = name,
        data = data,
        node_modules = node_modules,
        tags = tags,
        **kwargs,
    )
