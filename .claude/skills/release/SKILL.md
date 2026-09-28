---
name: release
description: Cut a powerpipe CLI patch release from a v{maj}.{min}.x release branch - the release workflow (which also builds the dashboard UI), verification and release notes, the merge-back, main and powerpipe.io PRs, and the hand-over to the Pipes release.
---

# Release powerpipe

Release branches are `v{maj}.{min}.x` (e.g. `v1.5.x`). A release is a tag cut from the branch head by the
`01 - Powerpipe: Release` workflow. Never tag through the GitHub Releases UI: it creates the tag but
does not build the dashboard UI assets, run goreleaser, publish binaries, or open the homebrew-tap PR.

Placeholders: `{x.y.z}` is the new version (e.g. `1.5.6`), `{prev}` the previous tag (e.g. `v1.5.5`,
from `gh release list --repo turbot/powerpipe --limit 1 --exclude-pre-releases`), `{x}-{y}-{z}` and `{xyz}` the same version
dash-separated (`1-5-6`) and with no separators (`156`).

Every PR you open below is opened under your own `gh` auth and needs a teammate's approval before it
merges. The exceptions: the `turbot/homebrew-tap` PR, which the workflow opens and merges itself, and
the powerpipe.io PR, whose `main` requires no review.

Powerpipe ships two artifacts under one version: the CLI binary, and the dashboard UI built from
`ui/dashboard/`. The release workflow builds both in one run (no separate FDW-style pre-step, and
`go.mod`'s only steampipe reference is the plugin-sdk library, not a release dependency) — but the
hand-over to Pipes needs both version fields, and the UI one is easy to leave pointed at `develop`.

## 1. Confirm what is being released

```bash
gh api 'repos/turbot/powerpipe/compare/{prev}...v{maj}.{min}.x' \
  -q '.commits[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"'
```

- Every intended fix is in the list. For a security fix, credit the commit that actually changed the
  dependency, not an adjacent PR titled like a batch dep bump.
- `CHANGELOG.md` on `v{maj}.{min}.x` has an entry for `v{x.y.z}` in the style of earlier entries, dated
  today, committed with the message `v{x.y.z}` (via a PR into the release branch).
- Open the release issue from `.github/ISSUE_TEMPLATE/release_issue.md`: title `Powerpipe v{x.y.z}`,
  label `release`. Its checklist names older workflow numbering; follow this skill's order instead, and
  only tick its manual checks (Homebrew install, Linux install script, update check) once you've run them.

## 2. Dispatch the release workflow

```bash
gh workflow run 01-powerpipe-release.yaml --repo turbot/powerpipe --ref v{maj}.{min}.x \
  -f environment='Final (RC and final release)' -f version={x.y.z} -f confirmDevelop=false
sleep 15; gh run list --repo turbot/powerpipe --workflow 01-powerpipe-release.yaml --limit 3 --json databaseId,createdAt,headBranch
gh run watch --repo turbot/powerpipe <run-id>
```

Take the run whose `createdAt` is after your dispatch; if none is, list again. `version` has no `v` prefix;
the workflow adds it. `confirmDevelop` is required but not read by any job; pass `false`.
`Development (alpha)` / `Development (beta)` are for pre-release test builds only.

The `build_assets` job builds `ui/dashboard/` and runs its unit tests before `build_and_release` tags
and ships the CLI — a dashboard UI failure blocks the whole release, not just the UI. Later jobs open
and merge the homebrew-tap PR and dispatch `12 - Test: Linux Distros (Post-release)`.

## 3. Verify and publish the release notes

```bash
gh release view v{x.y.z} --repo turbot/powerpipe
gh pr list --repo turbot/homebrew-tap --state merged --limit 5
gh run list --repo turbot/powerpipe --workflow 12-test-post-release-linux-distros.yaml --limit 1
```

The release has its binaries, the homebrew-tap PR for `{x.y.z}` is merged, and the post-release test run
passed. The release is published with an empty body (goreleaser's changelog is disabled): edit it and
paste in the `CHANGELOG.md` entry for `v{x.y.z}`.

## 4. PRs

1. `v{maj}.{min}.x` into `develop`, titled `Merge branch 'v{maj}.{min}.x' into develop`. If the branches
   conflict (usually `go.mod`/`go.sum`), open it from a branch cut off `develop` that merges
   `origin/v{maj}.{min}.x` with the conflicts resolved (keep the higher version of each dependency), then `go mod tidy`.
2. `v{maj}.{min}.x` into `main`, titled `Release Powerpipe v{x.y.z}`, label `release`. Never merge `main` into
   the release branch: if they conflict, open it from a branch cut off `v{maj}.{min}.x` that merges
   `origin/main`, keeping `main`'s action SHA pins and the release branch's `go-version`. Body:
   ```
   ## Release Issue
   [Powerpipe v{x.y.z}](<release issue URL>)

   ## Checklist
   - [ ] Confirmed that version has been correctly upgraded.
   ```
3. `turbot/powerpipe.io`: branch `pp-{xyz}` off `main`, add
   `content/changelog/<YYYYMMDD>-powerpipe-cli-v{x}-{y}-{z}.md` (flat directory, no year subfolder —
   unlike steampipe.io):
   ```
   ---
   title: Powerpipe CLI v{x.y.z} - <short summary>
   publishedAt: "<YYYY-MM-DD>T10:00:00"
   permalink: powerpipe-cli-v{x}-{y}-{z}
   tags: cli
   ---
   ```
   Body matches the `CHANGELOG.md` entry. PR title `Powerpipe CLI v{x.y.z}`, base `main`.
   Once merged, run the `Trigger Vercel deploy (prod)` workflow from `main` and check the page loads.
   Vercel blocks the deploy if the merge commit's author isn't a member of the team that owns the
   project — if it doesn't go live, a Vercel team member needs to push again or approve access.

## 5. Hand over

Give the new version to whoever runs the Turbot Pipes release, as both `powerpipeCliVersion` and
`powerpipeDashboardUiVersion` — check no Pipes workflow still points the UI version at `develop` from
interim UI work. Tick the release issue's checklist and close it.
