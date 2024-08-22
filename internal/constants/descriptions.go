package constants

const (
	PowerpipeShortDescription = "Powerpipe: Dashboards for DevOps."
	PowerpipeLongDescription  = `Powerpipe: Dashboards for DevOps.
	
Visualize cloud configurations. Assess security posture against a
massive library of benchmarks. Build custom dashboards with code.
	
Common commands:
	
  # Install a mod from the hub - https://hub.powerpipe.io
  powerpipe mod init
  powerpipe mod install github.com/turbot/steampipe-mod-aws-compliance

  # View dashboards at http://localhost:9033
  powerpipe server

  # List and run benchmarks in the terminal
  powerpipe benchmark list
  powerpipe benchmark run aws_compliance.benchmark.cis_v140
	
Documentation available at https://powerpipe.io/docs`
)
