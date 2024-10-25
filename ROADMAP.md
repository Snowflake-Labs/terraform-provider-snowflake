# Our roadmap

## (25.10.2024) Project state overview

### Goals

Since the last update we have focused on:

* [reducing the feature gap](#reducing-the-feature-gap) (focusing on the Snowflake essential GA resources)
* redesigning identifiers (check [\#3045](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/3045) and [identifiers_rework_design_decisions](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/identifiers_rework_design_decisions.md))
* reworking the provider's configuration (the doc/discussion will be shared when ready)
* researching the object renaming in our provider (the doc will be shared when ready)

These steps were all needed to get us closer to the first stable version of the provider which... is really close. In the next 1-2 months we want to:

* [wrap up the functional scope](#wrap-up-the-functional-scope) (not all the objects will be declared stable, more details below)
* [prepare for the V1 release](#prepare-for-the-v1-release)
* [prepare some basic performance benchmarks](#prepare-some-basic-performance-benchmarks) (especially, after a few major changes to the resources logic)
* [improve/update the documentation](#improveupdate-the-documentation)
* [run a closed early adopter program](#run-a-closed-early-adopter-program) to verify the readiness of the provider to enter a stable V1

If there won't be any major obstacles or critical issues we aim to release V1 on **December 9th**. To better understand its scope, please check the ["what is V1?"](#what-is-v1) section.

#### reducing the feature gap

During the last six months, we have been tackling objects from the [essential](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/ESSENTIAL_GA_OBJECTS.MD) and [remaining](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/REMAINING_GA_OBJECTS.MD) object lists. We’ve been aligning the implementation, adding missing attributes, and fixing known issues of the chosen objects (full list below). We had to make design decisions that sometimes were not only dictated by our engineering assessments but also by the limitations of Terraform and the underlying [SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2). The main decisions are listed inside the repository in the [Design decisions before v1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#design-decisions-before-v1) (we will validate if all essential ones are present there before releasing V1).

#### wrap up the functional scope

It’s about finishing the redesign of objects we want to declare stable. This mainly affects tables and accounts, but it also involves small alterations in other objects (which will be listed in the migration guide as usual).

As shown [below](#which-resources-will-be-declared-stable), all but one of the [essential](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/ESSENTIAL_GA_OBJECTS.MD) objects and a few of the [remaining](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/REMAINING_GA_OBJECTS.MD) objects made the cut.

#### prepare for the V1 release

This is mainly cleaning up the repository but also activities around the release:

* removing deprecated resources
* marking the resources as [preview features](#preview-resourcesdatasources)
* removing deprecated attributes
* potentially renaming some configuration options
* summarizing migration guidelines between v0.x.x and v1.0.0

#### prepare some basic performance benchmarks

During the resources redesign we introduced multiple changes that may affect the performance. Namely:

* more SQL statements are run (`SHOW`, `DESCRIBE`, and `SHOW PARAMETERS` when needed)
* the state we save is bigger because of the `show_output`, `describe_output`, and `parameters`.

We observed that our customers tend to have lots of objects in single terraform deployments. ​​This leads to longer planning and execution times. To be able to guide “what is too much”, we need to perform tests with more objects on our end.

#### improve/update the documentation

We greatly improved the docs and the transparency of the project. However, there are still topics that need our attention (e.g. adding a migration guide directly to the [registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs), adding missing design decisions like granting ownership, or adding more guides \- similar to [identifiers rework](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/guides/identifiers) \- like importing existing infrastructure into Snowflake).

#### run a closed early adopter program

We planned V1 to be as close as possible to the latest 0.x.x version before the V1 release. However, some changes/migrations are still expected. To improve confidence, we have decided to provide early V1 binaries to early adopters. We are still actively recruiting customers; please reach out to your Snowflake Account Manager at the earliest if you would like to participate. The program runs from mid-November to mid-December.

#### what is V1?

The first major version, V1, marks the first step in getting to GA by providing stable versions to customers who use the provider. We hope to have all our current customers migrate to V1. The provider's Product and Engineering teams will be available for migration or any other questions, as we believe this migration is key in preparing our customers for seamless GA adoption.

From the engineering point of view, the provider will be in the stable version, but it will still stay in the Snowflake-Labs GitHub organization. We plan to change that and move it to the official snowflakedb org so that it gets the official Snowflake support. This will be a necessary step to reach the GA.

#### which resources will be declared stable

Check [this list](v1-preparations/LIST_OF_STABLE_RESOURCES_FOR_V1.md) for details.

#### preview resources/datasources

On our road to V1, we went through the resources, starting with the most used ones. We did not cover all of them (as described above). Because of that, in the newest [v0.97.0](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.97.0/docs) version of the provider, we have multiple resources that were not redesigned/fixed.

We discussed two main options: removing them from 1.0.0 or marking them as preview features. We were mostly worried that removing resources would prevent the majority of our users from migrating to the stable version. On the other hand, we know they are not ready so we don’t want to declare them as stable.

After consideration, we decided to leave them as preview features that need to be **explicitly enabled by the user**. This way, we are not reducing the provider's functionality between v0.x.x and v1.0.0 and leave the possibility to use them while accepting the limitations they have. However, these resources will be subject to change after V1. They should be treated as [Snowflake Preview Features](https://docs.snowflake.com/en/release-notes/preview-features) so changes to their schemas (breaking changes included\!) may be introduced even without bumping the major version of the provider.

#### “attachment” resources clarification

During our road to V1 we tried to limit the number of resources needed to be configured in order to manage the given Snowflake object correctly. Because of that, we moved [Snowflake parameters](https://docs.snowflake.com/en/sql-reference/parameters) handling directly to the given object’s resource (check [this](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#snowflake-parameters)). We did that to other types of properties too (e.g. we changed the logic for public keys handling in the [snowflake_user](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.95.0/docs/resources/user#rsa_public_key) resource, so that [snowflake_user_public_keys](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.94.1/docs/resources/user_public_keys) is no longer compatible with it).

Still, these “attachment” objects serve a specific use case (i.e. the main object is not managed by Terraform but part of the object may be). It opened a question for the future not only because of the aforementioned use case but also because of a wider perspective on the default resource behavior. For example, a resource monitor can be attached to a warehouse only by a user with an ACCOUNTADMIN role (check [\#3019](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3019)). Some of our users would like to provision warehouses separately from assigning resource monitors but the caveat here is that leaving the resource monitor empty in the resource config will currently remove any assigned resources. Handling this would require adding a separate attachment resource and allowing a conditional change in behavior for empty assignments in the main object.

The topic is wide. For the V1, we decided to keep most of the attachment resources as [preview features](#preview-resourcesdatasources) and we will discuss the need for handling the use cases described in this section as a separate topic after V1.

#### which resources will be left as preview features

Check [this list](v1-preparations/LIST_OF_PREVIEW_FEATURES_FOR_V1.md) for details.

#### which resources will be removed

Check [this list](v1-preparations/LIST_OF_REMOVED_RESOURCES_FOR_V1.md) for details.

#### roadmap short after V1

Right after V1, we would like to focus on helping all of you with the migration. We will prioritize it so we encourage you to approach us with any issues you might have.

In the meantime, if we have enough time, we want to prioritize redesigning the object marked as preview features. Currently, stages and shares open the list.

#### next year priorities

This is only a general overview of the next year and may be subject to change:

* Graduate out of Snowflake-Labs into the official snowflakedb organization
* GA of the Snowflake Terraform Provider
* Research performance improvements (optimize Snowflake invocations)
* Grants improvements
* Redesign remaining GA objects
* Design transition to the [plugin framework](https://developer.hashicorp.com/terraform/plugin/framework)
* Introduce Terraform modules

## (05.05.2024) Roadmap Overview

### Goals
Since the last update we have focused on:
- [Finishing the SDK rewrite](#finishing-sdk-rewrite).
- [Redesigning grants](#redesigning-grants) (check announcements: [discussions/1890#discussioncomment-9071073](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/1890#discussioncomment-9071073), [discussions/2235](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/2235), and [discussions/2736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/2736)).
- Improving the provider’s stability (by [stabilizing the tests](#tests-stabilization), solving new incoming issues on a daily basis, and [introducing repository-wide fixes to multiple objects](#resolving-existing-issues)).
- Preparing the scope for the V1 (more below). Part of [supporting-all-snowflake-ga-features](#supporting-all-snowflake-ga-features).
- Raising the transparency of the project (this roadmap, [contribution guidelines](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/CONTRIBUTING.md), [old issues cleanup](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/2755), and [FAQ](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/CREATING_ISSUES.md#faq)).

The primary goals we are working on currently are:
- Introducing support for the fundamental GA features and improving the existing objects (resolving existing provider issues included). Continuation of [supporting-all-snowflake-ga-features](#supporting-all-snowflake-ga-features).
- Reworking identifiers.

The more concrete topics we are currently dealing with are presented in the following three sections: current, upcoming, and next.

|                      Current                       |                             Upcoming                             |                  Next                   |
|:--------------------------------------------------:|:----------------------------------------------------------------:|:---------------------------------------:|
| Preparing essential GA objects for the provider V1 | Preparing rest of the fundamental GA objects for the provider V1 |    Official Snowflake public preview    |
|                 Identifiers rework                 |                 Support object renaming properly                 |        Enable Snowflake support         |
|                                                    |                 Provider’s configuration rework                  | Support for the public preview features |
|                                                    |                      Prepare the V1 release                      |                                         |

#### Current (expected mid/late this year)
##### Preparing essential GA objects for the provider V1
As we stated in the [previous entry](#supporting-all-snowflake-ga-features) we want to inspect all the existing objects to find missing parameters and flaws in their designs. We gathered a list of objects we believe are most crucial, and we will address them first. The list is available [here](v1-preparations/ESSENTIAL_GA_OBJECTS.MD). After them, we will address the ones described in the [following entry](#preparing-the-rest-of-the-fundamental-ga-objects-for-the-provider-v1).

##### Identifiers rework
([previous entry](#identifiers-rework-1)) Identifiers were recently the second, next to the Grants, most common error source in users’ configurations. We want to make interaction with them easier (at least to the extent we have control of).

#### Upcoming (expected likely in late Fall)
##### Preparing the rest of the fundamental GA objects for the provider V1
This will be the continuation of [Preparing essential GA objects for the provider V1](#preparing-essential-ga-objects-for-the-provider-v1). It will address objects listed [here](v1-preparations/REMAINING_GA_OBJECTS.MD).

##### Support object renaming properly
Object renaming is a topic that arises in different contexts like renaming a database, column, or schema object to name a few. The renaming topic was brought up a long time ago, e.g. in [#420](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/420), [#753](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/753), and [this forum entry](https://community.snowflake.com/s/question/0D5Do00000KWFhqKAH/how-to-rename-snowflake-database-on-terraform-with-the-snowflake-provider). We want to address the renaming in general before the stable V1.

##### Provider’s configuration rework
([previous entry](#providers-configuration-rework-1)) It is one of the last moments before going V1 to make incompatible changes in the provider. The current configuration contains many deprecated parameters, inconsistencies with the documentation, and other design flaws. We want to address it.

##### Prepare the V1 release
This will be the moment to validate our V1 efforts by checking if everything was implemented and making the migration for all of you as smooth as possible. This includes:
- Listing of all breaking changes
- Summarizing the migration notes
- Communicating the V1 release in detail
- Describing the new release cycle post-V1
- And many more…

**Important** We plan to introduce the changes before the V1 to allow you to migrate most of the objects before the official release. Because we are still not providing the backward bugfixes, it's always best to bump the provider version with the new releases (following the [migration guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#migration-guide)).

#### Next
- Official Snowflake public preview
- Enable Snowflake support
- Support for the public preview features

## (19.01.2024) Roadmap Overview
### Goals
The primary goals we are working on currently are:
- Adding missing and updating existing functionalities (resources and data sources);
- Resolving existing provider issues;
- Improving provider’s stability.

We believe fulfilling these goals will help us reach V1 with a stable, reliable, and functional provider. The more concrete topics we are currently dealing with are presented in the following three sections: current, upcoming, and next.

|                                                                                     Current                                                                                      |                                                                                              Upcoming                                                                                               |                                                                         Next                                                                          |
|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------------------------------------------------------------------------------------------------------------------------------------------------:|
|                                                                                 Redesign Grants                                                                                  | Design proper resources for the majority of Snowflake objects. Support all Snowflake GA features, starting with the most critical resources like databases, schemas, tables, tasks, and warehouses. | Rework provider’s configuration. Covers current configurations deprecated parameters, inconsistencies with the documentation, and other design flaws. |
| Finish SDK rewrite. Migrate existing resources and data sources to the new SDK, aiding in safer and more extendable generation of SQL statements executed against Snowflake API. |                                                                                         Rework identifiers                                                                                          |                                     Stabilization of tests, ensuring quicker development and stability assurance.                                     |
|                                                     Address open issues in the repo repository, focusing on critical issues.                                                     |                                                              Address open issues in the repo repository, focusing on critical issues.                                                               |                                                                                                                                                       |

#### Current
##### Redesigning Grants
Grants proved to be one of the most common pain points for the provider’s users. We have been focusing on designing the proper resources for the past few weeks. The development is in progress, but more topics still need our attention (like granting ownership, and imported privileges, to name a few).

##### Finishing SDK rewrite
Last year, we changed the approach to generating the SQL statements executed against Snowflake API. The previous, old implementation was error-prone and hard to maintain. We are concluding migrating existing resources and data sources to the new SDK we are developing. It has already proved to be safer and more extendable.

##### Resolving existing issues
Having the ~470 open issues in the repository is not fun. We want to reduce that number drastically. We have recently taken multiple different steps to achieve it:
- We respond to most of the incoming issues faster.
- We classified and prioritized the existing issues. We picked the resources that were causing the most trouble for our users. We will focus first on resource monitors, databases, and tasks. At the same time, we introduce improvements in reporting errors and handle common pitfalls globally.
- We plan to close the issues regarding ancient provider versions. There will be a separate announcement about it.

#### Upcoming
##### Supporting all Snowflake GA features
Eventually, we want to support all Snowflake features. We first want to support all the GA ones. It does not only mean that we will add the missing resources; we will also carefully inspect the existing ones to find missing parameters and flaws in their designs. We will start with the most critical resources like databases, schemas, tables, tasks, and warehouses.

##### Identifiers rework
Identifiers were recently the second, next to the Grants, most common error source in users’ configurations. We want to make interaction with them easier (at least to the extent we have control of).

##### Increasing transparency and involving the community in discussions
We are actively being asked about the state of the development, plans for introducing new resources, and design decisions. This roadmap is one of the many steps we are willing to take to be more transparent to our users.

#### Next
##### Provider’s configuration rework
It is one of the last moments before going V1 to make incompatible changes in the provider. The current configuration contains many deprecated parameters, inconsistencies with the documentation, and other design flaws. We want to address it.

##### Tests stabilization
We are extensively testing our provider. We rely on our tests when introducing new features. Unfortunately, historically, testing was not the biggest concern in the project; many tests are missing, and existing ones are not always correct. Having reliable test sets is essential for quicker development and stability assurance.
