##################################
### on object to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_standard_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_standard_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name   = snowflake_role.test.name
  outbound_privileges = "COPY"
  on {
    object_type = "SCHEMA"
    object_name = "\"${snowflake_standard_database.test.name}\".\"${snowflake_schema.test.name}\""
  }
}

##################################
### on object to database role
##################################

resource "snowflake_standard_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_standard_database.test.name
}

resource "snowflake_standard_database_role" "test" {
  name     = "test_database_role"
  database = snowflake_standard_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  database_role_name  = "\"${snowflake_standard_database_role.test.database}\".\"${snowflake_standard_database_role.test.name}\""
  outbound_privileges = "REVOKE"
  on {
    object_type = "SCHEMA"
    object_name = "\"${snowflake_standard_database.test.name}\".\"${snowflake_schema.test.name}\""
  }
}

##################################
### on all tables in database to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_standard_database" "test" {
  name = "test_database"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_standard_database.test.name
    }
  }
}

##################################
### on all tables in schema to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_standard_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_standard_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    all {
      object_type_plural = "TABLES"
      in_schema          = "\"${snowflake_standard_database.test.name}\".\"${snowflake_schema.test.name}\""
    }
  }
}

##################################
### on future tables in database to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_standard_database" "test" {
  name = "test_database"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    future {
      object_type_plural = "TABLES"
      in_database        = snowflake_standard_database.test.name
    }
  }
}

##################################
### on future tables in schema to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_standard_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_standard_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    future {
      object_type_plural = "TABLES"
      in_schema          = "\"${snowflake_standard_database.test.name}\".\"${snowflake_schema.test.name}\""
    }
  }
}

##################################
### RoleBasedAccessControl (RBAC example)
##################################

resource "snowflake_role" "test" {
  name = "role"
}

resource "snowflake_standard_database" "test" {
  name = "database"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    object_type = "DATABASE"
    object_name = snowflake_standard_database.test.name
  }
}

resource "snowflake_grant_account_role" "test" {
  role_name = snowflake_role.test.name
  user_name = "username"
}

provider "snowflake" {
  profile = "default"
  alias   = "secondary"
  role    = snowflake_role.test.name
}

## With ownership on the database, the secondary provider is able to create schema on it without any additional privileges.
resource "snowflake_schema" "test" {
  depends_on = [snowflake_grant_ownership.test, snowflake_grant_account_role.test]
  provider   = snowflake.secondary
  database   = snowflake_standard_database.test.name
  name       = "schema"
}
