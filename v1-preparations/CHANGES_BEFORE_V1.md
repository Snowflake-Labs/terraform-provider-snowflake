# Changes before v1

This document is a changelog of resources and datasources as part of the https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1.
Each provider version lists changes made in resources and datasources definitions during v1 preparations, like added, modified and removed fields.

## Default values
For any resource that went through the rework as part of the [resource preparation for V1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1),
the behaviour for default values may change from the previous one. 

In the past, the provider copied defaults from Snowflake, creating a tight coupling between them. 
However, this approach posed a challenge as the defaults on the Snowflake side could change and vary between accounts based on their configurations.

Now, whenever the value is not specified in the configuration, we let the Snowflake fill out the default value for a given field
(if there is one). Using such defaults may lead to non-idempotent cases where the same configuration may 
create a resource with slightly different configuration in Snowflake (depending on the Snowflake Edition and Version, 
current account configuration, and most-likely other factors). That is why we recommend setting optional fields where
you want to ensure that the specified value has been set on the Snowflake side.

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
