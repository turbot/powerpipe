load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

setup() {
  export POWERPIPE_INSTALL_DIR="$BATS_TMPDIR/powerpipe-tag-filter"
  mkdir -p "$POWERPIPE_INSTALL_DIR/config"
  export LOCAL_DB_URL="duckdb://$POWERPIPE_INSTALL_DIR/tag_filter.duckdb"
}

@test "benchmark tag filter with deprecated=true includes only deprecated controls" {
  cd "$MODS_DIR/tag_filtering_mod"
  run powerpipe benchmark run tag_filtering_benchmark --tag deprecated=true --progress=false --output=json --database="$LOCAL_DB_URL" --mod-install=false
  assert_success
  filtered_output="$(echo "$output" | sed -n '/^{/,$p')"
  echo "$filtered_output" > "$BATS_TMPDIR/tag_filter_true.json"
  run jd "$TEST_DATA_DIR/expected_tag_filter_true.json" "$BATS_TMPDIR/tag_filter_true.json"
  echo "$filtered_output"
  assert_success
  cd -
}

@test "benchmark tag filter with deprecated!=true includes non-matching and missing tags" {
  cd "$MODS_DIR/tag_filtering_mod"
  run powerpipe benchmark run tag_filtering_benchmark --tag deprecated!=true --progress=false --output=json --database="$LOCAL_DB_URL" --mod-install=false
  # current behavior excludes missing tags; this is expected to fail until fixed
  filtered_output="$(echo "$output" | sed -n '/^{/,$p')"
  echo "$filtered_output" > "$BATS_TMPDIR/tag_filter_not_true.json"
  run jd "$TEST_DATA_DIR/expected_tag_filter_not_true.json" "$BATS_TMPDIR/tag_filter_not_true.json"
  echo "$filtered_output"
  assert_success
  cd -
}
