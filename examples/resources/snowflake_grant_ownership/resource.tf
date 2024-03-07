##################################
### on object to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name   = snowflake_role.test.name
  outbound_privileges = "COPY"
  on {
    object_type = "SCHEMA"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
  }
}

##################################
### on object to database role
##################################

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_database.test.name
}

resource "snowflake_database_role" "test" {
  name     = "test_database_role"
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  database_role_name  = "\"${snowflake_database_role.test.database}\".\"${snowflake_database_role.test.name}\""
  outbound_privileges = "REVOKE"
  on {
    object_type = "SCHEMA"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
  }
}

##################################
### on all tables in database to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    all {
      plural_object_type = "TABLES"
      in_database        = snowflake_database.test.name
    }
  }
}

##################################
### on all tables in schema to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    all {
      plural_object_type = "TABLES"
      in_schema          = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
    }
  }
}

##################################
### on future tables in database to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    future {
      plural_object_type = "TABLES"
      in_database        = snowflake_database.test.name
    }
  }
}

##################################
### on future tables in schema to account role
##################################

resource "snowflake_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    future {
      plural_object_type = "TABLES"
      in_schema          = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
    }
  }
}

