"""web_assets

This is a simple macro for "web_assets" that echo's "js_library"
This is kept seperate so that users can override it with gazelle's map_kind directive
"""
load("@aspect_rules_js//js:defs.bzl", "js_library")

def web_assets(**kwargs):
    js_library(**kwargs)
