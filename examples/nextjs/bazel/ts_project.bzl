load("@aspect_rules_ts//ts:defs.bzl", _ts_project = "ts_project")

def ts_project(**kwargs):
    """Provides defaults for ts_project"""

    declaration = kwargs.pop("declaration", True)
    resolve_json_module = kwargs.pop("resolve_json_module", True)
    tsconfig = kwargs.pop("tsconfig", "//:tsconfig")


    _ts_project(
        declaration = declaration,
        resolve_json_module = resolve_json_module,
        tsconfig = tsconfig,
        **kwargs,
    )
