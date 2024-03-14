<a href="https://powerpipe.io?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme"><img width="67%" src="https://powerpipe.io/images/powerpipe_wordmark.svg"></a>

[![mods](https://img.shields.io/badge/mods-52-blue)](https://hub.steampipe.io/mods?objectives=dashboard?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![slack](https://img.shields.io/badge/slack-2695-blue)](https://turbot.com/community/join?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) &nbsp;
[![maintained by](https://img.shields.io/badge/maintained%20by-Turbot-blue)](https://turbot.com?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme)

[Powerpipe](https://powerpipe.io) is **Dashboards for DevOps**. Use it to visualize any data source, and run compliance benchmarks and controls, for effective decision-making and ongoing compliance monitoring.

**Dashboards and reports**. High-level dashboards provide a quick overview. Use them to highlight misconfigurations and hotspots. [Filter](https://powerpipe.io/docs/run/benchmark/benchmark-dashboard#filtering--grouping), pivot, and [snapshot](https://powerpipe.io/docs/run/snapshots) results.

**Benchmarks**. We offer [5,000+ open-source controls](https://hub.powerpipe.io) from CIS, NIST, PCI, HIPAA, FedRamp and more. [Run instantly on your machine](https://powerpipe.io/docs#run-security-and-compliance-benchmarks) or as part of your deployment pipeline.

**Relationship diagrams**. The only dashboarding tool designed from the ground up to [visualize DevOps data](https://powerpipe.io/docs#visualize-cloud-infrastructure). Explore your cloud, understand relationships, drill down to the details.

**Code, not clicks**. Our dashboards are [code](https://powerpipe.io/docs/build): version-controlled, composable, shareable, easy to edit — designed for the way you work. 

## Demo time!

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

Now, [set up and visualize your first dashboard →](https://powerpipe.io/docs)

## Powerpipe mods: dashboards and benchmarks

Powerpipe [mods](https://hub.powerpipe.io) are sets of pre-built dashboards that visualize your resources and benchmarks that check your cloud resources for compliance. Ready to use mods are available for [AWS](https://hub.powerpipe.io/?q=aws), [Azure](https://hub.powerpipe.io/?q=azure), [GCP](https://hub.powerpipe.io/?q=gcp), [GitHub](https://hub.powerpipe.io/?q=github), [Kubernetes](https://hub.powerpipe.io/?q=kubernetes), [Terraform](https://hub.powerpipe.io/?q=terraform), [M365](https://hub.powerpipe.io/mods/turbot/microsoft365_compliance) and much more to cover common use cases for [security & compliance](https://hub.powerpipe.io/?objectives=compliance), [cost management](https://hub.powerpipe.io/?objectives=cost), [shift-left scanning](https://hub.powerpipe.io/?categories=iac), and [asset inventory and insights](https://hub.powerpipe.io/?objectives=dashboard).

Many mods are written for Steampipe and its plugin ecosystem. But Powerpipe is database-agnostic, and we provide samples for dashboards that use other data sources via [Postgres, SQLite, DuckDB, and MySQL](https://hub.powerpipe.io/?engines=postgres,duckdb,sqlite,mysql).


## Open source and contributing
This repository is published under the [AGPL 3.0](https://www.gnu.org/licenses/agpl-3.0.html) license. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). Contributors must sign our [Contributor License Agreement](https://turbot.com/open-source#cla) as part of their first pull request. We look forward to collaborating with you!

[Powerpipe](https://powerpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Developing

<details>
<summary>Developing Powerpipe</summary>

Prerequisites:

- [Golang](https://golang.org/doc/install) Version 1.21 or higher.

Clone `github.com/powerpipe` and `github.com/turbot/pipe-fittings` repositories:

```sh
git clone git@github.com:turbot/powerpipe
git clone git@github.com:turbot/pipe-fittings
cd powerpipe
```

Build will build powerpipe binary in /usr/local/bin/ unless `OUTPUT_DIR` is specified:

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

## Turbot Pipes

 Bring your team to [Turbot Pipes](https://turbot.com/pipes?utm_id=gspreadme&utm_source=github&utm_medium=repo&utm_campaign=github&utm_content=readme) to use Powerpipe together in the cloud.

## Get involved

**[Join #powerpipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:
* [Powerpipe](https://github.com/turbot/powerpipe/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

