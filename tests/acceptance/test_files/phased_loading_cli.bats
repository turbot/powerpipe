load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# Comprehensive CLI tests for phased/lazy loading
# Tests verify that lazy loading produces identical output to eager loading
# across all CLI commands and output formats.

setup() {
  export PHASED_LOADING_MOD="$MODS_DIR/phased_loading_comparison_mod"
  export EMPTY_MOD="$MODS_DIR/empty_mod"
  export VARIABLE_TAGS_MOD="$MODS_DIR/variable_tags_mod"
}

# Helper function to sort JSON array by qualified_name for consistent comparison
sort_json_by_name() {
  jq -S 'sort_by(.qualified_name)'
}

# ============================================
# Dashboard List - Output Format Tests
# ============================================

@test "dashboard list - JSON output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output json
  assert_success
  eager_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$eager_output" > "$BATS_TMPDIR/eager_dash_json.json"

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success
  lazy_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$lazy_output" > "$BATS_TMPDIR/lazy_dash_json.json"

  # Compare using jd for semantic JSON comparison
  run jd "$BATS_TMPDIR/eager_dash_json.json" "$BATS_TMPDIR/lazy_dash_json.json"
  echo "JSON diff: $output"
  assert_success

  cd -
}

@test "dashboard list - plain output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output plain
  assert_success
  eager_output="$output"

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output plain
  assert_success
  lazy_output="$output"

  assert_equal "$eager_output" "$lazy_output"

  cd -
}

@test "dashboard list - pretty output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output pretty
  assert_success
  eager_output="$output"

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output pretty
  assert_success
  lazy_output="$output"

  assert_equal "$eager_output" "$lazy_output"

  cd -
}

@test "dashboard list - includes all metadata fields" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Verify dashboard_with_tags has all expected fields
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_tags\") | has(\"qualified_name\", \"resource_name\", \"title\", \"tags\")'"
  assert_success

  # Verify tags are populated
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_tags\") | .tags.service == \"test_service\"'"
  assert_success

  # Verify title is populated
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_tags\") | .title == \"Dashboard With Tags\"'"
  assert_success

  cd -
}

# ============================================
# Benchmark List - Output Format Tests
# ============================================

@test "benchmark list - JSON output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe benchmark list --output json
  assert_success
  eager_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$eager_output" > "$BATS_TMPDIR/eager_bench_json.json"

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success
  lazy_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$lazy_output" > "$BATS_TMPDIR/lazy_bench_json.json"

  # Compare
  run jd "$BATS_TMPDIR/eager_bench_json.json" "$BATS_TMPDIR/lazy_bench_json.json"
  echo "JSON diff: $output"
  assert_success

  cd -
}

@test "benchmark list - plain output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe benchmark list --output plain
  assert_success
  eager_output="$output"

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output plain
  assert_success
  lazy_output="$output"

  assert_equal "$eager_output" "$lazy_output"

  cd -
}

@test "benchmark list - includes all metadata fields" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Verify benchmark_with_tags has expected fields
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"benchmark_with_tags\") | has(\"qualified_name\", \"resource_name\", \"title\", \"tags\")'"
  assert_success

  # Verify tags populated
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"benchmark_with_tags\") | .tags.service == \"test_service\"'"
  assert_success

  cd -
}

# ============================================
# Query List Tests
# ============================================

@test "query list - JSON output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe query list --output json
  # Query list may be empty in this mod, but should still succeed
  eager_status="$status"
  eager_output=$(echo "$output" | sed -n '/^\[/,$p' | jq -S '.' 2>/dev/null || echo "[]")
  echo "$eager_output" > "$BATS_TMPDIR/eager_query_json.json"

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe query list --output json
  lazy_status="$status"
  lazy_output=$(echo "$output" | sed -n '/^\[/,$p' | jq -S '.' 2>/dev/null || echo "[]")
  echo "$lazy_output" > "$BATS_TMPDIR/lazy_query_json.json"

  # Both should have same status
  assert_equal "$eager_status" "$lazy_status"

  # Compare if both succeeded
  if [ "$eager_status" -eq 0 ] && [ "$lazy_status" -eq 0 ]; then
    run jd "$BATS_TMPDIR/eager_query_json.json" "$BATS_TMPDIR/lazy_query_json.json"
    echo "Query list diff: $output"
    assert_success
  fi

  cd -
}

@test "control list - JSON output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe control list --output json
  assert_success
  eager_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$eager_output" > "$BATS_TMPDIR/eager_ctrl_json.json"

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe control list --output json
  assert_success
  lazy_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$lazy_output" > "$BATS_TMPDIR/lazy_ctrl_json.json"

  # Compare
  run jd "$BATS_TMPDIR/eager_ctrl_json.json" "$BATS_TMPDIR/lazy_ctrl_json.json"
  echo "Control list diff: $output"
  assert_success

  cd -
}

# ============================================
# Empty Mod Edge Case Tests
# ============================================

@test "empty mod - dashboard list returns empty array" {
  cd "$EMPTY_MOD"

  run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')
  count=$(echo "$json_output" | jq 'length')
  assert_equal "$count" "0"

  cd -
}

@test "empty mod - benchmark list returns empty array" {
  cd "$EMPTY_MOD"

  run powerpipe benchmark list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')
  count=$(echo "$json_output" | jq 'length')
  assert_equal "$count" "0"

  cd -
}

@test "empty mod - query list returns empty array" {
  cd "$EMPTY_MOD"

  run powerpipe query list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')
  count=$(echo "$json_output" | jq 'length')
  assert_equal "$count" "0"

  cd -
}

# ============================================
# Variable Tags Mod Tests
# ============================================

@test "variable tags mod - tags are resolved (not raw variables)" {
  cd "$VARIABLE_TAGS_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Tags should NOT contain unresolved variable references like ${var.xxx}
  has_unresolved=$(echo "$json_output" | jq '[.[].tags | values[] | select(type == "string" and contains("${"))] | length')
  assert_equal "$has_unresolved" "0"

  cd -
}

@test "variable tags mod - variable tags have resolved values" {
  cd "$VARIABLE_TAGS_MOD"

  # Test with eager loading to get baseline
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output json
  assert_success
  eager_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$eager_output" > "$BATS_TMPDIR/eager_vartags.json"

  # Test with lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success
  lazy_output=$(echo "$output" | sed -n '/^\[/,$p' | sort_json_by_name)
  echo "$lazy_output" > "$BATS_TMPDIR/lazy_vartags.json"

  # Compare
  run jd "$BATS_TMPDIR/eager_vartags.json" "$BATS_TMPDIR/lazy_vartags.json"
  echo "Variable tags diff: $output"
  assert_success

  cd -
}

@test "variable tags mod - literal tags preserved" {
  cd "$VARIABLE_TAGS_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Dashboard with literal tags should have expected values
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_literal_tags\") | .tags.type == \"literal\"'"
  assert_success

  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_literal_tags\") | .tags.category == \"test\"'"
  assert_success

  cd -
}

# ============================================
# Count Verification Tests
# ============================================

@test "dashboard count identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output json
  assert_success
  eager_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success
  lazy_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  echo "Dashboard count - Eager: $eager_count, Lazy: $lazy_count"
  assert_equal "$eager_count" "$lazy_count"

  cd -
}

@test "benchmark count identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe benchmark list --output json
  assert_success
  eager_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success
  lazy_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  echo "Benchmark count - Eager: $eager_count, Lazy: $lazy_count"
  assert_equal "$eager_count" "$lazy_count"

  cd -
}

@test "control count identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe control list --output json
  assert_success
  eager_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe control list --output json
  assert_success
  lazy_count=$(echo "$output" | sed -n '/^\[/,$p' | jq 'length')

  echo "Control count - Eager: $eager_count, Lazy: $lazy_count"
  assert_equal "$eager_count" "$lazy_count"

  cd -
}

# ============================================
# Tag and Title Consistency Tests
# These verify that resources with tags/titles have them correctly populated
# ============================================

@test "tagged dashboards have consistent tags between eager and lazy" {
  cd "$PHASED_LOADING_MOD"

  # Get tags from eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output json
  assert_success
  eager_tags=$(echo "$output" | sed -n '/^\[/,$p' | jq -S '[.[] | select(.tags != null) | {name: .resource_name, tags: .tags}] | sort_by(.name)')
  echo "$eager_tags" > "$BATS_TMPDIR/eager_tags.json"

  # Get tags from lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success
  lazy_tags=$(echo "$output" | sed -n '/^\[/,$p' | jq -S '[.[] | select(.tags != null) | {name: .resource_name, tags: .tags}] | sort_by(.name)')
  echo "$lazy_tags" > "$BATS_TMPDIR/lazy_tags.json"

  # Compare
  run jd "$BATS_TMPDIR/eager_tags.json" "$BATS_TMPDIR/lazy_tags.json"
  echo "Tags diff: $output"
  assert_success

  cd -
}

@test "tagged benchmarks have consistent tags between eager and lazy" {
  cd "$PHASED_LOADING_MOD"

  # Get tags from eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe benchmark list --output json
  assert_success
  eager_tags=$(echo "$output" | sed -n '/^\[/,$p' | jq -S '[.[] | select(.tags != null) | {name: .resource_name, tags: .tags}] | sort_by(.name)')
  echo "$eager_tags" > "$BATS_TMPDIR/eager_bench_tags.json"

  # Get tags from lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success
  lazy_tags=$(echo "$output" | sed -n '/^\[/,$p' | jq -S '[.[] | select(.tags != null) | {name: .resource_name, tags: .tags}] | sort_by(.name)')
  echo "$lazy_tags" > "$BATS_TMPDIR/lazy_bench_tags.json"

  # Compare
  run jd "$BATS_TMPDIR/eager_bench_tags.json" "$BATS_TMPDIR/lazy_bench_tags.json"
  echo "Tags diff: $output"
  assert_success

  cd -
}

@test "all dashboards have titles in lazy mode" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Check no dashboard has empty/null title
  missing_titles=$(echo "$json_output" | jq '[.[] | select(.title == null or .title == "")] | length')
  assert_equal "$missing_titles" "0"

  cd -
}

@test "all benchmarks have titles in lazy mode" {
  cd "$PHASED_LOADING_MOD"

  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Check no benchmark has empty/null title
  missing_titles=$(echo "$json_output" | jq '[.[] | select(.title == null or .title == "")] | length')
  assert_equal "$missing_titles" "0"

  cd -
}
