[<picture><source media="(prefers-color-scheme: dark)" srcset="https://powerpipe.io/images/powerpipe_logo.svg"><source media="(prefers-color-scheme: light)" srcset="https://powerpipe.io/images/powerpipe_logo_darkmode.svg"><img width="67%" alt="Powerpipe Logo" src="https://steampipe.io/images/powerpipe_logo.svg"></picture>](https://powerpipe.io?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme)

[![plugins](https://img.shields.io/badge/apis_supported-137-blue)](https://hub.steampipe.io/plugins?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp; 
[![benchmarks](https://img.shields.io/badge/controls-4733-blue)](https://hub.steampipe.io/mods?objectives=compliance?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![dashboards](https://img.shields.io/badge/dashboards-708-blue)](https://hub.steampipe.io/mods?objectives=dashboard?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![slack](https://img.shields.io/badge/slack-1959-blue)](https://turbot.com/community/join?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![maintained by](https://img.shields.io/badge/maintained%20by-Turbot-blue)](https://turbot.com?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme)

Powerpipe is **dashboards and benchmarks as code**. Use it to visualize any data source, and run compliance benchmarks and controls, for effective decision-making and ongoing compliance monitoring.

**Dashboards**. Powerpipe runs dashboards that present data in visual and interactive ways.

**Benchmarks**. Powerpipe can also run benchmarks, which are packages of controls to check your infrastucture against compliance standards.

**Modularity and customization**. Dashboards and benchmarks are built with [HCL and SQL](https://steampipe.io/blog/dashboards-as-code), and are available for [remixing and reuse](https://steampipe.io/blog/remixing-dashboards). 

## Demo Time!

**[Watch on YouTube →](https://www.youtube.com/watch?v=TBD)**

[![Powerpipe demo](https://powerpipe.io/images/powerpipe_hero_video_thumbnail.png)](https://www.youtube.com/watch?v=TBD)

## Getting Started

The <a href="https://powerpipe.io/downloads?utm_id=gfpreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme">downloads</a> page shows you how, but tl;dr:

Linux or WSL

```sh
sudo /bin/sh -c "$(curl -fsSL https://powerpipe.io/install/powerpipe.sh)"
```

MacOS

```sh
brew tap turbot/tap
brew install powerpipe
```

## Powerpipe: Dashboards and benchmarks

Powerpipe [mods](https://hub.powerpipe.io.io/mods?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) which are sets of **dashboards** that visualize your resources and **benchmarks** that check your cloud resources for compliance.

<table>
  <tr>
   <td><b>Compliance</b></td>
   <td>Check AWS, Azure, GCP, etc for compliance with HIPAA, PCI, etc
  </tr>
  <tr>
   <td><b>Cost</b></td>
   <td>Review what AWS, Azure, GCP, and other clouds are costing you</td>
  </tr>
  <tr>
   <td><b>Insights</b></td>
   <td>Visualize cloud resources with charts, tables, and interactive widgets</td>
  </tr>
  <tr>
   <td><b>Security</b></td>
   <td>Use CIS, NIST, FedRAMP etc to assess the security of AWS, Azure, GCP, etc</td>
  </tr>
  <tr>
   <td><b>Tags</b></td>
   <td>Verify the consistency of tags applied to AWS, Azure, and GCP resources</td>
  </tr>
  <tr>
   <td><b>Your mod</b></td>
   <td>Build your own <a href="https://powerpoint.io/docs/build?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme">benchmarks and dashboards</a></td>
  </tr>
 </table>


Dashboards and benchmarks use SQL to gather data and HCL to flow the data into [dashboard widgets](https://steampipe.io/blog/dashboards-as-code?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) and [benchmark controls](https://steampipe.io/blog/release-0-11-0?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme#composable-mods). You can use the existing suites of benchmarks and dashboards, or build derivative versions, or create your own. 

### Get started with dashboards and benchmarks

<details>
<summary>Install the Net Insights mod</summary>
<br/>
The <a href="https://hub.powerpipe.io/mods/turbot/net_insights?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme">Net Insights</a> mod works with the Net plugin shown above. To run it, first clone its repo and change to that directory.

```sh
git clone https://github.com/turbot/powerpipe-mod-net-insights
cd powerpipe-mod-net-insights
```
</details>

<details>
<br/>
<summary>Run benchmarks in the CLI</summary>

All the benchmarks:

```sh
powerpipe benchmark run all
```

A single benchmark:

```sh
powerpipe benchmark run benchmark.dns_best_practices
```

A single control:

```sh
powerpipe control  run control.dns_ns_name_valid
```
 
Available <a href="docs/reference/cli/benchmark#powerpipe-benchmark-run">formats</a> include JSON, CSV, HTML, and ASFF. 

You can use <a href="https://steampipe.io/docs/develop/writing-control-output-templates?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme">custom output templates</a> to create new output formats.
</details>

<details>
<summary>Run benchmarks as dashboards</summary>
<br/>
Launch the dashboard server: `powerpipe server`, then open `http://localhost:9033` in your browser. 

The home page lists available dashboards. Click `DNS Best Practices` to view that dashboard.

Note that the default domains are `microsoft.com` and `github.com`. You can <a href="https://hub.steampipe.io/mods/turbot/net_insights?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme#configuration">change those defaults</a> to check other domains.
</details>

<details>
<summary>Use dashboards to explore your resources</summary>
<br/>
Dashboards use charts, tables, and interactive <a href="https://powerpipe.io/docs/powerpipe-hcl?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme#dashboards">widgets</a> to help you explore and visualize your resources. 

The <a href="https://hub.powerpipe.io/mods/turbot/aws_insights?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme">AWS Insights</a> mod, for example, provides dozens of dashboards that exercise the full set of widgets. To use these dashboards, first install the <a href="https://hub.powerpipe.io/plugins/turbot/aws?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme">AWS plugin</a> and <a href="https://hub.powerpipe.io/plugins/turbot/aws?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme#configuration">authenticate</a>. Then clone `AWS Insights`, change to its directory, launch `steampipe dashboard`, and open `localhost:9033`.
</details>


## Developing

<details>
<summary>Developing Powerpipe</summary>

Prerequisites:

- [Golang](https://golang.org/doc/install) Version 1.21 or higher.

Clone `github.com/powerpipe` and `github.com/turbot/pipe-fittings` repositories:

```sh
git clone https://github.com/turbot/powerpipe
git clone https://github.com/turbot/pipe-fittings
cd powerpipe
```

The build lands in `/usr/local/bin/` unless `OUTPUT_DIR` is specified.

```sh
make
```

Check the version:
```sh
powerpipe --version
```
```
Powerpipe v0.1.0-local.1
```

</details>

## Open Source & Contributing
This repository is published under the [AGPL 3.0](https://www.gnu.org/licenses/agpl-3.0.html) license. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). Contributors must sign our [Contributor License Agreement](https://turbot.com/open-source#cla) as part of their first pull request. We look forward to collaborating with you!

[Powerpipe](https://powerpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #powerpipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:
* [Powerpipe](https://github.com/turbot/powerpipe/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

