#!/bin/bash
# Runs the acceptance tests in Docker, as CI does, without changing this machine.
#
#   tests/acceptance/run-docker.sh                       # the test files that CI runs
#   tests/acceptance/run-docker.sh check.bats mod.bats   # some test files
#
# Each test file runs in its own container, with its own Steampipe service, on a copy of this
# checkout: the files git tracks, with their uncommitted changes, and new files that git does
# not ignore. (run-local.sh instead stops every Steampipe service on the machine.)
#
# ACCEPTANCE_JOBS     how many test files run at the same time (default: AGENT_FLEET_CPU_SHARE, or 4)
# ACCEPTANCE_LOG_DIR  where the log of each test file goes (default: a new temporary directory)
# SPIPETOOLS_TOKEN    passed to the tests if set, as in CI; without it the tests that use
#                     Turbot Pipes skip
set -euo pipefail

root=$(cd "$(dirname "$0")/../.." && pwd)
workflow=$root/.github/workflows/11-test-acceptance.yaml

if [ ! -x "$root/tests/acceptance/lib/bats-core/bin/bats" ]; then
  echo "The bats submodules are missing. Run: git submodule update --init" >&2
  exit 2
fi

if [ $# -gt 0 ]; then
  files=("$@")
else
  # the test blocks of the CI matrix: a block that is commented out there does not run here either,
  # and a line that is not a comment or a block stops the run, so that the two lists cannot drift apart
  block_re="^[[:space:]]*-[[:space:]]*(\"([A-Za-z0-9_-]+)\"|'([A-Za-z0-9_-]+)'|([A-Za-z0-9_-]+))[[:space:]]*(#.*)?$"
  files=()
  while IFS= read -r line; do
    if [[ $line =~ ^[[:space:]]*(#.*)?$ ]]; then
      continue
    fi
    if [[ ! $line =~ $block_re ]]; then
      echo "Cannot read this line of the test_block matrix in $workflow:" >&2
      echo "$line" >&2
      exit 2
    fi
    files+=("${BASH_REMATCH[2]}${BASH_REMATCH[3]}${BASH_REMATCH[4]}.bats")
  done < <(sed -n '/test_block:/,/runs-on:/p' "$workflow" | sed '1d;$d')
  if [ ${#files[@]} -eq 0 ]; then
    echo "Found no test blocks in $workflow" >&2
    exit 2
  fi
fi
for file in "${files[@]}"; do
  if [ ! -f "$root/tests/acceptance/test_files/$file" ]; then
    echo "No such test file: tests/acceptance/test_files/$file" >&2
    exit 2
  fi
done

jobs=${ACCEPTANCE_JOBS:-${AGENT_FLEET_CPU_SHARE:-4}}
logs=${ACCEPTANCE_LOG_DIR:-$(mktemp -d "${TMPDIR:-/tmp}/powerpipe-acceptance.XXXXXX")}
mkdir -p "$logs"
logs=$(cd "$logs" && pwd)
run_id=powerpipe-acceptance-$$-$RANDOM

cleanup() {
  # remove the containers until none is left: a container that `docker run --rm` is still
  # removing keeps the image in use
  for _ in 1 2 3 4 5 6 7 8 9 10; do
    containers=$(docker ps -aq --filter "label=$run_id")
    if [ -z "$containers" ]; then
      break
    fi
    # shellcheck disable=SC2086 # one container ID per word
    docker rm -f $containers >/dev/null 2>&1 || true
    sleep 1
  done
  docker image rm "$run_id" >/dev/null 2>&1 || true
}
trap cleanup EXIT

# on Ctrl-C or kill, stop xargs and the test files it started, so that no test file starts after
# the cleanup; the cleanup then removes their containers
stop() {
  trap - INT TERM
  if [ -n "${runner:-}" ]; then
    # shellcheck disable=SC2046 # one process ID per word
    kill "$runner" $(pgrep -P "$runner" 2>/dev/null) 2>/dev/null || true
    # quietly: bash reports a background job that a signal ends
    wait "$runner" 2>/dev/null || true
  fi
  exit "$1"
}
trap 'stop 130' INT
trap 'stop 143' TERM

echo "Building the test image from $root"
cd "$root"
if ! git ls-files -z --cached --others --exclude-standard --deduplicate \
  | while IFS= read -r -d '' path; do if [ -e "$path" ]; then printf '%s\0' "$path"; fi; done \
  | COPYFILE_DISABLE=1 tar --exclude .git --null --files-from - -cf - \
  | docker build --tag "$run_id" --file tests/acceptance/Dockerfile - >"$logs/build.log" 2>&1; then
  tail -n 40 "$logs/build.log" >&2
  echo "The test image did not build. Log: $logs/build.log" >&2
  exit 1
fi

# runs one test file in its own container, and records its exit code
run_file() {
  local file=$1 code=0
  docker run --rm --label "$run_id" --name "$run_id-${file%.bats}" --env SPIPETOOLS_TOKEN "$run_id" \
    bash -c 'steampipe service start >/dev/null && timeout 15m tests/acceptance/run.sh "$1"' _ "$file" \
    >"$logs/${file%.bats}.log" 2>&1 || code=$?
  echo "$code" >"$logs/${file%.bats}.code"
  echo "finished $file (exit $code)"
}
export -f run_file
export run_id logs

echo "Running ${#files[@]} test files, $jobs at a time. Logs: $logs"
# in the background, so that a signal stops the script at once (bash runs a trap only after the
# command in the foreground ends)
# shellcheck disable=SC2016 # $1 is the argument of the inner bash
printf '%s\n' "${files[@]}" | xargs -n 1 -P "$jobs" bash -c 'run_file "$1"' _ &
runner=$!
wait "$runner" || true
runner=

failed=0
echo
for file in "${files[@]}"; do
  name=${file%.bats}
  code=$(cat "$logs/$name.code" 2>/dev/null || echo "none")
  ok=$(grep -c '^ok ' "$logs/$name.log" || true)
  skipped=$(grep -c '^ok .* # skip' "$logs/$name.log" || true)
  not_ok=$(grep -c '^not ok ' "$logs/$name.log" || true)
  if [ "$code" = "0" ]; then
    printf 'pass  %-36s %4s passed, %s skipped\n' "$file" "$((ok - skipped))" "$skipped"
  else
    failed=$((failed + 1))
    printf 'FAIL  %-36s %4s passed, %s skipped, %s failed, exit %s\n' "$file" "$((ok - skipped))" "$skipped" "$not_ok" "$code"
    grep '^not ok ' "$logs/$name.log" | sed 's/^/        /' || true
  fi
done
echo
if [ "$failed" -gt 0 ]; then
  echo "$failed of ${#files[@]} test files failed. Logs: $logs"
  exit 1
fi
echo "All ${#files[@]} test files passed. Logs: $logs"
