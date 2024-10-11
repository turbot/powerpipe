load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

@test "connecting to a duckdb backend" {
  cd $MODS_DIR/duckdb_mod

  # run a powerpipe query connecting to a local duckdb database
  run powerpipe query run query.total_employee --database duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb --output csv
  echo $output
  # check output
  assert_equal "$output" "$(cat $TEST_DATA_DIR/0)"
}

@test "using json extension(casting) with duckdb backend" {
  cd $MODS_DIR/duckdb_mod

  # run a powerpipe query(which performs a JSON casting operation) connecting to a local duckdb database
  run powerpipe query run query.json_casting --database duckdb:///$MODS_DIR/duckdb_mod/employee.duckdb --output csv
  echo $output
  # check output
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_duckdb_backend_json_casting.csv)"
}
