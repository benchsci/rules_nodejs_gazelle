#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

# Set by GH actions, see
# https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
NAME=rules_nodejs_gazelle
COMMIT=$(git rev-parse HEAD)
PREFIX=${NAME}-${COMMIT}
SHA=$(git archive --format=tar --prefix=${PREFIX}/ ${COMMIT} | gzip | shasum -a 256 | awk '{print $1}')
INTEGRITY=$(git archive --format=tar --prefix=${PREFIX}/ ${COMMIT} | gzip | openssl dgst -sha256 -binary | openssl base64 -A)

cat <<EOF
MODULE.bazel setup:

\`\`\`starlark
bazel_dep(name = "com_github_benchsci_rules_nodejs_gazelle", version = "0.0.0", repo_name = "com_github_benchsci_rules_nodejs_gazelle")

archive_override(
    module_name = "com_github_benchsci_rules_nodejs_gazelle",
    integrity = "sha256-${INTEGRITY}",
    strip_prefix = "${PREFIX}",
    urls = ["https://github.com/benchsci/rules_nodejs_gazelle/archive/$COMMIT.tar.gz"],
)
\`\`\`

WORKSPACE setup:

\`\`\`starlark
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
http_archive(
    name = "com_github_benchsci_rules_nodejs_gazelle",
    sha256 = "${SHA}",
    strip_prefix = "${PREFIX}",
    url = "https://github.com/benchsci/rules_nodejs_gazelle/archive/$COMMIT.tar.gz",
)
\`\`\`
EOF
