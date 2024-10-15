load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

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

# database specified in mod definition - so the mod level database gets the highest precedence
@test "database specified in mod definition" {
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