#!/usr/bin/env sh
# Validate published tutorial YAML examples with grcli's Gemara validator.
#
# Usage: validate-examples.sh [GEMARA_REF]
#
# GEMARA_REF defaults to latest. GEMARA_EXAMPLES_SPEC_DIR may point to a
# local Gemara checkout when validating a branch or another non-release ref.

set -eu

# Prefer Git for Windows' POSIX utilities over Windows system utilities.
case "$(uname -s 2>/dev/null || true)" in
  MINGW*|MSYS*|CYGWIN*) export PATH="/usr/bin:/bin:${PATH}" ;;
esac

GEMARA_REF="${1:-latest}"
SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
TUTORIALS_DIR="${TUTORIALS_DIR:-${SCRIPT_DIR}/../../tutorials}"
SPEC_DIR="${GEMARA_EXAMPLES_SPEC_DIR:-}"
INVALID_FIXTURE="${TUTORIALS_DIR}/testdata/invalid-guidance-catalog.yaml"

if ! command -v grcli >/dev/null 2>&1; then
  echo "ERROR: grcli is required to validate tutorial examples. Install grcli and add it to PATH." >&2
  exit 1
fi

if ! command -v cue >/dev/null 2>&1; then
  echo "ERROR: cue is required by grcli validate. Install CUE and add it to PATH." >&2
  exit 1
fi

if [ -n "$SPEC_DIR" ] && [ ! -d "$SPEC_DIR" ]; then
  echo "ERROR: Gemara spec directory does not exist: $SPEC_DIR" >&2
  exit 1
fi

validate_file() {
  file="$1"
  if [ -n "$SPEC_DIR" ]; then
    grcli validate --file "$file" --spec "$SPEC_DIR"
  else
    grcli validate --file "$file" --gemara-version "$GEMARA_REF"
  fi
}

if [ ! -d "$TUTORIALS_DIR" ]; then
  echo "ERROR: tutorials directory does not exist: $TUTORIALS_DIR" >&2
  exit 1
fi

examples=$(find "$TUTORIALS_DIR" -type f \( -name '*.yaml' -o -name '*.yml' \) \
  ! -path "${TUTORIALS_DIR}/testdata/*" | sort)
if [ -z "$examples" ]; then
  echo "ERROR: no tutorial YAML examples found under $TUTORIALS_DIR" >&2
  exit 1
fi

echo "  > Validating tutorial examples against Gemara ${GEMARA_REF} ..."
printf '%s\n' "$examples" | while IFS= read -r yaml_file; do
  [ -n "$yaml_file" ] || continue
  validate_file "$yaml_file"
done

if [ ! -f "$INVALID_FIXTURE" ]; then
  echo "ERROR: invalid-fixture smoke test file is missing: $INVALID_FIXTURE" >&2
  exit 1
fi

echo "  > Checking that grcli rejects a schema-invalid GuidanceCatalog ..."
smoke_output=$(mktemp)
trap 'rm -f "$smoke_output"' EXIT HUP INT TERM
if validate_file "$INVALID_FIXTURE" >"$smoke_output" 2>&1; then
  cat "$smoke_output"
  echo "ERROR: grcli accepted the intentionally invalid tutorial fixture" >&2
  exit 1
fi

# A missing cue executable, malformed input, or another setup failure is not
# evidence that the schema validator rejected the fixture. Require grcli to
# report an actual cue vet validation failure.
if ! grep -Fq 'cue vet exited with:' "$smoke_output"; then
  cat "$smoke_output" >&2
  echo "ERROR: grcli failed before reporting a CUE schema validation error" >&2
  exit 1
fi

cat "$smoke_output"
echo "  > Smoke test passed: the invalid fixture was rejected by CUE validation."
