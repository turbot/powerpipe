<a href="https://powerpipe.io?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme"><img width="67%" src="https://powerpipe.io/images/powerpipe_wordmark.svg"></a>

[![mods](https://img.shields.io/badge/mods-52-blue)](https://hub.steampipe.io/mods?objectives=dashboard?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![slack](https://img.shields.io/badge/slack-2297-blue)](https://turbot.com/community/join?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![maintained by](https://img.shields.io/badge/maintained%20by-Turbot-blue)](https://turbot.com?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme)

Powerpipe is **dashboards and benchmarks as code**. Use it to visualize any data source, and run compliance benchmarks and controls, for effective decision-making and ongoing compliance monitoring.

**Benchmarks** - [5,000+ open-source controls](https://hub.powerpipe.io) from CIS, NIST, PCI, HIPAA, FedRamp and more. [Run instantly on your machine](https://powerpipe.io/docs#run-security-and-compliance-benchmarks) or as part of your deployment pipeline.

**Relationship Diagrams** - The only dashboarding tool designed from the ground up to [visualize DevOps data](https://powerpipe.io/docs#visualize-cloud-infrastructure). Explore your cloud, understand relationships, drill down to the details.

**Dashboards & Reports** - High-level dashboards provide a quick overview. Use them to highlight misconfigurations and hotspots. [Filter](https://powerpipe.io/docs/run/benchmark/benchmark-dashboard#filtering--grouping), pivot, and [snapshot](https://powerpipe.io/docs/run/snapshots) results.

**Code, not clicks** - Our dashboards are [code](https://powerpipe.io/docs/powerpipe-hcl): version-controlled, composable, shareable, easy to edit — designed for the way you work. [Join](https://turbot.com/community/join) our open-source community!

## Demo Time!

**[Watch on YouTube →](https://www.youtube.com/watch?v=-h6RSpvR0FE)**

[![Powerpipe demo](https://powerpipe.io/images/powerpipe_hero_video_thumbnail.png)](https://www.youtube.com/watch?v=-h6RSpvR0FE)

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

## Your first dashboard

Dashboards use charts, tables, and interactive <a href="https://powerpipe.io/docs/powerpipe-hcl?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme#dashboards">widgets</a> to help you explore and visualize your resources. For example, the <a href="RNAcentral">RNAcentrals</a> mod visualizes a dataset of RNA types. To run the RNACentral dashboard:

```sh
mkdir rna && cd rna
powerpipe mod install https://github.com/turbot/powerpipe-mod-rnacentral
powerpipe server --database postgres://reader:NWDMCE5xdipIjRrp@hh-pgsql-public.ebi.ac.uk:5432/pfmegrnargs
```

View your dashboard at http://localhost:9033

![rnacentral](./images/rnacentral.png)

## Dashboards and Benchmarks with Steampipe

Many Powerpipe [mods](https://hub.powerpipe.io.io/mods?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) includes **benchmarks** that check your cloud resources for compliance. The [Net Insights](https://hub.powerpipe.io/mods/turbot/net_insights?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) mod provides configuration, compliance and security controls to validate security best practices for DNS records. To run the [SSL/TLS Best Practices benchmark](https://hub.powerpipe.io/mods/turbot/net_insights/controls/benchmark.ssl_best_practices):

```sh
# Net insights uses the net plugin from Steampipe
steampipe plugin install net
steampipe service start

# Install the Powerpipe mod
mkdir net && cd net
powerpipe mod install https://github.com/turbot/powerpipe-mod-net-insights

# Start the dashboard server
powerpipe server
```

View your dashboard at http://localhost:9033

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

