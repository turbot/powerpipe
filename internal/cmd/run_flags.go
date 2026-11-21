package cmd

// runModFlagsDeprecated is the warning for --mod-install and --pull on a run command. A run never installs or
// updates mods (#1211); both flags are still accepted, so a script that passes them keeps working.
const runModFlagsDeprecated = "a run does not install or update mods: use 'powerpipe mod install' or 'powerpipe mod update' first"
