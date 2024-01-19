# Our roadmap

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
