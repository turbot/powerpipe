load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"
source "$FILE_PATH/snapshot_delete.sh"

# Tests for tests/acceptance/snapshot_delete.sh, which the snapshot tests use to delete the
# snapshot they publish.
#
# No network and no Pipes token: a stub `curl` first on PATH plays the Pipes API. It logs its
# arguments and standard input in $STUB_DIR/curl.log and answers with the HTTP status in
# $STUB_DIR/status (204 when the file is missing). An empty status file is a request that
# gets no answer: curl prints 000.

PRINTED="https://pipes.turbot.com/org/turbot-ops/workspace/clitesting/powerpipe/snapshot/snap_d4k8s2v8p1q3x7n0a5bg_1a2b3c4d5e6f7g8h9i0j"
API="https://pipes.turbot.com/api/v0/org/turbot-ops/workspace/clitesting/snapshot/snap_d4k8s2v8p1q3x7n0a5bg_1a2b3c4d5e6f7g8h9i0j"

function setup() {
  STUB_DIR="$(mktemp -d)"
  export STUB_DIR
  : >"$STUB_DIR/curl.log"
  mkdir -p "$STUB_DIR/bin"
  cat >"$STUB_DIR/bin/curl" <<'STUB'
#!/bin/bash
echo "ARGS $*" >>"$STUB_DIR/curl.log"
if [[ " $* " == *" -K - "* ]]; then
  echo "STDIN $(cat)" >>"$STUB_DIR/curl.log"
fi
status=204
[[ -f "$STUB_DIR/status" ]] && status="$(cat "$STUB_DIR/status")"
[[ -n "$status" ]] || status=000
printf '%s' "$status"
STUB
  chmod +x "$STUB_DIR/bin/curl"
  PATH="$STUB_DIR/bin:$PATH"
  export SPIPETOOLS_TOKEN="test-token"
}

function teardown() {
  rm -rf "$STUB_DIR"
}

@test "snapshot delete url - drops the powerpipe segment and goes through /api/v0 to the snapshot path" {
  run snapshot_delete_url "$PRINTED"
  assert_success
  assert_output "$API"
}

@test "snapshot delete url - takes the org and workspace from the printed url" {
  run snapshot_delete_url "https://pipes.turbot.com/org/acme/workspace/ws1/powerpipe/snapshot/snap_x"
  assert_success
  assert_output "https://pipes.turbot.com/api/v0/org/acme/workspace/ws1/snapshot/snap_x"
}

@test "snapshot delete url - keeps a host with a port and ignores a query string and a trailing slash" {
  run snapshot_delete_url "http://localhost:3000/org/acme/workspace/ws1/powerpipe/snapshot/snap_x/?a=b"
  assert_success
  assert_output "http://localhost:3000/api/v0/org/acme/workspace/ws1/snapshot/snap_x"
}

@test "snapshot delete url - a url that is not a snapshot url fails and says so" {
  run snapshot_delete_url "https://pipes.turbot.com/org/turbot-ops/workspace/clitesting"
  assert_failure
  assert_output --partial "cannot build a snapshot delete URL"
}

@test "snapshot delete url - an empty url fails" {
  run snapshot_delete_url ""
  assert_failure
}

@test "snapshot delete - sends a DELETE to the api snapshot path with the token on stdin" {
  run snapshot_delete "$PRINTED"
  assert_success
  run cat "$STUB_DIR/curl.log"
  assert_output --partial "-X DELETE"
  assert_output --partial "$API"
  assert_output --partial "STDIN header = \"Authorization: Bearer test-token\""
  # the token is not an argument, so it is not in a process list
  run grep '^ARGS.*test-token' "$STUB_DIR/curl.log"
  assert_failure
}

@test "snapshot delete - any 2xx passes" {
  echo 200 >"$STUB_DIR/status"
  run snapshot_delete "$PRINTED"
  assert_success
}

@test "snapshot delete - a 404 fails and names the url and the status" {
  echo 404 >"$STUB_DIR/status"
  run snapshot_delete "$PRINTED"
  assert_failure
  assert_output --partial "DELETE $API answered HTTP 404"
}

@test "snapshot delete - a 429 fails and names the url and the status" {
  echo 429 >"$STUB_DIR/status"
  run snapshot_delete "$PRINTED"
  assert_failure
  assert_output --partial "DELETE $API answered HTTP 429"
}

@test "snapshot delete - no answer from the api fails" {
  : >"$STUB_DIR/status"
  run snapshot_delete "$PRINTED"
  assert_failure
  assert_output --partial "DELETE $API answered HTTP 000"
}

@test "snapshot delete - a url it cannot build sends no request" {
  run snapshot_delete "https://pipes.turbot.com/somewhere/else"
  assert_failure
  run cat "$STUB_DIR/curl.log"
  assert_output ""
}

@test "snapshot delete - called last in a bats teardown it fails the test" {
  # the snapshot tests delete in teardown: prove a failed delete there fails the test
  # built with printf: a line starting @test inside a heredoc would be read as a test of this file
  {
    printf 'source "%s/snapshot_delete.sh"\n' "$FILE_PATH"
    printf 'function teardown() { snapshot_delete "%s"; }\n' "$PRINTED"
    printf '%s\n' '@test "inner" { true; }'
  } >"$STUB_DIR/inner.bats"
  echo 404 >"$STUB_DIR/status"
  run bats "$STUB_DIR/inner.bats"
  assert_failure
  assert_output --partial "answered HTTP 404"
}
