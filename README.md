# Snowflake Terraform Provider

> ⚠️ **Please note**: If you believe you have found a security issue, _please responsibly disclose_ by contacting us at [triage-terraformprovider-dl@snowflake.com](mailto:triage-terraformprovider-dl@snowflake.com).

> ⚠️ **Disclaimer**: the project is still in the 0.x.x version, which means it’s still in the experimental phase (check [Go module versioning](https://go.dev/doc/modules/version-numbers#v0-number) for more details). It can be used in production but makes no stability or backward compatibility guarantees. We do not provide backward bug fixes and, therefore, always suggest using the newest version. We are providing only limited support for the provider; priorities will be assigned on a case-by-case basis.
>
> Our main current goals are stabilization, addressing existing issues, and providing the missing features (prioritizing the GA features; supporting PrPr and PuPr features are not high priorities now).
>
> With all that in mind, we aim to reach V1 with a stable, reliable, and functional provider. V1 will be free of all the above limitations.

----

![.github/workflows/ci.yml](https://github.com/Snowflake-Labs/terraform-provider-snowflake/workflows/.github/workflows/ci.yml/badge.svg)

This is a terraform provider for managing [Snowflake](https://www.snowflake.com/) resources.

## Table of contents
<!-- TOC -->
* [Snowflake Terraform Provider](#snowflake-terraform-provider)
  * [Table of contents](#table-of-contents)
  * [Getting started](#getting-started)
  * [Migration guide](#migration-guide)
  * [Roadmap](#roadmap)
  * [Getting Help](#getting-help)
  * [Would you like to create an issue?](#would-you-like-to-create-an-issue)
  * [Additional debug logs for `snowflake_grant_privileges_to_role` resource](#additional-debug-logs-for-snowflake_grant_privileges_to_role-resource)
  * [Additional SQL Client configuration](#additional-sql-client-configuration)
  * [Contributing](#contributing)
  * [Releases](#releases)
<!-- TOC -->

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

## Migration guide

Please check the [migration guide](./MIGRATION_GUIDE.md) when changing the version of the provider.

## Roadmap

Check [Roadmap](./ROADMAP.md).

## Getting Help

Some links that might help you:

- The [introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0) shows how to set up your Snowflake account for Terraform (service user, role, authentication, etc) and how to create your first resources in Terraform.
- The [docs on the Terraform registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest) are a complete reference of all resources and data sources supported and contain more advanced examples.
- The [discussions area](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions) of this repo, we use this forum to discuss new features and changes to the provider.
- **If you are an enterprise customer**, reach out to your account team. This helps us prioritize issues.
- The [issues section](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues) might already have an issue addressing your question.

## Would you like to create an issue?
If you would like to create a GitHub issue, please read our [guide](./CREATING_ISSUES.md) first.
It contains useful links, FAQ, and commonly known issues with solutions that may already solve your case.

## Additional debug logs for `snowflake_grant_privileges_to_role` resource
Set environment variable `SF_TF_ADDITIONAL_DEBUG_LOGGING` to a non-empty value. Additional logs will be visible with `sf-tf-additional-debug` prefix, e.g.:
```text
2023/12/08 12:58:22.497078 sf-tf-additional-debug [DEBUG] Creating new client from db
```

## Additional SQL Client configuration
Currently underlying sql [gosnowflake](https://github.com/snowflakedb/gosnowflake) driver is wrapped with [instrumentedsql](https://github.com/luna-duclos/instrumentedsql). In order to use raw [gosnowflake](https://github.com/snowflakedb/gosnowflake) driver, set environment variable `SF_TF_NO_INSTRUMENTED_SQL` to a non-empty value.

By default, the underlying driver is set to error level logging. It can be changed by setting `driver_tracing` field in the configuration to one of (from most to least verbose):
- `trace`
- `debug`
- `info`
- `print`
- `warning`
- `error`
- `fatal`
- `panic`

Read more in [provider configuration docs](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs#schema).

## Contributing

Check [Contributing](./CONTRIBUTING.md).

## Releases

Releases will be performed as needed, typically once every 2 weeks.

Releases are published to [the terraform registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest). Each change has its own release notes (e.g. https://github.com/Snowflake-Labs/terraform-provider-snowflake/releases/tag/v0.89.0) and migration guide if needed (e.g. https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#v0880--v0890).
