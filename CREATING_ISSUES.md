# Creating GitHub issues

> **Note:** Part of this guide was moved to the following documentsdas:
> - [FAQ](./FAQ.md)
> - [Known issues](./KNOWN_ISSUES.md).
>
> Their subsections were not removed to preserve already existing links to them.

* [1. Check the documentation for a given resource](#1-check-the-documentation-for-a-given-resource)
* [2. Check the existing GitHub issues.](#2-check-the-existing-github-issues)
* [3. Go through the frequently asked questions and commonly known issues.](#3-go-through-the-frequently-asked-questions-and-commonly-known-issues)
* [4. Check the official Snowflake documentation](#4-check-the-official-snowflake-documentation)
* [5. Choose the correct template and use as much information as possible.](#5-choose-the-correct-template-and-use-as-much-information-as-possible)

This guide was made to aid with creating the GitHub issues, so you can maximize your chances of getting help as quickly as possible.
To correctly report the issue, we suggest going through the following steps.

### 1. Check the documentation for a given resource
Some resources have known limitations, they are usually described in the documentation for a given resource at the top ([example](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/account)).
Those limitations are mostly coming out of limited support for certain commands or Snowflake side or limited output from SHOW/DESC commands.
Please, make sure that the issue you are experiencing is not related to one of the limitations.
Otherwise, follow the guidelines provided with the limitation.

### 2. Check the existing GitHub issues.
Please, go to the [provider’s issues page](https://github.com/snowflakedb/terraform-provider-snowflake/issues) and search if there is any other similar issue to the one you would like to report.
This helps us to keep all relevant information in one place, including any potential workarounds.
If you are unsure how to do it, [here](https://docs.github.com/en/issues/tracking-your-work-with-issues/filtering-and-searching-issues-and-pull-requests) is a quick guide showing different filtering options available on GitHub.
It’s good to search by keywords (like **IMPORTED_PRIVILEGES**) or affected resource names (like **snowflake_database**) for quick and effective results.
Remember to search through open and closed issues, because there may be a chance we have already fixed the issue, and it’s working in the latest version of the provider.

### 3. Go through the frequently asked questions and commonly known issues.
We’ve put together a list of frequently asked questions ([FAQ](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/FAQ.md)) and [commonly known issues](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/KNOWN_ISSUES.md).
Please check, If the answer you're looking for is in one of the lists.

### 4. Check the official Snowflake documentation
It’s common for some Snowflake objects to have special cases of usage in some scenarios.
Mostly, they are not validated in the provider for various reasons.
All of them can result in errors during `terraform plan` or `terraform apply`.
For those reasons, it’s worth to check [the official Snowflake documentation](https://docs.snowflake.com/) before assuming it’s a provider issue.
Depending on the situation where an error occurred a corresponding [SQL command](https://docs.snowflake.com/en/sql-reference-commands) should be looked at.
Especially take a closer look at the “usage notes” section ([for example](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#usage-notes)) where all the special cases should be listed.

### 5. Choose the correct template and use as much information as possible.
Currently, we have a few predefined templates for different types of GitHub issues.
Choose the appropriate one for your use case ([create an issue](https://github.com/snowflakedb/terraform-provider-snowflake/issues/new/choose)).
Remember to provide as much information as possible when filling out the forms, especially category and object types which appear in almost every template.
That way we will be able to categorize the issues and plan future improvements. When filling out corresponding templates you need to remember the following:
- Bug - It’s important to know the root cause of the issue, that is why we encourage you to fill out the optional fields If you think they can be essential in the analysis. That way we will be able to answer or fix the issue without asking for additional context.
- General Usage - Like in the case of bugs, any additional context can speed up the process.
- Documentation - If there’s an error somewhere in the documentation, please check the related parts. For example, an error in the documentation for stage could be also found in the dependent resources like external tables.
- Feature Request - Before filling out the feature request, please familiarize yourself with the publicly available [roadmap](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md) in case the problem will be resolved by upcoming plans. Also, it would be helpful to reference the roadmap item if the proposals are closely related. That way, we can take a closer look when doing the planned task.

## FAQ

Moved to [FAQ.md](./FAQ.md).

## Commonly known issues

Moved to [KNOWN_ISSUES.md](./KNOWN_ISSUES.md).
