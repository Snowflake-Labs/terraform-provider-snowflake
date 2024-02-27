# Migration guide

This document is meant to help you migrate your Terraform config to the new newest version. In migration guides, we will only
describe deprecations or breaking changes and help you to change your configuration to keep the same (or similar) behavior
across different versions.

## v0.86.0 ➞ v0.87.0
### Provider configuration changes

#### **IMPORTANT** *(bug fix)* Configuration hierarchy
There were several issues reported about the configuration hierarchy, e.g. [#2294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2294) and [#2242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2242).
In fact, the order of precedence described in the docs was not followed. This have led to the incorrect behavior.

After migrating to this version, the hierarchy from the docs should be followed:
```text
The Snowflake provider will use the following order of precedence when determining which credentials to use:
1) Provider Configuration
2) Environment Variables
3) Config File
```

**BEWARE**: your configurations will be affected with that change because they may have been leveraging the incorrect configurations precedence. Please be sure to check all the configurations before running terraform.

### snowflake_failover_group resource changes
#### *(bug fix)* ACCOUNT PARAMETERS is returned as PARAMETERS from SHOW FAILOVER GROUPS
Longer context in [#2517](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2517).
After this change, one apply may be required to update the state correctly for failover group resources using `ACCOUNT PARAMETERS`.

### snowflake_schema resource changes
#### *(behavior change)* Schema `data_retention_days`
To make `snowflake_schema.data_retention_days` truly optional field (previously it was producing plan every time when no value was set),
we added `-1` as a possible value as set it as default. That got rid of the unexpected plans when no value is set and added possibility to use default value assigned by Snowflake (see [the data retention period](https://docs.snowflake.com/en/user-guide/data-time-travel#data-retention-period)).

## v0.85.0 ➞ v0.86.0
### snowflake_table_constraint resource changes

#### *(behavior change)* NOT NULL removed from possible types
The `type` of the constraint was limited back to `UNIQUE`, `PRIMARY KEY`, and `FOREIGN KEY`.
The reason for that is, that syntax for Out-of-Line constraint ([docs](https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint#out-of-line-unique-primary-foreign-key)) does not contain `NOT NULL`.
It is noted as a behavior change but in some way it is not; with the previous implementation it did not work at all with `type` set to `NOT NULL` because the generated statement was not a valid Snowflake statement.

We will consider adding `NOT NULL` back because it can be set by `ALTER COLUMN columnX SET NOT NULL`, but first we want to revisit the whole resource design.

#### *(behavior change)* table_id reference
The docs were inconsistent. Example prior to 0.86.0 version showed using the `table.id` as the `table_id` reference. The description of the `table_id` parameter never allowed such a value (`table.id` is a `|`-delimited identifier representation and only the `.`-separated values were listed in the docs: https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.85.0/docs/resources/table_constraint#required. The misuse of `table.id` parameter will result in error after migrating to 0.86.0. To make the config work, please remove and reimport the constraint resource from the state as described in [resource migration doc](./docs/technical-documentation/resource_migration.md).

After discussions in [#2535](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2535) we decided to provide a temporary workaround in 0.87.0 version, so that the manual migration is not necessary. It allows skipping the migration and jumping straight to 0.87.0 version. However, the temporary workaround will be gone in one of the future versions. Please adjust to the newly suggested reference with the new resources you create.

### snowflake_external_function resource changes

#### *(behavior change)* return_null_allowed default is now true
The `return_null_allowed` attribute default value is now `true`. This is a behavior change because it was `false` before. The reason it was changed is to match the expected default value in the [documentation](https://docs.snowflake.com/en/sql-reference/sql/create-external-function#optional-parameters) `Default: The default is NULL (i.e. the function can return NULL values).`

#### *(behavior change)* comment is no longer required
The `comment` attribute is now optional. It was required before, but it is not required in Snowflake API.

### snowflake_external_functions data source changes

#### *(behavior change)* schema is now required with database
The `schema` attribute is now required with `database` attribute to match old implementation `SHOW EXTERNAL FUNCTIONS IN SCHEMA "<database>"."<schema>"`. In the future this may change to make schema optional.

## vX.XX.X -> v0.85.0

### Migration from old (grant) resources to new ones

In recent changes, we introduced a new grant resources to replace the old ones.
To aid with the migration, we wrote a guide to show one of the possible ways to migrate deprecated resources to their new counter-parts.
As the guide is more general and applies to every version (and provider), we moved it [here](./docs/technical-documentation/resource_migration.md).

### snowflake_procedure resource changes
#### *(deprecation)* return_behavior
`return_behavior` parameter is deprecated because it is also deprecated in the Snowflake API.

### snowflake_function resource changes
#### *(behavior change)* return_type
`return_type` has become force new because there is no way to alter it without dropping and recreating the function.

## v0.84.0 ➞ v0.85.0

### snowflake_notification_integration resource changes
#### *(behavior change)* notification_provider
`notification_provider` becomes required and has three possible values `AZURE_STORAGE_QUEUE`, `AWS_SNS`, and `GCP_PUBSUB`.
It is still possible to set it to `AWS_SQS` but because there is no underlying SQL, so it will result in an error.
Attributes `aws_sqs_arn` and `aws_sqs_role_arn` will be ignored.
Computed attributes `aws_sqs_external_id` and `aws_sqs_iam_user_arn` won't be updated.

#### *(behavior change)* force new for multiple attributes
Force new was added for the following attributes (because no usable SQL alter statements for them):
- `azure_storage_queue_primary_uri`
- `azure_tenant_id`
- `gcp_pubsub_subscription_name`
- `gcp_pubsub_topic_name`

#### *(deprecation)* direction
`direction` parameter is deprecated because it is added automatically on the SDK level.

#### *(deprecation)* type
`type` parameter is deprecated because it is added automatically on the SDK level (and basically it's always `QUEUE`).

## v0.73.0 ➞ v0.74.0
### Provider configuration changes

In this change we have done a provider refactor to make it more complete and customizable by supporting more options that
were already available in Golang Snowflake driver. This lead to several attributes being added and a few deprecated.
We will focus on the deprecated ones and show you how to adapt your current configuration to the new changes.

#### *(rename)* username ➞ user

```terraform
provider "snowflake" {
  # before
  username = "username"

  # after
  user = "username"
}
```

#### *(structural change)* OAuth API

```terraform
provider "snowflake" {
  # before
  browser_auth        = false
  oauth_access_token  = "<access_token>"
  oauth_refresh_token = "<refresh_token>"
  oauth_client_id     = "<client_id>"
  oauth_client_secret = "<client_secret>"
  oauth_endpoint      = "<endpoint>"
  oauth_redirect_url  = "<redirect_uri>"

  # after
  authenticator = "ExternalBrowser"
  token         = "<access_token>"
  token_accessor {
    refresh_token   = "<refresh_token>"
    client_id       = "<client_id>"
    client_secret   = "<client_secret>"
    token_endpoint  = "<endpoint>"
    redirect_uri    = "<redirect_uri>"
  }
}
```

#### *(remove redundant information)* region

Specifying a region is a legacy thing and according to https://docs.snowflake.com/en/user-guide/admin-account-identifier
you can specify a region as a part of account parameter. Specifying account parameter with the region is also considered legacy,
but with this approach it will be easier to convert only your account identifier to the new preferred way of specifying account identifier.

```terraform
provider "snowflake" {
  # before
  region = "<cloud_region_id>"

  # after
  account = "<account_locator>.<cloud_region_id>"
}
```

#### *(todo)* private key path

```terraform
provider "snowflake" {
  # before
  private_key_path = "<filepath>"

  # after
  private_key = file("<filepath>")
}
```

#### *(rename)* session_params ➞ params

```terraform
provider "snowflake" {
  # before
  session_params = {}

  # after
  params = {}
}
```

#### *(behavior change)* authenticator (JWT)

Before the change `authenticator` parameter did not have to be set for private key authentication and was deduced by the provider. The change is a result of the introduced configuration alignment with an underlying [gosnowflake driver](https://github.com/snowflakedb/gosnowflake). The authentication type is required there, and it defaults to user+password one. From this version, set `authenticator` to `JWT` explicitly.
