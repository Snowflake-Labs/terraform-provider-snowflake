
# Identifiers rework

## Table of contents
<!-- TOC -->
* [Identifiers rework](#identifiers-rework)
  * [Table of contents](#table-of-contents)
  * [Topics](#topics)
    * [New identifier parser](#new-identifier-parser)
    * [Using the recommended format for account identifiers](#using-the-recommended-format-for-account-identifiers)
    * [Better handling for identifiers with arguments](#better-handling-for-identifiers-with-arguments)
    * [Quoting differences](#quoting-differences)
    * [New computed fully qualified name field in resources](#new-computed-fully-qualified-name-field-in-resources)
    * [New resource identifier format](#new-resource-identifier-format)
    * [Known limitations and identifier recommendations](#known-limitations-and-identifier-recommendations)
    * [New identifier conventions](#new-identifier-conventions)
  * [Next steps](#next-steps)
  * [Conclusions](#conclusions)
<!-- TOC -->

This document summarises work done in the [identifiers rework](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework) and future plans for further identifier improvements. 
But before we dive into results and design decisions, here’s the list of reasons why we decided to rework the identifiers in the first place:
- Common issues with identifiers with arguments (identifiers for functions, procedures, and external functions).
- Meaningless error messages whenever an invalid identifier is specified.
- Inconsistencies in quotes causing differences in Terraform plans.
- The inconvenience of specifying fully qualified names in certain resource fields (e.g. object name in privilege-granting resources).
- Mixed usage of account identifier formats across resources.

Now, knowing the issues we wanted to solve, we would like to present the changes and design decisions we made.

## Topics

### New identifier parser
To resolve many of our underlying problems with parsing identifiers, we decided to go with the new one that will be able to correctly parse fully qualified names of objects. 
In addition to a better parsing function, we made sure it will return user-friendly error messages that will be able to find the root cause of a problem when specifying invalid identifiers. 
Previously, the error looked like [this](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2091).

### Using the recommended format for account identifiers
Previously, the use of account identifiers was mixed across the resources, in many cases causing confusion ([commonly known issues reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/CREATING_ISSUES.md#incorrect-account-identifier-snowflake_databasefrom_share)). 
Some of them required an account locator format (that was not fully supported and is currently deprecated), and some of the new recommended ones. 
We decided to unify them and use the new account identifier format everywhere.

### Better handling for identifiers with arguments
Previously, the handling of identifiers with arguments was not done fully correctly, causing many issues and confusion on how to use them ([commonly known issues reference](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/CREATING_ISSUES.md#granting-on-functions-or-procedures)).
The main pain point was using them with privilege-granting resources. To address this we had to make two steps. 
The first one was adding a dedicated representation of an identifier containing arguments and using it in our SDK. 
The second one was additional parsing for the output of SHOW GRANTS in our SDK which was only necessary for functions, 
procedures, and external functions that returned non-valid identifier formats.

### Quoting differences
There are many reported issues on identifier quoting and how it is inconsistent across resources and causes plan diffs to enforce certain format (e.g. [#2982](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2982), [#2236](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2236)). 
To address that, we decided to add diff suppress on identifier fields that ignore changes related to differences in quotes. 
The main root cause of such differences was that Snowflake has specific rules when a given identifier (or part of an identifier) is quoted and when it’s not. 
The diff suppression should make those rules irrelevant whenever identifiers in your Terraform configuration contain quotes or not.

### New computed fully qualified name field in resources
With the combination of quotes, old parsing methods, and other factors, it was a struggle to specify the fully qualified name of an object needed (e.g. [#2164](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2164), [#2754](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2754)). 
Now, with v0.95.0, every resource that represents an object in Snowflake (e.g. user, role), and not an association (e.g. grants) will have a new computed field named `fully_qualified_name`. 
With the new computed field, it will be much easier to use resources requiring fully qualified names, for examples of usage head over to the [documentation for granting privileges to account role](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_privileges_to_account_role).

### New resource identifier format
This will be a small shift in the identifier representation for resources. The general rule will be now that:
- If a resource can only be described with the Snowflake identifier, the fully qualified name will be put into the resource identifier. Previously, it was almost the same, except it was separated by pipes, and it was not a valid identifier.
- If a resource cannot be described only by a single Snowflake identifier, then the resource identifier will be a pipe-separated text of all parts needed to identify a given resource ([example](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_privileges_to_account_role#import)). Mind that this approach is not compliant with identifiers containing pipes, but this approach is a middle ground between an easy-to-specify separator and a character that shouldn’t be that common in the identifier (it was previously used for all identifiers).

### Known limitations and identifier recommendations
The main limitations around identifiers are strictly connected to what characters are used. Here’s a list of recommendations on which characters should be generally avoided when specifying identifiers:
- Avoid dots ‘.’ inside identifiers. It’s the main separator between identifier parts and although we are handling dots inside identifiers, there may be cases where it’s impossible to parse the identifier correctly.
- Avoid pipes ‘|’ inside identifiers. It’s the separator for our more complex resource identifiers that could make our parser split the resource identifier into the wrong parts.
- Avoid parentheses ‘(’ and ‘)’ when specifying identifiers for functions, procedures, external functions. Parentheses as part of their identifiers could potentially make our parser split the identifier into wrong parts causing issues.

As a general recommendation, please lean toward simple names without any special characters, and if word separation is needed, use underscores. 
This also applies to other “identifiers” like column names in tables or argument names in functions. 
If you are currently using complex identifiers, we recommend considering migration to simpler identifiers for a more straightforward and less error-prone experience.
Also, we want to make it clear that every field specifying identifier (or its part, e.g. `name`, `database`, `schema`) are always case-sensitive. By specifying
identifiers with lowercase characters in Terraform, you also have to refer to them with lowercase names in quotes in Snowflake. 
For example, by specifying an account role with `name = "test"` to check privileges granted to the role in Snowflake, you have to call:
```sql
show grants to role "test";
show grants to role test; -- this won't work, because unquoted identifiers are converted to uppercase according to https://docs.snowflake.com/en/sql-reference/identifiers-syntax#label-identifier-casing
```

### New identifier conventions
Although, we are closing the identifiers rework, some resources won’t have the mentioned improvements. 
They were mostly applied to the objects that were already prepared for v1 ([essential objects](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/ESSENTIAL_GA_OBJECTS.MD)). 
The remaining resources (and newly created ones) will receive these improvements [during v1 preparation](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1) following our internal guidelines that contain those new rules regarding identifiers. 
No matter if the resource has been refactored or not, the same recommendations mentioned above apply.

## Next steps
While we have completed the identifiers rework for now, we plan to revisit these topics in the future to ensure continued improvements. 
In the upcoming phases, we will focus on addressing the following key areas:
- Implementing better validations for identifiers.
- Providing support for new identifier formats in our resources (e.g. [instance roles](https://docs.snowflake.com/en/sql-reference/snowflake-db-classes#instance-roles)).

## Conclusions
We have concluded the identifiers rework, implementing significant improvements to address common issues and inconsistencies in identifier handling. 
Moving forward, we aim to continue enhancing our identifier functionalities to provide a smoother experience.
We value your feedback on the recent changes made to the identifiers. Please share your thoughts and suggestions to help us refine our identifier management further.
