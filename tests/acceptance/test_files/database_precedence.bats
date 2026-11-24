load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# These set of tests are skipped locally
# To run these tests locally set the SPIPETOOLS_TOKEN env var.
# These tests will be skipped locally unless the below env var is set.
# The "permitted root paths" tests at the end of the file use only local files,
# so they run without the token.

function setup() {
  if [[ -z "${SPIPETOOLS_TOKEN}" && "${BATS_TEST_DESCRIPTION}" != "permitted root paths"* ]]; then
    skip
  fi
}

function teardown() {
  rm -f "$POWERPIPE_INSTALL_DIR/config/permitted_paths_duckdb.ppc"
}

# no database specified in mod or within mod resources - so the default steampipe connection
# gets the highest precedence
@test "no database specified in mod or within resource" {

  # checkout the mod with no database specified in mod.pp
  cd $MODS_DIR/mod_with_no_db

  # run a powerpipe query to verify that the default steampipe connection is used
  run powerpipe query run query.steampipe_db_query --output csv
  echo $output

  # check output that the defasult steampipe connection is used
  assert_output --partial "cache_enabled"
}

# database(sqlite) specified in mod definition(connection ref) - so the mod level database gets the highest precedence
@test "database specified in mod definition(connection ref)" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  cat $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db

  # run a powerpipe query to verify that the mod level database(sqlite) is used
  run powerpipe query run query.sqlite_db_query --output csv
  echo $output

  # check output that the mod level database is used
  assert_output --partial "Total Albums"
}

# TODO - remove this test once deprecated --database flag is removed
# no database specified in mod or within mod resources
# database(sqlite) specified in --database argument - so that gets the highest precedence
@test "database specified in --database arg(connection string)" {

  # checkout the mod with no database specified in mod.pp
  cd $MODS_DIR/mod_with_no_db

  # run a powerpipe query to verify that the database(sqlite) specified in --database argument is used
  run powerpipe query run query.sqlite_db_query --database sqlite:///$MODS_DIR/sqlite_mod/chinook.db --output csv
  echo $output

  # check output that the database specified in --database argument is used
  assert_output --partial "Total Albums"
}

# TODO - remove this test once deprecated --database flag is removed
# database(sqlite) specified in mod definition
# database(duckdb) specified in --database argument(connection string) - so that gets the highest precedence
@test "database specified in mod definition and --database arg(connection string)" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db

  # run a powerpipe query to verify that the database(duckdb) specified in --database argument is used
  run powerpipe query run query.duckdb_db_query --database duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb --output csv
  echo $output

  # check output that the database specified in --database argument is used
  assert_output --partial "Total Employees"
}

# TODO - remove this test once deprecated --database flag is removed
# database(sqlite) specified in mod definition
# database(duckdb) specified in --database argument(connection ref) - so that gets the highest precedence
@test "database specified in mod definition and --database arg(connection ref)" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  file_name = "$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db

  # run a powerpipe query to verify that the database(duckdb) specified in --database argument is used
  run powerpipe query run query.duckdb_db_query --database connection.duckdb.employees --output csv
  echo $output

  # check output that the database specified in --database argument is used over the one specified in mod definition
  assert_output --partial "Total Employees"
}

# database specified in mod definition through a var
# default value of database var is sqlite(connection ref)
# so the mod level database gets the highest precedence
@test "database specified through variable in mod" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db_var

  # run a powerpipe query to verify that the database specified through variable in mod is used
  run powerpipe query run query.sqlite_db_query --output csv
  echo $output

  # check output that the database specified through default value of variable in mod is used
  assert_output --partial "Total Albums"
}

# database specified in mod definition through a var
# default value of database var is sqlite(connection ref)
# database also specified at runtime through --var argument
# so the database passed through --var gets the highest precedence
@test "database specified through variable in mod and passed through --var arg" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  file_name = "$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db_var

  # run a powerpipe query to verify that the database specified through variable in mod is used
  run powerpipe query run query.duckdb_db_query --output csv --var database=connection.duckdb.employees
  echo $output

  # check output that the database specified through default value of variable in mod is used
  assert_output --partial "Total Employees"
}

# database specified in mod definition through a var
# default value of database var is sqlite(connection ref)
# database also specified at runtime through  powerpipe.ppvars file
# so the database passed through powerpipe.ppvars file gets the highest precedence
@test "database specified through variable in mod and passed through .ppvars file" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  file_name = "$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF

  # write the .ppvars file with the database variable
  cat << EOF > $MODS_DIR/mod_with_db_var/powerpipe.ppvars
database=connection.duckdb.employees
EOF

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db_var

  # run a powerpipe query to verify that the database specified through variable in mod is used
  run powerpipe query run query.duckdb_db_query --output csv
  echo $output

  # check output that the database specified through default value of variable in mod is used
  assert_output --partial "Total Employees"

  # cleanup the .ppvars file
  rm $MODS_DIR/mod_with_db_var/powerpipe.ppvars
}

# database specified in mod require block
@test "database specified in dependency mod require block" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # checkout the mod with database specified in mod.pp dependant mod require block
  cd $MODS_DIR/mod_with_db_in_require

  # run a powerpipe query to verify that the database specified through mod require block is used
  run powerpipe query run mod_with_db_var.query.sqlite_db_query --output csv
  echo $output

  # check output that the database specified through default value of mod require block is used
  assert_output --partial "Total Albums"
}

# database specified in mod require through a var
@test "database specified through variable in dependency mod require block" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  file_name = "$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF

  # checkout the mod with database specified in mod.pp dependant mod require block
  cd $MODS_DIR/mod_with_db_in_require_var

  # run a powerpipe query to verify that the database specified through mod require block is used
  run powerpipe query run mod_with_db_var.query.sqlite_db_query --output csv
  echo $output

  # check output that the database specified through default value of mod require block is used
  assert_output --partial "Total Albums"
}

# database specified in mod - implicit workspace
# database specified in mod definition(implicit workspace) - so the mod level database gets the highest precedence
# and should query the pipes workspace
@test "database specified in mod definition(implicit workspace)" {
  # checkout the mod with implicit workspace database specified in mod.pp
  cd $MODS_DIR/mod_with_db_implicit_workspace

  # run a powerpipe query to verify that the pipes workspace mentioned in mod is used
  run powerpipe query run query.pipes_workspace_query --output csv --pipes-token $SPIPETOOLS_TOKEN
  echo $output

  # check output that the pipes workspace is queried
  assert_output --partial "turbot-silverwater"
}

# database specified in mod definition through a var
# default value of database var is implicit workspace
# so the mod level database gets the highest precedence
@test "database specified through variable(implicit workspace) in mod" {

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db_var_implicit_workspace

  # run a powerpipe query to verify that the pipes workspace specified through variable in mod is used
  run powerpipe query run query.pipes_workspace_query --output csv --pipes-token $SPIPETOOLS_TOKEN
  echo $output

  # check output that the pipes workspace specified through default value of variable in mod is used
  assert_output --partial "turbot-silverwater"
}

# database specified in mod definition through a var
# default value of database var is implicit workspace
# database also specified at runtime through  vars argument
# so the pipes workspace specified in --var gets the highest precedence
@test "database specified through variable in mod and passed through --var(implicit workspace) arg" {

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db_var_implicit_workspace

  # run a powerpipe query to verify that the pipes workspace specified through --var in mod is used
  run powerpipe query run query.pipes_workspace_query --output csv --pipes-token $SPIPETOOLS_TOKEN --var database="turbot-ops/clitesting"
  echo $output

  # check output that the pipes workspace specified through --var in mod is used
  assert_output --partial "turbot-silverwater"
}

# database specified in resource
# database specified in resource - so the resource level database gets the highest precedence
@test "database specified in resource" {
    # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  file_name = "$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF

  # checkout the mod with database specified in resource
  cd $MODS_DIR/mod_with_db_in_resource

  # run a powerpipe query to verify that the database specified in resource is used
  run powerpipe query run query.sqlite_db_query --output csv
  echo $output

  # check output that the database specified in resource is used
  assert_output --partial "Total Albums"

  # run a powerpipe query to verify that the database specified in resource is used
  run powerpipe query run query.duckdb_db_query --output csv
  echo $output

  # check output that the database specified in resource is used
  assert_output --partial "Total Employees"
}

# database specified in mod
# database specified in resource
# so the database specified in resource gets the highest precedence
@test "database specified in mod and also in resource" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  file_name = "$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF

  # checkout the mod with database specified in resource and in mod
  cd $MODS_DIR/mod_with_db_in_mod_and_resource

  # run a powerpipe query to verify that the database specified in resource is used
  run powerpipe query run query.sqlite_db_query --output csv
  echo $output

  # check output that the database specified in resource is used
  assert_output --partial "Total Albums"

  # run a powerpipe query to verify that the database specified in resource is used
  run powerpipe query run query.duckdb_db_query --output csv
  echo $output

  # check output that the database specified in resource is used
  assert_output --partial "Total Employees"
}

# TODO - remove this test once deprecated --database flag is removed
# database specified in mod
# database specified in resource
# so the database specified in resource gets the highest precedence
@test "database specified in mod and also in resource, also passed through --database arg" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  file_name = "$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  file_name = "$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF

  # checkout the mod with database specified in resource and in mod
  cd $MODS_DIR/mod_with_db_in_mod_and_resource

  # run a powerpipe query to verify that the database specified in resource is used
  run powerpipe query run query.sqlite_db_query  --database connection.steampipe.default --output csv
  echo $output

  # check output that the database specified in resource is used
  assert_output --partial "Total Albums"

  # run a powerpipe query to verify that the database specified in resource is used
  run powerpipe query run query.duckdb_db_query --database connection.steampipe.default --output csv
  echo $output

  # check output that the database specified in resource is used
  assert_output --partial "Total Employees"
}

# database specified in mod through connection ref, search_path_prefix specified in connection
# so the search_path_prefix specified in connection gets the highest precedence
@test "database specified in mod through connection ref, search_path_prefix specified in connection" {
  # add the steampipe connection with search_path_prefix defined
  # write the steampipe connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/with_search_path_prefix.ppc
connection "steampipe" "with_search_path_prefix" {
  host     = "localhost"
  port     = 9193
  db       = "steampipe"
  username = "steampipe"
  search_path_prefix = ["foo"]
}
EOF

  # checkout the mod with database(connection ref) specified in mod.pp, connection has search_path_prefix defined
  cd $MODS_DIR/mod_with_db_steampipe_search_path

  # run a powerpipe query to verify that the search_path specified in mod is used
  run powerpipe query run query.search_path --output csv
  echo $output

  # check output that the search_path specified in mod is used
  assert_output --partial "foo"
}

# database specified in mod through connection ref, search_path_prefix specified in connection
# database specified by --var in runtime
# so the search_path_prefix specified in --var gets the highest precedence
@test "database specified in mod through connection ref, search_path_prefix specified in connection, database specified by --var in runtime" {
  # add the steampipe connection with search_path_prefix defined
  # write the steampipe connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/with_search_path_prefix_bar.ppc
connection "steampipe" "with_search_path_prefix_bar" {
  host     = "localhost"
  port     = 9193
  db       = "steampipe"
  username = "steampipe"
  search_path_prefix = ["bar"]
}
EOF

  # checkout the mod with database(connection ref) specified in mod.pp, connection has search_path_prefix defined
  cd $MODS_DIR/mod_with_db_steampipe_search_path

  # run a powerpipe query to verify that the search_path specified in --var is used
  run powerpipe query run query.search_path --output csv --var database=connection.steampipe.with_search_path_prefix_bar
  echo $output

  # check output that the search_path specified in mod is used
  assert_output --partial "bar"
}

# database specified in mod through connection ref, search_path_prefix specified in connection
# database specified in query resource
# so the database/connection specified in resource gets the highest precedence
@test "database specified in mod through connection ref, search_path_prefix specified in connection, database specified in query resource" {
  # add the steampipe connection with search_path_prefix defined
  # write the steampipe connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/with_search_path_prefix_bar.ppc
connection "steampipe" "with_search_path_prefix_bar" {
  host     = "localhost"
  port     = 9193
  db       = "steampipe"
  username = "steampipe"
  search_path_prefix = ["bar"]
}
EOF

  # checkout the mod with database(connection ref) specified in mod.pp, connection has search_path_prefix defined
  cd $MODS_DIR/mod_with_resource_db_steampipe_search_path

  # run a powerpipe query to verify that the search_path/database specified in query is used
  run powerpipe query run query.search_path --output csv
  echo $output

  # check output that the search_path specified in query is used
  assert_output --partial "bar"
}

# database specified in mod through connection ref, search_path_prefix specified in connection
# database specified by --database in runtime
# so the search_path_prefix specified in --database gets the highest precedence
@test "database specified in mod through connection ref, search_path_prefix specified in connection, database specified by --database in runtime" {
  # add the steampipe connection with search_path_prefix defined
  # write the steampipe connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/with_search_path_prefix_bar.ppc
connection "steampipe" "with_search_path_prefix_bar" {
  host     = "localhost"
  port     = 9193
  db       = "steampipe"
  username = "steampipe"
  search_path_prefix = ["bar"]
}
EOF

  # checkout the mod with database(connection ref) specified in mod.pp, connection has search_path_prefix defined
  cd $MODS_DIR/mod_with_db_steampipe_search_path

  # run a powerpipe query to verify that the search_path specified in --var is used
  run powerpipe query run query.search_path --output csv --database connection.steampipe.with_search_path_prefix_bar
  echo $output

  # check output that the search_path specified in mod is used
  assert_output --partial "bar"
}

# POWERPIPE_PERMITTED_ROOT_PATHS keeps the files that a mod names inside the folders it lists.
# Each test below writes its own mod into a temp folder and sets the setting to that folder.

# write a mod with one query whose title reads the file named by $1
function write_file_read_mod() {
  mkdir -p "$BATS_TEST_TMPDIR/mod"
  cat << EOF > "$BATS_TEST_TMPDIR/mod/mod.pp"
mod "permitted_paths_file_mod" {
  title = "Permitted paths file mod"
}

query "file_title" {
  title = file("$1")
  sql   = "select 1 as one"
}
EOF
}

# write a mod whose database is $1 (an HCL expression), with a query that runs against it
function write_database_mod() {
  mkdir -p "$BATS_TEST_TMPDIR/mod"
  cat << EOF > "$BATS_TEST_TMPDIR/mod/mod.pp"
mod "permitted_paths_db_mod" {
  title    = "Permitted paths database mod"
  database = $1
}

query "count_rows" {
  sql = "select 1 as one"
}
EOF
}

@test "permitted root paths: file() read outside the permitted folder fails and names the setting" {
  echo "outside the folder" > "$BATS_TEST_TMPDIR/outside.txt"
  write_file_read_mod "$BATS_TEST_TMPDIR/outside.txt"
  cd "$BATS_TEST_TMPDIR/mod"

  POWERPIPE_PERMITTED_ROOT_PATHS="$BATS_TEST_TMPDIR/mod" run powerpipe query show query.file_title --output json
  echo $output

  assert_failure
  assert_output --partial "POWERPIPE_PERMITTED_ROOT_PATHS"
  assert_output --partial "outside.txt"
}

@test "permitted root paths: the same file() read works when the setting is unset" {
  echo "outside the folder" > "$BATS_TEST_TMPDIR/outside.txt"
  write_file_read_mod "$BATS_TEST_TMPDIR/outside.txt"
  cd "$BATS_TEST_TMPDIR/mod"

  run powerpipe query show query.file_title --output json
  echo $output

  assert_success
  assert_output --partial "outside the folder"
}

@test "permitted root paths: file() read inside the permitted folder works" {
  mkdir -p "$BATS_TEST_TMPDIR/mod"
  echo "inside the folder" > "$BATS_TEST_TMPDIR/mod/inside.txt"
  write_file_read_mod "$BATS_TEST_TMPDIR/mod/inside.txt"
  cd "$BATS_TEST_TMPDIR/mod"

  POWERPIPE_PERMITTED_ROOT_PATHS="$BATS_TEST_TMPDIR/mod" run powerpipe query show query.file_title --output json
  echo $output

  assert_success
  assert_output --partial "inside the folder"
}

@test "permitted root paths: mod database that is a DuckDB file outside the permitted folder fails and names the setting" {
  write_database_mod "\"duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb\""
  cd "$BATS_TEST_TMPDIR/mod"

  POWERPIPE_PERMITTED_ROOT_PATHS="$BATS_TEST_TMPDIR/mod" run powerpipe query run query.count_rows --output csv
  echo $output

  assert_failure
  assert_output --partial "POWERPIPE_PERMITTED_ROOT_PATHS"
  assert_output --partial "employee.duckdb"
}

@test "permitted root paths: mod database that is a DuckDB file loads when the setting is unset" {
  write_database_mod "\"duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb\""
  cd "$BATS_TEST_TMPDIR/mod"

  run powerpipe query run query.count_rows --output csv
  echo $output

  assert_success
  assert_output --partial "one"
}

@test "permitted root paths: mod database that is a SQLite file outside the permitted folder fails and names the setting" {
  write_database_mod "\"sqlite:///$MODS_DIR/sqlite_mod/chinook.db\""
  cd "$BATS_TEST_TMPDIR/mod"

  POWERPIPE_PERMITTED_ROOT_PATHS="$BATS_TEST_TMPDIR/mod" run powerpipe query run query.count_rows --output csv
  echo $output

  assert_failure
  assert_output --partial "POWERPIPE_PERMITTED_ROOT_PATHS"
  assert_output --partial "chinook.db"
}

@test "permitted root paths: mod database that is a SQLite file loads when the setting is unset" {
  write_database_mod "\"sqlite:///$MODS_DIR/sqlite_mod/chinook.db\""
  cd "$BATS_TEST_TMPDIR/mod"

  run powerpipe query run query.count_rows --output csv
  echo $output

  assert_success
  assert_output --partial "one"
}

@test "permitted root paths: SQLite file as a query database outside the permitted folder fails and names the setting" {
  mkdir -p "$BATS_TEST_TMPDIR/mod"
  cat << EOF > "$BATS_TEST_TMPDIR/mod/mod.pp"
mod "permitted_paths_query_db_mod" {
  title = "Permitted paths query database mod"
}

query "count_albums" {
  database = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
  sql      = "select count(*) as albums from albums"
}
EOF
  cd "$BATS_TEST_TMPDIR/mod"

  POWERPIPE_PERMITTED_ROOT_PATHS="$BATS_TEST_TMPDIR/mod" run powerpipe query run query.count_albums --output csv
  echo $output

  assert_failure
  assert_output --partial "POWERPIPE_PERMITTED_ROOT_PATHS"
  assert_output --partial "chinook.db"

  # the same mod loads when the setting is unset (show does not connect to the default database)
  run powerpipe query show query.count_albums --output json
  echo $output

  assert_success
  assert_output --partial "chinook.db"
}

@test "permitted root paths: .ppc DuckDB connection outside the permitted folder fails and names the setting" {
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/permitted_paths_duckdb.ppc
connection "duckdb" "permitted_paths_employees" {
  file_name = "$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF
  write_database_mod "connection.duckdb.permitted_paths_employees"
  cd "$BATS_TEST_TMPDIR/mod"

  POWERPIPE_PERMITTED_ROOT_PATHS="$BATS_TEST_TMPDIR/mod" run powerpipe query run query.count_rows --output csv
  echo $output

  assert_failure
  assert_output --partial "POWERPIPE_PERMITTED_ROOT_PATHS"
  assert_output --partial "employee.duckdb"

  # the same connection works when the setting is unset
  run powerpipe query run query.count_rows --output csv
  echo $output

  assert_success
  assert_output --partial "one"
}

@test "permitted root paths: a relative entry stops powerpipe at start and names the entry" {
  write_database_mod "\"sqlite:///$MODS_DIR/sqlite_mod/chinook.db\""
  cd "$BATS_TEST_TMPDIR/mod"

  POWERPIPE_PERMITTED_ROOT_PATHS="relative/folder" run powerpipe query run query.count_rows --output csv
  echo $output

  assert_failure
  assert_output --partial "POWERPIPE_PERMITTED_ROOT_PATHS"
  assert_output --partial "relative/folder"
}
