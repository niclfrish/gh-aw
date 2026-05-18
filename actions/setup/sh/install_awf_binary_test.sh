#!/bin/bash
set +o histexpand

# Test script for install_awf_binary.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_PATH="$SCRIPT_DIR/install_awf_binary.sh"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

print_result() {
  local test_name="$1"
  local result="$2"

  TESTS_RUN=$((TESTS_RUN + 1))

  if [ "$result" = "PASS" ]; then
    echo -e "${GREEN}✓ PASS${NC}: $test_name"
    TESTS_PASSED=$((TESTS_PASSED + 1))
  else
    echo -e "${RED}✗ FAIL${NC}: $test_name"
    TESTS_FAILED=$((TESTS_FAILED + 1))
  fi
}

test_script_syntax() {
  echo ""
  echo "Test 1: Verify script syntax"

  if bash -n "$SCRIPT_PATH" 2>/dev/null; then
    print_result "Script syntax is valid" "PASS"
  else
    print_result "Script has syntax errors" "FAIL"
  fi
}

test_fallback_on_missing_release() {
  echo ""
  echo "Test 2: Fallback to latest release when pinned checksums return 404"

  local tmpdir fakebin output
  tmpdir=$(mktemp -d)
  fakebin="$tmpdir/fakebin"
  output="$tmpdir/output.log"
  mkdir -p "$fakebin"

  cat > "$fakebin/uname" <<'EOF'
#!/bin/bash
if [ "$1" = "-s" ]; then
  echo "Plan9"
elif [ "$1" = "-m" ]; then
  echo "x86_64"
else
  /usr/bin/uname "$@"
fi
EOF

  cat > "$fakebin/node" <<'EOF'
#!/bin/bash
echo "v18.20.0"
EOF

  cat > "$fakebin/curl" <<'EOF'
#!/bin/bash
output_file=""
write_out=""
args=("$@")

while [ $# -gt 0 ]; do
  case "$1" in
    -o)
      output_file="$2"
      shift 2
      ;;
    -w)
      write_out="$2"
      shift 2
      ;;
    *)
      shift
      ;;
  esac
done

url="${args[${#args[@]}-1]}"

case "$url" in
  *"/releases/download/v0.25.28/checksums.txt")
    if [ "$write_out" = "%{http_code}" ]; then
      printf "404"
      exit 0
    fi
    exit 22
    ;;
  *"/releases/latest")
    if [ "$write_out" = "%{url_effective}" ]; then
      printf "https://github.com/github/gh-aw-firewall/releases/tag/v0.25.40"
      exit 0
    fi
    exit 0
    ;;
  *"/releases/download/v0.25.40/checksums.txt")
    if [ -n "$output_file" ] && [ "$output_file" != "/dev/null" ]; then
      cat > "$output_file" <<'CHECKSUMS'
dummy awf-bundle.js
dummy awf-linux-x64
CHECKSUMS
    fi
    if [ "$write_out" = "%{http_code}" ]; then
      printf "200"
    fi
    exit 0
    ;;
  *)
    exit 1
    ;;
esac
EOF

  chmod +x "$fakebin/uname" "$fakebin/node" "$fakebin/curl"

  if PATH="$fakebin:$PATH" bash "$SCRIPT_PATH" "v0.25.28" >"$output" 2>&1; then
    print_result "Script should fail on unsupported OS in test harness" "FAIL"
  elif grep -q "Falling back to latest AWF release: v0.25.40" "$output" && \
       grep -q "Using fallback AWF version: v0.25.40" "$output"; then
    print_result "Script falls back to latest release after 404" "PASS"
  else
    print_result "Script did not use fallback release after 404" "FAIL"
  fi

  rm -rf "$tmpdir"
}

echo "=== Testing install_awf_binary.sh ==="
echo "Script: $SCRIPT_PATH"

test_script_syntax
test_fallback_on_missing_release

echo ""
echo "=== Test Summary ==="
echo "Tests run: $TESTS_RUN"
echo -e "${GREEN}Tests passed: $TESTS_PASSED${NC}"
if [ $TESTS_FAILED -gt 0 ]; then
  echo -e "${RED}Tests failed: $TESTS_FAILED${NC}"
  exit 1
else
  echo -e "${GREEN}All tests passed!${NC}"
  exit 0
fi
