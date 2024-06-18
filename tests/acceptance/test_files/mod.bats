load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"


# operation: install; pull-mode: default (latest); top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod (pull mode default) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: default (latest); top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod (pull mode default) - top already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: default (latest); top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but does not meet version constraints; expected: update
@test "install mod (pull mode default) - top level not meet version constraints" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the version of top level mod
  $TEST_SCRIPTS_DIR/update_top_level_mod_version.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}


function setup() {
  # create the work folder to run the tests
  tmp_dir="$(mktemp -d)"
  mkdir -p "${tmp_dir}"
}

function teardown() {
  # cleanup the work folder
  rm -rf "${tmp_dir}"
}