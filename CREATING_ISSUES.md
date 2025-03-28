# Creating GitHub issues

* [Creating GitHub issues](#creating-github-issues)
  * [1. Check the existing GitHub issues.](#1-check-the-existing-github-issues)
  * [2. Go through the frequently asked questions and commonly known issues.](#2-go-through-the-frequently-asked-questions-and-commonly-known-issues)
  * [3. Check the official Snowflake documentation](#3-check-the-official-snowflake-documentation)
  * [4. Choose the correct template and use as much information as possible.](#4-choose-the-correct-template-and-use-as-much-information-as-possible)
* [FAQ](#faq)
  * [When will the Snowflake feature X be available in the provider?](#when-will-the-snowflake-feature-x-be-available-in-the-provider)
  * [When will my bug report be fixed/released?](#when-will-my-bug-report-be-fixedreleased)
  * [How to migrate from version X to Y?](#how-to-migrate-from-version-x-to-y)
  * [What are the current/future plans for the provider?](#what-are-the-currentfuture-plans-for-the-provider)
  * [How can I contribute?](#how-can-i-contribute)
  * [How can I debug the issue myself?](#how-can-i-debug-the-issue-myself)
  * [How can I import already existing Snowflake infrastructure into Terraform?](#how-can-i-import-already-existing-snowflake-infrastructure-into-terraform)
  * [What identifiers are valid inside the provider and how to reference one resource inside the other one?](#what-identifiers-are-valid-inside-the-provider-and-how-to-reference-one-resource-inside-the-other-one)
* [Commonly known issues](#commonly-known-issues)
  * [Old Terraform CLI version](#old-terraform-cli-version)
  * [Errors with connection to Snowflake](#errors-with-connection-to-snowflake)
  * [How to set up the connection with the private key?](#how-to-set-up-the-connection-with-the-private-key)
  * [Incorrect identifier (index out of bounds) (even with the old error message)](#incorrect-identifier-index-out-of-bounds-even-with-the-old-error-message)
  * [Incorrect account identifier (snowflake_database.from_share)](#incorrect-account-identifier-snowflake_databasefrom_share)
  * [Granting on Functions or Procedures](#granting-on-functions-or-procedures)
  * [Infinite diffs, empty privileges, errors when revoking on grant resources](#infinite-diffs-empty-privileges-errors-when-revoking-on-grant-resources)
  * [Granting PUBLIC role fails](#granting-public-role-fails)
  * [Issues with grant_ownership resource](#issues-with-grant_ownership)
  * [Using QUOTED_IDENTIFIERS_IGNORE_CASE with the provider](#using-quoted_identifiers_ignore_case-with-the-provider)
  * [Experiencing Go related issues (e.g., using Suricata-based firewalls, like AWS Network Firewall, with >=v1.0.4 version of the provider)](#experiencing-go-related-issues-eg-using-suricata-based-firewalls-like-aws-network-firewall-with-v104-version-of-the-provider)

This guide was made to aid with creating the GitHub issues, so you can maximize your chances of getting help as quickly as possible.
To correctly report the issue, we suggest going through the following steps.

### 1. Check the existing GitHub issues.
Please, go to the [provider’s issues page](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues) and search if there is any other similar issue to the one you would like to report.
This helps us to keep all relevant information in one place, including any potential workarounds.
If you are unsure how to do it, [here](https://docs.github.com/en/issues/tracking-your-work-with-issues/filtering-and-searching-issues-and-pull-requests) is a quick guide showing different filtering options available on GitHub.
It’s good to search by keywords (like **IMPORTED_PRIVILEGES**) or affected resource names (like **snowflake_database**) for quick and effective results.
Remember to search through open and closed issues, because there may be a chance we have already fixed the issue, and it’s working in the latest version of the provider.

### 2. Go through the frequently asked questions and commonly known issues.
We’ve put together a list of frequently asked questions ([FAQ](#faq)) and [commonly known issues](#commonly-known-issues).
Please check, If the answer you're looking for is in one of the lists.

### 3. Check the official Snowflake documentation
It’s common for some Snowflake objects to have special cases of usage in some scenarios.
Mostly, they are not validated in the provider for various reasons.
All of them can result in errors during `terraform plan` or `terraform apply`.
For those reasons, it’s worth to check [the official Snowflake documentation](https://docs.snowflake.com/) before assuming it’s a provider issue.
Depending on the situation where an error occurred a corresponding [SQL command](https://docs.snowflake.com/en/sql-reference-commands) should be looked at.
Especially take a closer look at the “usage notes” section ([for example](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#usage-notes)) where all the special cases should be listed.

### 4. Choose the correct template and use as much information as possible.
Currently, we have a few predefined templates for different types of GitHub issues.
Choose the appropriate one for your use case ([create an issue](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/new/choose)).
Remember to provide as much information as possible when filling out the forms, especially category and object types which appear in almost every template.
That way we will be able to categorize the issues and plan future improvements. When filling out corresponding templates you need to remember the following:
- Bug - It’s important to know the root cause of the issue, that is why we encourage you to fill out the optional fields If you think they can be essential in the analysis. That way we will be able to answer or fix the issue without asking for additional context.
- General Usage - Like in the case of bugs, any additional context can speed up the process.
- Documentation - If there’s an error somewhere in the documentation, please check the related parts. For example, an error in the documentation for stage could be also found in the dependent resources like external tables.
- Feature Request - Before filling out the feature request, please familiarize yourself with the publicly available [roadmap](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md) in case the problem will be resolved by upcoming plans. Also, it would be helpful to reference the roadmap item if the proposals are closely related. That way, we can take a closer look when doing the planned task.

## FAQ
### When will the Snowflake feature X be available in the provider?
It depends on the status of the feature. Snowflake marks features as follows:
- Private Preview (PrPr)
- Public Preview (PuPr)
- Generally Available (GA)

Currently, our main focus is on making the provider stable with the most stable GA features,
but please take a closer look at our recently updated [roadmap](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#05052024-roadmap-overview)
which describes our priorities for the next quarters.

### When will my bug report be fixed/released?
Our team is checking daily incoming GitHub issues. The resolution depends on the complexity and the topic of a given issue, but the general rules are:
- If the issue is easy enough, we tend to answer it immediately and provide fix depending on the issue and our current workload.
- If the issue needs more insight, we tend to reproduce the issue usually in the matter of days and answer/fix it right away (also very dependent on our current workload).
- If the issue is a part of the incoming topic on the [roadmap](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md), we postpone it to resolve it with the related tasks.

The releases are usually happening once every two weeks, mostly done on Wednesday.

### How to migrate from version X to Y?
As noted at the top of our [README](https://github.com/Snowflake-Labs/terraform-provider-snowflake?tab=readme-ov-file#snowflake-terraform-provider),
the project is still experimental and breaking change may occur. We try to minimize such changes, but with some of the changes required for version 1.0.0, it’s inevitable.
Because of that, whenever we introduce any breaking change, we add it to the [migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md).
It’s a document containing every breaking change (starting from around v0.73.0) with additional hints on how to migrate resources between the versions.

### What are the current/future plans for the provider?
Our current plans are documented in the publicly available [roadmap](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md) that you can find in our repository.
We will be updating it to keep you posted on what’s coming for the provider.

### How can I contribute?
If you would like to contribute to the project, please follow our [contribution guidelines](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/CONTRIBUTING.md).

### How can I debug the issue myself?
The provider is simply an abstraction issuing SQL commands through the Go Snowflake driver, so most of the errors will be connected to incorrectly built or executed SQL statements.
To see what SQLs are being run you have to set more verbose logging check the [section below](#how-can-i-turn-on-logs).
To confirm the correctness of the SQLs, refer to the [official Snowflake documentation](https://docs.snowflake.com/).

### How can I turn on logs?
The provider offers two main types of logging:
- Terraform execution (check [Terraform Debugging Documentation](https://www.terraform.io/internals/debugging)) - you can set it through the `TF_LOG` environment variable, e.g.: `TF_LOG=DEBUG`; it will make output of the Terraform execution more verbose.
- Snowflake communication (using the logs from the underlying [Go Snowflake driver](https://github.com/snowflakedb/gosnowflake)) - you can set it directly in the provider config ([`driver_tracing`](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/1.0.3/docs#driver_tracing-3) attribute), by `SNOWFLAKE_DRIVER_TRACING` environmental variable (e.g. `SNOWFLAKE_DRIVER_TRACING=info`), or by `drivertracing` field in the TOML file. To see the communication with Snowflake (including the SQL commands run) we recommend setting it to `info`.

As driver logs may seem cluttered, to locate the SQL commands run, search for:
- (preferred) `--terraform_provider_usage_tracking`
- `msg="Query:`
- `msg="Exec:`

### How can I import already existing Snowflake infrastructure into Terraform?
Please refer to [this document](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/guides/resource_migration.md#3-two-options-from-here)
as it describes different approaches of importing the existing Snowflake infrastructure into Terraform as configuration.
One thing worth noting is that some approaches can be automated by scripts interacting with Snowflake and generating needed configuration blocks,
which is highly recommended for large-scale migrations.

### What identifiers are valid inside the provider and how to reference one resource inside the other one?
Please refer to [this document](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/guides/identifiers_rework_design_decisions.md)
- For the recommended identifier format, take a look at the ["Known limitations and identifier recommendations"](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/guides/identifiers_rework_design_decisions.md#known-limitations-and-identifier-recommendations) section.
- For a new way of referencing object identifiers in resources, take a look at the ["New computed fully qualified name field in resources" ](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/guides/identifiers_rework_design_decisions.md#new-computed-fully-qualified-name-field-in-resources) section.

## Commonly known issues
### Old Terraform CLI version
**Problem:** Sometimes you can get errors similar to:
```text
│ Error: Provider produced invalid plan
│
│ Provider "registry.terraform.io/snowflake-labs/snowflake" planned an invalid value for snowflake_schema_grant.schema_grant.on_all: planned value cty.False for a
│ non-computed attribute.
│
│ This is a bug in the provider, which should be reported in the provider's own issue tracker.
│
```
[GitHub issue reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2347)
**Solution:** You have to be using at least 1.1.5 version of the Terraform CLI.

### Errors with connection to Snowflake
**Problem**: If you are getting connection errors with Snowflake error code, similar to this one:
```text
│
│ Error: open snowflake connection: 390144 (08004): JWT token is invalid.
│
```

**Related issues**: [Experiencing Go related issues (e.g., using Suricata-based firewalls, like AWS Network Firewall, with >=v1.0.4 version of the provider)](#experiencing-go-related-issues-eg-using-suricata-based-firewalls-like-aws-network-firewall-with-v104-version-of-the-provider)

**Solution**: Go to the [official Snowflake documentation](https://docs.snowflake.com/en/user-guide/key-pair-auth-troubleshooting#list-of-errors) and search by error code (390144 in this case).

[GitHub issue reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2432#issuecomment-1915074774)

**Problem**: Getting `Error: 260000: account is empty` error with non-empty `account` configuration after upgrading to v1, with the same provider configuration which worked up to v0.100.0

**Solution**: `account` configuration [has been removed in v1.0.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#removed-deprecated-objects). Please specify your organization name and account name separately as mentioned in the [migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#removed-deprecated-objects):
* `account_name` (`accountname` if you're sourcing it from `config` TOML)
* `organization_name` (`organizationname` if you're sourcing it from `config` TOML)

GitHub issue reference: [#3198](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3198), [#3308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3308)

### How to set up the connection with the private key?
**Problem:** From the version v0.78.0, we introduced a lot of provider configuration changes. One of them was deprecating `private_key_path` in favor of `private_key`.

[GitHub issue reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2489), [Migration Guide reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#v0730--v0740)

**Solution:** Use a non-deprecated `private_key` field with the use of the [file](https://developer.hashicorp.com/terraform/language/functions/file) function to pass the private key.

### Incorrect identifier (index out of bounds) (even with the old error message)
**Problem:** When getting stack traces similar to:
```text
│ panic: runtime error: index out of range [2] with length 2
│
│ goroutine 61 [running]:
│ github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake.SchemaObjectIdentifierFromQualifiedName({0x140001c2870?, 0x103adf987?})
│ github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake/identifier.go:58 +0x174
```

[GitHub issue reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2224)

**Solution:** Some fields may expect different types of identifiers, when in doubt check [our documentation](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs) for the field or the [official Snowflake documentation](https://docs.snowflake.com/) what type of identifier is needed.

### Incorrect identifier type (panic: interface conversion)
**Problem** When getting stack traces similar to:
```text
panic: interface conversion: sdk.ObjectIdentifier is sdk.AccountObjectIdentifier, not sdk.DatabaseObjectIdentifier
```

[GitHub issue reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2779)

**Solution:** Some fields may expect different types of identifiers, when in doubt check [our documentation](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs) for the field or the [official Snowflake documentation](https://docs.snowflake.com/) what type of identifier is needed. Quick reference:
- AccountObjectIdentifier - `<name>`
- DatabaseObjectIdentifier - `<database>.<name>`
- SchemaObjectIdentifier - `<database>.<schema>.<name>`
- TableColumnIdentifier - `<database>.<schema>.<table>.<name>`

### Incorrect account identifier (snowflake_database.from_share)
**Problem:** From 0.87.0 version, we are quoting incoming external account identifier correctly, which may break configurations that specified account identifier as `<org_name>.<acc_name>` that worked previously by accident.

[GitHub issue reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2590)

**Solution:** As specified in the [migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#behavior-change-external-object-identifier-changes), use account locator instead.

### Granting on Functions or Procedures
**Problem:** Right now, when granting any privilege on Function or Procedure with this or similar configuration:

```terraform
resource "snowflake_grant_privileges_to_account_role" "grant_on_procedure" {
  privileges        = ["USAGE"]
  account_role_name = snowflake_account_role.name
  on_schema_object {
    object_type = "PROCEDURE"
    object_name = "\"${snowflake_database.database.name}\".\"${snowflake_schema.schema.name}\".\"${snowflake_procedure_sql.procedure.name}\""
  }
}
```

You may encounter the following error:
```text
│ Error: 090208 (42601): Argument types of function 'procedure_name' must be
│ specified.
```

**Related issues:** [#2375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2375), [#2922](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2922)

**Solution:** Specify the arguments in the `object_name`:

```terraform
resource "snowflake_grant_privileges_to_account_role" "grant_on_procedure" {
  privileges        = ["USAGE"]
  account_role_name = snowflake_account_role.name
  on_schema_object {
    object_type = "PROCEDURE"
    object_name = "\"${snowflake_database.database.name}\".\"${snowflake_schema.schema.name}\".\"${snowflake_procedure_sql.procedure.name}\"(NUMBER, VARCHAR)"
  }
}
```

If you manage the procedure in Terraform, you can use `fully_qualified_name` field:

```terraform
resource "snowflake_grant_privileges_to_account_role" "grant_on_procedure" {
  privileges        = ["USAGE"]
  account_role_name = snowflake_account_role.name
  on_schema_object {
    object_type = "PROCEDURE"
    object_name = snowflake_procedure_sql.procedure_name.fully_qualified_name
  }
}
```

### Infinite diffs, empty privileges, errors when revoking on grant resources
**Problem:** If you encountered one of the following issues:
- Issue with revoking: `Error: An error occurred when revoking privileges from an account role.
- Plan in every `terraform plan` run (mostly empty privileges)
It's possible that the `object_type` you specified is "incorrect."
Let's say you would like to grant `SELECT` on event table. In Snowflake, it's possible to specify
`TABLE` object type instead of dedicated `EVENT TABLE` one. As `object_type` is one of the fields
we filter on, it needs to exactly match with the output provided by `SHOW GRANTS` command.

**Related issues:** [#2749](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2749), [#2803](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2803)

**Solution:** Here's a list of things that may help with your issue:
- Firstly, check if the privilege has been granted in Snowflake. If it is, it means the configuration is correct (or at least compliant with Snowflake syntax).
- When granting `IMPORTED PRIVILEGES` on `SNOWFLAKE` database/application, use `object_type = "DATABASE"`.
- Run `SHOW GRANTS` command with the right filters to find the granted privilege and check what is the object type returned of that command. If it doesn't match the one you have in your configuration, then follow those steps:
  - Use state manipulation (no revoking)
    - Remove the resource from your state using `terraform state rm`.
    - Change the `object_type` to correct value.
    - Import the state from Snowflake using `terraform import`.
  - Remove the grant configuration and after `terraform apply` put it back with the correct `object_type` (requires revoking).

### Granting PUBLIC role fails
**Problem:** When you try granting PUBLIC role, like:
```terraform
resource "snowflake_account_role" "any_role" {
  name = "ANY_ROLE"
}

resource "snowflake_grant_account_role" "this_is_a_bug" {
  parent_role_name = snowflake_account_role.any_role.name
  role_name        = "PUBLIC"
}
```
Terraform may fail with:
```
╷
│ Error: Provider produced inconsistent result after apply
│
│ When applying changes to snowflake_grant_account_role.this_is_a_bug, provider "provider["registry.terraform.io/snowflake-labs/snowflake"]" produced an
│ unexpected new value: Root object was present, but now absent.
│
│ This is a bug in the provider, which should be reported in the provider's own issue tracker.
╵
```

**Related issues:** [#3001](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3001), [#2848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2848)

**Solution:** This happens, because the PUBLIC role is a "pseudo-role" (see [docs](https://docs.snowflake.com/en/user-guide/security-access-control-overview#system-defined-roles)) that is already assigned to every role and user, so there is no need to grant it through Terraform. If you have an issue with removing the resources please use `terraform state rm <id>` to remove the resource from the state (and you can safely remove the configuration).

### Issues with grant_ownership

Please read our [guide for grant_ownership](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/guides/grant_ownership_common_use_cases) resource.
It contains common use cases and issues that you may encounter when dealing with ownership transfers.

### Using QUOTED_IDENTIFIERS_IGNORE_CASE with the provider

**Problem:** When `QUOTED_IDENTIFIERS_IGNORE_CASE` parameter is set to true, but resource identifier fields are filled with lowercase letters,
during `terrform apply` they may fail with the `Error: Provider produced inconsistent result after apply` error (removing themselves from the state in the meantime).

**Related issues:** [#2967](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2967)

**Solution:** Either turn off the parameter or adjust your configurations to use only upper-cased names for identifiers and import back the resources.

### Experiencing Go related issues (e.g., using Suricata-based firewalls, like AWS Network Firewall, with >=v1.0.4 version of the provider)

**Problem:** The communication from the provider is dropped, because of the firewalls that are incompatible with the `tlskyber` setting introduced in [Go v1.23](https://go.dev/doc/godebug#go-123).

**Related issues:** [#3421](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3421)

**Solution:** [Solution described in the migration guide for v1.0.3 to v1.0.4 upgrade](./MIGRATION_GUIDE.md#new-go-version-and-conflicts-with-suricata-based-firewalls-like-aws-network-firewall).
