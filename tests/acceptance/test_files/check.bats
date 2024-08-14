load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

@test "verify powerpipe benchmark exitCode - no control alarms or errors" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe benchmark run all_controls_ok
  assert_equal $status 0
  cd -
}

@test "verify powerpipe benchmark exitCode - with controls in error" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe benchmark run control_check_rendering_benchmark
  assert_equal $status 2
  cd -
}

# @test "verify powerpipe benchmark exitCode - runtime error(insufficient args)" {
#   cd $FUNCTIONALITY_TEST_MOD
#   run powerpipe benchmark
#   assert_equal $status 254
#   cd -
# }

# TODO: Implement STEAMPIPE_DISPLAY_WIDTH in powerpipe to test the below tests
# https://github.com/turbot/powerpipe/issues/154

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

@test "powerpipe control run - output csv - no header" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run sample_control_mixed_results_1 --output=csv --progress=false --header=false
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_check_csv_noheader.csv)"
  cd -
}

@test "powerpipe control run - output csv(check tags and dimensions sorting)" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run sample_control_sorted_tags_and_dimensions --output=csv --progress=false
  assert_equal "$output" "$(cat $TEST_DATA_DIR/expected_check_csv_sorted_tags.csv)"
  cd -
}

@test "powerpipe control run - output json" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run sample_control_mixed_results_1 --output json --progress=false --export output.json
  output=""
  run jd "$TEST_DATA_DIR/expected_check_json.json" output.json
  echo $output
  assert_success
  rm -f output.json
  cd -
}

@test "powerpipe control run - export csv" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run sample_control_mixed_results_1 --export test.csv --progress=false
  assert_equal "$(cat test.csv)" "$(cat $TEST_DATA_DIR/expected_check_csv.csv)"
  rm -f test.csv
  cd -
}

@test "powerpipe benchmark run - export csv (control re-used/ multiple parents)" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe benchmark run control_reused --export test.csv --progress=false
  assert_equal "$(cat test.csv)" "$(cat $TEST_DATA_DIR/expected_check_csv_multiple_parents.csv)"
  rm -f test.csv
  cd -
}

@test "powerpipe benchmark run - export json (control re-used/ multiple parents)" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe benchmark run control_reused --export test.json --progress=false
  assert_equal "$(cat test.json)" "$(cat $TEST_DATA_DIR/expected_check_json_multiple_parents.json)"
  rm -f test.json
  cd -
}

@test "powerpipe benchmark run - dry run (control re-used/ multiple parents)" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe benchmark run control_reused --dry-run
  assert_success
  cd -
}

@test "powerpipe control run - export csv - pipe separator" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run sample_control_mixed_results_1 --export test.csv --separator="|" --progress=false
  assert_equal "$(cat test.csv)" "$(cat $TEST_DATA_DIR/expected_check_csv_pipe_separator.csv)"
  rm -f test.csv
  cd -
}

@test "powerpipe control run - export csv(check tags and dimensions sorting)" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run sample_control_sorted_tags_and_dimensions --export test.csv --progress=false
  assert_equal "$(cat test.csv)" "$(cat $TEST_DATA_DIR/expected_check_csv_sorted_tags.csv)"
  rm -f test.csv
  cd -
}

@test "powerpipe control run - export json" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run control.sample_control_mixed_results_1 --export test.json --progress=false
  output=""
  run jd "$TEST_DATA_DIR/expected_check_json.json" test.json
  echo $output
  assert_success
  rm -f test.json
  cd -
}

@test "powerpipe control run - export html" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run control.sample_control_mixed_results_1 --export test.html --progress=false
  
  # checking for OS type, since sed command is different for linux and OSX
  # removing the 642nd line, since it contains file locations and timestamps
  if [[ "$OSTYPE" == "darwin"* ]]; then
    run sed -i ".html" "642d" test.html
    run sed -i ".html" "642d" test.html
    run sed -i ".html" "642d" test.html
  else
    run sed -i "642d" test.html
    run sed -i "642d" test.html
    run sed -i "642d" test.html
  fi

  assert_equal "$(cat test.html)" "$(cat $TEST_DATA_DIR/expected_check_html.html)"
  rm -rf test.html*
  cd -
}

@test "powerpipe control run - export md" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run control.sample_control_mixed_results_1 --export test.md --progress=false
  
  # checking for OS type, since sed command is different for linux and OSX
  # removing the 42nd line, since it contains file locations and timestamps
  if [[ "$OSTYPE" == "darwin"* ]]; then
    run sed -i ".md" "42d" test.md
  else
    run sed -i "42d" test.md
  fi

  assert_equal "$(cat test.md)" "$(cat $TEST_DATA_DIR/expected_check_markdown.md)"
  rm -rf test.md*
  cd -
}

@test "powerpipe control run - export nunit3" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run control.sample_control_mixed_results_1 --export test.xml --progress=false

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 6th line, since it contains duration, and duration will be different in each run
  if [[ "$OSTYPE" == "darwin"* ]]; then
    run sed -i ".xml" "6d" test.xml
  else
    run sed -i "6d" test.xml
  fi

  assert_equal "$(cat test.xml)" "$(cat $TEST_DATA_DIR/expected_check_nunit3.xml)"
  rm -f test.xml*
  cd -
}

@test "powerpipe control run - export snapshot" {
  cd $CONTROL_RENDERING_TEST_MOD
  run powerpipe control run control.sample_control_mixed_results_1 --export test.pps --progress=false

  # get the patch diff between the two snapshots
  run jd -f patch $TEST_DATA_DIR/expected_check_snapshot.pps test.pps

  # run the script to evaluate the patch
  # returns nothing if there is no diff(except start_time, end_time & search_path)
  diff=$($FILE_PATH/json_patch.sh $output)
  echo $diff
  rm -f test.pps

  # check if there is no diff returned by the script
  assert_equal "$diff" ""
  cd -
}

# testing the check summary output feature in powerpipe
@test "benchmark run summary output" {
  cd $FUNCTIONALITY_TEST_MOD
  run powerpipe benchmark run benchmark.control_summary_benchmark

  echo $output

  # TODO: Find a way to store the output in a file and match it with the 
  # expected file. For now the work-around is to check whether the output
  # contains `summary`
  assert_output --partial 'Summary'
}
