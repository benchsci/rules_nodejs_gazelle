load("@aspect_rules_js//js:defs.bzl", "js_library")

def web_assets(**kwargs):
    """Maps web_assets to js_library"""
    js_library(**kwargs)
