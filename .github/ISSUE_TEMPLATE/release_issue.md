---
name: Powerpipe Release
about: Powerpipe Release
title: "Powerpipe v<INSERT_VERSION_HERE>"
labels: release
---

#### Changelog

[Powerpipe v<INSERT_VERSION_HERE> Changelog](https://github.com/turbot/powerpipe/blob/v<INSERT_VERSION_HERE>/CHANGELOG.md)

## Checklist

### Pre-release checks

- [ ] The release PR's page shows its checks passed and what changes
- [ ] Update check is working

### Release Powerpipe

- [ ] A pipelings release approver approves the release PR and merges it. The merge builds, tags, publishes the GitHub Release with its notes, updates Homebrew, writes the changelog page on powerpipe.io, starts the Linux distro test and merges `main` back into `develop`.

### Post-release checks

- [ ] Test Linux install script
- [ ] Test Homebrew install, and that the dashboard loads in it
- [ ] Check the changelog page is live on powerpipe.io
- [ ] Tell the Turbot Pipes release the new tag: Pipes pins the CLI and the dashboard UI by it
