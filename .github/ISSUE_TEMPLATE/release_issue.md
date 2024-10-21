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

- [ ] All PR acceptance test pass in `powerpipe`
- [ ] Update check is working
- [ ] Powerpipe version is correct
- [ ] Powerpipe Changelog updated and reviewed

### Release Powerpipe

- [ ] Merge the release PR
- [ ] Trigger the `02 - Powerpipe: Release` workflow. This will create the release build.

### Post-release checks

- [ ] Update Changelog in the Release page (copy and paste from CHANGELOG.md)
- [ ] Test Linux install script
- [ ] Test Homebrew install
- [ ] Release branch merged to `develop`
- [ ] Raise Changelog update to `powerpipe.io`, get it reviewed.
- [ ] Merge Changelog update to `powerpipe.io`.