load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

@test "list with no mods installed" {
  run powerpipe mod list
  assert_output 'name
local'
}

@test "install latest(--force)" {
  run powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance --force
  assert_output --partial 'Installed 1 mod:

local
└── github.com/turbot/steampipe-mod-aws-compliance'
  # need the check the version from mod.sp file as well
}

@test "install latest and then run install" {
  powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance --force
  run powerpipe mod install
  assert_output 'All mods are up to date'
}

@test "install old version when latest already installed" {
  powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance --force
  run powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance@0.1
  assert_output '
Downgraded 1 mod:

local
└── github.com/turbot/steampipe-mod-aws-compliance@v0.1.0'
}

@test "install mod version, remove .cache file and then run install" {
  # install particular mod version, remove .mod.cache.json file and run mod install
  powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance@0.1 --force
  rm -rf .mod.cache.json
  run powerpipe mod install

  # should install the same cached version
  # better message
  assert_output '
Installed 1 mod:

local
└── github.com/turbot/steampipe-mod-aws-compliance@v0.1.0'
}

@test "install a mod with protocol in url" {
  run powerpipe mod install https://github.com/turbot/steampipe-mod-hackernews-insights@0.3.0 --force
  # should install with the protocol in the url prefix
  assert_output '
Installed 1 mod:

local
└── github.com/turbot/steampipe-mod-hackernews-insights@v0.3.0'
}

# Installed 4 mods:

# local
# └── github.com/pskrbasu/steampipe-mod-top-level@v3.0.0
#     ├── github.com/pskrbasu/steampipe-mod-dependency-1@v4.0.0
#     └── github.com/pskrbasu/steampipe-mod-dependency-2@v3.0.0
#         └── github.com/pskrbasu/steampipe-mod-dependency-1@v3.0.0
@test "complex mod dependency resolution - test tree structure" {
  run powerpipe mod install github.com/pskrbasu/steampipe-mod-top-level
  # test the tree structure output
  assert_output '
Installed 4 mods:

local
└── github.com/pskrbasu/steampipe-mod-top-level@v3.0.0
    ├── github.com/pskrbasu/steampipe-mod-dependency-1@v4.0.0
    └── github.com/pskrbasu/steampipe-mod-dependency-2@v3.0.0
        └── github.com/pskrbasu/steampipe-mod-dependency-1@v3.0.0'
}

@test "complex mod dependency resolution - test benchmark and controls resolution 1" {
  powerpipe mod install github.com/pskrbasu/steampipe-mod-top-level

  run powerpipe benchmark run top_level.benchmark.bm_version_dependency_mod_1 --output csv
  # check the output - benchmark should run the control and query from dependency mod 1 which will
  # have the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | 4.0    | 4.0      | alarm  |
# +--------+----------+--------+
  assert_output 'group_id,title,description,control_id,control_title,control_description,reason,resource,status,severity
top_level.benchmark.bm_version_dependency_mod_1,Benchmark version dependency mod 1,,dependency_1.control.version,,,4.0,4.0,alarm,'
}

@test "complex mod dependency resolution - test benchmark and controls resolution 2" {
  powerpipe mod install github.com/pskrbasu/steampipe-mod-top-level

  run powerpipe benchmark run top_level.benchmark.bm_version_dependency_mod_2 --output csv
  # check the output - benchmark should run the control and query from dependency mod 2 which will
  # have the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | 3.0    | 3.0      | ok     |
# +--------+----------+--------+
  assert_output 'group_id,title,description,control_id,control_title,control_description,reason,resource,status,severity
top_level.benchmark.bm_version_dependency_mod_2,Benchmark version dependency mod 2,,dependency_2.control.version,,,3.0,3.0,ok,'
}

function teardown() {
  rm -rf .powerpipe/
  rm -rf .mod.cache.json
  rm -rf mod.sp
}

function setup() {
  cd $FILE_PATH/test_data/mods/mod_install
}
