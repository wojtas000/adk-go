#!/usr/bin/env bash
# Fails when internal/version/version.go does not match the release tag.
#
# Run against a checked-out release tag, it guarantees the version constant was
# bumped before the release was cut. The core module is tagged vX.Y.Z, so the
# constant must equal the tag with the leading v stripped.
#
# Submodules (plugin/agentanalytics/vX.Y.Z) are versioned independently and do
# not carry this constant, so any tag that is not a bare vX.Y.Z is skipped.
#
# Usage: check-version.sh <tag>
set -euo pipefail

TAG="${1:?usage: check-version.sh <tag>}"
DIR="$(cd "$(dirname "$0")" && pwd)"

case "$TAG" in
  v[0-9]*) ;;  # core module release: vX.Y.Z
  *)
    echo "Tag '$TAG' is not a core module tag (vX.Y.Z); skipping the version check."
    exit 0
    ;;
esac

want="${TAG#v}"
got="$("$DIR/version.sh" get)"

if [ "$got" != "$want" ]; then
  cat >&2 <<MSG
Release tag is '$TAG' but internal/version/version.go declares Version = "$got".

Before releasing v$want, bump the constant so it matches the tag. The
"Prepare release" workflow (Actions -> Prepare release) does this for you: it
opens a PR setting the version, and the release tag should be created on the
merge commit.
MSG
  exit 1
fi

echo "internal/version/version.go matches the release tag ($TAG)."
