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

# @test "powerpipe check long control title" {
#   cd $CONTROL_RENDERING_TEST_MOD
#   export STEAMPIPE_DISPLAY_WIDTH=100
#   run powerpipe check control.control_long_title --progress=false --theme=plain
#   assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_long_title.txt)"
#   cd -
# }

# @test "powerpipe check short control title" {
#   cd $CONTROL_RENDERING_TEST_MOD
#   export STEAMPIPE_DISPLAY_WIDTH=100
#   run powerpipe check control.control_short_title --progress=false --theme=plain
#   assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_short_title.txt)"
#   cd -
# }

# @test "powerpipe check unicode control title" {
#   cd $CONTROL_RENDERING_TEST_MOD
#   export STEAMPIPE_DISPLAY_WIDTH=100
#   run powerpipe check control.control_unicode_title --progress=false --theme=plain
#   assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_unicode_title.txt)"
#   cd -
# }

# @test "powerpipe check reasons(very long, very short, unicode)" {
#   cd $CONTROL_RENDERING_TEST_MOD
#   export STEAMPIPE_DISPLAY_WIDTH=100
#   run powerpipe check control.control_long_short_unicode_reasons --progress=false --theme=plain
#   assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_reasons.txt)"
#   cd -
# }

# @test "powerpipe check control with all possible statuses(10 OK, 5 ALARM, 2 ERROR, 1 SKIP and 3 INFO)" {
#   cd $CONTROL_RENDERING_TEST_MOD
#   export STEAMPIPE_DISPLAY_WIDTH=100
#   run powerpipe check control.sample_control_mixed_results_1 --progress=false --theme=plain
#   assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_mixed_results.txt)"
#   cd -
# }

# @test "powerpipe check control with all resources in ALARM" {
#   cd $CONTROL_RENDERING_TEST_MOD
#   export STEAMPIPE_DISPLAY_WIDTH=100
#   run powerpipe check control.sample_control_all_alarms --progress=false --theme=plain
#   assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_all_alarm.txt)"
#   cd -
# }

# @test "powerpipe check control with blank dimension" {
#   cd $BLANK_DIMENSION_VALUE_TEST_MOD
#   export STEAMPIPE_DISPLAY_WIDTH=100
#   run powerpipe check all --progress=false --theme=plain
#   assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_blank_dimension.txt)"
#   cd -
# }

@test "powerpipe check - output csv - no header" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe check control.sample_control_mixed_results_1 --output=csv --progress=false --header=false
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_check_csv_noheader.csv)"
  cd -
}

@test "powerpipe check - output csv(check tags and dimensions sorting)" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe check control.sample_control_sorted_tags_and_dimensions --output=csv --progress=false
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_check_csv_sorted_tags.csv)"
  cd -
}

@test "powerpipe check - output json" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe check control.sample_control_mixed_results_1 --output json --progress=false --export output.json
  output=""
  run jd "$TEST_DATA_DIR/expected_check_json.json" output.json
  echo $output
  assert_success
  rm -f output.json
  cd -
}

@test "powerpipe check - export csv" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe check control.sample_control_mixed_results_1 --export test.csv --progress=false
  assert_equal "$(cat test.csv)" "$(cat $TEST_DATA_DIR/expected_check_csv.csv)"
  rm -f test.csv
  cd -
}

@test "powerpipe check - export csv - pipe separator" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe check control.sample_control_mixed_results_1 --export test.csv --separator="|" --progress=false
  assert_equal "$(cat test.csv)" "$(cat $TEST_DATA_DIR/expected_check_csv_pipe_separator.csv)"
  rm -f test.csv
  cd -
}

@test "powerpipe check - export csv(check tags and dimensions sorting)" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe check control.sample_control_sorted_tags_and_dimensions --export test.csv --progress=false
  assert_equal "$(cat test.csv)" "$(cat $TEST_DATA_DIR/expected_check_csv_sorted_tags.csv)"
  rm -f test.csv
  cd -
}
