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
  skip
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

@test "dashboard container show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard container show dashboard_sample_dashboard_1_anonymous_container_0 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 16th and 22nd lines, since they contain information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "16d;22d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "16d;22d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_container_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard card show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard card show sample_card_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 5th line, since it contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "5d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "5d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_card_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard image show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard image show sample_image_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 6th line, since it contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "6d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "6d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_image_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard text show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard text show sample_text_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 5th line, since it contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "5d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "5d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_text_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard chart show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard chart show sample_chart_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 6th line, since it contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "6d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "6d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_chart_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard flow show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard flow show sample_flow_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 9th line, since it contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "9d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "9d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_flow_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard graph show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard graph show sample_graph_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 14th, 25th and 32nd lines, since they contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "14d;25d;32d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "14d;25d;32d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_graph_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard hierarchy show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard hierarchy show sample_hierarchy_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 14th, 25th and 32nd lines, since they contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "14d;25d;32d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "14d;25d;32d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_hierarchy_show_output.json" output.json
  echo $output
  assert_success
}

@test "dashboard table show output" {
  cd $SIMPLE_MOD_DIR
  powerpipe dashboard table show sample_table_1 --output json > output.json

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 6th line, since it contains information which would differ in github runners
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # For macOS, adding a backup extension (.bak) and then removing it to mimic in-place editing without a backup
    run sed -i ".bak" -e "6d" output.json && rm output.json.bak
  else
    # For Linux, using in-place editing without a backup file directly
    run sed -i -e "6d" output.json
  fi

  run jd "$TEST_DATA_DIR/expected_dashboard_table_show_output.json" output.json
  echo $output
  assert_success

}

function teardown() {
  rm -f output.json
  cd -
}