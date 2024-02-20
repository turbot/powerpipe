# test_workspace_mod_var_set_from_explicit_ppvars

### Description

This mod is used to test variable resolution in a mod by passing the variable value in an explicit ppvars file. The mod has a default value of variable 'version' set.

### Usage

This mod is used in the tests in `mod_vars.bats` to simulate a scenario where the version defined in the mod is picked from the passed
variable value in an explicit ppvars file over the default value of variable 'version' set in the mod. 