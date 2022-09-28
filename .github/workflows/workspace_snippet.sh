#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

# Set by GH actions, see
# https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
NAME=rules_nodejs_gazelle
COMMIT=$(git rev-parse HEAD)
PREFIX=${NAME}-${COMMIT}
SHA=$(git archive --format=tar --prefix=${PREFIX}/ ${COMMIT} | gzip | shasum -a 256 | awk '{print $1}')

cat << EOF
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
