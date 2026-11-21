load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# Each test starts `powerpipe server` on its own copy of server_reload_mod, with an empty SQLite database, on a
# free port, changes the files the way a client does, and checks the load state at GET /api/v0/workspace.

WORKSPACE_LOAD="python3 -I $TEST_SCRIPTS_DIR/workspace_load.py"

function setup() {
  WS="$(mktemp -d)"
  cp -R "$MODS_DIR/server_reload_mod/." "$WS"
  touch "$WS/test.db"

  # an install folder of its own, with the database connection and a stand-in for the dashboard UI: these
  # tests use only the API, and the server starts only when the UI it finds has its version
  export POWERPIPE_INSTALL_DIR="$WS-install"
  mkdir -p "$POWERPIPE_INSTALL_DIR/config" "$POWERPIPE_INSTALL_DIR/dashboard/assets"
  cat << EOF > "$POWERPIPE_INSTALL_DIR/config/server_reload.ppc"
connection "sqlite" "server_reload" {
  file_name = "$WS/test.db"
}
EOF
  echo '<html></html>' > "$POWERPIPE_INSTALL_DIR/dashboard/assets/index.html"
  printf '{"version":"%s"}' "$(powerpipe --version | awk '{print $2}')" > "$POWERPIPE_INSTALL_DIR/dashboard/assets/versions.json"

  cd "$WS"
  PORT="$($WORKSPACE_LOAD free-port)"
  BASE="http://localhost:$PORT"
  powerpipe server --port "$PORT" > "$WS/server.log" 2>&1 &
  SERVER_PID=$!
  $WORKSPACE_LOAD wait-up "$BASE"
}

function teardown() {
  kill "$SERVER_PID" 2>/dev/null || true
  wait "$SERVER_PID" 2>/dev/null || true
  cd /
  rm -rf "$WS" "$WS-dep" "$WS-install"
}

# make_dependency writes a mod with one benchmark to $WS-dep, to install from the local folder
function make_dependency() {
  mkdir -p "$WS-dep"
  cat << EOF > "$WS-dep/mod.pp"
mod "server_reload_dep" {
  title = "Server reload dependency"
}

control "dep" {
  title = "Dependency"
  sql   = "select 'dep' as resource, 'ok' as status, 'from the dependency' as reason"
}

benchmark "dep" {
  title    = "Dependency"
  children = [control.dep]
}
EOF
}

@test "server: a started server has finished its first load" {
  run $WORKSPACE_LOAD state "$BASE"
  assert_success
  assert_equal "$(echo "$output" | jq -r '.loading')" "false"
  assert [ "$(echo "$output" | jq -r '.load_count')" -ge 1 ]
  assert_equal "$(echo "$output" | jq -r '.last_load.error // ""')" ""
}

@test "server: a run with --host uses a variable change once the load state shows it loaded" {
  printf 'who = "second"\n' > powerpipe.ppvars
  run $WORKSPACE_LOAD wait-loaded "$BASE" powerpipe.ppvars
  assert_success
  assert_equal "$(echo "$output" | jq -r '.error // ""')" ""

  run powerpipe benchmark run server_reload.benchmark.who --host "$BASE" --output snapshot
  assert_output --partial "who is second"
}

@test "server: a broken file gives a failed load, and runs keep the workspace as last loaded" {
  printf 'control "broken" {\n  sql = \n}\n' > broken.pp
  run $WORKSPACE_LOAD wait-loaded "$BASE" broken.pp
  assert_success
  assert_output --partial "broken.pp"

  run powerpipe benchmark run server_reload.benchmark.who --host "$BASE" --output snapshot
  assert_output --partial "who is first"

  # deleting the file changes its folder, which the load state counts
  rm broken.pp
  run $WORKSPACE_LOAD wait-loaded "$BASE" .
  assert_success
  assert_equal "$(echo "$output" | jq -r '.error // ""')" ""
}

@test "server: a run request during a reload is answered 503 with Retry-After, and --host waits for the reload" {
  # enough controls that a load takes seconds
  $WORKSPACE_LOAD generate many.pp 8000
  $WORKSPACE_LOAD wait-loaded "$BASE" many.pp

  printf '# a change\n' >> mod.pp
  run $WORKSPACE_LOAD probe-reload "$BASE" server_reload.benchmark.who
  assert_success
  assert_output --regexp '^503 [0-9]+$'

  printf '# another change\n' >> mod.pp
  $WORKSPACE_LOAD probe-reload "$BASE" server_reload.benchmark.who
  run powerpipe benchmark run server_reload.benchmark.who --host "$BASE" --output snapshot
  assert_output --partial "who is first"
}

@test "server: mod install reloads the server, and the load state covers the lock file" {
  make_dependency
  powerpipe mod install "$WS-dep" --mod-location "$WS"
  run $WORKSPACE_LOAD wait-loaded "$BASE" mod.pp .mod.cache.json
  assert_success
  assert_equal "$(echo "$output" | jq -r '.error // ""')" ""

  run powerpipe benchmark run server_reload_dep.benchmark.dep --host "$BASE" --output snapshot
  assert_output --partial "from the dependency"
}

@test "server: an install that changes nothing still lets a client's wait end" {
  make_dependency
  powerpipe mod install "$WS-dep" --mod-location "$WS"
  $WORKSPACE_LOAD wait-loaded "$BASE" mod.pp .mod.cache.json

  # the installer writes the lock file again, with the same content
  powerpipe mod install "$WS-dep" --mod-location "$WS"
  run $WORKSPACE_LOAD wait-loaded "$BASE" mod.pp .mod.cache.json
  assert_success
}

@test "server: mod uninstall reloads the server" {
  make_dependency
  powerpipe mod install "$WS-dep" --mod-location "$WS"
  $WORKSPACE_LOAD wait-loaded "$BASE" mod.pp .mod.cache.json

  powerpipe mod uninstall "$WS-dep" --mod-location "$WS"
  run $WORKSPACE_LOAD wait-loaded "$BASE" mod.pp .mod.cache.json
  assert_success

  run powerpipe benchmark run server_reload_dep.benchmark.dep --host "$BASE" --output snapshot
  assert_failure
}
