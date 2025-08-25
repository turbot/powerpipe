load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# Test setup - create test configuration files
setup() {
  # Create test directories
  mkdir -p $BATS_TEST_DIRNAME/test_configs/valid_config
  mkdir -p $BATS_TEST_DIRNAME/test_configs/invalid_config
  mkdir -p $BATS_TEST_DIRNAME/test_configs/empty_config
  
  # Create a simple configuration file that should be valid
  cat > $BATS_TEST_DIRNAME/test_configs/valid_config/test.ppc << 'EOF'
# This is a test configuration file
# It should be loaded when POWERPIPE_CONFIG_PATH is set correctly
EOF

  # Create an invalid configuration file
  cat > $BATS_TEST_DIRNAME/test_configs/invalid_config/invalid.ppc << 'EOF'
connection "invalid" "test" {
  invalid_argument = "this should cause an error"
}
EOF
}

# Test cleanup
teardown() {
  # Clean up test directories
  rm -rf $BATS_TEST_DIRNAME/test_configs
}

@test "POWERPIPE_CONFIG_PATH with valid configuration should work" {
  cd $MODS_DIR/functionality_test_mod
  
  # Set POWERPIPE_CONFIG_PATH to point to valid configuration
  export POWERPIPE_CONFIG_PATH="$BATS_TEST_DIRNAME/test_configs/valid_config"
  
  # Run a simple command that should work with valid config
  run powerpipe mod list
  echo "$output"
  
  # Should succeed
  assert_success
  assert_output --partial "functionality_test_mod"
}

@test "POWERPIPE_CONFIG_PATH with invalid configuration should fail appropriately" {
  cd $MODS_DIR/functionality_test_mod
  
  # Set POWERPIPE_CONFIG_PATH to point to invalid configuration
  export POWERPIPE_CONFIG_PATH="$BATS_TEST_DIRNAME/test_configs/invalid_config"
  
  # Run a command that loads configuration
  run powerpipe mod list
  echo "$output"
  
  # Should fail with configuration error
  assert_failure
  assert_output --partial "Failed to load"
}

@test "POWERPIPE_CONFIG_PATH with empty directory should use defaults" {
  cd $MODS_DIR/functionality_test_mod
  
  # Set POWERPIPE_CONFIG_PATH to point to empty directory
  export POWERPIPE_CONFIG_PATH="$BATS_TEST_DIRNAME/test_configs/empty_config"
  
  # Run a simple command
  run powerpipe mod list
  echo "$output"
  
  # Should succeed (falls back to default behavior)
  assert_success
  assert_output --partial "functionality_test_mod"
}

@test "POWERPIPE_CONFIG_PATH with multiple colon-separated paths should respect precedence" {
  cd $MODS_DIR/functionality_test_mod
  
  # Set POWERPIPE_CONFIG_PATH with multiple paths - valid first, then invalid
  export POWERPIPE_CONFIG_PATH="$BATS_TEST_DIRNAME/test_configs/valid_config:$BATS_TEST_DIRNAME/test_configs/invalid_config"
  
  # Run a command that loads configuration
  run powerpipe mod list
  echo "$output"
  
  # Should succeed because valid config has higher precedence
  assert_success
}

@test "POWERPIPE_CONFIG_PATH should work with relative paths" {
  cd $MODS_DIR/functionality_test_mod
  
  # Set POWERPIPE_CONFIG_PATH with relative path
  export POWERPIPE_CONFIG_PATH="../test_configs/valid_config"
  
  # Run a simple command
  run powerpipe mod list
  echo "$output"
  
  # Should succeed
  assert_success
  assert_output --partial "functionality_test_mod"
}

@test "POWERPIPE_CONFIG_PATH should work with tilde expansion" {
  cd $MODS_DIR/functionality_test_mod
  
  # Create config in home directory
  mkdir -p ~/test_powerpipe_config
  cat > ~/test_powerpipe_config/test.ppc << 'EOF'
# This is a test configuration file in home directory
EOF
  
  # Set POWERPIPE_CONFIG_PATH with tilde
  export POWERPIPE_CONFIG_PATH="~/test_powerpipe_config"
  
  # Run a simple command
  run powerpipe mod list
  echo "$output"
  
  # Should succeed
  assert_success
  assert_output --partial "functionality_test_mod"
  
  # Cleanup
  rm -rf ~/test_powerpipe_config
}

@test "POWERPIPE_CONFIG_PATH should be ignored when not set" {
  cd $MODS_DIR/functionality_test_mod
  
  # Unset POWERPIPE_CONFIG_PATH
  unset POWERPIPE_CONFIG_PATH
  
  # Run a simple command
  run powerpipe mod list
  echo "$output"
  
  # Should succeed with default behavior
  assert_success
  assert_output --partial "functionality_test_mod"
}
