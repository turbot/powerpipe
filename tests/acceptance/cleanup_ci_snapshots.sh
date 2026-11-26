#!/bin/bash
#
# Deletes the snapshots the acceptance tests leave in the CI workspace in Pipes.
#
# The snapshot tests delete their own snapshot in teardown, but a create that times out at
# the API gateway still lands in the workspace and gives the client no id, so no test can
# delete it. A scheduled workflow runs this script to remove every snapshot older than
# MAX_AGE_HOURS, whatever left it behind. A test run uses its snapshot for seconds, so the
# default of one hour never takes one from a run in progress (and a test's own delete
# tolerates a snapshot that is already gone).
#
# Environment:
#   PIPES_TOKEN       required. A Pipes token that can list and delete snapshots in the workspace.
#   CLEANUP_ORG       org handle of the CI workspace (default turbot-ops)
#   CLEANUP_WORKSPACE workspace handle of the CI workspace (default clitesting)
#   MAX_AGE_HOURS     delete snapshots created more than this many hours ago (default 1)
#   DRY_RUN           when "true", list what would be deleted and delete nothing
#   PIPES_API         API base URL (default https://pipes.turbot.com/api/v0)
#
# A request that Pipes answers with HTTP 429 is sent again, up to five tries in all. The wait
# before the next try is the Retry-After header's number of seconds when the answer has one,
# else 2, 4, 8 then 16 seconds, and never more than 60. A delete counts as failed only when its
# tries run out.
#
# Exit status: 0 when every old snapshot is gone, 1 when the list or a delete failed.

set -u

org="${CLEANUP_ORG:-turbot-ops}"
workspace="${CLEANUP_WORKSPACE:-clitesting}"
max_age_hours="${MAX_AGE_HOURS:-1}"
dry_run="${DRY_RUN:-false}"
api="${PIPES_API:-https://pipes.turbot.com/api/v0}"
page_size=100
max_pages=500
max_tries=5
max_wait=60

if [[ -z "${PIPES_TOKEN:-}" ]]; then
  echo "PIPES_TOKEN is not set" >&2
  exit 1
fi

if ! [[ "$max_age_hours" =~ ^[0-9]+$ ]] || [[ "$max_age_hours" -lt 1 ]]; then
  echo "MAX_AGE_HOURS must be a whole number of hours, at least 1 (got '$max_age_hours')" >&2
  exit 1
fi

for tool in curl jq; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "$tool is required" >&2
    exit 1
  fi
done

base="$api/org/$org/workspace/$workspace/snapshot"
cutoff=$(( $(date -u +%s) - max_age_hours * 3600 ))

# The seconds to wait before the next try after a 429: the Retry-After header's whole number
# when the answer has one (an HTTP date is not read), else a back-off that doubles from 2
# seconds, and never more than max_wait.
retry_wait() {
  local attempt="$1" wait="" value
  value=$(grep -i '^retry-after:' "$header_file" 2>/dev/null | tail -n 1 | cut -d: -f2- | tr -d '[:space:]')
  if [[ "$value" =~ ^[0-9]+$ ]]; then
    wait=$((10#$value))
  else
    wait=$((2 << (attempt - 1)))
  fi
  [[ "$wait" -le "$max_wait" ]] || wait="$max_wait"
  echo "$wait"
}

# Sends one request and prints its HTTP status. An HTTP 429 is sent again after a wait, up to
# max_tries tries; the status printed is the last answer's. The token goes in a curl config
# read from stdin, so it never shows in a process list.
api_call() {
  local method="$1" url="$2" attempt=1 status wait
  while :; do
    : >"$header_file"
    status=$(printf 'header = "Authorization: Bearer %s"\n' "$PIPES_TOKEN" |
      curl -sS --max-time 120 --config - -X "$method" -D "$header_file" -o "$body_file" -w '%{http_code}' "$url") || return $?
    if [[ "$status" != 429 || "$attempt" -ge "$max_tries" ]]; then
      echo "$status"
      return 0
    fi
    wait=$(retry_wait "$attempt")
    echo "HTTP 429 from $method ${url##*/}, waiting ${wait}s before try $((attempt + 1)) of $max_tries" >&2
    sleep "$wait"
    attempt=$((attempt + 1))
  done
}

body_file=$(mktemp)
header_file=$(mktemp)
ids_file=$(mktemp)
trap 'rm -f "$body_file" "$header_file" "$ids_file"' EXIT

echo "Listing snapshots in $org/$workspace older than $max_age_hours hours"

# Collect every page before deleting anything, so a failed list deletes nothing half-way
# and deletes do not shift the pages being read.
next_token=""
pages=0
total=0
while :; do
  pages=$((pages + 1))
  if [[ "$pages" -gt "$max_pages" ]]; then
    echo "Gave up after $max_pages pages of snapshots" >&2
    exit 1
  fi

  url="$base?limit=$page_size"
  if [[ -n "$next_token" ]]; then
    url="$url&next_token=$(jq -rn --arg t "$next_token" '$t|@uri')"
  fi

  status=$(api_call GET "$url") || { echo "List request failed" >&2; exit 1; }
  if [[ "$status" != 200 ]]; then
    echo "List request returned HTTP $status: $(head -c 300 "$body_file")" >&2
    exit 1
  fi

  jq -e '(.items // []) | type == "array"' "$body_file" >/dev/null 2>&1 || {
    echo "List response has no items array: $(head -c 300 "$body_file")" >&2
    exit 1
  }

  total=$((total + $(jq '(.items // []) | length' "$body_file")))

  # created_at is ISO 8601 UTC, possibly with fractional seconds. A snapshot with no
  # readable created_at is not old by any measure, so it is kept.
  jq -r --argjson cutoff "$cutoff" '
    (.items // [])[]
    | select((.state // "available") != "deleted")
    | select(
        ((.created_at // "") | sub("\\.[0-9]+"; "") | try fromdateiso8601 catch null) as $t
        | $t != null and $t < $cutoff
      )
    | .id' "$body_file" >>"$ids_file"

  next_token=$(jq -r '.next_token // ""' "$body_file")
  [[ -n "$next_token" ]] || break
done

old=$(wc -l <"$ids_file" | tr -d ' ')
echo "Found $total snapshots, $old older than $max_age_hours hours"

failed=0
deleted=0
while IFS= read -r id; do
  [[ -n "$id" ]] || continue
  # An id from the API is a handle like snap_abc123; refuse anything else in a URL.
  if ! [[ "$id" =~ ^[A-Za-z0-9_-]+$ ]]; then
    echo "Skipping snapshot with an unexpected id: $id" >&2
    failed=$((failed + 1))
    continue
  fi
  if [[ "$dry_run" == "true" ]]; then
    echo "Would delete $id"
    continue
  fi
  status=$(api_call DELETE "$base/$id") || status="000"
  case "$status" in
    # 404: already gone, which is what we wanted
    200 | 202 | 204 | 404) deleted=$((deleted + 1)) ;;
    *)
      echo "Delete of $id returned HTTP $status: $(head -c 300 "$body_file")" >&2
      failed=$((failed + 1))
      ;;
  esac
done <"$ids_file"

if [[ "$dry_run" == "true" ]]; then
  echo "Dry run: nothing deleted"
else
  echo "Deleted $deleted snapshots, $failed failed"
fi

[[ "$failed" -eq 0 ]]
