load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# Tests for tests/acceptance/cleanup_ci_snapshots.sh, the script the scheduled cleanup
# workflow runs to delete snapshots the snapshot tests leave in the CI workspace.
#
# No network and no Pipes token: a stub `curl` first on PATH plays the Pipes API. It
# serves the list from page1.json / page2.json in $STUB_DIR, answers each DELETE with the
# status in $STUB_DIR/delete_status.<id> (204 unless the id has a file there), and logs every
# request in $STUB_DIR/calls.log. A file named seq.<id> (or seq.list for the list) holds one
# status per line: each request takes the first line off it, and once it is empty the usual
# status applies. A 429 carries a Retry-After header when $STUB_DIR/retry_after exists. A stub
# `sleep` records the seconds it was asked for in $STUB_DIR/sleep.log and does not wait.

function ago() {
  # ago <hours>: an ISO 8601 UTC time that many hours in the past (GNU or BSD date)
  date -u -d "$1 hours ago" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v-"$1"H +%Y-%m-%dT%H:%M:%SZ
}

function setup() {
  STUB_DIR="$(mktemp -d)"
  export STUB_DIR
  : >"$STUB_DIR/calls.log"
  : >"$STUB_DIR/stdin.log"
  : >"$STUB_DIR/sleep.log"
  echo 200 >"$STUB_DIR/list_status"

  mkdir -p "$STUB_DIR/bin"
  cat >"$STUB_DIR/bin/curl" <<'STUB'
#!/bin/bash
# stand-in for curl: parses just the flags cleanup_ci_snapshots.sh uses
method=GET
out=""
hdr=""
url=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    -X) method="$2"; shift 2 ;;
    -o) out="$2"; shift 2 ;;
    -D) hdr="$2"; shift 2 ;;
    -w | --max-time | --config) shift 2 ;;
    -sS) shift ;;
    *) url="$1"; shift ;;
  esac
done
cat >>"$STUB_DIR/stdin.log"
echo "$method $url" >>"$STUB_DIR/calls.log"

# next_status <seq file> <default>: the first queued status, else the default
next_status() {
  if [[ -s "$1" ]]; then
    head -n 1 "$1"
    tail -n +2 "$1" >"$1.rest" && mv "$1.rest" "$1"
  else
    cat "$2"
  fi
}
respond_headers() {
  [[ -n "$hdr" ]] || return 0
  echo "HTTP/2 $1" >"$hdr"
  if [[ "$1" == 429 && -f "$STUB_DIR/retry_after" ]]; then
    printf 'Retry-After: %s\r\n' "$(cat "$STUB_DIR/retry_after")" >>"$hdr"
  fi
}

if [[ "$method" == "GET" ]]; then
  status=$(next_status "$STUB_DIR/seq.list" "$STUB_DIR/list_status")
  respond_headers "$status"
  if [[ "$status" == 429 ]]; then
    echo '{"error":{"title":"Rate limit exceeded."}}' >"$out"
  elif [[ "$status" != 200 ]]; then
    echo '{"error":{"title":"Gateway Timeout"}}' >"$out"
  elif [[ "$url" == *next_token=* ]]; then
    cat "$STUB_DIR/page2.json" >"$out"
  else
    cat "$STUB_DIR/page1.json" >"$out"
  fi
  printf '%s' "$status"
else
  id="${url##*/}"
  [[ -f "$STUB_DIR/delete_status.$id" ]] || echo 204 >"$STUB_DIR/delete_status.$id"
  status=$(next_status "$STUB_DIR/seq.$id" "$STUB_DIR/delete_status.$id")
  respond_headers "$status"
  echo '{}' >"$out"
  printf '%s' "$status"
fi
STUB
  chmod +x "$STUB_DIR/bin/curl"

  cat >"$STUB_DIR/bin/sleep" <<'STUB'
#!/bin/bash
echo "$1" >>"$STUB_DIR/sleep.log"
STUB
  chmod +x "$STUB_DIR/bin/sleep"

  export PATH="$STUB_DIR/bin:$PATH"
  export PIPES_TOKEN="tpt_stub_secret_value"
  export PIPES_API="https://pipes.example.test/api/v0"
  export MAX_AGE_HOURS=24
  unset DRY_RUN CLEANUP_ORG CLEANUP_WORKSPACE

  # page 1: one old, one recent, one old but already deleted; points on to page 2
  # page 2: one old that is missing created_at, one old with fractional seconds, one old
  cat >"$STUB_DIR/page1.json" <<EOF
{"items":[
  {"id":"snap_old1","created_at":"$(ago 48)","state":"available"},
  {"id":"snap_new1","created_at":"$(ago 2)","state":"available"},
  {"id":"snap_gone","created_at":"$(ago 72)","state":"deleted"}
],"next_token":"tok2"}
EOF
  cat >"$STUB_DIR/page2.json" <<EOF
{"items":[
  {"id":"snap_nodate","state":"available"},
  {"id":"snap_old2","created_at":"$(ago 30 | sed 's/Z$/.123456Z/')","state":"available"},
  {"id":"snap_old3","created_at":"$(ago 100)"}
]}
EOF
}

function teardown() {
  rm -rf "$STUB_DIR"
}

@test "snapshot cleanup - deletes every old snapshot across pages and keeps the rest" {
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_success

  run grep '^DELETE' "$STUB_DIR/calls.log"
  assert_output "DELETE https://pipes.example.test/api/v0/org/turbot-ops/workspace/clitesting/snapshot/snap_old1
DELETE https://pipes.example.test/api/v0/org/turbot-ops/workspace/clitesting/snapshot/snap_old2
DELETE https://pipes.example.test/api/v0/org/turbot-ops/workspace/clitesting/snapshot/snap_old3"

  # the second page was asked for with the first page's token
  run grep -c 'next_token=tok2' "$STUB_DIR/calls.log"
  assert_output "1"
}

@test "snapshot cleanup - sends the token in a header and never prints it" {
  run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_success
  refute_output --partial "tpt_stub_secret_value"

  run grep -c 'Authorization: Bearer tpt_stub_secret_value' "$STUB_DIR/stdin.log"
  assert_output "5"

  # the token never reaches curl's own arguments
  run grep 'tpt_stub_secret_value' "$STUB_DIR/calls.log"
  assert_failure
}

@test "snapshot cleanup - a delete that returns 404 counts as done" {
  echo 404 >"$STUB_DIR/delete_status.snap_old1"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_success
  assert_output --partial "Deleted 3 snapshots, 0 failed"
}

@test "snapshot cleanup - one failed delete still deletes the others, then fails" {
  echo 500 >"$STUB_DIR/delete_status.snap_old1"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_failure
  assert_output --partial "Delete of snap_old1 returned HTTP 500"
  assert_output --partial "Deleted 2 snapshots, 1 failed"

  run grep -c '^DELETE' "$STUB_DIR/calls.log"
  assert_output "3"
}

@test "snapshot cleanup - a failed list deletes nothing and fails" {
  echo 504 >"$STUB_DIR/list_status"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_failure
  assert_output --partial "List request returned HTTP 504"

  run grep -c '^DELETE' "$STUB_DIR/calls.log"
  assert_output "0"
}

@test "snapshot cleanup - a list answering the second page with junk deletes nothing" {
  echo '<html>bad gateway</html>' >"$STUB_DIR/page2.json"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_failure
  assert_output --partial "List response has no items array"

  run grep -c '^DELETE' "$STUB_DIR/calls.log"
  assert_output "0"
}

@test "snapshot cleanup - a workspace with no snapshots is a clean pass" {
  # the API may leave items out when the list is empty
  for body in '{}' '{"items":null}' '{"items":[]}'; do
    echo "$body" >"$STUB_DIR/page1.json"
    run $FILE_PATH/cleanup_ci_snapshots.sh
    echo "$output"
    assert_success
    assert_output --partial "Found 0 snapshots, 0 older than 24 hours"
  done

  run grep -c '^DELETE' "$STUB_DIR/calls.log"
  assert_output "0"
}

@test "snapshot cleanup - the default cutoff is one hour, so a run's leftovers go at the next cleanup" {
  unset MAX_AGE_HOURS
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_success
  assert_output --partial "older than 1 hours"

  # snap_new1 is two hours old: gone with the default, kept at the 24 hour cutoff above
  run grep -c 'snap_new1' "$STUB_DIR/calls.log"
  assert_output "1"
}

@test "snapshot cleanup - DRY_RUN lists what it would delete and deletes nothing" {
  DRY_RUN=true run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_success
  assert_output --partial "Would delete snap_old1"
  assert_output --partial "Dry run: nothing deleted"

  run grep -c '^DELETE' "$STUB_DIR/calls.log"
  assert_output "0"
}

@test "snapshot cleanup - MAX_AGE_HOURS moves the cutoff" {
  # only the 72 and 100 hour old snapshots are older than 60 hours (snap_gone is already deleted)
  MAX_AGE_HOURS=60 run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_success
  assert_output --partial "Deleted 1 snapshots, 0 failed"

  run grep '^DELETE' "$STUB_DIR/calls.log"
  assert_output "DELETE https://pipes.example.test/api/v0/org/turbot-ops/workspace/clitesting/snapshot/snap_old3"
}

@test "snapshot cleanup - no token fails before any request" {
  unset PIPES_TOKEN
  run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_failure
  assert_output --partial "PIPES_TOKEN is not set"

  run wc -c <"$STUB_DIR/calls.log"
  assert_output --partial "0"
}

@test "snapshot cleanup - a bad MAX_AGE_HOURS fails before any request" {
  MAX_AGE_HOURS=0 run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_failure
  assert_output --partial "MAX_AGE_HOURS must be a whole number"

  MAX_AGE_HOURS=soon run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_failure
}

@test "snapshot cleanup - an id that is not a plain handle is skipped, not put in a URL" {
  cat >"$STUB_DIR/page2.json" <<EOF
{"items":[{"id":"../../other","created_at":"$(ago 48)"}]}
EOF
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_failure
  assert_output --partial "Skipping snapshot with an unexpected id"

  run grep -c 'other' "$STUB_DIR/calls.log"
  assert_output "0"
}

@test "snapshot cleanup - a delete answered 429 and then 204 counts as deleted" {
  printf '429\n' >"$STUB_DIR/seq.snap_old1"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_success
  assert_output --partial "Deleted 3 snapshots, 0 failed"

  # snap_old1 was sent twice, the other two once
  run grep -c '^DELETE' "$STUB_DIR/calls.log"
  assert_output "4"
  run grep -c 'snap_old1$' "$STUB_DIR/calls.log"
  assert_output "2"

  # no Retry-After header, so the first back-off of 2 seconds
  run cat "$STUB_DIR/sleep.log"
  assert_output "2"
}

@test "snapshot cleanup - a delete answered 429 on every try fails after five tries and the rest still run" {
  echo 429 >"$STUB_DIR/delete_status.snap_old1"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_failure
  assert_output --partial "Delete of snap_old1 returned HTTP 429"
  assert_output --partial "Deleted 2 snapshots, 1 failed"

  run grep -c 'snap_old1$' "$STUB_DIR/calls.log"
  assert_output "5"
  run grep -c '^DELETE' "$STUB_DIR/calls.log"
  assert_output "7"

  # four waits between five tries, doubling from 2 seconds, and none after the last
  run cat "$STUB_DIR/sleep.log"
  assert_output "2
4
8
16"
}

@test "snapshot cleanup - a Retry-After header sets the wait, up to 60 seconds" {
  echo 7 >"$STUB_DIR/retry_after"
  printf '429\n' >"$STUB_DIR/seq.snap_old1"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_success
  run cat "$STUB_DIR/sleep.log"
  assert_output "7"

  # an hour is cut to 60 seconds
  : >"$STUB_DIR/sleep.log"
  echo 3600 >"$STUB_DIR/retry_after"
  printf '429\n' >"$STUB_DIR/seq.snap_old2"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_success
  run cat "$STUB_DIR/sleep.log"
  assert_output "60"
}

@test "snapshot cleanup - a Retry-After of 0 still waits one second" {
  echo 0 >"$STUB_DIR/retry_after"
  printf '429\n429\n' >"$STUB_DIR/seq.snap_old1"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_success
  assert_output --partial "Deleted 3 snapshots, 0 failed"

  # two 429s, each waited out for one second (not zero), then the delete lands
  run cat "$STUB_DIR/sleep.log"
  assert_output "1
1"
  run grep -c 'snap_old1$' "$STUB_DIR/calls.log"
  assert_output "3"

  # a zero written with leading zeros counts as 0 too
  : >"$STUB_DIR/sleep.log"
  echo 00 >"$STUB_DIR/retry_after"
  printf '429\n' >"$STUB_DIR/seq.snap_old2"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_success
  run cat "$STUB_DIR/sleep.log"
  assert_output "1"
}

@test "snapshot cleanup - a Retry-After that is not a number falls back to the back-off" {
  echo "Wed, 21 Oct 2026 07:28:00 GMT" >"$STUB_DIR/retry_after"
  printf '429\n429\n' >"$STUB_DIR/seq.snap_old1"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  assert_success
  run cat "$STUB_DIR/sleep.log"
  assert_output "2
4"
}

@test "snapshot cleanup - a list answered 429 and then 200 lists as normal" {
  printf '429\n' >"$STUB_DIR/seq.list"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_success
  assert_output --partial "Found 6 snapshots, 3 older than 24 hours"
  assert_output --partial "Deleted 3 snapshots, 0 failed"

  run cat "$STUB_DIR/sleep.log"
  assert_output "2"
}

@test "snapshot cleanup - a list answered 429 on every try deletes nothing and fails" {
  echo 429 >"$STUB_DIR/list_status"
  run $FILE_PATH/cleanup_ci_snapshots.sh
  echo "$output"
  assert_failure
  assert_output --partial "List request returned HTTP 429"

  run grep -c '^GET' "$STUB_DIR/calls.log"
  assert_output "5"
  run grep -c '^DELETE' "$STUB_DIR/calls.log"
  assert_output "0"
}
