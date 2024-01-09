load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

@test "powerpipe check exitCode - no control alarms or errors" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe check benchmark.all_controls_ok
  assert_equal $status 0
  cd -
}

@test "powerpipe check exitCode - with controls in error" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe check benchmark.control_check_rendering_benchmark
  assert_equal $status 2
  cd -
}

@test "powerpipe check exitCode - with controls in error(running multiple benchmarks together)" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe check benchmark.control_summary_benchmark benchmark.check_cache_benchmark
  assert_equal $status 2
  cd -
}

@test "powerpipe check exitCode - runtime error(insufficient args)" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe check
  assert_equal $status 254
  cd -
}

@test "steampipe check long control title" {
  cd $CONTROL_RENDERING_TEST_MOD
  export STEAMPIPE_DISPLAY_WIDTH=100
  run steampipe check control.control_long_title --progress=false --theme=plain
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_long_title.txt)"
  cd -
}

@test "steampipe check short control title" {
  cd $CONTROL_RENDERING_TEST_MOD
  export STEAMPIPE_DISPLAY_WIDTH=100
  run steampipe check control.control_short_title --progress=false --theme=plain
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_short_title.txt)"
  cd -
}

@test "steampipe check unicode control title" {
  cd $CONTROL_RENDERING_TEST_MOD
  export STEAMPIPE_DISPLAY_WIDTH=100
  run steampipe check control.control_unicode_title --progress=false --theme=plain
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_unicode_title.txt)"
  cd -
}

@test "steampipe check reasons(very long, very short, unicode)" {
  cd $CONTROL_RENDERING_TEST_MOD
  export STEAMPIPE_DISPLAY_WIDTH=100
  run steampipe check control.control_long_short_unicode_reasons --progress=false --theme=plain
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_reasons.txt)"
  cd -
}

@test "steampipe check control with all possible statuses(10 OK, 5 ALARM, 2 ERROR, 1 SKIP and 3 INFO)" {
  cd $CONTROL_RENDERING_TEST_MOD
  export STEAMPIPE_DISPLAY_WIDTH=100
  run steampipe check control.sample_control_mixed_results_1 --progress=false --theme=plain
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_mixed_results.txt)"
  cd -
}

@test "steampipe check control with all resources in ALARM" {
  cd $CONTROL_RENDERING_TEST_MOD
  export STEAMPIPE_DISPLAY_WIDTH=100
  run steampipe check control.sample_control_all_alarms --progress=false --theme=plain
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_all_alarm.txt)"
  cd -
}

@test "steampipe check control with blank dimension" {
  cd $BLANK_DIMENSION_VALUE_TEST_MOD
  export STEAMPIPE_DISPLAY_WIDTH=100
  run steampipe check all --progress=false --theme=plain
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_blank_dimension.txt)"
  cd -
}

@test "steampipe check - output csv - no header" {
  cd $CONTROL_RENDERING_TEST_MOD
  run steampipe check control.sample_control_mixed_results_1 --output=csv --progress=false --header=false
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_check_csv_noheader.csv)"
  cd -
}

@test "steampipe check - output csv(check tags and dimensions sorting)" {
  cd $CONTROL_RENDERING_TEST_MOD
  run steampipe check control.sample_control_sorted_tags_and_dimensions --output=csv --progress=false
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_check_csv_sorted_tags.csv)"
  cd -
}

