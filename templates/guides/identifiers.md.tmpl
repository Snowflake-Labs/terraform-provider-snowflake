---
page_title: "Identifiers rework"
subcategory: ""
description: |-

---
# Identifiers rework

## New computed fully qualified name field in resources

With the combination of quotes, old parsing methods, and other factors, it was a struggle to specify the fully qualified name of an object needed (e.g. [#2164](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2164), [#2754](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2754)). Now, with v0.95.0, every resource that represents an object in Snowflake (e.g. user, role), and not an association (e.g. grants) will have a new computed field named `fully_qualified_name`. With the new computed field, it will be much easier to use resources requiring fully qualified names, for examples of usage head over to the [documentation for granting privileges to account role](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_privileges_to_account_role).

For example, instead of writing

```object_name = “\”${snowflake_table.database}\”.\”${snowflake_table.schema}\”.\”${snowflake_table.name}\””```

now we can write

```object_name = snowflake_table.fully_qualified_name```

If you don't manage table in Terraform, you can construct the proper id yourself like before:  `"\"database_name\".\"schema_name\".\"table_name\""`

<!--- TODO: fill the rest of the document -->
