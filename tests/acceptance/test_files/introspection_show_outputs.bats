load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

@test "control show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe control show sample_control_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 7th and 23rd lines, since they contain information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "7d;23d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "7d;23d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_control_show_output.json" output.json
  echo $output
  assert_success
}

@test "query show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe query show sample_query_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 6th and 37th lines, since they contain information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "6d;37d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "6d;37d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_query_show_output.json" output.json
  echo $output
  assert_success
}

@test "variable show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe variable show sample_var_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 6th line, since it contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "6d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "6d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_variable_show_output.json" output.json
  echo $output
  assert_success
}

@test "benchmark show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe benchmark show sample_benchmark_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 9th and 23rd lines, since they contain information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "9d;23d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "9d;23d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_benchmark_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard show sample_dashboard_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 9th and 15th lines, since they contain information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "9d;15d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "9d;15d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_show_output.json" output.json
  echo $output
  assert_success
}

function teardown() {
  rm -f output.json
  cd -
}