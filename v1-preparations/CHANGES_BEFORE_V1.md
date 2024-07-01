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

## "Empty" values
The [Terraform SDK v2](https://github.com/hashicorp/terraform-plugin-sdk) that is currently used in our provider detects the presence of the attribute based on its non-zero Golang value. This means, that it is not possible to distinguish the removal of the value inside a config from setting it explicitely to a zero value, e.g. `0` for the numeric value (check [this thread](https://discuss.hashicorp.com/t/is-it-possible-to-differentiate-between-a-zero-value-and-a-removed-property-in-the-terraform-provider-sdk/43131)). Before we migrate to the new recommended [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) we want to handle such cases the same way inside the provider. It means that:
- boolean attributes will be migrated to the string attributes with two values: `"true"` and `"false"` settable in the config and the special third value `"default"` that will mean, that the given attribute is not set inside the config.
- integer values with the possible `0` value in Snowflake (e.g. `AUTO_SUSPEND` in [warehouse](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse)) will have a special default (usually a `-1`) assigned on the provider side when the config is left empty for them.
- string values with the possible empty (`""`) value (e.g. default for column value in a table) will have a special default `"<Snowflake Terraform Provider string default>"` that will be used for the empty config.
It won't be possible to use the above values directly (it will be for the string attributes) but users should be aware of them, because they may appear in the terraform plans.

## Snowflake parameters
[Snowflake parameters](https://docs.snowflake.com/en/sql-reference/parameters) have different types and hierarchies. In the earlier versions of the provider they have been handled non-intuitively by setting the deault values inside the provider (e.g. [#2356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356)). We want to change that. Because of that we decided to:
- make all parameters available for the given object available in the resource (without the need to use the `snowflake_object_parameter` resource which future will be discussed in the next few weeks)
- remove the default values from Snowflake parameters in every resource before the V1. This is an important **breaking change**. In the previous versions usually not setting the given parameter resulted in using the provider default. This was different from creating the same object without the parameter by hand (because Snowflake just takes the parameter from the hierarchy in such case).
- provider will identify both the internal and the external changes to these parameters on both `value` and `level` levels, e.g.:
  - setting the parameter inside the config and then manually unsetting it to the same value on the higher level will result in detecting a change
  - not setting the parameter inside the config and then manually changing the parameter on object level to the same value as the value one level higher in the hierarchy will result in detecting a change
- handle parameters as optional/computed values in the provider
- add, in all objects having at least one parameter, a special computed collection `parameters` containing all the values and levels of parameters (the result of `SHOW PARAMETERS IN <object> <name>`).

## Config values in the state
Currently, not setting a value for the given attribute inside the config results in populating this field in state with the value extracted from Snowflake (usually by running `SHOW`/`DESCRIBE`). This poses a challenge to identify if the change happened externally or is it just a default Snowflake value (multiple issues reported describe the issue with the infinite plans or weird drifts - this is one of the main reasons). With getting rid of the Snowflake defaults from the provider, it's not an easy thing to do in the currently used [Terraform SDK v2](https://github.com/hashicorp/terraform-plugin-sdk). We have considered and tested a variety of options, including custom diff suppression, setting these fields as optional and computed, and others, but there were smaller or bigger problems with these approaches. What we ended up with, and what will be a guideline for the V1 is:
- we do not fill the given attribute in the state if it is not present inside a config
- we encourage to always use the value directly if you don't want to depend on the Snowflake default (consult [default values](#default-values) section)
- this may result in change detection with migrations to the newer versions of the provider (because currently, the value was is independently of being present in the config or not and there is no way to deduce its presence in the automatic state migrations we can provide) - alternative would be to follow our [resource migration guide](../docs/technical-documentation/resource_migration.md)
- we will provide a `show_output` and `describe_output` in each resource (more in [Raw Snowflake output](#raw-snowflake-output) section)

## Raw Snowflake output
Because of the changes regarding [Config values in the state](#config-values-in-the-state) we want to still allow our users to get the value of the given attrobute, even when it is not set in the config. For each resource (and datasource) we will provide:
- `show_output` computed field, containing the response of `SHOW <object>` for the given managed object
- `describe_output` computed field, containing the response of `DESCRIBE <object> <name>` for the given managed object
- `parameters` computed field, containing all the values and levels of Snowflake parameters (the result of `SHOW PARAMETERS IN <object> <name>`)

This way, it is still possible to obtain the values in your configs, even without setting them directly for the given managed object.

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
