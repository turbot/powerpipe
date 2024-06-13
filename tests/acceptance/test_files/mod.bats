load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# operation: install
# pull-mode: default (latest)
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: no mods installed
# expected: install all mods
@test "install mod (pull mode default) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # check the stdout mod tree
  assert_output 'Initializing mod, created mod.pp.

Installed 3 mods:

local
└── github.com/pskrbasu/powerpipe-mod-1@v1.0.0
    └── github.com/pskrbasu/powerpipe-mod-2@v2.0.0
        └── github.com/pskrbasu/powerpipe-mod-3@v1.0.0'

  # check the folder structure (all 3 mods should be present and also check mod contents)
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: default (latest)
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top level mod already installed
# expected: all mods are up to date
@test "install mod (pull mode default) - top already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # check the stdout mod tree
  assert_output 'All targetted mods are up to date'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: default (latest)
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top installed but does not meet version constraints
# expected: update
@test "install mod (pull mode default) - top level not meet version constraints" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # update the version of top level mod
  JSON_FILE=".mod.cache.json"
  jq '.local["github.com/pskrbasu/powerpipe-mod-1"].version = "0.1.0"' "$JSON_FILE" > tmp.$$.json && mv tmp.$$.json "$JSON_FILE"

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # check the stdout mod tree
  assert_output '
Installed 3 mods:

local
└── github.com/pskrbasu/powerpipe-mod-1@v1.0.0
    └── github.com/pskrbasu/powerpipe-mod-2@v2.0.0
        └── github.com/pskrbasu/powerpipe-mod-3@v1.0.0'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: default (latest)
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top installed but retagged in repo
# expected: all mods are up to date
@test "install mod (pull mode default) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # update the commit of top level mod
  JSON_FILE=".mod.cache.json"
  jq '.local["github.com/pskrbasu/powerpipe-mod-1"].commit = "43746d1cd00bb9ecdccc9c6babe0a93bef4ee446"' "$JSON_FILE" > tmp.$$.json && mv tmp.$$.json "$JSON_FILE"

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # check the stdout mod tree
  assert_output 'All targetted mods are up to date'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: full
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: no mods installed
# expected: install all mods
@test "install mod (pull mode full) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  # check the stdout mod tree
  assert_output 'Initializing mod, created mod.pp.

Installed 3 mods:

local
└── github.com/pskrbasu/powerpipe-mod-1@v1.0.0
    └── github.com/pskrbasu/powerpipe-mod-2@v2.0.0
        └── github.com/pskrbasu/powerpipe-mod-3@v1.0.0'

  # check the folder structure (all 3 mods should be present and also check mod contents)
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: full
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top level mod already installed
# expected: all mods are up to date
@test "install mod (pull mode full) - top already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  # check the stdout mod tree
  assert_output 'All targetted mods are up to date'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: full
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top installed but does not meet version constraints
# expected: update
@test "install mod (pull mode full) - top level not meet version constraints" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # update the version of top level mod
  JSON_FILE=".mod.cache.json"
  jq '.local["github.com/pskrbasu/powerpipe-mod-1"].version = "0.1.0"' "$JSON_FILE" > tmp.$$.json && mv tmp.$$.json "$JSON_FILE"

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  # check the stdout mod tree
  assert_output '
Installed 3 mods:

local
└── github.com/pskrbasu/powerpipe-mod-1@v1.0.0
    └── github.com/pskrbasu/powerpipe-mod-2@v2.0.0
        └── github.com/pskrbasu/powerpipe-mod-3@v1.0.0'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: full
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top installed but retagged in repo
# expected: update
@test "install mod (pull mode full) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # update the commit of top level mod
  JSON_FILE=".mod.cache.json"
  jq '.local["github.com/pskrbasu/powerpipe-mod-1"].commit = "43746d1cd00bb9ecdccc9c6babe0a93bef4ee446"' "$JSON_FILE" > tmp.$$.json && mv tmp.$$.json "$JSON_FILE"

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  # check the stdout mod tree
  assert_output '
Upgraded 1 mod:

local
└── github.com/pskrbasu/powerpipe-mod-1@v1.0.0'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: full
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: new version available that satisfies constraint
# expected: update
@test "install mod (pull mode full) - new version available that satisfies constraint" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # update the version of top level mod
  JSON_FILE=".mod.cache.json"
  jq '.local["github.com/pskrbasu/powerpipe-mod-1"].version = "0.1.0"' "$JSON_FILE" > tmp.$$.json && mv tmp.$$.json "$JSON_FILE"

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull full

  # check the stdout mod tree
  assert_output '
Installed 3 mods:

local
└── github.com/pskrbasu/powerpipe-mod-1@v1.0.0
    └── github.com/pskrbasu/powerpipe-mod-2@v2.0.0
        └── github.com/pskrbasu/powerpipe-mod-3@v1.0.0'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: latest
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: no mods installed
# expected: install all mods
@test "install mod (pull mode latest) - no mods installed" {
  cd "$tmp_dir"
  # no mods installed

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull latest

  # check the stdout mod tree
  assert_output 'Initializing mod, created mod.pp.

Installed 3 mods:

local
└── github.com/pskrbasu/powerpipe-mod-1@v1.0.0
    └── github.com/pskrbasu/powerpipe-mod-2@v2.0.0
        └── github.com/pskrbasu/powerpipe-mod-3@v1.0.0'

  # check the folder structure (all 3 mods should be present and also check mod contents)
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: latest
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top level mod already installed
# expected: all mods are up to date
@test "install mod (pull mode latest) - top already installed" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull latest

  # check the stdout mod tree
  assert_output 'All targetted mods are up to date'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: latest
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top installed but does not meet version constraints
# expected: update
@test "install mod (pull mode latest) - top level not meet version constraints" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # update the version of top level mod
  JSON_FILE=".mod.cache.json"
  jq '.local["github.com/pskrbasu/powerpipe-mod-1"].version = "0.1.0"' "$JSON_FILE" > tmp.$$.json && mv tmp.$$.json "$JSON_FILE"

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1 --pull latest

  # check the stdout mod tree
  assert_output '
Installed 3 mods:

local
└── github.com/pskrbasu/powerpipe-mod-1@v1.0.0
    └── github.com/pskrbasu/powerpipe-mod-2@v2.0.0
        └── github.com/pskrbasu/powerpipe-mod-3@v1.0.0'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
}

# operation: install
# pull-mode: latest
# top-level-mod constraint: version
# l1 constraint: version
# l2 constraint: version
# scenario: top installed but retagged in repo
# expected: all mods are up to date
@test "install mod (pull mode latest) - top installed but retagged in repo" {
  cd "$tmp_dir"
  # mod already installed
  powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # update the commit of top level mod
  JSON_FILE=".mod.cache.json"
  jq '.local["github.com/pskrbasu/powerpipe-mod-1"].commit = "43746d1cd00bb9ecdccc9c6babe0a93bef4ee446"' "$JSON_FILE" > tmp.$$.json && mv tmp.$$.json "$JSON_FILE"

  # install mod
  run powerpipe mod install github.com/pskrbasu/powerpipe-mod-1

  # check the stdout mod tree
  assert_output 'All targetted mods are up to date'

  # check the folder structure (all 3 mods should be present and also check mod contents) - should be unchanged
  run ls .powerpipe/mods/github.com/pskrbasu/
  assert_output 'powerpipe-mod-1@v1.0.0
powerpipe-mod-2@v2.0.0
powerpipe-mod-3@v1.0.0'
  run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0/
  assert_output 'README.md
mod.sp
query.sp'

  # check lock file correct
  version=$(jq -r '.local["github.com/pskrbasu/powerpipe-mod-1"].version' .mod.cache.json) # check top level mod version - should be unchanged
  assert_equal "$version" "1.0.0"
  version_2=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"].["github.com/pskrbasu/powerpipe-mod-2"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_2" "2.0.0"
  version_3=$(jq -r '.["github.com/pskrbasu/powerpipe-mod-2@v2.0.0"].["github.com/pskrbasu/powerpipe-mod-3"].version' .mod.cache.json) # check dependency mod version - should be unchanged
  assert_equal "$version_3" "1.0.0"
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
