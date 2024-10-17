load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# These set of tests are skipped locally
# To run these tests locally set the SPIPETOOLS_TOKEN env var.
# These tests will be skipped locally unless the below env var is set.

function setup() {
  if [[ -z "${SPIPETOOLS_TOKEN}" ]]; then
    skip
  fi
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
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
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
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
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
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  connection_string = "duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb"
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
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
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
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  connection_string = "duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb"
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
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  connection_string = "duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb"
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

# database specified in mod require through a var
@test "database specified through variable in dependency mod require block" {
  skip "not working"
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  connection_string = "duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb"
}
EOF

  # checkout the mod with database specified in mod.pp dependant mod require block
  cd $MODS_DIR/mod_with_db_in_require

  # run a powerpipe query to verify that the database specified through mod require block is used
  run powerpipe query run query.sqlite_db_query --output csv
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
  assert_output --partial "redhood-aaa"
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
  assert_output --partial "redhood-aaa"
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
  assert_output --partial "redhood-aaa"
}

# test steampipe connection with implicit workspace works
@test "steampipe connection with implicit workspace" {
  skip "not working"
    # add the steampipe connection with pipes workspace
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/steampipe.ppc
connection "steampipe" "pipes" {
  workspace = "turbot-ops/clitesting"
}
EOF

  # checkout the mod with database specified in mod.pp
  cd $MODS_DIR/mod_with_db_var

  # run a powerpipe query to verify that the database specified through --var(connection ref) in mod is used
  run powerpipe query run query.pipes_workspace_query --output csv --var database=connection.steampipe.pipes --pipes-token $SPIPETOOLS_TOKEN
  echo $output

  # check output that the database specified through default value of variable in mod is used
  assert_output --partial "redhood-aaa"
}

# database specified in resource
# database specified in resource - so the resource level database gets the highest precedence
@test "database specified in resource" {
    # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  connection_string = "duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb"
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

# TODO - remove this test once deprecated --database flag is removed
# database specified in mod
# database specified in resource
# so the database specified in resource gets the highest precedence
@test "database specified in mod and also in resource" {
  # add the sqlite connection
  # write the sqlite connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/sqlite_conn.ppc
connection "sqlite" "albums" {
  connection_string = "sqlite:///$MODS_DIR/sqlite_mod/chinook.db"
}
EOF

  # add the duckdb connection
  # write the duckdb connection with $MODS_DIR placeholder directly into the config file
  cat << EOF > $POWERPIPE_INSTALL_DIR/config/duckdb_conn.ppc
connection "duckdb" "employees" {
  connection_string = "duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb"
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
