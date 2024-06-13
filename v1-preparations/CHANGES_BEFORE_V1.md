# Changes before v1

This document is a changelog of resources and datasources as part of the https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1. Each provider version lists changes made in resources and datasources definitions during v1 preparations, like added, modified and removed fields.

## v0.92.0 -> v0.93.0
### snowflake_saml_integration deprecation
Resource deprecated in favor of snowflake_saml2_integration.

This new resource has new fields:
- `allowed_user_domains`
- `allowed_email_patterns`
- `comment`

## v0.91.0 ➞ v0.92.0
### snowflake_scim_integration resource changes

New fields:
- `enabled`
- `sync_password`
- `comment`

Changed fields:
- `provisioner_role` renamed to `run_as_role`

Other changes:
- `scim_client` and `run_as_role` marked as `ForceNew`
