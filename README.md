# Snowflake Terraform Provider

> ⚠️ **Please note**: If you believe you have found a security issue, _please responsibly disclose_ by contacting us at [team-cloud-foundation-tools-dl@snowflake.com](mailto:team-cloud-foundation-tools-dl@snowflake.com).

----

![.github/workflows/ci.yml](https://github.com/Snowflake-Labs/terraform-provider-snowflake/workflows/.github/workflows/ci.yml/badge.svg)

This is a terraform provider for managing [Snowflake](https://www.snowflake.com/) resources.

## Table of contents
- [Snowflake Terraform Provider](#snowflake-terraform-provider)
  - [Table of contents](#table-of-contents)
  - [Getting started](#getting-started)
  - [SDK migration table](#sdk-migration-table)
  - [Getting Help](#getting-help)
  - [Additional debug logs for `snowflake_grant_privileges_to_role` resource](#additional-debug-logs-for-snowflake_grant_privileges_to_role-resource)
  - [Additional SQL Client configuration](#additional-sql-client-configuration)
  - [Contributing](#contributing)


## Getting started

> If you're still using the `chanzuckerberg/snowflake` source, see [Upgrading from CZI Provider](./CZI_UPGRADE.md) to upgrade to the current version.

Install the Snowflake Terraform provider by adding a requirement block and a provider block to your Terraform codebase:
```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "~> 0.61"
    }
  }
}

provider "snowflake" {
  account  = "abc12345" # the Snowflake account identifier
  user     = "johndoe"
  password = "v3ry$3cr3t"
  role     = "ACCOUNTADMIN"
}
```

For more information on provider configuration see the [provider docs on the Terraform registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs).

Don't forget to run `terraform init` and you're ready to go! 🚀

Start browsing the [registry docs](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs) to find resources and data sources to use.

## SDK migration table

This table represents the current state of SDK migration from pkg/snowflake to pkg/sdk package.
The goal of migration is to support every Snowflake feature with more type safe API and use it in every resource / datasource.

SDK migration status - indicates if given resource has been migrated into new SDK.

Integration status - indicates if given resource / datasource is using new SDK.

✅ - done<br>
❌ - not started<br>
👨‍💻 - in progress<br>
🟨 - partially done<br>


| Object Type                         | SDK migration status | Resource name                                  | Datasource name               | Integration status |
|-------------------------------------|----------------------|------------------------------------------------|-------------------------------|--------------------|
| Account                             | ✅                    | snowflake_account                              | snowflake_account             | ✅                  |
| Managed Account                     | ❌                    | snowflake_managed_account                      | snowflake_managed_account     | ❌                  |
| User                                | ✅                    | snowflake_user                                 | snowflake_user                | ✅                  |
| Database Role                       | ✅                    | snowflake_database_role                        | snowflake_database_role       | ✅                  |
| Role                                | ✅                    | snowflake_role                                 | snowflake_role                | 👨‍💻              |
| Grant Privilege to Application Role | ❌                    | snowflake_grant_privileges_to_application_role | snowflake_grants              | ❌                  |
| Grant Privilege to Database Role    | ✅                    | snowflake_grant_privileges_to_database_role    | snowflake_grants              | 👨‍💻              |
| Grant Privilege to Role             | ❌                    | snowflake_grant_privileges_to_role             | snowflake_grants              | ✅                  |
| Grant Role                          | ❌                    | snowflake_grant_role                           | snowflake_grants              | ❌                  |
| Grant Database Role                 | ✅                    | snowflake_grant_database_role                  | snowflake_grants              | ❌                  |
| Grant Application Role              | ❌                    | snowflake_grant_application_role               | snowflake_grants              | ❌                  |
| Grant Privilege to Share            | ✅                    | snowflake_grant_privileges_to_share            | snowflake_grants              | ❌                  |
| Grant Ownership                     | ✅                    | snowflake_grant_ownership                      | snowflake_grants              | ❌                  |
| API Integration                     | ❌                    | snowflake_api_integration                      | snowflake_integrations        | ❌                  |
| Notification Integration            | ❌                    | snowflake_notification_integration             | snowflake_integrations        | ❌                  |
| Security Integration                | ❌                    | snowflake_security_integration                 | snowflake_integrations        | ❌                  |
| Storage Integration                 | ❌                    | snowflake_storage_integration                  | snowflake_integrations        | ❌                  |
| Network Policy                      | ✅                    | snowflake_network_policy                       | snowflake_network_policy      | ✅                  |
| Password Policy                     | ✅                    | snowflake_password_policy                      | snowflake_password_policy     | ✅                  |
| Session Policy                      | ✅                    | snowflake_session_policy                       | snowflake_session_policy      | ❌                  |
| Replication Group                   | ❌                    | snowflake_replication_group                    | snowflake_replication_group   | ❌                  |
| Failover Group                      | ✅                    | snowflake_failover_group                       | snowflake_failover_group      | ✅                  |
| Connection                          | ❌                    | snowflake_connection                           | snowflake_connection          | ❌                  |
| Account Parameters                  | ✅                    | snowflake_account_parameter                    | snowflake_parameters          | ❌                  |
| Session Parameters                  | ✅                    | snowflake_session_parameter                    | snowflake_parameters          | ❌                  |
| Object Parameters                   | ✅                    | snowflake_object_parameter                     | snowflake_parameters          | ❌                  |
| Warehouse                           | ✅                    | snowflake_warehouse                            | snowflake_warehouse           | 🟨                 |
| Resource Monitor                    | ✅                    | snowflake_resource_monitor                     | snowflake_resource_monitor    | ✅                  |
| Database                            | ✅                    | snowflake_database                             | snowflake_database            | ✅                  |
| Schema                              | ✅                    | snowflake_schema                               | snowflake_schema              | ✅                  |
| Share                               | ✅                    | snowflake_share                                | snowflake_share               | ✅                  |
| Table                               | 👨‍💻                | snowflake_table                                | snowflake_table               | ❌                  |
| Dynamic Table                       | ✅                    | snowflake_dynamic_table                        | snowflake_dynamic_table       | ❌                  |
| External Table                      | ✅                    | snowflake_external_table                       | snowflake_external_table      | ❌                  |
| Event Table                         | ❌                    | snowflake_event_table                          | snowflake_event_table         | ❌                  |
| View                                | ❌                    | snowflake_view                                 | snowflake_view                | ❌                  |
| Materialized View                   | ❌                    | snowflake_materialized_view                    | snowflake_materialized_view   | ❌                  |
| Sequence                            | ❌                    | snowflake_sequence                             | snowflake_sequence            | ❌                  |
| Function                            | ❌                    | snowflake_function                             | snowflake_function            | ❌                  |
| External Function                   | ❌                    | snowflake_external_function                    | snowflake_external_function   | ❌                  |
| Stored Procedure                    | ❌                    | snowflake_stored_procedure                     | snowflake_stored_procedure    | ❌                  |
| Stream                              | ✅                    | snowflake_stream                               | snowflake_stream              | ✅                  |
| Task                                | ✅                    | snowflake_task                                 | snowflake_task                | ❌                  |
| Masking Policy                      | ✅                    | snowflake_masking_policy                       | snowflake_masking_policy      | ✅                  |
| Row Access Policy                   | ❌                    | snowflake_row_access_policy                    | snowflake_row_access_policy   | ❌                  |
| Tag                                 | ✅                    | snowflake_tag                                  | snowflake_tag                 | ❌                  |
| Secret                              | ❌                    | snowflake_secret                               | snowflake_secret              | ❌                  |
| Stage                               | ❌                    | snowflake_stage                                | snowflake_stage               | ❌                  |
| File Format                         | ✅                    | snowflake_file_format                          | snowflake_file_format         | ✅                  |
| Pipe                                | ✅                    | snowflake_pipe                                 | snowflake_pipe                | ✅                  |
| Alert                               | ✅                    | snowflake_alert                                | snowflake_alert               | ✅                  |
| Application                         | ❌                    | snowflake_application                          | snowflake_application         | ❌                  |
| Application Package                 | ❌                    | snowflake_application_package                  | snowflake_application_package | ❌                  |
| Application Role                    | ❌                    | snowflake_application_role                     | snowflake_application_role    | ❌                  |
| Streamlit                           | ❌                    | snowflake_streamlit                            | snowflake_streamlit           | ❌                  |
| Versioned Schema                    | ❌                    | snowflake_versioned_schema                     | snowflake_versioned_schema    | ❌                  |
| Tag Association                     | ❌                    | snowflake_tag_association                      | snowflake_tag_association     | ❌                  |


## Getting Help

Some links that might help you:

- The [introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0) shows how to set up your Snowflake account for Terraform (service user, role, authentication, etc) and how to create your first resources in Terraform.
- The [docs on the Terraform registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest) are a complete reference of all resources and data sources supported and contain more advanced examples.
- The [discussions area](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions) of this repo, we use this forum to discuss new features and changes to the provider.
- **If you are an enterprise customer**, reach out to your account team. This helps us prioritize issues.
- The [issues section](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues) might already have an issue addressing your question.

## Additional debug logs for `snowflake_grant_privileges_to_role` resource
Set environment variable `SF_TF_ADDITIONAL_DEBUG_LOGGING` to a non-empty value. Additional logs will be visible with `sf-tf-additional-debug` prefix, e.g.:
```text
2023/12/08 12:58:22.497078 sf-tf-additional-debug [DEBUG] Creating new client from db
```

## Additional SQL Client configuration
Currently underlying sql [gosnowflake](https://github.com/snowflakedb/gosnowflake) driver is wrapped with [instrumentedsql](https://github.com/luna-duclos/instrumentedsql). In order to use raw [gosnowflake](https://github.com/snowflakedb/gosnowflake) driver, set environment variable `SF_TF_NO_INSTRUMENTED_SQL` to a non-empty value.

By default, the underlying driver is set to error level logging. It can be changed by setting `SF_TF_GOSNOWFLAKE_LOG_LEVEL` to one of:
- `panic`
- `fatal`
- `error`
- `warn`
- `warning`
- `info`
- `debug`
- `trace`

*note*: It's possible it will be one of the provider config parameters in the future provider versions.

## Contributing

Cf. [Contributing](./CONTRIBUTING.md).
