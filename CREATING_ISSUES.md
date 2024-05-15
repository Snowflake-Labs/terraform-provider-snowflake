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
* [Commonly known issues](#commonly-known-issues)
  * [Old Terraform CLI version](#old-terraform-cli-version)
  * [Errors with connection to Snowflake](#errors-with-connection-to-snowflake)
  * [How to set up the connection with the private key?](#how-to-set-up-the-connection-with-the-private-key)
  * [Incorrect identifier (index out of bounds) (even with the old error message)](#incorrect-identifier-index-out-of-bounds-even-with-the-old-error-message)
  * [Incorrect account identifier (snowflake_database.from_share)](#incorrect-account-identifier-snowflake_databasefrom_share)

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
To see what SQLs are being run you have to set the `TF_LOG=DEBUG` environment variable. 
To confirm the correctness of the SQLs, refer to the [official Snowflake documentation](https://docs.snowflake.com/).

### How can I import already existing Snowflake infrastructure into Terraform?
Please refer to [this document](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/resource_migration.md#3-two-options-from-here) 
as it describes different approaches of importing the existing Snowflake infrastructure into Terrafrom as configuration. 
One thing worth noting is that some approaches can be automated by scripts interacting with Snowflake and generating needed configuration blocks, 
which is highly recommended for large-scale migrations.

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
**Solution**: Go to the [official Snowflake documentation](https://docs.snowflake.com/en/user-guide/key-pair-auth-troubleshooting#list-of-errors) and search by error code (390144 in this case).

[GitHub issue reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2432#issuecomment-1915074774)

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
- NewTableColumnIdentifier - `<database>.<schema>.<table>.<name>`

### Incorrect account identifier (snowflake_database.from_share)
**Problem:** From 0.87.0 version, we are quoting incoming external account identifier correctly, which may break configurations that specified account identifier as `<org_name>.<acc_name>` that worked previously by accident.

[GitHub issue reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2590)

**Solution:** As specified in the [migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#behavior-change-external-object-identifier-changes), use account locator instead.
