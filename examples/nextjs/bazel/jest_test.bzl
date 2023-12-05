load("@aspect_rules_jest//jest:defs.bzl", _jest_test = "jest_test")
load("//bazel:ts_project.bzl", "ts_project")

def jest_test(name = "", srcs = [], deps=[], data=[], **kwargs):
    """Provides defaults for jest_test"""

    node_modules = kwargs.pop("node_modules", "//:node_modules")
    tags = kwargs.pop("tags", ["jest"])

    # This ensures that the test is typechecked.
    ts_project(
        name = "%s_js" % name,
        srcs = srcs,
        deps = deps,
    )

    # However we run the test with the untranspiled srcs.
    # This is to ensure that the snapshot filenames match what would
    # be used when invoking the test without bazel.
    _jest_test(
        name = name,
        data = srcs + data + ["//:tsconfig"],
        node_modules = node_modules,
        tags = tags,
        **kwargs,
    )
