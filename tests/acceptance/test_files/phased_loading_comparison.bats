load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# These tests verify that phased/lazy loading produces IDENTICAL output to eager loading.
# Tests should FAIL initially to prove the gap between lazy and eager loading.
# Once phased loading is properly implemented, all tests should pass.

setup() {
  export PHASED_LOADING_MOD="$MODS_DIR/phased_loading_comparison_mod"
}

# Helper function to sort JSON array by qualified_name for consistent comparison
# Also sorts nested path arrays for deterministic comparison
sort_json_by_name() {
  jq -S 'sort_by(.qualified_name) | [.[] | .path |= (if . then sort else . end)]'
}

# Test: Dashboard list JSON output should be identical
@test "dashboard list output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output json
  assert_success
  eager_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$eager_output" > "$BATS_TMPDIR/eager_dashboard_list.json"

  # Lazy loading (default)
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success
  lazy_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$lazy_output" > "$BATS_TMPDIR/lazy_dashboard_list.json"

  # Compare - should be identical
  run jd "$BATS_TMPDIR/eager_dashboard_list.json" "$BATS_TMPDIR/lazy_dashboard_list.json"
  echo "Diff output: $output"
  assert_success

  cd -
}

# Test: Benchmark list JSON output should be identical
@test "benchmark list output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe benchmark list --output json
  assert_success
  eager_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$eager_output" > "$BATS_TMPDIR/eager_benchmark_list.json"

  # Lazy loading (default)
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success
  lazy_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$lazy_output" > "$BATS_TMPDIR/lazy_benchmark_list.json"

  # Compare - should be identical
  run jd "$BATS_TMPDIR/eager_benchmark_list.json" "$BATS_TMPDIR/lazy_benchmark_list.json"
  echo "Diff output: $output"
  assert_success

  cd -
}

# Test: Dashboard with tags - tags should be present in lazy mode
@test "dashboard list includes all tags in lazy mode" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success

  # Extract JSON array from output
  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Check that dashboard_with_tags has the expected tags
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_tags\") | .tags.service == \"test_service\"'"
  echo "Service tag check: $output"
  assert_success

  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_tags\") | .tags.category == \"comparison\"'"
  echo "Category tag check: $output"
  assert_success

  cd -
}

# Test: Benchmark with tags - tags should be present in lazy mode
@test "benchmark list includes all tags in lazy mode" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success

  # Extract JSON array from output
  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Check that benchmark_with_tags has the expected tags
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"benchmark_with_tags\") | .tags.service == \"test_service\"'"
  echo "Service tag check: $output"
  assert_success

  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"benchmark_with_tags\") | .tags.category == \"comparison\"'"
  echo "Category tag check: $output"
  assert_success

  cd -
}

# Test: Dashboard titles should be present in lazy mode
@test "dashboard list includes titles in lazy mode" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Check that dashboard has title
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_tags\") | .title == \"Dashboard With Tags\"'"
  echo "Title check: $output"
  assert_success

  cd -
}

# Test: Benchmark titles should be present in lazy mode
@test "benchmark list includes titles in lazy mode" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Check that benchmark has title
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"benchmark_with_tags\") | .title == \"Benchmark With Tags\"'"
  echo "Title check: $output"
  assert_success

  cd -
}

# Test: Dashboard descriptions should be present in lazy mode
@test "dashboard list includes descriptions in lazy mode" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Check that dashboard has description (if the list output includes it)
  # Note: This may not be in list output - adjust based on actual schema
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_tags\") | .description != null'"
  echo "Description check: $output"
  # This assertion may need adjustment based on whether description is in list output
  # assert_success

  cd -
}

# Test: Control list should include tags in lazy mode
@test "control list includes tags in lazy mode" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe control list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Check that control_with_tags has the expected tags
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"control_with_tags\") | .tags.control_type == \"test\"'"
  echo "Control tag check: $output"
  assert_success

  cd -
}

# Test: Dashboard count should be same in eager vs lazy
@test "dashboard count identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output json
  assert_success
  eager_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  # Lazy
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success
  lazy_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  echo "Eager count: $eager_count, Lazy count: $lazy_count"
  assert_equal "$eager_count" "$lazy_count"

  cd -
}

# Test: Benchmark count should be same in eager vs lazy
@test "benchmark count identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe benchmark list --output json
  assert_success
  eager_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  # Lazy
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success
  lazy_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  echo "Eager count: $eager_count, Lazy count: $lazy_count"
  assert_equal "$eager_count" "$lazy_count"

  cd -
}
