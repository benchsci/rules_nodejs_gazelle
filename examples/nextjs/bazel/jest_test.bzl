load("@aspect_rules_jest//jest:defs.bzl", _jest_test = "jest_test")

def jest_test(name = "", srcs = [], data=[], **kwargs):
    """Provides defaults for jest_test"""

    node_modules = kwargs.pop("node_modules", "//:node_modules")
    tags = kwargs.pop("tags", ["jest"])

    _jest_test(
        name = name,
        data = srcs + data + ["//:tsconfig"],
        node_modules = node_modules,
        tags = tags,
        **kwargs,
    )
