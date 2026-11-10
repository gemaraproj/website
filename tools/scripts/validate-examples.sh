#!/usr/bin/env sh
# validate-examples.sh – validate every tutorial YAML example against
#   a pinned Gemara schema release using `cue vet`.
#
# Usage:
#   validate-examples.sh [GEMARA_REF] [--smoke-test]
#
#   GEMARA_REF      Schema release tag (default: v1.6.0). Pass "latest" or any tag.
#   --smoke-test    Instead of scanning all tutorial YAMLs, validate only the
#                   intentionally-invalid fixtures in tutorials/testdata/ and
#                   exit non-zero when they unexpectedly pass (CI sanity check).
#
# Environment:
#   TUTORIALS_DIR   Override the tutorial root (default: relative to this script).
#
# The mapping from metadata.type → CUE schema definition is maintained
# here and must be updated when new artifact types are added to Gemara.

set -eu

# ---------------------------------------------------------------------------
# Platform path fixups (Windows Git for Windows)
#  - Prepend /usr/bin so GNU find is chosen over Windows' system32\find.exe
#  - Append $GOPATH/bin so the `cue` binary from `go install` is found.
# ---------------------------------------------------------------------------
case "$(uname -s 2>/dev/null || true)" in
  MINGW*|MSYS*|CYGWIN*)
    export PATH="/usr/bin:/bin:${PATH}"
    if command -v go >/dev/null 2>&1; then
      export PATH="${PATH}:$(go env GOPATH)/bin"
    fi
    ;;
esac

# ---------------------------------------------------------------------------
# Argument parsing
# ---------------------------------------------------------------------------
GEMARA_REF="v1.6.0"
SMOKE_TEST=0

for arg in "$@"; do
  case "$arg" in
    --smoke-test) SMOKE_TEST=1 ;;
    --*)          echo "Unknown flag: $arg" >&2; exit 1 ;;
    *)            GEMARA_REF="$arg" ;;
  esac
done

SCHEMA_MODULE="github.com/gemaraproj/gemara@${GEMARA_REF}"
SCRIPT_DIR="$(dirname "$0")"
TUTORIALS_DIR="${TUTORIALS_DIR:-${SCRIPT_DIR}/../../tutorials}"

# ---------------------------------------------------------------------------
# type_to_schema: map a Gemara metadata.type value to its #Definition name
# ---------------------------------------------------------------------------
type_to_schema() {
  case "$1" in
    AuditLog)          echo "#AuditLog" ;;
    CapabilityCatalog) echo "#CapabilityCatalog" ;;
    ControlCatalog)    echo "#ControlCatalog" ;;
    EnforcementLog)    echo "#EnforcementLog" ;;
    EvaluationLog)     echo "#EvaluationLog" ;;
    GuidanceCatalog)   echo "#GuidanceCatalog" ;;
    MappingDocument)   echo "#MappingDocument" ;;
    Policy)            echo "#Policy" ;;
    PrincipleCatalog)  echo "#PrincipleCatalog" ;;
    RiskCatalog)       echo "#RiskCatalog" ;;
    ThreatCatalog)     echo "#ThreatCatalog" ;;
    VectorCatalog)     echo "#VectorCatalog" ;;
    *)                 echo "" ;;
  esac
}

# ---------------------------------------------------------------------------
# validate_file: run `cue vet` on a single YAML file and print result.
# Returns 0 on pass, 1 on fail.
# ---------------------------------------------------------------------------
validate_file() {
  yaml_file="$1"
  expected_pass="$2"   # 1 = expect pass (normal examples), 0 = expect fail (invalid fixtures)

  artifact_type=$(awk '
    /^metadata:/ { in_meta=1 }
    in_meta && /^[[:space:]]+type:/ {
      sub(/^[[:space:]]+type:[[:space:]]*/, "")
      gsub(/["'"'"']/, "")
      print
      exit
    }
    /^[^[:space:]]/ && !/^metadata:/ { in_meta=0 }
  ' "$yaml_file")

  if [ -z "$artifact_type" ]; then
    echo "  SKIP  $yaml_file (no metadata.type found)"
    return 2
  fi

  schema=$(type_to_schema "$artifact_type")

  if [ -z "$schema" ]; then
    echo "  UNKN  $yaml_file (unrecognised type: $artifact_type)"
    return 3
  fi

  if cue vet -c -d "$schema" "$SCHEMA_MODULE" "$yaml_file" 2>&1; then
    echo "  PASS  $yaml_file ($artifact_type → $schema)"
    return 0
  else
    echo "  FAIL  $yaml_file ($artifact_type → $schema)"
    return 1
  fi
}

# ---------------------------------------------------------------------------
# Smoke test mode: validate fixtures in testdata/ — they must all FAIL cue vet
# ---------------------------------------------------------------------------
if [ "$SMOKE_TEST" -eq 1 ]; then
  echo "  > Smoke-testing invalid fixtures in ${TUTORIALS_DIR}/testdata/ against ${SCHEMA_MODULE} ..."
  echo "  > (All fixtures are expected to FAIL validation.)"
  echo ""

  smoke_pass=0
  smoke_fail=0

  while IFS= read -r yaml_file; do
    artifact_type=$(awk '
      /^metadata:/ { in_meta=1 }
      in_meta && /^[[:space:]]+type:/ {
        sub(/^[[:space:]]+type:[[:space:]]*/, "")
        gsub(/["'"'"']/, "")
        print
        exit
      }
      /^[^[:space:]]/ && !/^metadata:/ { in_meta=0 }
    ' "$yaml_file")

    if [ -z "$artifact_type" ]; then
      # Fixtures with no type still need a schema; try #GuidanceCatalog as default.
      artifact_type="GuidanceCatalog"
    fi

    schema=$(type_to_schema "$artifact_type")
    if [ -z "$schema" ]; then
      schema="#GuidanceCatalog"
    fi

    if cue vet -c -d "$schema" "$SCHEMA_MODULE" "$yaml_file" >/dev/null 2>&1; then
      echo "  UNEXPECTED PASS  $yaml_file — fixture should have failed!"
      smoke_fail=$((smoke_fail + 1))
    else
      echo "  CORRECTLY REJECTED  $yaml_file ($artifact_type → $schema)"
      smoke_pass=$((smoke_pass + 1))
    fi
  done <<EOF
$(find "${TUTORIALS_DIR}/testdata" -name '*.yaml' -o -name '*.yml' | sort)
EOF

  echo ""
  echo "  > Smoke results: ${smoke_pass} correctly rejected, ${smoke_fail} unexpectedly passed"

  if [ "$smoke_pass" -eq 0 ]; then
    echo "  > WARNING: No invalid fixtures found in tutorials/testdata/."
    exit 1
  fi

  if [ "$smoke_fail" -gt 0 ]; then
    echo "  > FAIL: ${smoke_fail} fixture(s) passed when they should have failed."
    echo "    The validator may not be working correctly, or the fixture is no longer invalid."
    exit 1
  fi

  echo "  > Smoke test passed: all fixtures correctly rejected."
  exit 0
fi

# ---------------------------------------------------------------------------
# Normal mode: validate all tutorial YAMLs (skip testdata/)
# ---------------------------------------------------------------------------
pass=0
fail=0
skip=0
unknown=0

echo "  > Validating tutorial examples against ${SCHEMA_MODULE} ..."
echo ""

while IFS= read -r yaml_file; do
  # Skip fixtures in testdata/ — those are intentionally invalid
  case "$yaml_file" in
    */testdata/*) continue ;;
  esac

  result=0
  validate_file "$yaml_file" 1 || result=$?

  case "$result" in
    0) pass=$((pass + 1)) ;;
    1) fail=$((fail + 1)) ;;
    2) skip=$((skip + 1)) ;;
    3) unknown=$((unknown + 1)) ;;
  esac
done <<EOF
$(find "$TUTORIALS_DIR" -name '*.yaml' -o -name '*.yml' | sort)
EOF

echo ""
echo "  > Results: ${pass} passed, ${fail} failed, ${skip} skipped, ${unknown} unknown type"

if [ "$unknown" -gt 0 ]; then
  echo "  > WARNING: ${unknown} file(s) had an unrecognised artifact type."
  echo "    Update the type_to_schema() map in validate-examples.sh."
fi

if [ "$fail" -gt 0 ]; then
  echo "  > FAIL: ${fail} example(s) did not validate. See errors above."
  exit 1
fi

echo "  > All examples validated successfully."
