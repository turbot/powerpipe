# Powerpipe

Powerpipe is "dashboards as code," running controls and benchmarks for compliance checks. It transforms complex data into visual, interactive dashboards, providing real-time insights for efficient decision-making and ongoing compliance monitoring, all in a streamlined, code-driven format.

**Running Controls**. Powerpipe can execute controls and benchmarks, which are predefined checks or tests against specific compliance or quality standards. This automation ensures that compliance is continually monitored and maintained.

**Real-Time Dashboards**. Powerpipe can run real-time dashboards that visually represent compliance and performance data. These dashboards provide immediate insights and allow for quick decision-making.

**Customizable and Interactive**. Each dashboard is fully customizable, catering to the specific needs and preferences of different teams or departments within an organization. The interactive nature of these dashboards allows for a deeper dive into the data, offering granular insights.

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

## Open Source & Contributing
This repository is published under the [AGPL 3.0](https://www.gnu.org/licenses/agpl-3.0.html) license. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). Contributors must sign our [Contributor License Agreement](https://turbot.com/open-source#cla) as part of their first pull request. We look forward to collaborating with you!

[Powerpipe](https://powerpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #powerpipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:
* [Powerpipe](https://github.com/turbot/powerpipe/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

