## Kind: environment

# Issue #1170: add-to-project / add-to-project red on develop

## What I was doing
Tracing why `.github/workflows/31-add-issues-to-pipeling-issue-tracker.yaml` fails on every `issues: opened` run.

## Finding
- The workflow only calls `turbot/steampipe-workflows/.github/workflows/assign-issue-to-pipeling-issue-tracker.yml@main` with `secrets: inherit`. Nothing in this repo sets the project URL, the token or the `fields` input.
- The real error in the log is `Could not resolve to a ProjectV2 with the number 31`. The `Unexpected input(s) 'fields'` line is only a warning. The failure is the `GH_PROJECT_PAT` secret (expired, revoked, or without access to org project 31), or the project having moved or been deleted.
- Run history: success through 2026-03-31, failure on every run since 2026-04-15, with no change to this repo's workflow file in between (last touched by #1034).
- The "token has only read permissions" in the issue refers to the job's GITHUB_TOKEN, which is not what the action uses; it uses `GH_PROJECT_PAT`.

## Ruled out
- A change to this repo's workflow file: it passes nothing that could fix the failure.
- Reading or listing secrets: this worker's token gets 403 on repo and org secrets, and that is outside the run's scope anyway.
- The workflow lives under `.github/workflows/`, and any edit would be parked.

## Next step
A person with org admin rights on turbot should (1) renew `GH_PROJECT_PAT` (org or repo secret, including for turbot/powerpipe) with `project` scope and access to org project 31, and (2) confirm project 31 still exists. If the "Pipeling issue tracker" project is retired, delete this workflow instead (should this be closed as won't fix?). The `fields` input to `actions/add-to-project@v1.0.2` is also unsupported and can be dropped from steampipe-workflows.
