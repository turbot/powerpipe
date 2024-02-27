[<picture><source media="(prefers-color-scheme: dark)" srcset="https://powerpipe.io/images/powerpipe_logo.svg"><source media="(prefers-color-scheme: light)" srcset="https://powerpipe.io/images/powerpipe_logo_darkmode.svg"><img width="67%" alt="Powerpipe Logo" src="https://steampipe.io/images/powerpipe_logo.svg"></picture>](https://powerpipe.io?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme)

[![plugins](https://img.shields.io/badge/apis_supported-137-blue)](https://hub.steampipe.io/plugins?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp; 
[![benchmarks](https://img.shields.io/badge/controls-4733-blue)](https://hub.steampipe.io/mods?objectives=compliance?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![dashboards](https://img.shields.io/badge/dashboards-708-blue)](https://hub.steampipe.io/mods?objectives=dashboard?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![slack](https://img.shields.io/badge/slack-1959-blue)](https://turbot.com/community/join?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![maintained by](https://img.shields.io/badge/maintained%20by-Turbot-blue)](https://turbot.com?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme)

Powerpipe is **dashboards and benchmarks as code**. Use it to visualize any data source, and run compliance benchmarks and controls, for effective decision-making and ongoing compliance monitoring.

**Dashboards**. Powerpipe runs dashboards that present data in visual and interactive ways.

**Benchmarks**. Powerpipe can also run benchmarks, which are packages of controls to check your infrastructure against compliance standards.

**Modularity and customization**. Dashboards and benchmarks are built with [HCL and SQL](https://steampipe.io/blog/dashboards-as-code), and are available for [remixing and reuse](https://steampipe.io/blog/remixing-dashboards). 

## Demo Time!

**[Watch on YouTube →](https://www.youtube.com/watch?v=TBD)**

[![Powerpipe demo](https://powerpipe.io/images/powerpipe_hero_video_thumbnail.png)](https://www.youtube.com/watch?v=TBD)

## Install Powerpipe

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

## Run a dashboard

Dashboards use charts, tables, and interactive <a href="https://powerpipe.io/docs/powerpipe-hcl?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme#dashboards">widgets</a> to help you explore and visualize your resources. For example, the <a href="RNAcentral">RNAcentrals</a> mod visualizes a dataset of RNA types. To run the RNACentral dashboard:


1. git clone https://github.com/turbot/powerpipe-mod-rnacentral

1. cd powerpipe-mod-rnacentral

1. powerpipe server --database postgres://reader:NWDMCE5xdipIjRrp@hh-pgsql-public.ebi.ac.uk:5432/pfmegrnargs

1. open localhost:9033 in a browser

That's it! Here's the dashboard.

![rnacentral](./images/rnacentral.png)


## Run a benchmark

Many Powerpipe [mods](https://hub.powerpipe.io.io/mods?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) includes **benchmarks** that check your cloud resources for compliance. The [Net Insights](https://hub.powerpipe.io/mods/turbot/net_insights?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) mod provides configuration, compliance and security controls to validate security best practices for DNS records. To run the [SSL/TLS Best Practices benchmark](https://hub.powerpipe.io/mods/turbot/net_insights/controls/benchmark.ssl_best_practices):

1. git clone https://github.com/turbot/powerpipe-mod-net-insights

1. cd powerpipe-mod-net-insights

1. powerpipe server

1. open localhost:9033 in a browser
``

Here's the dashboard.

![net insights dashboard](./images/net_insights_dashboard.png)

You can run benchmarks as dashboards too! Here's the console ouput.

```hcl
powerpipe benchmark run benchmark.ssl_certificate_best_practices --output=brief
```

![net insights console](./images/net_insights_console.png)

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

