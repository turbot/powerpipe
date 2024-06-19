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

# operation: install; pull-mode: default (latest); top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repos; expected: all mods are up to date
@test "install mod (pull mode default) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the commit hash of top level mod to simulate retagging
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

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

# operation: install; pull-mode: full; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod (pull mode full) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod (pull mode full) - top already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but does not meet version constraints; expected: update
@test "install mod (pull mode full) - top level not meet version constraints" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the version of top level mod
  $TEST_SCRIPTS_DIR/update_top_level_mod_version.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repo; expected: update
@test "install mod (pull mode full) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the commit hash of top level mod to simulate retagging
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: new version available that satisfies constraint; expected: update
@test "install mod (pull mode full) - new version available that satisfies constraint" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the version of top level mod
  $TEST_SCRIPTS_DIR/update_top_level_mod_version.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod (pull mode latest) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod (pull mode latest) - top already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but does not meet version constraints; expected: update
@test "install mod (pull mode latest) - top level not meet version constraints" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the version of top level mod
  $TEST_SCRIPTS_DIR/update_top_level_mod_version.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repos; expected: all mods are up to date
@test "install mod (pull mode latest) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the commit hash of top level mod to simulate retagging
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod (pull mode development) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod (pull mode development) - top already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but does not meet version constraints; expected: update
@test "install mod (pull mode development) - top level not meet version constraints" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the version of top level mod
  $TEST_SCRIPTS_DIR/update_top_level_mod_version.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repos; expected: all mods are up to date
@test "install mod (pull mode development) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the commit hash of top level mod to simulate retagging
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod (pull mode minimal) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod (pull mode minimal) - top already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but does not meet version constraints; expected: update
@test "install mod (pull mode minimal) - top level not meet version constraints" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the version of top level mod
  $TEST_SCRIPTS_DIR/update_top_level_mod_version.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: version; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repos; expected: all mods are up to date
@test "install mod (pull mode minimal) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  
  # update the commit hash of top level mod to simulate retagging
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: default (latest); top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from branch (pull mode default) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: default (latest); top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from branch (pull mode default) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: default (latest); top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top installed but new commit is available; expected: update top level mod
@test "install mod from branch (pull mode default) - top installed but new commit is available" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  
  # update the commit hash of top level mod to simulate old commit
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from branch (pull mode full) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from branch (pull mode full) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top installed but new commit is available; expected: update top level mod
@test "install mod from branch (pull mode full) - top installed but new commit is available" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  
  # update the commit hash of top level mod to simulate old commit
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from branch (pull mode latest) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from branch (pull mode latest) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top installed but new commit is available; expected: update top level mod
@test "install mod from branch (pull mode latest) - top installed but new commit is available" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  
  # update the commit hash of top level mod to simulate old commit
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from branch (pull mode development) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from branch (pull mode development) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top installed but new commit is available; expected: update top level mod
@test "install mod from branch (pull mode development) - top installed but new commit is available" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  
  # update the commit hash of top level mod to simulate old commit
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from branch (pull mode minimal) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods_branch.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from branch (pull mode minimal) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: branch; l1 constraint: version; l2 constraint: version; scenario: top installed but new commit is available; expected:all mods are up to date
@test "install mod from branch (pull mode minimal) - top installed but new commit is available" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new

  
  # update the commit hash of top level mod to simulate old commit
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1#new --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure_branch.txt)"
}

# operation: install; pull-mode: default (latest); top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from tag (pull mode default) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: default (latest); top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from tag (pull mode default) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: default (latest); top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repo; expected: all mods are up to date
@test "install mod from tag (pull mode default) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  
  # update the commit hash of top level mod to simulate retagged
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from tag (pull mode full) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from tag (pull mode full) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: full; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repo; expected: update top level mod
@test "install mod from tag (pull mode full) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  
  # update the commit hash of top level mod to simulate retagged
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from tag (pull mode latest) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from tag (pull mode latest) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: latest; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repo; expected: all mods are up to date
@test "install mod from tag (pull mode latest) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  
  # update the commit hash of top level mod to simulate retagged
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from tag (pull mode development) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from tag (pull mode development) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: development; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repo; expected: all mods are up to date
@test "install mod from tag (pull mode development) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  
  # update the commit hash of top level mod to simulate retagged
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from tag (pull mode minimal) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/installed_3_mods.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top level mod already installed; expected: all mods are up to date
@test "install mod from tag (pull mode minimal) - top level mod already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: minimal; top-level-mod constraint: tag; l1 constraint: version; l2 constraint: version; scenario: top installed but retagged in repo; expected: all mods are up to date
@test "install mod from tag (pull mode minimal) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0

  
  # update the commit hash of top level mod to simulate retagged
  $TEST_SCRIPTS_DIR/update_top_level_mod_commit.sh

  # run install command
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1@v1.0.0 --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
  # check the folder structure (all 3 mods should be present and also check mod contents)
  echo "Verifying the mod folder structure"
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"
}

# operation: install; pull-mode: default(minimal); top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from local file (pull mode default) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed

  

  
  # check the partial output since the absolute path is not known
  echo "Verifying the mod tree output"
  assert_output  --partial "$(cat $TEST_DATA_DIR/installed_3_mods_partial.txt)"

  
}

# operation: install; pull-mode: default(minimal); top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: top level mod installed; expected: all mods are up to date
@test "install mod from local file (pull mode default) - top level mod installed" {
  cd "$tmp_dir"
  # no mods installed
  powerpipe mod install $MODS_DIR/test_mods/mod_installed

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
}

# operation: install; pull-mode: full; top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from local file (pull mode full) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed --pull full

  

  
  # check the partial output since the absolute path is not known
  echo "Verifying the mod tree output"
  assert_output  --partial "$(cat $TEST_DATA_DIR/installed_3_mods_partial.txt)"

  
}

# operation: install; pull-mode: full; top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: top level mod installed; expected: all mods are up to date
@test "install mod from local file (pull mode full) - top level mod installed" {
  cd "$tmp_dir"
  # no mods installed
  powerpipe mod install $MODS_DIR/test_mods/mod_installed

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed --pull full

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
}

# operation: install; pull-mode: latest; top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from local file (pull mode latest) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed --pull latest

  

  
  # check the partial output since the absolute path is not known
  echo "Verifying the mod tree output"
  assert_output  --partial "$(cat $TEST_DATA_DIR/installed_3_mods_partial.txt)"

  
}

# operation: install; pull-mode: latest; top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: top level mod installed; expected: all mods are up to date
@test "install mod from local file (pull mode latest) - top level mod installed" {
  cd "$tmp_dir"
  # no mods installed
  powerpipe mod install $MODS_DIR/test_mods/mod_installed

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed --pull latest

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
}

# operation: install; pull-mode: development; top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from local file (pull mode development) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed --pull development

  

  
  # check the partial output since the absolute path is not known
  echo "Verifying the mod tree output"
  assert_output  --partial "$(cat $TEST_DATA_DIR/installed_3_mods_partial.txt)"

  
}

# operation: install; pull-mode: development; top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: top level mod installed; expected: all mods are up to date
@test "install mod from local file (pull mode development) - top level mod installed" {
  cd "$tmp_dir"
  # no mods installed
  powerpipe mod install $MODS_DIR/test_mods/mod_installed

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed --pull development

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
}

# operation: install; pull-mode: minimal; top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: no mods installed; expected: install all mods
@test "install mod from local file (pull mode minimal) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed
  

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed --pull minimal

  

  
  # check the partial output since the absolute path is not known
  echo "Verifying the mod tree output"
  assert_output  --partial "$(cat $TEST_DATA_DIR/installed_3_mods_partial.txt)"

  
}

# operation: install; pull-mode: minimal; top-level-mod constraint: file; l1 constraint: version; l2 constraint: version; scenario: top level mod installed; expected: all mods are up to date
@test "install mod from local file (pull mode minimal) - top level mod installed" {
  cd "$tmp_dir"
  # no mods installed
  powerpipe mod install $MODS_DIR/test_mods/mod_installed

  

  # run install command
  run powerpipe mod install $MODS_DIR/test_mods/mod_installed --pull minimal

  
  # check the stdout mod tree
  echo "Verifying the mod tree output"
  assert_output "$(cat $TEST_DATA_DIR/mod_up_to_date.txt)"

  

  
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