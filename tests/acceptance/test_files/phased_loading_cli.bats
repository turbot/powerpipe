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
# Also sorts nested path arrays for deterministic comparison
sort_json_by_name() {
  jq -S 'sort_by(.qualified_name) | [.[] | .path |= (if . then sort else . end)]'
}

# Helper function to strip update notification banners from plain text output
# The banner appears intermittently and interferes with output comparison
strip_update_banner() {
  grep -v -E '^(\+---|\|.*version.*available|\|.*powerpipe.io|\|[[:space:]]*\|$)' | sed '/^$/d'
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
  eager_output=$(echo "$output" | strip_update_banner)

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output plain
  assert_success
  lazy_output=$(echo "$output" | strip_update_banner)

  assert_equal "$eager_output" "$lazy_output"

  cd -
}

@test "dashboard list - pretty output identical eager vs lazy" {
  cd "$PHASED_LOADING_MOD"

  # Eager loading
  POWERPIPE_WORKSPACE_PRELOAD=true run powerpipe dashboard list --output pretty
  assert_success
  eager_output=$(echo "$output" | strip_update_banner)

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe dashboard list --output pretty
  assert_success
  lazy_output=$(echo "$output" | strip_update_banner)

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
  eager_output=$(echo "$output" | strip_update_banner)

  # Lazy loading
  POWERPIPE_WORKSPACE_PRELOAD=false run powerpipe benchmark list --output plain
  assert_success
  lazy_output=$(echo "$output" | strip_update_banner)

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

# Variable tags that reference locals/variables (tags = local.common_tags) are now supported
# in lazy loading. The eval context is built with variables and locals at workspace initialization.
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

# ============================================
# Full Qualified Name Format Tests
# These verify that list commands output full qualified names (mod.type.name)
# to ensure backward compatibility with v1.4.2 behavior
# ============================================

@test "dashboard list - NAME column shows full qualified name in plain output" {
  cd "$PHASED_LOADING_MOD"

  run powerpipe dashboard list --output plain
  assert_success

  # Strip any update banners and check that NAME column contains full qualified names
  plain_output=$(echo "$output" | strip_update_banner)

  # Verify the output contains the full qualified name format (mod.dashboard.name)
  echo "$plain_output" | grep -q "phased_loading_comparison.dashboard.dashboard_with_tags"
  assert_success

  echo "$plain_output" | grep -q "phased_loading_comparison.dashboard.dashboard_without_tags"
  assert_success

  # Verify it does NOT show just the short name without the mod prefix
  # The NAME column should never be just "dashboard_with_tags" without the mod prefix
  if echo "$plain_output" | grep -E "^[[:space:]]*phased_loading_comparison[[:space:]]+dashboard_with_tags[[:space:]]*$" > /dev/null; then
    fail "NAME column shows short name instead of full qualified name"
  fi

  cd -
}

@test "dashboard list - qualified_name has correct format in JSON output" {
  cd "$PHASED_LOADING_MOD"

  run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Verify qualified_name matches the format: mod_name.dashboard.short_name
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_with_tags\") | .qualified_name == \"phased_loading_comparison.dashboard.dashboard_with_tags\"'"
  assert_success

  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"dashboard_without_tags\") | .qualified_name == \"phased_loading_comparison.dashboard.dashboard_without_tags\"'"
  assert_success

  # Verify all dashboards have qualified_name in the correct format (mod.dashboard.name)
  run bash -c "echo '$json_output' | jq -e '[.[] | .qualified_name | test(\"^phased_loading_comparison\\\\.dashboard\\\\.[a-z_]+\$\")] | all'"
  assert_success

  cd -
}

@test "benchmark list - NAME column shows full qualified name in plain output" {
  cd "$PHASED_LOADING_MOD"

  run powerpipe benchmark list --output plain
  assert_success

  plain_output=$(echo "$output" | strip_update_banner)

  # Verify the output contains the full qualified name format (mod.benchmark.name)
  echo "$plain_output" | grep -q "phased_loading_comparison.benchmark.benchmark_with_tags"
  assert_success

  cd -
}

@test "benchmark list - qualified_name has correct format in JSON output" {
  cd "$PHASED_LOADING_MOD"

  run powerpipe benchmark list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Verify qualified_name matches the format: mod_name.benchmark.short_name
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"benchmark_with_tags\") | .qualified_name == \"phased_loading_comparison.benchmark.benchmark_with_tags\"'"
  assert_success

  # Verify all benchmarks have qualified_name in the correct format
  run bash -c "echo '$json_output' | jq -e '[.[] | .qualified_name | test(\"^phased_loading_comparison\\\\.benchmark\\\\.[a-z_]+\$\")] | all'"
  assert_success

  cd -
}

@test "control list - NAME column shows full qualified name in plain output" {
  cd "$PHASED_LOADING_MOD"

  run powerpipe control list --output plain
  assert_success

  plain_output=$(echo "$output" | strip_update_banner)

  # Verify the output contains the full qualified name format (mod.control.name)
  echo "$plain_output" | grep -q "phased_loading_comparison.control.control_with_tags"
  assert_success

  cd -
}

@test "control list - qualified_name has correct format in JSON output" {
  cd "$PHASED_LOADING_MOD"

  run powerpipe control list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')

  # Verify qualified_name matches the format: mod_name.control.short_name
  run bash -c "echo '$json_output' | jq -e '.[] | select(.resource_name == \"control_with_tags\") | .qualified_name == \"phased_loading_comparison.control.control_with_tags\"'"
  assert_success

  # Verify all controls have qualified_name in the correct format
  run bash -c "echo '$json_output' | jq -e '[.[] | .qualified_name | test(\"^phased_loading_comparison\\\\.control\\\\.[a-z_]+\$\")] | all'"
  assert_success

  cd -
}

@test "dashboard show - works with full qualified name from list output" {
  cd "$PHASED_LOADING_MOD"

  # Get the qualified name from list output
  run powerpipe dashboard list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')
  qualified_name=$(echo "$json_output" | jq -r '.[] | select(.resource_name == "dashboard_with_tags") | .qualified_name')

  # Verify we got the expected qualified name
  assert_equal "$qualified_name" "phased_loading_comparison.dashboard.dashboard_with_tags"

  # Verify show command works with this qualified name
  run powerpipe dashboard show "$qualified_name"
  assert_success

  # Verify the show output contains expected data
  echo "$output" | grep -q "Dashboard With Tags"
  assert_success

  cd -
}

@test "benchmark show - works with full qualified name from list output" {
  cd "$PHASED_LOADING_MOD"

  # Get the qualified name from list output
  run powerpipe benchmark list --output json
  assert_success

  json_output=$(echo "$output" | sed -n '/^\[/,$p')
  qualified_name=$(echo "$json_output" | jq -r '.[] | select(.resource_name == "benchmark_with_tags") | .qualified_name')

  # Verify we got the expected qualified name
  assert_equal "$qualified_name" "phased_loading_comparison.benchmark.benchmark_with_tags"

  # Verify show command works with this qualified name
  run powerpipe benchmark show "$qualified_name"
  assert_success

  # Verify the show output contains expected data
  echo "$output" | grep -q "Benchmark With Tags"
  assert_success

  cd -
}
