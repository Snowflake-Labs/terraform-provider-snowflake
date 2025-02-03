---
page_title: "Grant ownership - common use cases"
subcategory: ""
description: |-

---
# Grant ownership - common use cases

This guide is a follow-up for the [grant_ownership resource overview](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/b863d2e79ae6ae021552c4348e3012b8053ede17/docs/technical-documentation/resource_migration.md) document.
These examples should guide you through Snowflake role management in Terraform with the use of grant_ownership resource.
Here's a list of grant ownership common use cases:

- [Basic RBAC example](#basic-rbac-example)
- [Granting ownership with a less privileged role (granting MANAGED ACCESS)](#granting-ownership-with-a-less-privileged-role-granting-managed-access)
- [Modifying objects you don't own after transferring the ownership](#modifying-objects-you-dont-own-after-transferring-the-ownership)
- [Fixing the state after using a less privileged role in grant_ownership resource](#fixing-the-state-after-using-a-less-privileged-role-in-grant_ownership-resource)

This list may be further extended with more cases; please approach us through [GitHub issue](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/new?template=03-documentation.yml)
if you would like to see any others or contribute ([contribution guidelines](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/b863d2e79ae6ae021552c4348e3012b8053ede17/CONTRIBUTING.md)).

### Basic RBAC example
Here's an easy example of using RBAC (Role-based Access Control). Of course, there are many ways to perform RBAC, and here, we are not proposing any
option over the other. It is only supposed to show, more or less, how the grant_ownership could be used in such a scenario.
The approach is depending on the use case and should be first consulted with a Snowflake Account Manager before creating any role-based system right away.

Keep in mind that this example uses highly privileged role (ACCOUNTADMIN) and for lower privileges roles, you should look into
other examples to see what else is needed to perform the same actions.

#### First deployment
This configuration imitates the "main" Terraform deployment that manages the account objects

```terraform
provider "snowflake" {
  role = "ACCOUNTADMIN"
  # ...
}

resource "snowflake_account_role" "team_a" {
  name = "TEAM_A_ROLE"
}

resource "snowflake_account_role" "team_b" {
  name = "TEAM_B_ROLE"
}

# Make <team_a_user> able to use the TEAM_A_ROLE
resource "snowflake_grant_account_role" "grant_team_a_role" {
  role_name = snowflake_account_role.team_a.name
  user_name = "<team_a_user>"
}

# Make <team_b_user> able to use the TEAM_B_ROLE
resource "snowflake_grant_account_role" "grant_team_b_role" {
  role_name = snowflake_account_role.team_b.name
  user_name = "<team_b_user>"
}

resource "snowflake_database" "team_a_database" {
  name = "TEST_DATABASE"
}

resource "snowflake_grant_ownership" "grant_team_a_database" {
  account_role_name = snowflake_account_role.team_a.name
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.team_a_database.name
  }
}
```

#### Second deployment
If the second deployment uses different user, then the TEST_A_ROLE should be granted to that user in the first deployment first.
By using our ownership of the TEST_DATABASE, we can manage its further access to other teams.

```terraform
provider "snowflake" {
  role = "TEAM_A_ROLE"
  # ...
}

resource "snowflake_schema" "team_b_schema" {
  database = "TEST_DATABASE"
  name = "TEAM_B_SCHEMA"
}

resource "snowflake_grant_privileges_to_account_role" "grant_access_to_database" {
  account_role_name = "TEAM_B_ROLE"
  privileges = ["USAGE"]
  on_account_object {
    object_type = "DATABASE"
    object_name = "TEST_DATABASE"
  }
}

resource "snowflake_grant_privileges_to_account_role" "grant_access_to_schema" {
  account_role_name = "TEAM_B_ROLE"
  privileges = ["USAGE"]
  on_schema {
    schema_name = snowflake_schema.team_b_schema.fully_qualified_name
  }
}

resource "snowflake_grant_privileges_to_account_role" "grant_privileges_to_team_b" {
  account_role_name = "TEAM_B_ROLE"
  privileges = ["USAGE", "CREATE TABLE", "CREATE VIEW"]
  on_schema {
    schema_name = snowflake_schema.team_b_schema.fully_qualified_name
  }
}
```

Then a team using TEAM_B_ROLE can take it from here and create all the tables / views they need (in the worksheet SQL or in any other way).
Just to confirm the above configuration work, you can use the following script:

```snowflake
USE ROLE TEAM_B_ROLE;
USE DATABASE TEST_DATABASE;
USE SCHEMA TEAM_B_SCHEMA;
CREATE TABLE TEST_TABLE(N INT);
-- Has only privilege to create tables and views, so the following command will fail:
CREATE TASK TEST_TASK SCHEDULE = '60 MINUTES' AS SELECT CURRENT_TIMESTAMP;
```

### Granting ownership with a less privileged role (granting MANAGED ACCESS)

This example shows how a less privileged role can be used to transfer ownership of the objects they currently own.
Read more in the [official Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#access-control-requirements).
For this setup, the necessary objects were created by running:

```snowflake
USE ROLE ACCOUNTADMIN;
CREATE ROLE LESS_PRIVILEGED_ROLE;
CREATE ROLE ANOTHER_LESS_PRIVILEGED_ROLE;
GRANT ROLE LESS_PRIVILEGED_ROLE TO USER '<your_terraform_user>';
GRANT CREATE DATABASE, MANAGE GRANTS ON ACCOUNT TO ROLE LESS_PRIVILEGED_ROLE;
```

and after the initial, setup the following configuration can be tested:

```terraform
provider "snowflake" {
  role = "LESS_PRIVILEGED_ROLE"
}

resource "snowflake_database" "test_database" {
  name = "TEST_DATABASE"
}

resource "snowflake_grant_ownership" "grant_ownership_to_another_role" {
  account_role_name = "ANOTHER_LESS_PRIVILEGED_ROLE"
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.test_database.name
  }
}
```

The ownership transfer is possible because here you have both:
- Ownership of the created above database.
- MANAGE GRANTS privilege on the currently used role.

Once the ownership is taken away, you still must be able to take the ownership back to the original role, so that
the Terraform is able to perform successful delete operation once the resource is removed from the configuration.
If you used a less privileged role to grant ownership, [here's an example](#fixing-the-state-after-using-a-less-privileged-role-in-grant_ownership-resource) of how the errors may look like and how to fix them.

That being said, granting ownership would be still possible without MANAGE GRANTS, but you wouldn't be able to grant
the ownership back to the original role. This is a common mistake when dealing with ownership transfers. With Terraform, you have to think
about ownership when it's taken away from the current role, and what will happen when you would like to take it back.

Currently, the least privileged role that is able to transfer ownership has to have at least MANAGE GRANTS privilege.
In the future, we are planning to support other mechanisms that would allow you to use roles without MANAGE GRANTS.
However, other assumptions would be imposed, e.g., that the current user is granted to the role it transfers the ownership to.

### Modifying objects you don't own after transferring the ownership

By transferring ownership of an object to another role, you are limiting currently used role's access control on this object.
This doesn't play well with Terraform ideology that the resource "owns" its part on the infrastructure
and should be able to make changes on that object to eventually match the configuration with the state on the Snowflake side.
By limiting privileges for that resource to make changes on the Snowflake side, you may encounter strange errors related to limited access.
You can commonly encounter this when there will be a need for updating an object after its ownership was transferred to another role. Note that
every object has its access requirements and privileges needed to perform certain actions could be different in your case.
The example was also done on ACCOUNTADMIN role, which means depending on the use case; additional privileges could be necessary for a given action to run successfully.
Imagine you have the following configuration, and you want to change the comment parameter of the database:

```terraform
provider "snowflake" {
  role = "ACCOUNTADMIN"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_account_role" "test" {
  name = "test_role"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}
```

Then, some day you would like to change the comment property of the database like so:

```terraform
# ...

resource "snowflake_database" "test" {
  name = "test_database"
  comment = "new comment"
}

# ...
```

With the current setup, you will encounter the following error (or similar one):
```text
│ Error: 003001 (42501): SQL access control error:
│ Insufficient privileges to operate on database 'test_database'
```

This happened, because now, you don't own this database, and your current role cannot perform any actions on it.
To let the current role modify the database it doesn't own you possibly have a few choices.
1. One of the possible options is to grant the currently used role with necessary privilege (we chose this one in the examples below).
2. Another one could be to create a hierarchy of roles that would possibly allow you to possess certain privileges to the database.

There are possibly more paths that lead to the same place, but to keep it simple, we focus on less extreme cases.

Also, keep in mind that the currently used role has MANAGE GRANTS privilege which makes it easier.
Currently, using less privileged roles (minimum is having MANAGE GRANTS privilege) is not possible.
It will be available once more functionalities are added to the resource.

Going back to the example, firstly, you have to revert the database change and grant the correct privilege (MODIFY) to be able to set the comment on the database.

```terraform
# ...

resource "snowflake_database" "test" {
  name = "test_database"
  # comment = "new comment"
}

resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "ACCOUNTADMIN"
  privileges = [ "MODIFY" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}

resource "snowflake_grant_ownership" "test" {
  depends_on = [ snowflake_grant_privileges_to_account_role.test ]

  # ...
}

# ...
```

After that, you should be able to set the comment of your database, here's how the complete configuration should look like:

```terraform
provider "snowflake" {
  role = "ACCOUNTADMIN"
}

resource "snowflake_database" "test" {
  name = "test_database"
  comment = "new comment"
}

resource "snowflake_account_role" "test" {
  name = "test_role"
}

resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "ACCOUNTADMIN"
  privileges = [ "MODIFY" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}

resource "snowflake_grant_ownership" "test" {
  depends_on = [ snowflake_grant_privileges_to_account_role.test ]

  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}
```

This shows that using ownership transfer (either in provider or only in Snowflake)
requires pre-planning on the overall access architecture and foresight in possible incoming changes.
Otherwise, It may be challenging to introduce certain changes afterward.

### Fixing the state after using a less privileged role in grant_ownership resource

Here's a short example showing how this could look like. Firstly, let's prepare a few objects on the Snowflake side:

```snowflake
CREATE ROLE CREATE_DATABASES_ROLE;
CREATE ROLE ANOTHER_ROLE;
GRANT CREATE DATABASE ON ACCOUNT TO ROLE CREATE_DATABASES_ROLE;
GRANT ROLE CREATE_DATABASES_ROLE TO USER "<terraform_user_name>";
```

then run the following configuration:

```terraform
provider "snowflake" {
  role = "CREATE_DATABASES_ROLE"
}

resource "snowflake_database" "test" {
  name = "TEST_DATABASE_NAME"
}

resource "snowflake_grant_ownership" "transfer_ownership" {
  account_role_name = "ANOTHER_ROLE"
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}
```

After the first apply the provider will raise a warning and an error:
```text
╷
│ Warning: Failed to retrieve grants. Marking the resource as removed.
│ 
│   with snowflake_grant_ownership.transfer_ownership,
│   on main.tf line 18, in resource "snowflake_grant_ownership" "transfer_ownership":
│   18: resource "snowflake_grant_ownership" "transfer_ownership" {
│ 
│ Id: 
│ Error: [errors.go:22] object does not exist or not authorized
╵
╷
│ Error: Provider produced inconsistent result after apply
│ 
│ When applying changes to snowflake_grant_ownership.transfer_ownership, provider
│ "provider[\"registry.terraform.io/snowflake-labs/snowflake\"]" produced an unexpected new value: Root object was present, but now absent.
│ 
│ This is a bug in the provider, which should be reported in the provider's own issue tracker.
╵
```

What happened is after ownership transfers, the current role lost the ability to confirm that the ownership is granted to the correct role.
Because of that, the grant_ownership resource produces inconsistent results and database resource removed itself from the state because
from its perspective the database wasn't created (it couldn't find the database by calling SHOW DATABASES).

To fix this issue, you have to firstly grant the ownership back to the original role. You have to do this from a role
that has at least MANAGE GRANTS privilege (e.g. ACCOUNTADMIN or a custom role with this privilege).

```snowflake
GRANT OWNERSHIP ON DATABASE TEST_DATABASE_NAME TO ROLE CREATE_DATABASES_ROLE;
```

then you have to adjust the configuration, so the ownership is commented out (or completely removed), and import the database resource.

```shell
terraform import snowflake_database.test '"TEST_DATABASE_NAME"'
```

At this point, your configuration should look similar to this:

```terraform
provider "snowflake" {
  role = "CREATE_DATABASES_ROLE"
}

resource "snowflake_database" "test" {
  name = "TEST_DATABASE_NAME"
}

# resource "snowflake_grant_ownership" "transfer_ownership" {
#   account_role_name = "ANOTHER_ROLE"
#   on {
#     object_type = "DATABASE"
#     object_name = snowflake_database.test.name
#   }
# }
```

After running `terraform plan` you should see no changes planned from the provider side,
and you can start over from this point to grant ownership again, but now apply it using one of the provided examples.
