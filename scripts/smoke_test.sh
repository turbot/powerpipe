#!/bin/sh
# This is a script with set of commands to smoke test a powerpipe build.
# The plan is to gradually add more tests to this script.

/usr/local/bin/steampipe service start # start the steampipe service
/usr/local/bin/steampipe plugin install net # install the net plugin

/usr/local/bin/powerpipe --version # check version

/usr/local/bin/powerpipe mod install github.com/turbot/steampipe-mod-net-insights # verify mod install 
/usr/local/bin/powerpipe mod list # verify mod list
/usr/local/bin/powerpipe mod uninstall github.com/turbot/steampipe-mod-net-insights # verify mod uninstall
/usr/local/bin/powerpipe mod list # verify mod list after uninstalling

/usr/local/bin/powerpipe mod install github.com/turbot/steampipe-mod-net-insights # re-install for other tests

/usr/local/bin/powerpipe control list # verify control list
/usr/local/bin/powerpipe control run net_insights.control.dns_mx_at_least_two # verify control run

# the file path is different for darwin and linux
if [ "$(uname -s)" = "Darwin" ]; then
  /usr/local/bin/powerpipe control run net_insights.control.dns_mx_at_least_two --export /Users/runner/query.pps # verify file export
  cat /Users/runner/query.pps | jq '.end_time' # verify file created is readable
else
  /usr/local/bin/powerpipe control run net_insights.control.dns_mx_at_least_two --export /home/steampipe/query.pps # verify file export
  cat /home/steampipe/query.pps | jq '.end_time' # verify file created is readable
fi
