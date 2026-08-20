#!/usr/bin/env bash
# Reads or writes the core module's version constant.
#
# The value lives in internal/version/version.go as `const Version = "X.Y.Z"`,
# and is the single source of truth tagged for llm request tagging. Both the
# release workflows and a maintainer working locally go through this script so
# the constant is never edited by hand and drifts from the release tag.
#
# Usage:
#   version.sh get              # prints the current version, e.g. 2.0.0
#   version.sh set X.Y.Z[-pre]  # rewrites the constant, prints the new value
set -euo pipefail

ROOT="$(git -C "$(dirname "$0")" rev-parse --show-toplevel)"
FILE="$ROOT/internal/version/version.go"

# The numeric core is required; an optional -prerelease and +build are allowed
# so a release candidate (2.1.0-rc1) can be prepared the same way. A release tag
# for the core module is still vX.Y.Z.
SEMVER='^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$'

usage() {
  echo "usage: version.sh get | set <X.Y.Z[-pre]>" >&2
  exit 2
}

current() {
  # Extract the value from the single `const Version = "..."` line.
  sed -n 's/^const Version = "\(.*\)"$/\1/p' "$FILE"
}

case "${1:-}" in
  get)
    [ -f "$FILE" ] || { echo "$FILE not found" >&2; exit 1; }
    v="$(current)"
    [ -n "$v" ] || { echo "could not find the Version constant in $FILE" >&2; exit 1; }
    echo "$v"
    ;;
  set)
    new="${2:-}"
    [ -n "$new" ] || usage
    echo "$new" | grep -Eq "$SEMVER" || { echo "not a valid version: '$new'" >&2; exit 1; }
    [ -f "$FILE" ] || { echo "$FILE not found" >&2; exit 1; }
    # Rewrite only the constant line. The version has no characters that are
    # special to sed's replacement text, and SEMVER above guarantees that.
    sed -i -E 's/^const Version = ".*"$/const Version = "'"$new"'"/' "$FILE"
    grep -qx "const Version = \"$new\"" "$FILE" || {
      echo "failed to update $FILE" >&2
      exit 1
    }
    echo "$new"
    ;;
  *)
    usage
    ;;
esac
