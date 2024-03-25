#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

# Set by GH actions, see
# https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
NAME=rules_nodejs_gazelle
TAG=${GITHUB_REF_NAME}
# The prefix is chosen to match what GitHub generates for source archives
PREFIX="${NAME}-${TAG:1}"
ARCHIVE="${NAME}-$TAG.tar.gz"
git archive --format=tar --prefix="${PREFIX}/" "${TAG}" | gzip >"$ARCHIVE"
SHA=$(shasum -a 256 "$ARCHIVE" | awk '{print $1}')

cat <<EOF
MODULE.bazel setup:

\`\`\`starlark
bazel_dep(name = "com_github_benchsci_rules_nodejs_gazelle", version = "${TAG:1}", repo_name = "com_github_benchsci_rules_nodejs_gazelle")

\`\`\`

WORKSPACE setup:

\`\`\`starlark
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
http_archive(
    name = "com_github_benchsci_rules_nodejs_gazelle",
    sha256 = "${SHA}",
    strip_prefix = "${PREFIX}",
    url = "https://github.com/benchsci/rules_nodejs_gazelle/releases/download/${TAG}/${ARCHIVE}",
)
\`\`\`
EOF
