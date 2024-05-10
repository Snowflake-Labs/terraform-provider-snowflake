# Our roadmap

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
