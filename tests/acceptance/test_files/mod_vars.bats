load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

### workspace mod tests ###
# The following set of tests use a workspace mod(the mod is committed) with no dependencies and the
# variable has a default. The variable is then set from the command line, an auto ppvars file, an explicit ppvars 
# file, or an ENV var. The tests below check that the variable is resolved correctly in each of these cases.

@test "test variable resolution in workspace mod set from command line(--var)" {
  cd $FILE_PATH/test_data/mods/test_workspace_mod_var_set_from_command_line

  run powerpipe query run query.version --output csv --var version="v5.0.0"
  # check the output - query should use the value of variable set from the command line
  # --var flag ("v5.0.0") which will give the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | v5.0.0 | v5.0.0   | ok     |
# +--------+----------+--------+
  assert_output 'reason,resource,status
v5.0.0,v5.0.0,ok'
}

@test "test variable resolution in workspace mod set from powerpipe.ppvars file" {
  cd $FILE_PATH/test_data/mods/test_workspace_mod_var_set_from_powerpipe.ppvars

  run powerpipe query run query.version --output csv
  # check the output - query should use the value of variable set from the powerpipe ppvars
  # file ("v7.0.0") which will give the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | v7.0.0 | v7.0.0   | ok     |
# +--------+----------+--------+
  assert_output 'reason,resource,status
v7.0.0,v7.0.0,ok'
}

@test "test variable resolution in workspace mod set from *.auto.ppvars file" {
  cd $FILE_PATH/test_data/mods/test_workspace_mod_var_set_from_auto.ppvars

  run powerpipe query run query.version --output csv
  # check the output - query should use the value of variable set from the auto ppvars
  # file ("v7.0.0") which will give the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | v7.0.0 | v7.0.0   | ok     |
# +--------+----------+--------+
  assert_output 'reason,resource,status
v7.0.0,v7.0.0,ok'
}

@test "test variable resolution in workspace mod set from explicit ppvars file" {
  cd $FILE_PATH/test_data/mods/test_workspace_mod_var_set_from_explicit_ppvars

  run powerpipe query run query.version --output csv --var-file='deps.ppvars'
  # check the output - query should use the value of variable set from the explicit ppvars
  # file ("v8.0.0") which will give the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | v8.0.0 | v8.0.0   | ok     |
# +--------+----------+--------+
  assert_output 'reason,resource,status
v8.0.0,v8.0.0,ok'
}

@test "test variable resolution in workspace mod set from ENV" {
  cd $FILE_PATH/test_data/mods/test_workspace_mod_var_set_from_command_line
  export PP_VAR_version=v9.0.0
  run powerpipe query run query.version --output csv
  # check the output - query should use the value of variable set from the ENV var
  # PP_VAR_version ("v9.0.0") which will give the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | v9.0.0 | v9.0.0   | ok     |
# +--------+----------+--------+
  assert_output 'reason,resource,status
v9.0.0,v9.0.0,ok'
}

# ### dependency mod tests ###
# # The following set of tests use a dependency mod(the mod is committed) that has a variable dependency but the
# # variable does not have a default. This means that the variable must be set from the command
# # line, an auto ppvars file, an explicit ppvars file, or an ENV var. The tests below check that
# # the variable is resolved correctly in each of these cases.

# @test "test variable resolution in dependency mod set from command line(--var)" {
#   cd $FILE_PATH/test_data/mods/test_dependency_mod_var_set_from_command_line

#   run powerpipe query run dependency_vars_1.query.version --output csv --var dependency_vars_1.version="v5.0.0"
#   # check the output - query should use the value of variable set from the command line
#   # --var flag ("v5.0.0") which will give the output:
# # +--------+----------+--------+
# # | reason | resource | status |
# # +--------+----------+--------+
# # | v5.0.0 | v5.0.0   | ok     |
# # +--------+----------+--------+
#   assert_output 'reason,resource,status
# v5.0.0,v5.0.0,ok'
# }

# @test "test variable resolution in dependency mod set from steampipe.spvars file" {
#   cd $FILE_PATH/test_data/mods/test_dependency_mod_var_set_from_steampipe.spvars

#   run steampipe query dependency_vars_1.query.version --output csv
#   # check the output - query should use the value of variable set from the steampipe.spvars
#   # file ("v7.0.0") which will give the output:
# # +--------+----------+--------+
# # | reason | resource | status |
# # +--------+----------+--------+
# # | v7.0.0 | v7.0.0   | ok     |
# # +--------+----------+--------+
#   assert_output 'reason,resource,status
# v7.0.0,v7.0.0,ok'
# }

# @test "test variable resolution in dependency mod set from *.auto.spvars spvars file" {
#   cd $FILE_PATH/test_data/mods/test_dependency_mod_var_set_from_auto.spvars

#   run steampipe query dependency_vars_1.query.version --output csv
#   # check the output - query should use the value of variable set from the *.auto.spvars 
#   # file ("v8.0.0") which will give the output:
# # +--------+----------+--------+
# # | reason | resource | status |
# # +--------+----------+--------+
# # | v8.0.0 | v8.0.0   | ok     |
# # +--------+----------+--------+
#   assert_output 'reason,resource,status
# v8.0.0,v8.0.0,ok'
# }

# ### precedence tests ###

@test "test variable resolution precedence in workspace mod set from powerpipe.ppvars and *.auto.ppvars file" {
  cd $FILE_PATH/test_data/mods/test_workspace_mod_var_precedence_set_from_both_ppvars

  run powerpipe query run query.version --output csv
  # check the output - query should use the value of variable set from the  *.auto.ppvars("v8.0.0") file over 
  # powerpipe.ppvars("v7.0.0") which will give the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | v8.0.0 | v8.0.0   | ok     |
# +--------+----------+--------+
  assert_output 'reason,resource,status
v8.0.0,v8.0.0,ok'
}

@test "test variable resolution precedence in workspace mod set from powerpipe.ppvars and ENV" {
  cd $FILE_PATH/test_data/mods/test_workspace_mod_var_set_from_powerpipe.ppvars
  export PP_VAR_version=v9.0.0
  run powerpipe query run query.version --output csv
  # check the output - query should use the value of variable set from the powerpipe.ppvars("v7.0.0") file over 
  # ENV("v9.0.0") which will give the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | v7.0.0 | v7.0.0   | ok     |
# +--------+----------+--------+
  assert_output 'reason,resource,status
v7.0.0,v7.0.0,ok'
}

@test "test variable resolution precedence in workspace mod set from command line(--var) and powerpipe.ppvars file and *.auto.ppvars file" {
  cd $FILE_PATH/test_data/mods/test_workspace_mod_var_precedence_set_from_both_ppvars

  run powerpipe query run query.version --output csv --var version="v5.0.0"
  # check the output - query should use the value of variable set from the command line --var flag("v5.0.0") over 
  # powerpipe.ppvars("v7.0.0") and *.auto.ppvars file("v8.0.0") which will give the output:
# +--------+----------+--------+
# | reason | resource | status |
# +--------+----------+--------+
# | v5.0.0 | v5.0.0   | ok     |
# +--------+----------+--------+
  assert_output 'reason,resource,status
v5.0.0,v5.0.0,ok'
}

