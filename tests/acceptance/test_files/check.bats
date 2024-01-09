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
