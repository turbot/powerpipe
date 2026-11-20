---
name: release
description: Release powerpipe through the shared release workflows - the prep file, the cut and the release PR, the approval, watching the ship (tag, Homebrew, powerpipe.io page, distro test, merge-back), a re-ship, and the hand-over to the Pipes release.
---

# Release powerpipe

powerpipe releases through Turbot's shared release workflows (turbot/shipyard-workflows, called at `v0`). Nobody types a
version, pushes a tag, edits `CHANGELOG.md` or opens a Homebrew or powerpipe.io PR by hand. Every PR into `develop` that
ships carries a `Changelog:` line, and those lines are the release's changelog.

Powerpipe ships two things under one version: the CLI binary, and the dashboard UI built from `ui/dashboard/` and embedded
in it. `make dashboard_assets_archive` builds the UI before goreleaser, in the release PR's results and in the ship build.

## 1. Start and prep

`release-start.yml` files a prep job (an issue labelled `agent:release`) on Tuesdays at 09:00 UTC when a shipped change
waits, or when run by hand on `develop`. Dependency bumps alone start nothing, unless one fixes an advisory.

Write the prep file with the skill the prep job names, for the commit it names. Then a person runs
`prep.yml` with the file's content, and `release.yml` on `develop` with `commit` (the full id) and `prep-run` (the prep run's
id). The cut makes `release/<version>` with the CHANGELOG section and opens the release PR into `main`.

## 2. The release PR

`release-pr.yml` tests the candidate (`make release_test`), builds a snapshot, installs it and writes the record. With a
SHIP record, the record marks the PR ready for review and starts `ask.yml`, which asks in #powerpipe within a minute.
`main` accepts the merge only when the record's `shipyard/release` status passes.
One approver of `turbot/pipelings-release-approvers` who drove none of its changes approves it. Merge with a merge commit.

## 3. Watch the ship

`ship.yml` runs on the merge:

```
ship (build at the merge) ──► tag (tag, GitHub Release with the archives, merge-back PR into develop)
                                ├──► homebrew (turbot/homebrew-tap: powerpipe.rb and powerpipe@<version>.rb)
                                ├──► page (turbot/powerpipe.io: content/changelog/powerpipe-cli-v<x-y-z>.md)
                                └──► distros (12 - Test: Linux Distros (Post-release))
```

```bash
gh run list --repo turbot/powerpipe --workflow ship.yml --limit 1
gh release view v{x.y.z} --repo turbot/powerpipe
brew update && brew install turbot/tap/powerpipe && powerpipe --version
```

It is done when Homebrew installs the new version, `powerpipe server` loads a dashboard, and the page is live at
`https://powerpipe.io/changelog/powerpipe-cli-v{x}-{y}-{z}`. A failed `homebrew` or `page` job is finished by running it
again; it writes nothing twice. If the page job passed and the page is still not live, run the `Trigger Vercel deploy (prod)`
workflow in turbot/powerpipe.io from `main`.

**A release merged but not shipped** (the ship job itself is broken): `main` takes only release PRs, so run `release.yml`
on `develop` with `reship` (the version), `commit` (the tip of `main`) and `fixes` (the merge commits of the fix on
`develop`). A fix that changes `.github/workflows/` cannot go in `fixes`: the cut says how to pick it by hand.

## 4. Hand over

Tell whoever runs the Turbot Pipes release that the new version is live, with its tag (`v1.6.1`). Pipes pins both the
CLI and the dashboard UI by that tag, and the new version ships in the next Pipes release.
