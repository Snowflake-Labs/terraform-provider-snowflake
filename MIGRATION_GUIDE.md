# Migration guide

This document is meant to help you migrate your Terraform config to the new newest version. In migration guides, we will only
describe deprecations or breaking changes and help you to change your configuration to keep the same (or similar) behavior
across different versions.

> [!TIP]
> We highly recommend upgrading the versions one by one instead of bulk upgrades.

## v0.97.0 ➞ v0.98.0

### snowflake_streamsdata source changes
New filtering options:
- `like`
- `in`
- `starts_with`
- `limit`
- `with_describe`

New output fields
- `show_output`
- `describe_output`

Breaking changes:
- `database` and `schema` are right now under `in` field
- `streams` field now organizes output of show under `show_output` field and the output of describe under `describe_output` field.

Please adjust your Terraform configuration files.

### *(behavior change)* Provider configuration rework
On our road to v1, we have decided to rework configuration to address the most common issues (see a [roadmap entry](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#providers-configuration-rework)). We have created a list of topics we wanted to address before v1. We will prepare an announcement soon. The following subsections describe the things addressed in the v0.98.0.

#### *(behavior change)* changed behavior of some fields
For the fields that are not deprecated, we focused on improving validations and documentation. Also, we adjusted some fields to match our [driver's](https://github.com/snowflakedb/gosnowflake) defaults. Specifically:
- Relaxed validations for enum fields like `protocol` and `authenticator`. Now, the case on such fields is ignored.
- `user`, `warehouse`, `role` - added a validation for an account object identifier
- `validate_default_parameters`, `client_request_mfa_token`, `client_store_temporary_credential`, `ocsp_fail_open`,  - to easily handle three-value logic (true, false, unknown) in provider's config, type of these fields was changed from boolean to string. For more details about default values, please refer to the [changes before v1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#default-values) document.
- `client_ip` - added a validation for an IP address
- `port` - added a validation for a port number
- `okta_url`, `token_accessor.token_endpoint`, `client_store_temporary_credential` - added a validation for a URL address
- `login_timeout`, `request_timeout`, `jwt_expire_timeout`, `client_timeout`, `jwt_client_timeout`, `external_browser_timeout` - added a validation for setting this value to at least `0`
- `authenticator` - added a possibility to configure JWT flow with `SNOWFLAKE_JWT` (formerly, this was upported with `JWT`); the previous value `JWT` was left for compatibility, but will be removed before v1

### *(behavior change)* handling copy_grants
Currently, resources like `snowflake_view`, `snowflake_stream_on_table`, `snowflake_stream_on_external_table` and `snowflake_stream_on_directory_table`  support `copy_grants` field corresponding with `COPY GRANTS` during `CREATE`. The current behavior is that, when a change leading for recreation is detected (meaning a change that can not be handled by ALTER, but only by `CREATE OR REPLACE`), `COPY GRANTS` are used during recreation when `copy_grants` is set to `true`. Changing this field without changes in other field results in a noop because in this case there is no need to recreate a resource.

### *(new feature)* recovering stale streams
Starting from this version, the provider detects stale streams for `snowflake_stream_on_table`, `snowflake_stream_on_external_table` and `snowflake_stream_on_directory_table` and recreates them (optionally with `copy_grants`) to recover them. To handle this correctly, a new computed-only field `stale` has been added to these resource, indicating whether a stream is stale.

### *(new feature)* snowflake_stream_on_directory_table and snowflake_stream_on_view resource
Continuing changes made in [v0.97](#v0960--v0970), the new resource `snowflake_stream_on_directory_table` and `snowflake_stream_on_view` have been introduced to replace the previous `snowflake_stream` for streams on directory tables and streams on views.

To use the new `stream_on_directory_table`, change the old `stream` from
```terraform
resource "snowflake_stream" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  on_stage    = snowflake_stage.stage.fully_qualified_name

  comment = "A stream."
}
```

to

```terraform
resource "snowflake_stream_on_directory_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  stage             = snowflake_stage.stage.fully_qualified_name

  comment = "A stream."
}
```

To use the new `stream_on_view`, change the old `stream` from
```terraform
resource "snowflake_stream" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  on_view    = snowflake_view.view.fully_qualified_name

  comment = "A stream."
}
```

to

```terraform
resource "snowflake_stream_on_view" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  view             = snowflake_view.view.fully_qualified_name

  comment = "A stream."
}
```

Then, follow our [Resource migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md).

### *(new feature)* Secret resources
Added a new secrets resources for managing secrets.
We decided to split each secret flow into individual resources.
This segregation was based on the secret flows in CREATE SECRET. i.e.:
- `snowflake_secret_with_client_credentials`
- `snowflake_secret_with_authorization_code_grant`
- `snowflake_secret_with_basic_authentication`
- `snowflake_secret_with_generic_string`


See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/create-secret).

## v0.96.0 ➞ v0.97.0

### *(new feature)* snowflake_stream_on_table, snowflake_stream_on_external_table resource

To enhance clarity and functionality, the new resources `snowflake_stream_on_table` and `snowflake_stream_on_external_table` have been introduced to replace the previous `snowflake_stream`. Recognizing that the old resource carried multiple responsibilities within a single entity, we opted to divide it into more specialized resources.
The newly introduced resources are aligned with the latest Snowflake documentation at the time of implementation, and adhere to our [new conventions](#general-changes).
This segregation was based on the object on which the stream is created. The mapping between SQL statements and the resources is the following:
- `ON TABLE <table_name>` -> `snowflake_stream_on_table`
- `ON EXTERNAL TABLE <external_table_name>` -> `snowflake_stream_on_external_table` (this was previously not supported)

The resources for streams on directory tables and streams on views will be implemented in the future releases.

To use the new `stream_on_table`, change the old `stream` from
```terraform
resource "snowflake_stream" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  on_table    = snowflake_table.table.fully_qualified_name
  append_only = true

  comment = "A stream."
}
```

to

```terraform
resource "snowflake_stream_on_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  table             = snowflake_table.table.fully_qualified_name
  append_only       = "true"

  comment = "A stream."
}
```


Then, follow our [Resource migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md).

### *(new feature)* new snowflake_service_user and snowflake_legacy_service_user resources

Release v0.95.0 introduced reworked `snowflake_user` resource. As [noted](#note-user-types), the new `SERVICE` and `LEGACY_SERVICE` user types were not supported.

This release introduces two new resources to handle these new user types: `snowflake_service_user` and `snowflake_legacy_service_user`.

Both resources have schemas almost identical to the `snowflake_user` resource with the following exceptions:
- `snowflake_service_user` does not contain the following fields (because they are not supported for the user of type `SERVICE` in Snowflake):
  - `password`
  - `first_name`
  - `middle_name`
  - `last_name`
  - `must_change_password`
  - `mins_to_bypass_mfa`
  - `disable_mfa`
- `snowflake_legacy_service_user` does not contain the following fields (because they are not supported for the user of type `LEGACY_SERVICE` in Snowflake):
  - `first_name`
  - `middle_name`
  - `last_name`
  - `mins_to_bypass_mfa`
  - `disable_mfa`

`snowflake_users` datasource was adjusted to handle different user types and `type` field was added to the `describe_output`.

If you used to manage service or legacy service users through `snowflake_user` resource (e.g. using `lifecycle.ignore_changes`) or `snowflake_unsafe_execute`, please migrate to the new resources following [our guidelines on resource migration](docs/technical-documentation/resource_migration.md).

E.g. change the old config from:

```terraform
resource "snowflake_user" "service_user" {
  lifecycle {
    ignore_changes = [user_type]
  }

  name         = "Snowflake Service User"
  login_name   = "service_user"
  email        = "service_user@snowflake.example"

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."
}
```

to

```
resource "snowflake_service_user" "service_user" {
  name         = "Snowflake Service User"
  login_name   = "service_user"
  email        = "service_user@snowflake.example"

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."
}

```

Then, follow our [resource migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md).

Connected issues: [#2951](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2951)

## v0.95.0 ➞ v0.96.0

### snowflake_masking_policies data source changes
New filtering options:
- `in`
- `limit`
- `with_describe`

New output fields
- `show_output`
- `describe_output`

Breaking changes:
- `database` and `schema` are right now under `in` field
- `masking_policies` field now organizes output of show under `show_output` field and the output of describe under `describe_output` field.

Please adjust your Terraform configuration files.

### snowflake_masking_policy resource changes
New fields:
  - `show_output` field that holds the response from SHOW MASKING POLICIES.
  - `describe_output` field that holds the response from DESCRIBE MASKING POLICY.

#### *(breaking change)* Renamed fields in snowflake_masking_policy resource
Renamed fields:
  - `masking_expression` to `body`
Please rename these fields in your configuration files. State will be migrated automatically.

#### *(breaking change)* Removed fields from snowflake_masking_policy resource
Removed fields:
- `or_replace`
- `if_not_exists`
The value of these field will be removed from the state automatically.

#### *(breaking change)* Adjusted schema of arguments/signature
The field `signature` is renamed to `arguments` to be consistent with other resources.
Now, arguments are stored without nested `column` field. Please adjust that in your configs, like in the example below. State is migrated automatically.

The old configuration looks like this:
```
  signature {
    column {
      name = "val"
      type = "VARCHAR"
    }
  }
```

The new configuration looks like this:
```
  argument {
    name = "val"
    type = "VARCHAR"
  }
```

#### *(breaking change)* Identifiers related changes
During [identifiers rework](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework) we decided to
migrate resource ids from pipe-separated to regular Snowflake identifiers (e.g. `<database_name>|<schema_name>` -> `"<database_name>"."<schema_name>"`). Importing resources also needs to be adjusted (see [example](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/row_access_policy#import)).

Also, we added diff suppress function that prevents Terraform from showing differences, when only quoting is different.

No change is required, the state will be migrated automatically.

#### *(behavior change)* Boolean type changes
To easily handle three-value logic (true, false, unknown) in provider's configs, type of `exempt_other_policies` was changed from boolean to string.

For more details about default values, please refer to the [changes before v1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#default-values) document.

### *(breaking change)* resource_monitor resource
Removed fields:
- `set_for_account` (will be settable on account resource, right now, the preferred way is to set it through unsafe_execute resource)
- `warehouses` (can be set on warehouse resource, optionally through unsafe_execute resource only if the warehouse is not managed by Terraform)
- `suspend_triggers` (now, `suspend_trigger` should be used)
- `suspend_immediate_triggers` (now, `suspend_immediate_trigger` should be used)

### *(breaking change)* resource_monitor data source
Changes:
- New filtering option `like`
- Now, the output of `SHOW RESOURCE MONITORS` is now inside `resource_monitors.*.show_output`. Here's the list of currently available fields:
    - `name`
    - `credit_quota`
    - `used_credits`
    - `remaining_credits`
    - `level`
    - `frequency`
    - `start_time`
    - `end_time`
    - `suspend_at`
    - `suspend_immediate_at`
    - `created_on`
    - `owner`
    - `comment`

### snowflake_row_access_policies data source changes
New filtering options:
- `in`
- `limit`
- `with_describe`

New output fields
- `show_output`
- `describe_output`

Breaking changes:
- `database` and `schema` are right now under `in` field
- `row_access_policies` field now organizes output of show under `show_output` field and the output of describe under `describe_output` field.

Please adjust your Terraform configuration files.

### snowflake_row_access_policy resource changes
New fields:
  - `show_output` field that holds the response from SHOW ROW ACCESS POLICIES.
  - `describe_output` field that holds the response from DESCRIBE ROW ACCESS POLICY.

#### *(breaking change)* Renamed fields in snowflake_row_access_policy resource
Renamed fields:
  - `row_access_expression` to `body`
Please rename these fields in your configuration files. State will be migrated automatically.

#### *(breaking change)* Adjusted schema of arguments/signature
The field `signature` is renamed to `arguments` to be consistent with other resources.
Now, arguments are stored as a list, instead of a map. Please adjust that in your configs. State is migrated automatically. Also, this means that order of the items matters and may be adjusted.


The old configuration looks like this:
```
  signature = {
    A = "VARCHAR",
    B = "VARCHAR"
  }
```

The new configuration looks like this:
```
  argument {
    name = "A"
    type = "VARCHAR"
  }
  argument {
    name = "B"
    type = "VARCHAR"
  }
```

Argument names are now case sensitive. All policies created previously in the provider have upper case argument names. If you used lower case before, please adjust your configs. Values in the state will be migrated to uppercase automatically.

#### *(breaking change)* Adjusted behavior on changing name
Previously, after changing `name` field, the resource was recreated. Now, the object is renamed with `RENAME TO`.

#### *(breaking change)* Mitigating permadiff on `body`
Previously, `body` of a policy was compared as a raw string. This led to permament diff because of leading newlines (see https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2053).

Now, similarly to handling statements in other resources, we replace blank characters with a space. The provider can cause false positives in cases where a change in case or run of whitespace is semantically significant.

#### *(breaking change)* Identifiers related changes
During [identifiers rework](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework) we decided to
migrate resource ids from pipe-separated to regular Snowflake identifiers (e.g. `<database_name>|<schema_name>` -> `"<database_name>"."<schema_name>"`). Importing resources also needs to be adjusted (see [example](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/row_access_policy#import)).

Also, we added diff suppress function that prevents Terraform from showing differences, when only quoting is different.

No change is required, the state will be migrated automatically.

## v0.94.x ➞ v0.95.0

### *(breaking change)* database roles data source; field rename, schema structure changes, and adding missing filtering options

- `database` renamed to `in_database`
- Added `like` and `limit` filtering options
- `SHOW DATABASE ROLES` output is now put inside `database_roles.*.show_output`. Here's the list of currently available fields:
    - `created_on`
    - `name`
    - `is_default`
    - `is_current`
    - `is_inherited`
    - `granted_to_roles`
    - `granted_to_database_roles`
    - `granted_database_roles`
    - `owner`
    - `comment`
    - `owner_role_type`

### snowflake_views data source changes
New filtering options:
- `in`
- `like`
- `starts_with`
- `limit`
- `with_describe`

New output fields
- `show_output`
- `describe_output`

Breaking changes:
- `database` and `schema` are right now under `in` field
- `views` field now organizes output of show under `show_output` field and the output of describe under `describe_output` field.

### snowflake_view resource changes
New fields:
  - `row_access_policy`
  - `aggregation_policy`
  - `change_tracking`
  - `is_recursive`
  - `is_temporary`
  - `data_metric_schedule`
  - `data_metric_function`
  - `column`
- added `show_output` field that holds the response from SHOW VIEWS.
- added `describe_output` field that holds the response from DESCRIBE VIEW. Note that one needs to grant sufficient privileges e.g. with [grant_ownership](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_ownership) on the tables used in this view. Otherwise, this field is not filled.

#### *(breaking change)* Removed fields from snowflake_view resource
Removed fields:
- `or_replace` - `OR REPLACE` is added by the provider automatically when `copy_grants` is set to `"true"`
- `tag` - Please, use [tag_association](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/tag_association) instead.
The value of these field will be removed from the state automatically.

#### *(breaking change)* Required warehouse
For this resource, the provider now uses [policy references](https://docs.snowflake.com/en/sql-reference/functions/policy_references) which requires a warehouse in the connection. Please, make sure you have either set a DEFAULT_WAREHOUSE for the user, or specified a warehouse in the provider configuration.

### Identifier changes

#### *(breaking change)* resource identifiers for schema and streamlit
During [identifiers rework](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework) we decided to
migrate resource ids from pipe-separated to regular Snowflake identifiers (e.g. `<database_name>|<schema_name>` -> `"<database_name>"."<schema_name>"`).
Exception to that rule will be identifiers that consist of multiple parts (like in the case of [grant_privileges_to_account_role](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_privileges_to_account_role#import)'s resource id).
The change was applied to already refactored resources (only in the case of `snowflake_schema` and `snowflake_streamlit` this will be a breaking change, because the rest of the objects are single part identifiers in the format of `<name>`):
- `snowflake_api_authentication_integration_with_authorization_code_grant`
- `snowflake_api_authentication_integration_with_client_credentials`
- `snowflake_api_authentication_integration_with_jwt_bearer`
- `snowflake_oauth_integration_for_custom_clients`
- `snowflake_oauth_integration_for_partner_applications`
- `snowflake_external_oauth_integration`
- `snowflake_saml2_integration`
- `snowflake_scim_integration`
- `snowflake_database`
- `snowflake_shared_database`
- `snowflake_secondary_database`
- `snowflake_account_role`
- `snowflake_network_policy`
- `snowflake_warehouse`

No change is required, the state will be migrated automatically.
The rest of the objects will be changed when working on them during [v1 object preparations](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1).

#### *(breaking change)* diff suppress for identifier quoting
(The same set of resources listed above was adjusted)
To prevent issues like [this one](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2982), we added diff suppress function that prevents Terraform from showing differences,
when only quoting is different. In some cases, Snowflake output (mostly from SHOW commands) was dictating which field should be additionally quoted and which shouldn't, but that should no longer be the case.
Like in the change above, the rest of the objects will be changed when working on them during [v1 object preparations](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1).

### New `fully_qualified_name` field in the resources.
We added a new `fully_qualified_name` to snowflake resources. This should help with referencing other resources in fields that expect a fully qualified name. For example, instead of
writing

```object_name = “\”${snowflake_table.database}\”.\”${snowflake_table.schema}\”.\”${snowflake_table.name}\””```

 now we can write

```object_name = snowflake_table.fully_qualified_name```

See more details in [identifiers guide](./docs/guides/identifiers.md#new-computed-fully-qualified-name-field-in-resources).

See [example usage](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_privileges_to_account_role).

Some of the resources are excluded from this change:
- deprecated resources
  - `snowflake_database_old`
  - `snowflake_oauth_integration`
  - `snowflake_saml_integration`
- resources for which fully qualified name is not appropriate
  - `snowflake_account_parameter`
  - `snowflake_account_password_policy_attachment`
  - `snowflake_network_policy_attachment`
  - `snowflake_session_parameter`
  - `snowflake_table_constraint`
  - `snowflake_table_column_masking_policy_application`
  - `snowflake_tag_masking_policy_association`
  - `snowflake_tag_association`
  - `snowflake_user_password_policy_attachment`
  - `snowflake_user_public_keys`
  - grant resources

#### *(breaking change)* removed `qualified_name` from `snowflake_masking_policy`, `snowflake_network_rule`, `snowflake_password_policy` and `snowflake_table`
Because of introducing a new `fully_qualified_name` field for all of the resources, `qualified_name` was removed from `snowflake_masking_policy`, `snowflake_network_rule`,  `snowflake_password_policy` and `snowflake_table`. Please adjust your configurations. State is automatically migrated.

### snowflake_stage resource changes

#### *(bugfix)* Correctly handle renamed/deleted stage

Correctly handle the situation when stage was rename/deleted externally (earlier it resulted in a permanent loop). No action is required on the user's side.

Connected issues: [#2972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2972)

### snowflake_table resource changes

#### *(bugfix)* Handle data type diff suppression better for text and number

Data types are not entirely correctly handled inside the provider (read more e.g. in [#2735](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2735)). It will be still improved with the upcoming function, procedure, and table rework. Currently, diff suppression was fixed for text and number data types in the table resource with the following assumptions/limitations:
- for numbers the default precision is 38 and the default scale is 0 (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-numeric#number))
- for number types the following types are treated as synonyms: `NUMBER`, `DECIMAL`, `NUMERIC`, `INT`, `INTEGER`, `BIGINT`, `SMALLINT`, `TINYINT`, `BYTEINT`
- for text the default length is 16777216 (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-text#varchar))
- for text types the following types are treated as synonyms: `VARCHAR`, `CHAR`, `CHARACTER`, `STRING`, `TEXT`
- whitespace and casing is ignored
- if the type arguments cannot be parsed the defaults are used and therefore diff may be suppressed unexpectedly (please report such cases)

No action is required on the user's side.

Connected issues: [#3007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3007)

### snowflake_user resource changes

Because of the multiple changes in the resource, the easiest migration way is to follow our [migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md) to perform zero downtime migration. Alternatively, it is possible to follow some pointers below. Either way, familiarize yourself with the resource changes before version bumping. Also, check the [design decisions](./v1-preparations/CHANGES_BEFORE_V1.md).

#### *(breaking change)* user parameters added to snowflake_user resource

On our road to V1 we changed the approach to Snowflake parameters on the object level; now, we add them directly to the resource. This is a **breaking change** because now:
- Leaving the config empty does not set the default value on the object level but uses the one from hierarchy on Snowflake level instead (so after version bump, the diff running `UNSET` statements is expected).
- This change is not compatible with `snowflake_object_parameter` - you have to set the parameter inside `snowflake_user` resource **IF** you manage users through terraform **AND** you want to set the parameter on the user level.

For more details, check the [Snowflake parameters](./v1-preparations/CHANGES_BEFORE_V1.md#snowflake-parameters).

The following set of [parameters](https://docs.snowflake.com/en/sql-reference/parameters) was added to the `snowflake_user` resource:
 - [ABORT_DETACHED_QUERY](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query)
 - [AUTOCOMMIT](https://docs.snowflake.com/en/sql-reference/parameters#autocommit)
 - [BINARY_INPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#binary-input-format)
 - [BINARY_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#binary-output-format)
 - [CLIENT_MEMORY_LIMIT](https://docs.snowflake.com/en/sql-reference/parameters#client-memory-limit)
 - [CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX](https://docs.snowflake.com/en/sql-reference/parameters#client-metadata-request-use-connection-ctx)
 - [CLIENT_PREFETCH_THREADS](https://docs.snowflake.com/en/sql-reference/parameters#client-prefetch-threads)
 - [CLIENT_RESULT_CHUNK_SIZE](https://docs.snowflake.com/en/sql-reference/parameters#client-result-chunk-size)
 - [CLIENT_RESULT_COLUMN_CASE_INSENSITIVE](https://docs.snowflake.com/en/sql-reference/parameters#client-result-column-case-insensitive)
 - [CLIENT_SESSION_KEEP_ALIVE](https://docs.snowflake.com/en/sql-reference/parameters#client-session-keep-alive)
 - [CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY](https://docs.snowflake.com/en/sql-reference/parameters#client-session-keep-alive-heartbeat-frequency)
 - [CLIENT_TIMESTAMP_TYPE_MAPPING](https://docs.snowflake.com/en/sql-reference/parameters#client-timestamp-type-mapping)
 - [DATE_INPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#date-input-format)
 - [DATE_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#date-output-format)
 - [ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION](https://docs.snowflake.com/en/sql-reference/parameters#enable-unload-physical-type-optimization)
 - [ERROR_ON_NONDETERMINISTIC_MERGE](https://docs.snowflake.com/en/sql-reference/parameters#error-on-nondeterministic-merge)
 - [ERROR_ON_NONDETERMINISTIC_UPDATE](https://docs.snowflake.com/en/sql-reference/parameters#error-on-nondeterministic-update)
 - [GEOGRAPHY_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#geography-output-format)
 - [GEOMETRY_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#geometry-output-format)
 - [JDBC_TREAT_DECIMAL_AS_INT](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-treat-decimal-as-int)
 - [JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-treat-timestamp-ntz-as-utc)
 - [JDBC_USE_SESSION_TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-use-session-timezone)
 - [JSON_INDENT](https://docs.snowflake.com/en/sql-reference/parameters#json-indent)
 - [LOCK_TIMEOUT](https://docs.snowflake.com/en/sql-reference/parameters#lock-timeout)
 - [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#log-level)
 - [MULTI_STATEMENT_COUNT](https://docs.snowflake.com/en/sql-reference/parameters#multi-statement-count)
 - [NOORDER_SEQUENCE_AS_DEFAULT](https://docs.snowflake.com/en/sql-reference/parameters#noorder-sequence-as-default)
 - [ODBC_TREAT_DECIMAL_AS_INT](https://docs.snowflake.com/en/sql-reference/parameters#odbc-treat-decimal-as-int)
 - [QUERY_TAG](https://docs.snowflake.com/en/sql-reference/parameters#query-tag)
 - [QUOTED_IDENTIFIERS_IGNORE_CASE](https://docs.snowflake.com/en/sql-reference/parameters#quoted-identifiers-ignore-case)
 - [ROWS_PER_RESULTSET](https://docs.snowflake.com/en/sql-reference/parameters#rows-per-resultset)
 - [S3_STAGE_VPCE_DNS_NAME](https://docs.snowflake.com/en/sql-reference/parameters#s3-stage-vpce-dns-name)
 - [SEARCH_PATH](https://docs.snowflake.com/en/sql-reference/parameters#search-path)
 - [SIMULATED_DATA_SHARING_CONSUMER](https://docs.snowflake.com/en/sql-reference/parameters#simulated-data-sharing-consumer)
 - [STATEMENT_QUEUED_TIMEOUT_IN_SECONDS](https://docs.snowflake.com/en/sql-reference/parameters#statement-queued-timeout-in-seconds)
 - [STATEMENT_TIMEOUT_IN_SECONDS](https://docs.snowflake.com/en/sql-reference/parameters#statement-timeout-in-seconds)
 - [STRICT_JSON_OUTPUT](https://docs.snowflake.com/en/sql-reference/parameters#strict-json-output)
 - [TIMESTAMP_DAY_IS_ALWAYS_24H](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-day-is-always-24h)
 - [TIMESTAMP_INPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-input-format)
 - [TIMESTAMP_LTZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-ltz-output-format)
 - [TIMESTAMP_NTZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-ntz-output-format)
 - [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-output-format)
 - [TIMESTAMP_TYPE_MAPPING](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-type-mapping)
 - [TIMESTAMP_TZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-tz-output-format)
 - [TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#timezone)
 - [TIME_INPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#time-input-format)
 - [TIME_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#time-output-format)
 - [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#trace-level)
 - [TRANSACTION_ABORT_ON_ERROR](https://docs.snowflake.com/en/sql-reference/parameters#transaction-abort-on-error)
 - [TRANSACTION_DEFAULT_ISOLATION_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#transaction-default-isolation-level)
 - [TWO_DIGIT_CENTURY_START](https://docs.snowflake.com/en/sql-reference/parameters#two-digit-century-start)
 - [UNSUPPORTED_DDL_ACTION](https://docs.snowflake.com/en/sql-reference/parameters#unsupported-ddl-action)
 - [USE_CACHED_RESULT](https://docs.snowflake.com/en/sql-reference/parameters#use-cached-result)
 - [WEEK_OF_YEAR_POLICY](https://docs.snowflake.com/en/sql-reference/parameters#week-of-year-policy)
 - [WEEK_START](https://docs.snowflake.com/en/sql-reference/parameters#week-start)
 - [ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR](https://docs.snowflake.com/en/sql-reference/parameters#enable-unredacted-query-syntax-error)
 - [NETWORK_POLICY](https://docs.snowflake.com/en/sql-reference/parameters#network-policy)
 - [PREVENT_UNLOAD_TO_INTERNAL_STAGES](https://docs.snowflake.com/en/sql-reference/parameters#prevent-unload-to-internal-stages)

Connected issues: [#2938](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2938)

#### *(breaking change)* Changes in sensitiveness of name, login_name, and display_name

According to https://docs.snowflake.com/en/sql-reference/functions/all_user_names#usage-notes, `NAME`s are not considered sensitive data and `LOGIN_NAME`s are. Previous versions of the provider had this the other way around. In this version, `name` attribute was unmarked as sensitive, whereas `login_name` was marked as sensitive. This may break your configuration if you were using `login_name`s before e.g. in a `for_each` loop.

The `display_name` attribute was marked as sensitive. It defaults to `name` if not provided on Snowflake side. Because `name` is no longer sensitive, we also change the setting for the `display_name`.

Connected issues: [#2662](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2662), [#2668](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2668).

#### *(bugfix)* Correctly handle `default_warehouse`, `default_namespace`, and `default_role`

During the [identifiers rework](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework), we generalized how we compute the differences correctly for the identifier fields (read more in [this document](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/identifiers_rework_design_decisions.md)). Proper suppressor was applied to `default_warehouse`, `default_namespace`, and `default_role`. Also, all these three attributes were corrected (e.g. handling spaces/hyphens in names).

Connected issues: [#2836](https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2836), [#2942](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2942)

#### *(bugfix)* Correctly handle failed update

Not every attribute can be updated in the state during read (like `password` in the `snowflake_user` resource). In situations where update fails, we may end up with an incorrect state (read more in https://github.com/hashicorp/terraform-plugin-sdk/issues/476). We use a deprecated method from the plugin SDK, and now, for partially failed updates, we preserve the resource's previous state. It fixed this kind of situations for `snowflake_user` resource.

Connected issues: [#2970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2970)

#### *(breaking change)* Handling default secondary roles

Old field `default_secondary_roles` was removed in favour of the new, easier, `default_secondary_roles_option` because the only possible options that can be currently set are `('ALL')` and `()`.  The logic to handle set element changes was convoluted and error-prone. Additionally, [bcr 2024_07](https://docs.snowflake.com/en/release-notes/bcr-bundles/2024_07/bcr-1692) complicated the matter even more.

Now:
- the default value is `DEFAULT` - it falls back to Snowflake default (so `()` before and `('ALL')` after the BCR)
- to explicitly set to `('ALL')` use `ALL`
- to explicitly set to `()` use `NONE`

While migrating, the old `default_secondary_roles` will be removed from the state automatically and `default_secondary_roles_option` will be constructed based on the previous value (in some cases apply may be necessary).

Connected issues: [#3038](https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/3038)

#### *(breaking change)* Attributes changes

Attributes that are no longer computed:
- `login_name`
- `display_name`
- `disabled`
- `default_role`

New fields:
- `middle_name`
- `days_to_expiry`
- `mins_to_unlock`
- `mins_to_bypass_mfa`
- `disable_mfa`
- `default_secondary_roles_option`
- `show_output` - holds the response from `SHOW USERS`. Remember that the field will be only recomputed if one of the user attributes is changed.
- `parameters` - holds the response from `SHOW PARAMETERS IN USER`.

Removed fields:
- `has_rsa_public_key`
- `default_secondary_roles` - replaced with `default_secondary_roles_option`

Default changes:
- `must_change_password`
- `disabled`

Type changes:
- `must_change_password`: bool -> string (To easily handle three-value logic (true, false, unknown) in provider's configs, read more in https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/751239b7d2fee4757471db6c03b952d4728ee099/v1-preparations/CHANGES_BEFORE_V1.md?plain=1#L24)
- `disabled`: bool -> string (To easily handle three-value logic (true, false, unknown) in provider's configs, read more in https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/751239b7d2fee4757471db6c03b952d4728ee099/v1-preparations/CHANGES_BEFORE_V1.md?plain=1#L24)

#### *(breaking change)* refactored snowflake_users datasource
> **IMPORTANT NOTE:** when querying users you don't have permissions to, the querying options are limited.
You won't get almost any field in `show_output` (only empty or default values), the DESCRIBE command cannot be called, so you have to set `with_describe = false`.
Only `parameters` output is not affected by the lack of privileges.

Changes:
- account checking logic was entirely removed
- `pattern` renamed to `like`
- `like`, `starts_with`, and `limit` filters added
- `SHOW USERS` output is enclosed in `show_output` field inside `users` (all the previous fields in `users` map were removed)
- Added outputs from **DESC USER** and **SHOW PARAMETERS IN USER** (they can be turned off by declaring `with_describe = false` and `with_parameters = false`, **they're turned on by default**).
  The additional parameters call **DESC USER** (with `with_describe` turned on) and **SHOW PARAMETERS IN USER** (with `with_parameters` turned on) **per user** returned by **SHOW USERS**.
  The outputs of both commands are held in `users` entry, where **DESC USER** is saved in the `describe_output` field, and **SHOW PARAMETERS IN USER** in the `parameters` field.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

Connected issues: [#2902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2902)

#### *(breaking change)* snowflake_user_public_keys usage with snowflake_user

`snowflake_user_public_keys` is a resource allowing to set keys for the given user. Before this version, it was possible to have `snowflake_user` and `snowflake_user_public_keys` used next to each other.
Because the logic handling the keys in `snowflake_user` was fixed, it is advised to use `snowflake_user_public_keys` only when user is not managed through terraform. Having both resources configured for the same user will result in improper behavior.

To migrate, in case of having two resources:
- copy the keys to `rsa_public_key` and `rsa_public_key2` in `snowflake_user`
- remove `snowflake_user_public_keys` from state (following https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md#resource-migration)
- remove `snowflake_user_public_keys` from config

#### *(breaking change)* snowflake_network_policy_attachment usage with snowflake_user

`snowflake_network_policy_attachment` changes are similar to the changes to `snowflake_user_public_keys` above. It is advised to use `snowflake_network_policy_attachment` only when user is not managed through terraform. Having both resources configured for the same user will result in improper behavior.

To migrate, in case of having two resources:
- copy network policy to [network_policy](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.95.0/docs/resources/user#network_policy) attribute in the `snowflake_user` resource
- remove `snowflake_network_policy_attachment` from state (following https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md#resource-migration)
- remove `snowflake_network_policy_attachment` from config

References: [#3048](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/3048), [#3058](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3058)

#### *(note)* snowflake_user_password_policy_attachment and other user policies

`snowflake_user_password_policy_attachment` is not addressed in the current version.
Attaching other user policies is not addressed in the current version.

Both topics will be addressed in the following versions.

#### *(note)* user types

`service` and `legacy_service` user types are currently not supported. They will be supported in the following versions as separate resources (namely `snowflake_service_user` and `snowflake_legacy_service_user`).

## v0.94.0 ➞ v0.94.1
### changes in snowflake_schema

In order to avoid dropping `PUBLIC` schemas, we have decided to use `ALTER` instead of `OR REPLACE` during creation. In the future we are planning to use `CREATE OR ALTER` when it becomes available for schems.

## v0.93.0 ➞ v0.94.0
### *(breaking change)* changes in snowflake_scim_integration

In order to fix issues in v0.93.0, when a resource has Azure scim client, `sync_password` field is now set to `default` value in the state. State will be migrated automatically.

### *(breaking change)* refactored snowflake_schema resource

Renamed fields:
- renamed `is_managed` to `with_managed_access`
- renamed `data_retention_days` to `data_retention_time_in_days`

Please rename these fields in your configuration files. State will be migrated automatically.

Removed fields:
- `tag`
The value of this field will be removed from the state automatically. Please, use [tag_association](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/tag_association) instead.

New fields:
- the following set of [parameters](https://docs.snowflake.com/en/sql-reference/parameters) was added:
    - `max_data_extension_time_in_days`
    - `external_volume`
    - `catalog`
    - `replace_invalid_characters`
    - `default_ddl_collation`
    - `storage_serialization_policy`
    - `log_level`
    - `trace_level`
    - `suspend_task_after_num_failures`
    - `task_auto_retry_attempts`
    - `user_task_managed_initial_warehouse_size`
    - `user_task_timeout_ms`
    - `user_task_minimum_trigger_interval_in_seconds`
    - `quoted_identifiers_ignore_case`
    - `enable_console_output`
    - `pipe_execution_paused`
- added `show_output` field that holds the response from SHOW SCHEMAS.
- added `describe_output` field that holds the response from DESCRIBE SCHEMA. Note that one needs to grant sufficient privileges e.g. with [grant_ownership](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_ownership) on all objects in the schema. Otherwise, this field is not filled.
- added `parameters` field that holds the response from SHOW PARAMETERS IN SCHEMA.

We allow creating and managing `PUBLIC` schemas now. When the name of the schema is `PUBLIC`, it's created with `OR_REPLACE`. Please be careful with this operation, because you may experience data loss. `OR_REPLACE` does `DROP` before `CREATE`, so all objects in the schema will be dropped and this is not visible in Terraform plan. To restore data-related objects that might have been accidentally or intentionally deleted, pleas read about [Time Travel](https://docs.snowflake.com/en/user-guide/data-time-travel). The alternative is to import `PUBLIC` schema manually and then manage it with Terraform. We've decided this based on [#2826](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2826).

#### *(behavior change)* Boolean type changes
To easily handle three-value logic (true, false, unknown) in provider's configs, type of `is_transient` and `with_managed_access` was changed from boolean to string.

Terraform should recreate resources for configs lacking `is_transient` (`DROP` and then `CREATE` will be run underneath). To prevent this behavior, please set `is_transient` field.
For more details about default values, please refer to the [changes before v1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#default-values) document.

Terraform should perform an action for configs lacking `with_managed_access` (`ALTER SCHEMA DISABLE MANAGED ACCESS` will be run underneath which should not affect the Snowflake object, because `MANAGED ACCESS` is not set by default)
### *(breaking change)* refactored snowflake_schemas datasource
Changes:
- `database` is removed and can be specified inside `in` field.
- `like`, `in`, `starts_with`, and `limit` fields enable filtering.
- SHOW SCHEMAS output is enclosed in `show_output` field inside `schemas`.
- Added outputs from **DESC SCHEMA** and **SHOW PARAMETERS IN SCHEMA** (they can be turned off by declaring `with_describe = false` and `with_parameters = false`, **they're turned on by default**).
  The additional parameters call **DESC SCHEMA** (with `with_describe` turned on) and **SHOW PARAMETERS IN SCHEMA** (with `with_parameters` turned on) **per schema** returned by **SHOW SCHEMAS**.
  The outputs of both commands are held in `schemas` entry, where **DESC SCHEMA** is saved in the `describe_output` field, and **SHOW PARAMETERS IN SCHEMA** in the `parameters` field.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

### *(new feature)* new snowflake_account_role resource

Already existing `snowflake_role` was deprecated in favor of the new `snowflake_account_role`. The old resource got upgraded to
have the same features as the new one. The only difference is the deprecation message on the old resource.

New fields:
- added `show_output` field that holds the response from SHOW ROLES. Remember that the field will be only recomputed if one of the fields (`name` or `comment`) are changed.

### *(breaking change)* refactored snowflake_roles data source

Changes:
- New `in_class` filtering option to filter out roles by class name, e.g. `in_class = "SNOWFLAKE.CORE.BUDGET"`
- `pattern` was renamed to `like`
- output of SHOW is enclosed in `show_output`, so before, e.g. `roles.0.comment` is now `roles.0.show_output.0.comment`

### *(new feature)* new snowflake_account_role resource

Already existing `snowflake_role` was deprecated in favor of the new `snowflake_account_role`. The old resource got upgraded to
have the same features as the new one. The only difference is the deprecation message on the old resource.

New fields:
- added `show_output` field that holds the response from SHOW ROLES. Remember that the field will be only recomputed if one of the fields (`name` or `comment`) are changed.

### *(breaking change)* refactored snowflake_roles data source

Changes:
- New `in_class` filtering option to filter out roles by class name, e.g. `in_class = "SNOWFLAKE.CORE.BUDGET"`
- `pattern` was renamed to `like`
- output of SHOW is enclosed in `show_output`, so before, e.g. `roles.0.comment` is now `roles.0.show_output.0.comment`

### *(new feature)* snowflake_streamlit resource
Added a new resource for managing streamlits. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/create-streamlit). In this resource, we decided to split `ROOT_LOCATION` in Snowflake to two fields: `stage` representing stage fully qualified name and `directory_location` containing a path within this stage to root location.

### *(new feature)* snowflake_streamlits datasource
Added a new datasource enabling querying and filtering stremlits. Notes:
- all results are stored in `streamlits` field.
- `like`, `in`, and `limit` fields enable streamlits filtering.
- SHOW STREAMLITS output is enclosed in `show_output` field inside `streamlits`.
- Output from **DESC STREAMLIT** (which can be turned off by declaring `with_describe = false`, **it's turned on by default**) is enclosed in `describe_output` field inside `streamlits`.
  The additional parameters call **DESC STREAMLIT** (with `with_describe` turned on) **per streamlit** returned by **SHOW STREAMLITS**.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

### *(new feature)* refactored snowflake_network_policy resource

No migration required.

New behavior:
- `name` is no longer marked as ForceNew parameter. When changed, now it will perform ALTER RENAME operation, instead of re-creating with the new name.
- Additional validation was added to `blocked_ip_list` to inform about specifying `0.0.0.0/0` ip. More details in the [official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-network-policy#usage-notes).

New fields:
- `show_output` and `describe_output` added to hold the results returned by `SHOW` and `DESCRIBE` commands. Those fields will only be recomputed when specified fields change

### *(new feature)* snowflake_network_policies datasource

Added a new datasource enabling querying and filtering network policies. Notes:
- all results are stored in `network_policies` field.
- `like` field enables filtering.
- SHOW NETWORK POLICIES output is enclosed in `show_output` field inside `network_policies`.
- Output from **DESC NETWORK POLICY** (which can be turned off by declaring `with_describe = false`, **it's turned on by default**) is enclosed in `describe_output` field inside `network_policies`.
  The additional parameters call **DESC NETWORK POLICY** (with `with_describe` turned on) **per network policy** returned by **SHOW NETWORK POLICIES**.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

### *(fix)* snowflake_warehouse resource

Because of the issue [#2948](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2948), we are relaxing the validations for the Snowflake parameter values. Read more in [CHANGES_BEFORE_V1.md](v1-preparations/CHANGES_BEFORE_V1.md#validations).

## v0.92.0 ➞ v0.93.0

### general changes

With this change we introduce the first resources redesigned for the V1. We have made a few design choices that will be reflected in these and in the further reworked resources. This includes:
- Handling the [default values](./v1-preparations/CHANGES_BEFORE_V1.md#default-values).
- Handling the ["empty" values](./v1-preparations/CHANGES_BEFORE_V1.md#empty-values).
- Handling the [Snowflake parameters](./v1-preparations/CHANGES_BEFORE_V1.md#snowflake-parameters).
- Saving the [config values in the state](./v1-preparations/CHANGES_BEFORE_V1.md#config-values-in-the-state).
- Providing a ["raw Snowflake output"](./v1-preparations/CHANGES_BEFORE_V1.md#empty-values) for the managed resources.

They are all described in short in the [changes before v1 doc](./v1-preparations/CHANGES_BEFORE_V1.md). Please familiarize yourself with these changes before the upgrade.

### old grant resources removal
Following the [announcement](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/2736) we have removed the old grant resources. The two resources [snowflake_role_ownership_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role_ownership_grant) and [snowflake_user_ownership_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/user_ownership_grant) were not listed in the announcement, but they were also marked as deprecated ones. We are removing them too to conclude the grants redesign saga.

### *(new feature)* Api authentication resources
Added new api authentication resources, i.e.:
- `snowflake_api_authentication_integration_with_authorization_code_grant`
- `snowflake_api_authentication_integration_with_client_credentials`
- `snowflake_api_authentication_integration_with_jwt_bearer`

See reference [doc](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth).

### *(new feature)* snowflake_oauth_integration_for_custom_clients and snowflake_oauth_integration_for_partner_applications resources

To enhance clarity and functionality, the new resources `snowflake_oauth_integration_for_custom_clients` and `snowflake_oauth_integration_for_partner_applications` have been introduced
to replace the previous `snowflake_oauth_integration`. Recognizing that the old resource carried multiple responsibilities within a single entity, we opted to divide it into two more specialized resources.
The newly introduced resources are aligned with the latest Snowflake documentation at the time of implementation, and adhere to our [new conventions](#general-changes).
This segregation was based on the `oauth_client` attribute, where `CUSTOM` corresponds to `snowflake_oauth_integration_for_custom_clients`,
while other attributes align with `snowflake_oauth_integration_for_partner_applications`.

### *(new feature)* snowflake_security_integrations datasource
Added a new datasource enabling querying and filtering all types of security integrations. Notes:
- all results are stored in `security_integrations` field.
- `like` field enables security integrations filtering.
- SHOW SECURITY INTEGRATIONS output is enclosed in `show_output` field inside `security_integrations`.
- Output from **DESC SECURITY INTEGRATION** (which can be turned off by declaring `with_describe = false`, **it's turned on by default**) is enclosed in `describe_output` field inside `security_integrations`.
  **DESC SECURITY INTEGRATION** returns different properties based on the integration type. Consult the documentation to check which ones will be filled for which integration.
  The additional parameters call **DESC SECURITY INTEGRATION** (with `with_describe` turned on) **per security integration** returned by **SHOW SECURITY INTEGRATIONS**.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

### snowflake_external_oauth_integration resource changes

#### *(behavior change)* Renamed fields
Renamed fields:
- `type` to `external_oauth_type`
- `issuer` to `external_oauth_issuer`
- `token_user_mapping_claims` to `external_oauth_token_user_mapping_claim`
- `snowflake_user_mapping_attribute` to `external_oauth_snowflake_user_mapping_attribute`
- `scope_mapping_attribute` to `external_oauth_scope_mapping_attribute`
- `jws_keys_urls` to `external_oauth_jws_keys_url`
- `rsa_public_key` to `external_oauth_rsa_public_key`
- `rsa_public_key_2` to `external_oauth_rsa_public_key_2`
- `blocked_roles` to `external_oauth_blocked_roles_list`
- `allowed_roles` to `external_oauth_allowed_roles_list`
- `audience_urls` to `external_oauth_audience_list`
- `any_role_mode` to `external_oauth_any_role_mode`
- `scope_delimiter` to `external_oauth_scope_delimiter`
to align with Snowflake docs. Please rename this field in your configuration files. State will be migrated automatically.

#### *(behavior change)* Force new for multiple attributes after removing from config
Conditional force new was added for the following attributes when they are removed from config. There are no alter statements supporting UNSET on these fields.
- `external_oauth_rsa_public_key`
- `external_oauth_rsa_public_key_2`
- `external_oauth_scope_mapping_attribute`
- `external_oauth_jws_keys_url`
- `external_oauth_token_user_mapping_claim`

#### *(behavior change)* Conflicting fields
Fields listed below can not be set at the same time in Snowflake. They are marked as conflicting fields.
- `external_oauth_jws_keys_url` <-> `external_oauth_rsa_public_key`
- `external_oauth_jws_keys_url` <-> `external_oauth_rsa_public_key_2`
- `external_oauth_allowed_roles_list` <-> `external_oauth_blocked_roles_list`

#### *(behavior change)* Changed diff suppress for some fields
The fields listed below had diff suppress which removed '-' from strings. Now, this behavior is removed, so if you had '-' in these strings, please remove them. Note that '-' in these values is not allowed by Snowflake.
- `external_oauth_snowflake_user_mapping_attribute`
- `external_oauth_type`
- `external_oauth_any_role_mode`

### *(new feature)* snowflake_saml2_integration resource

The new `snowflake_saml2_integration` is introduced and deprecates `snowflake_saml_integration`. It contains new fields
and follows our new conventions making it more stable. The old SAML integration wasn't changed, so no migration needed,
but we recommend to eventually migrate to the newer counterpart.

### snowflake_scim_integration resource changes
#### *(behavior change)* Changed behavior of `sync_password`

Now, the `sync_password` field will set the state value to `default` whenever the value is not set in the config. This indicates that the value on the Snowflake side is set to the Snowflake default.

> [!WARNING]
> This change causes issues for Azure scim client (see [#2946](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2946)). The workaround is to remove the resource from the state with `terraform state rm`, add `sync_password = true` to the config, and import with `terraform import "snowflake_scim_integration.test" "aad_provisioning"`. After these steps, there should be no errors and no diff on this field. This behavior is fixed in v0.94 with state upgrader.


#### *(behavior change)* Renamed fields

Renamed field `provisioner_role` to `run_as_role` to align with Snowflake docs. Please rename this field in your configuration files. State will be migrated automatically.

#### *(new feature)* New fields
Fields added to the resource:
- `enabled`
- `sync_password`
- `comment`

#### *(behavior change)* Changed behavior of `enabled`
New field `enabled` is required. Previously the default value during create in Snowflake was `true`. If you created a resource with Terraform, please add `enabled = true` to have the same value.

#### *(behavior change)* Force new for multiple attributes
ForceNew was added for the following attributes (because there are no usable SQL alter statements for them):
- `scim_client`
- `run_as_role`

### snowflake_warehouse resource changes

Because of the multiple changes in the resource, the easiest migration way is to follow our [migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md) to perform zero downtime migration. Alternatively, it is possible to follow some pointers below. Either way, familiarize yourself with the resource changes before version bumping. Also, check the [design decisions](./v1-preparations/CHANGES_BEFORE_V1.md).

#### *(potential behavior change)* Default values removed
As part of the [redesign](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1) we are removing the default values for attributes having their defaults on Snowflake side to reduce coupling with the provider (read more in [default values](./v1-preparations/CHANGES_BEFORE_V1.md#default-values)). Because of that the following defaults were removed:
- `comment` (previously `""`)
- `enable_query_acceleration` (previously `false`)
- `query_acceleration_max_scale_factor` (previously `8`)
- `warehouse_type` (previously `"STANDARD"`)
- `max_concurrency_level` (previously `8`)
- `statement_queued_timeout_in_seconds` (previously `0`)
- `statement_timeout_in_seconds` (previously `172800`)

**Beware!** For attributes being Snowflake parameters (in case of warehouse: `max_concurrency_level`, `statement_queued_timeout_in_seconds`, and `statement_timeout_in_seconds`), this is a breaking change (read more in [Snowflake parameters](./v1-preparations/CHANGES_BEFORE_V1.md#snowflake-parameters)). Previously, not setting a value for them was treated as a fallback to values hardcoded on the provider side. This caused warehouse creation with these parameters set on the warehouse level (and not using the Snowflake default from hierarchy; read more in the [parameters documentation](https://docs.snowflake.com/en/sql-reference/parameters)). To keep the previous values, fill in your configs to the default values listed above.

All previous defaults were aligned with the current Snowflake ones, however it's not possible to distinguish between filled out value and no value in the automatic state upgrader. Therefore, if the given attribute is not filled out in your configuration, terraform will try to perform update after the change (to UNSET the given attribute to the Snowflake default); it should result in no changes on Snowflake object side, but it is required to make Terraform state aligned with your config. **All** other optional fields that were not set inside the config at all (because of the change in handling state logic on our provider side) will follow the same logic. To avoid the need for the changes, fill out the default fields in your config. Alternatively, run `terraform apply`; no further changes should be shown as a part of the plan.

#### *(note)* Automatic state migrations
There are three migrations that should happen automatically with the version bump:
- incorrect `2XLARGE`, `3XLARGE`, `4XLARGE`, `5XLARGE`, `6XLARGE` values for warehouse size are changed to the proper ones
- deprecated `wait_for_provisioning` attribute is removed from the state
- old empty resource monitor attribute is cleaned (earlier it was set to `"null"` string)

#### *(fix)* Warehouse size UNSET

Before the changes, removing warehouse size from the config was not handled properly. Because UNSET is not supported for warehouse size (check the [docs](https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse#properties-parameters) - usage notes for unset) and there are multiple defaults possible, removing the size from config will result in the resource recreation.

#### *(behavior change)* Validation changes
As part of the [redesign](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1) we are adjusting validations or removing them to reduce coupling between Snowflake and the provider. Because of that the following validations were removed/adjusted/added:
- `max_cluster_count` - adjusted: added higher bound (10) according to Snowflake docs
- `min_cluster_count` - adjusted: added higher bound (10) according to Snowflake docs
- `auto_suspend` - adjusted: added `0` as valid value
- `warehouse_size` - adjusted: removed incorrect `2XLARGE`, `3XLARGE`, `4XLARGE`, `5XLARGE`, `6XLARGE` values
- `resource_monitor` - added: validation for a valid identifier (still subject to change during [identifiers rework](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework))
- `max_concurrency_level` - added: validation according to MAX_CONCURRENCY_LEVEL parameter docs
- `statement_queued_timeout_in_seconds` - added: validation according to STATEMENT_QUEUED_TIMEOUT_IN_SECONDS parameter docs
- `statement_timeout_in_seconds` - added: validation according to STATEMENT_TIMEOUT_IN_SECONDS parameter docs

#### *(behavior change)* Deprecated `wait_for_provisioning` field removed
`wait_for_provisioning` field was deprecated a long time ago. It's high time it was removed from the schema.

#### *(behavior change)* `query_acceleration_max_scale_factor` conditional logic removed
Previously, the `query_acceleration_max_scale_factor` was depending on `enable_query_acceleration` parameter, but it is not required on Snowflake side. After migration, `terraform plan` should suggest changes if `enable_query_acceleration` was earlier set to false (manually or from default) and if `query_acceleration_max_scale_factor` was set in config.

#### *(behavior change)* `initially_suspended` forceNew removed
Previously, the `initially_suspended` attribute change caused the resource recreation. This attribute is used only during creation (to create suspended warehouse). There is no reason to recreate the whole object just to have initial state changed.

#### *(behavior change)* Boolean type changes
To easily handle three-value logic (true, false, unknown) in provider's configs, type of `auto_resume` and `enable_query_acceleration` was changed from boolean to string. This should not require updating existing configs (boolean/int value should be accepted and state will be migrated to string automatically), however we recommend changing config values to strings. Terraform should perform an action for configs lacking `auto_resume` or `enable_query_acceleration` (`ALTER WAREHOUSE UNSET AUTO_RESUME` and/or `ALTER WAREHOUSE UNSET ENABLE_QUERY_ACCELERATION` will be run underneath which should not affect the Snowflake object, because `auto_resume` and `enable_query_acceleration` are false by default).

#### *(note)* `resource_monitor` validation and diff suppression
`resource_monitor` is an identifier and handling logic may be still slightly changed as part of https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework. It should be handled automatically (without needed manual actions on user side), though, but it is not guaranteed.

#### *(behavior change)* snowflake_warehouses datasource
- Added `like` field to enable warehouse filtering
- Added missing fields returned by SHOW WAREHOUSES and enclosed its output in `show_output` field.
- Added outputs from **DESC WAREHOUSE** and **SHOW PARAMETERS IN WAREHOUSE** (they can be turned off by declaring `with_describe = false` and `with_parameters = false`, **they're turned on by default**).
  The additional parameters call **DESC WAREHOUSE** (with `with_describe` turned on) and **SHOW PARAMETERS IN WAREHOUSE** (with `with_parameters` turned on) **per warehouse** returned by **SHOW WAREHOUSES**.
  The outputs of both commands are held in `warehouses` entry, where **DESC WAREHOUSE** is saved in the `describe_output` field, and **SHOW PARAMETERS IN WAREHOUSE** in the `parameters` field.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

You can read more in ["raw Snowflake output"](./v1-preparations/CHANGES_BEFORE_V1.md#empty-values).

### *(new feature)* new database resources
As part of the [preparation for v1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1), we split up the database resource into multiple ones:
- Standard database - can be used as `snowflake_database` (replaces the old one and is used to create databases with optional ability to become a primary database ready for replication)
- Shared database - can be used as `snowflake_shared_database` (used to create databases from externally defined shares)
- Secondary database - can be used as `snowflake_secondary_database` (used to create replicas of databases from external sources)

All the field changes in comparison to the previous database resource are:
- `is_transient`
    - in `snowflake_shared_database`
        - removed: the field is removed from `snowflake_shared_database` as it doesn't have any effect on shared databases.
- `from_database` - database cloning was entirely removed and is not possible by any of the new database resources.
- `from_share` - the parameter was moved to the dedicated resource for databases created from shares `snowflake_shared_database`. Right now, it's a text field instead of a map. Additionally, instead of legacy account identifier format we're expecting the new one that with share looks like this: `<organization_name>.<account_name>.<share_name>`. For more information on account identifiers, visit the [official documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier).
- `from_replication` - the parameter was moved to the dedicated resource for databases created from primary databases `snowflake_secondary_database`
- `replication_configuration` - renamed: was renamed to `configuration` and is only available in the `snowflake_database`. Its internal schema changed that instead of list of accounts, we expect a list of nested objects with accounts for which replication (and optionally failover) should be enabled. More information about converting between both versions [here](#resource-renamed-snowflake_database---snowflake_database_old). Additionally, instead of legacy account identifier format we're expecting the new one that looks like this: `<organization_name>.<account_name>` (it will be automatically migrated to the recommended format by the state upgrader). For more information on account identifiers, visit the [official documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier).
- `data_retention_time_in_days`
  - in `snowflake_shared_database`
      - removed: the field is removed from `snowflake_shared_database` as it doesn't have any effect on shared databases.
  - in `snowflake_database` and `snowflake_secondary_database`
    - adjusted: now, it uses different approach that won't set it to -1 as a default value, but rather fills the field with the current value from Snowflake (this still can change).
- added: The following set of [parameters](https://docs.snowflake.com/en/sql-reference/parameters) was added to every database type:
    - `max_data_extension_time_in_days`
    - `external_volume`
    - `catalog`
    - `replace_invalid_characters`
    - `default_ddl_collation`
    - `storage_serialization_policy`
    - `log_level`
    - `trace_level`
    - `suspend_task_after_num_failures`
    - `task_auto_retry_attempts`
    - `user_task_managed_initial_warehouse_size`
    - `user_task_timeout_ms`
    - `user_task_minimum_trigger_interval_in_seconds`
    - `quoted_identifiers_ignore_case`
    - `enable_console_output`

The split was done (and will be done for several objects during the refactor) to simplify the resource on maintainability and usage level.
Its purpose was also to divide the resources by their specific purpose rather than cramping every use case of an object into one resource.

### *(behavior change)* Resource renamed snowflake_database -> snowflake_database_old
We made a decision to use the existing `snowflake_database` resource for redesigning it into a standard database.
The previous `snowflake_database` was renamed to `snowflake_database_old` and the current `snowflake_database`
contains completely new implementation that follows our guidelines we set for V1.
When upgrading to the 0.93.0 version, the automatic state upgrader should cover the migration for databases that didn't have the following fields set:
- `from_share` (now, the new `snowflake_shared_database` should be used instead)
- `from_replica` (now, the new `snowflake_secondary_database` should be used instead)
- `replication_configuration`

For configurations containing `replication_configuraiton` like this one:
```terraform
resource "snowflake_database" "test" {
  name = "<name>"
  replication_configuration {
    accounts = ["<account_locator>", "<account_locator_2>"]
    ignore_edition_check = true
  }
}
```

You have to transform the configuration into the following format (notice the change from account locator into the new account identifier format):
```terraform
resource "snowflake_database" "test" {
  name = "%s"
  replication {
    enable_to_account {
      account_identifier = "<organization_name>.<account_name>"
      with_failover      = false
    }
    enable_to_account {
      account_identifier = "<organization_name_2>.<account_name_2>"
      with_failover      = false
    }
  }
  ignore_edition_check = true
}
```

If you had `from_database` set, you should follow our [resource migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md) to remove
the database from state to later import it in the newer version of the provider.
Otherwise, it may cause issues when migrating to v0.93.0.
For now, we're dropping the possibility to create a clone database from other databases.
The only way will be to clone a database manually and import it as `snowflake_database`, but if
cloned databases diverge in behavior from standard databases, it may cause issues.

For databases with one of the fields mentioned above, manual migration will be needed.
Please refer to our [migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md) to perform zero downtime migration.

If you would like to upgrade to the latest version and postpone the upgrade, you still have to perform the manual migration
to the `snowflake_database_old` resource by following the [zero downtime migrations document](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md).
The only difference would be that instead of writing/generating new configurations you have to just rename the existing ones to contain `_old` suffix.

### *(behavior change)* snowflake_databases datasource
- `terse` and `history` fields were removed.
- `replication_configuration` field was removed from `databases`.
- `pattern` was replaced by `like` field.
- Additional filtering options added (`limit`).
- Added missing fields returned by SHOW DATABASES and enclosed its output in `show_output` field.
- Added outputs from **DESC DATABASE** and **SHOW PARAMETERS IN DATABASE** (they can be turned off by declaring `with_describe = false` and `with_parameters = false`, **they're turned on by default**).
The additional parameters call **DESC DATABASE** (with `with_describe` turned on) and **SHOW PARAMETERS IN DATABASE** (with `with_parameters` turned on) **per database** returned by **SHOW DATABASES**.
The outputs of both commands are held in `databases` entry, where **DESC DATABASE** is saved in the `describe_output` field, and **SHOW PARAMETERS IN DATABASE** in the `parameters` field.
It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

## v0.89.0 ➞ v0.90.0
### snowflake_table resource changes
#### *(behavior change)* Validation to column type added
While solving issue [#2733](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2733) we have introduced diff suppression for `column.type`. To make it work correctly we have also added a validation to it. It should not cause any problems, but it's worth noting in case of any data types used that the provider is not aware of.

### snowflake_procedure resource changes
#### *(behavior change)* Validation to arguments type added
Diff suppression for `arguments.type` is needed for the same reason as above for `snowflake_table` resource.

### tag_masking_policy_association resource changes
Now the `tag_masking_policy_association` resource will only accept fully qualified names separated by dot `.` instead of pipe `|`.

Before
```terraform
resource "snowflake_tag_masking_policy_association" "name" {
    tag_id            = snowflake_tag.this.id
    masking_policy_id = snowflake_masking_policy.example_masking_policy.id
}
```

After
```terraform
resource "snowflake_tag_masking_policy_association" "name" {
    tag_id            = "\"${snowflake_tag.this.database}\".\"${snowflake_tag.this.schema}\".\"${snowflake_tag.this.name}\""
    masking_policy_id = "\"${snowflake_masking_policy.example_masking_policy.database}\".\"${snowflake_masking_policy.example_masking_policy.schema}\".\"${snowflake_masking_policy.example_masking_policy.name}\""
}
```

It's more verbose now, but after identifier rework it should be similar to the previous form.

## v0.88.0 ➞ v0.89.0
#### *(behavior change)* ForceNew removed
The `ForceNew` field was removed in favor of in-place Update for `name` parameter in:
- `snowflake_file_format`
- `snowflake_masking_policy`
So from now, these objects won't be re-created when the `name` changes, but instead only the name will be updated with `ALTER .. RENAME TO` statements.

## v0.87.0 ➞ v0.88.0
### snowflake_procedure resource changes
#### *(behavior change)* Execute as validation added
From now on, the `snowflake_procedure`'s `execute_as` parameter allows only two values: OWNER and CALLER (case-insensitive). Setting other values earlier resulted in falling back to the Snowflake default (currently OWNER) and creating a permadiff.

### snowflake_grants datasource changes
`snowflake_grants` datasource was refreshed as part of the ongoing [Grants Redesign](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#redesigning-grants).

#### *(behavior change)* role fields renames
To be aligned with the convention in other grant resources, `role` was renamed to `account_role` for the following fields:
- `grants_to.role`
- `grants_of.role`
- `future_grants_to.role`.

To migrate simply change `role` to `account_role` in the aforementioned fields.

#### *(behavior change)* grants_to.share type change
`grants_to.share` was a text field. Because Snowflake introduced new syntax `SHOW GRANTS TO SHARE <share_name> IN APPLICATION PACKAGE <app_package_name>` (check more in the [docs](https://docs.snowflake.com/en/sql-reference/sql/show-grants#variants)) the type was changed to object. To migrate simply change:
```terraform
data "snowflake_grants" "example_to_share" {
  grants_to {
    share = "some_share"
  }
}
```
to
```terraform
data "snowflake_grants" "example_to_share" {
  grants_to {
    share {
      share_name = "some_share"
    }
  }
}
```
Note: `in_application_package` is not yet supported.

#### *(behavior change)* future_grants_in.schema type change
`future_grants_in.schema` was an object field allowing to set required `schema_name` and optional `database_name`. Our strategy is to be explicit, so the schema field was changed to string and fully qualified name is expected. To migrate change:
```terraform
data "snowflake_grants" "example_future_in_schema" {
  future_grants_in {
    schema {
      database_name = "some_database"
      schema_name   = "some_schema"
    }
  }
}
```
to
```terraform
data "snowflake_grants" "example_future_in_schema" {
  future_grants_in {
    schema = "\"some_database\".\"some_schema\""
  }
}
```
#### *(new feature)* grants_to new options
`grants_to` was enriched with three new options:
- `application`
- `application_role`
- `database_role`

No migration work is needed here.

#### *(new feature)* grants_of new options
`grants_to` was enriched with two new options:
- `database_role`
- `application_role`

No migration work is needed here.

#### *(new feature)* future_grants_to new options
`future_grants_to` was enriched with one new option:
- `database_role`

No migration work is needed here.

#### *(documentation)* improvements
Descriptions of attributes were altered. More examples were added (both for old and new features).

## v0.86.0 ➞ v0.87.0
### snowflake_database resource changes
#### *(behavior change)* External object identifier changes

Previously, in `snowflake_database` when creating a database form share, it was possible to provide `from_share.provider`
in the format of `<org_name>.<account_name>`. It worked even though we expected account locator because our "external" identifier wasn't quoting its string representation.
To be consistent with other identifier types, we quoted the output of "external" identifiers which makes such configurations break
(previously, they were working "by accident"). To fix it, the previous format of `<org_name>.<account_name>` has to be changed
to account locator format `<account_locator>` (mind that it's now case-sensitive). The account locator can be retrieved by calling `select current_account();` on the sharing account.
In the future we would like to eventually come back to the `<org_name>.<account_name>` format as it's recommended by Snowflake.

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

### snowflake_database, snowflake_schema, and snowflake_table resource changes
#### *(behavior change)* Database `data_retention_time_in_days` + Schema `data_retention_days` + Table `data_retention_time_in_days`
For context [#2356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356).
To make data retention fields truly optional (previously they were producing plan every time when no value was set),
we added `-1` as a possible value, and it is set as default. That got rid of the unexpected plans when no value is set and added possibility to use default value assigned by Snowflake (see [the data retention period](https://docs.snowflake.com/en/user-guide/data-time-travel#data-retention-period)).

### snowflake_table resource changes
#### *(behavior change)* Table `data_retention_days` field removed in favor of `data_retention_time_in_days`
For context [#2356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356).
To define data retention days for table `data_retention_time_in_days` should be used as deprecated `data_retention_days` field is being removed.

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

### snowflake_stage resource changes

#### *(behavior change/regression)* copy_options
Setting `copy_options` to `ON_ERROR = 'CONTINUE'` would result in a permadiff. Use `ON_ERROR = CONTINUE` (without single quotes) or bump to v0.89.0 in which the behavior was fixed.

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
