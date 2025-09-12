<a href="https://powerpipe.io"><img width="67%" src="https://powerpipe.io/images/powerpipe_wordmark_white_outline.svg"></a>

[![mods](https://img.shields.io/endpoint?url=https://turbot.com/api/badge-stats?stat=mods)](https://hub.powerpipe.io) &nbsp;
[![slack](https://img.shields.io/endpoint?url=https://turbot.com/api/badge-stats?stat=slack)](https://turbot.com/community/join) &nbsp;
[![maintained by](https://img.shields.io/badge/maintained%20by-Turbot-blue)](https://turbot.com)

## Dashboards for DevOps

Use [Powerpipe](https://powerpipe.io) to visualize any data source and run compliance benchmarks and controls that enable effective decision-making and ongoing compliance monitoring.

**Dashboards and reports**. High-level dashboards provide a quick overview. Use them to highlight misconfigurations and hotspots. [Filter](https://powerpipe.io/docs/run/benchmark/benchmark-dashboard#filtering--grouping), pivot, and [snapshot](https://powerpipe.io/docs/run/snapshots) results.

**Benchmarks**. We offer [5,000+ open-source controls](https://hub.powerpipe.io) from CIS, NIST, PCI, HIPAA, FedRamp and more. [Run instantly on your machine](https://powerpipe.io/docs#run-security-and-compliance-benchmarks) or as part of your deployment pipeline.

**Relationship diagrams**. The only dashboarding tool designed from the ground up to [visualize DevOps data](https://powerpipe.io/docs#visualize-cloud-infrastructure). Explore your cloud, understand relationships, drill down to the details.

**Code, not clicks**. Our dashboards are [code](https://powerpipe.io/docs/build): version-controlled, composable, shareable, easy to edit — designed for the way you work. 

## Demo time!

**[Watch on YouTube →](https://www.youtube.com/watch?v=-h6RSpvR0FE)**

<a href="https://www.youtube.com/watch?v=-h6RSpvR0FE"><img alt="powerpipe demo" width=500 src="https://powerpipe.io/images/powerpipe_hero_video_thumbnail.png"></a>

## Documentation

See the [documentation](https://powerpipe.io/docs) for:

- [Viewing dashboards](https://powerpipe.io/docs/run/dashboard)
- [Running benchmarks](https://powerpipe.io/docs/run/benchmark)
- [CLI commands](https://powerpipe.io/docs/reference/cli)
- [HCL reference](https://powerpipe.io/docs/powerpipe-hcl)
- [Configuration](https://powerpipe.io/docs/reference/config-files)
- [Building mods](https://powerpipe.io/docs/build)


## Install Powerpipe

Install Powerpipe from the [downloads](https://powerpipe.io/downloads) page:

```sh
# MacOS
brew install turbot/tap/powerpipe
```

```sh
# Linux or Windows (WSL2)
sudo /bin/sh -c "$(curl -fsSL https://powerpipe.io/install/powerpipe.sh)"
```

Now, [set up and visualize your first dashboard →](https://powerpipe.io/docs)

## Powerpipe mods: dashboards and benchmarks

Powerpipe [mods](https://hub.powerpipe.io) are sets of pre-built dashboards that visualize your resources and benchmarks that check your cloud resources for compliance. Ready to use mods are available for [AWS](https://hub.powerpipe.io/?q=aws), [Azure](https://hub.powerpipe.io/?q=azure), [GCP](https://hub.powerpipe.io/?q=gcp), [GitHub](https://hub.powerpipe.io/?q=github), [Kubernetes](https://hub.powerpipe.io/?q=kubernetes), [Terraform](https://hub.powerpipe.io/?q=terraform), [M365](https://hub.powerpipe.io/mods/turbot/microsoft365_compliance) and much more to cover common use cases for [security & compliance](https://hub.powerpipe.io/?objectives=compliance), [cost management](https://hub.powerpipe.io/?objectives=cost), [shift-left scanning](https://hub.powerpipe.io/?categories=iac), and [asset inventory and insights](https://hub.powerpipe.io/?objectives=dashboard).

Many mods are written for Steampipe and its plugin ecosystem. But Powerpipe is database-agnostic, and we provide samples for dashboards that use other data sources via [Postgres, SQLite, DuckDB, and MySQL](https://hub.powerpipe.io/?engines=postgres,duckdb,sqlite,mysql).

## Developing

If you want to help develop the core Powerpipe binary, these are the steps to build it.

<details>
<summary>Clone</summary>

Clone [github.com/powerpipe](https://github.com/turbot/powerpipe) and [github.com/turbot/pipe-fittings](https://github.com/turbot/pipe-fittings):

```sh
git clone git@github.com:turbot/powerpipe
git clone git@github.com:turbot/pipe-fittings
```

</details>

<details>
<summary>Build</summary>

```sh
cd powerpipe
make
```

The Powerpipe binary lands in `/usr/local/bin` unless you specify an alternate `OUTPUT_DIR`.

</details>

<details>
<summary>Check the install</summary>

```sh
powerpipe --version

powerpipe --help
```
</details>

If you're interested in developing [Powerpipe mods](https://hub.powerpipe.io), see our [documentation for mod developers](https://powerpipe.io/docs/build).

## Turbot Pipes

Bring your team to [Turbot Pipes](https://turbot.com/pipes) to use Powerpipe together in the cloud. In a Pipes workspace you can use [Steampipe](https://github.com/turbot/steampipe) for data access, Powerpipe to visualize query results, and [Flowpipe](https://github.com/turbot/flowpipe) to automate workflow.

## Open source and contributing

This repository is published under the [AGPL 3.0](https://www.gnu.org/licenses/agpl-3.0.html) license. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). Contributors must sign our [Contributor License Agreement](https://turbot.com/open-source#cla) as part of their first pull request. We look forward to collaborating with you!

[Powerpipe](https://powerpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get involved

**[Join #powerpipe on Slack →](https://turbot.com/community/join)**
