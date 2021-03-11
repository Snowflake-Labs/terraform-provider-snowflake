# Release

## Running a release

**Note: releases can only be done by those with keybase pgp keys whitelisted in the terraform registry.**

Releases are done by [goreleaser](https://goreleaser.com/) and run by our make files. There two goreleaser configs, `.goreleaser.yml` for regular releases and `.goreleaser.prerelease.yml` for doing prereleases (for testing).

Releases are [published to the terraform registry](https://registry.terraform.io/providers/chanzuckerberg/snowflake/latest), which requires that releases by signed.

## Adding a new releaser

To set up a new person for releasing, there are a few steps–

1. releaser: a [keybase account](https://keybase.io/) and a workstation set up with their [Keybase app](https://keybase.io/download).
2. releaser: a pgp key - `keybase pgp gen`
3. releaser: export public key.
   1. If you have a single key in keybase–
      1. `keybase pgp export`
   2. If you have more than one key–
      1. `keybase pgp export` to find id if key you want to export (keep this for later)
      2. `keybase pgp export -q KEY_ID`
4. github admin for chanzuckerberg: take public key exported above and add it [in the registry](https://registry.terraform.io/settings/gpg-keys)
5. releaser: set `KEYBASE_KEY_ID` environment variable
6. releaser: run `make release-prerelease` to test that releases are working correctly
7. releaser: run `make release` to release for real
