resource "snowflake_database" "db" {
  name = "database"
}

resource "snowflake_role" "db_role" {
  name = "role_name"
}

##################################
### on account privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["CREATE DATABASE", "CREATE USER"]
  role_name  = snowflake_role.db_role.name
  on_account = true
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  role_name         = snowflake_role.db_role.name
  on_account        = true
  all_privileges    = true
  with_grant_option = true
}

# all privileges + grant option + always apply
resource "snowflake_grant_privileges_to_account_role" "example" {
  role_name         = snowflake_role.db_role.name
  on_account        = true
  always_apply      = true
  all_privileges    = true
  with_grant_option = true
}

##################################
### on account object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["CREATE SCHEMA", "CREATE DATABASE ROLE"]
  role_name  = snowflake_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  role_name = snowflake_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
  all_privileges    = true
  with_grant_option = true
}

# all privileges + grant option + always apply
resource "snowflake_grant_privileges_to_account_role" "example" {
  role_name = snowflake_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
  always_apply      = true
  all_privileges    = true
  with_grant_option = true
}

##################################
### schema privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["MODIFY", "CREATE TABLE"]
  role_name  = snowflake_role.db_role.name
  on_schema {
    schema_name = "\"${snowflake_database.db.name}\".\"my_schema\"" # note this is a fully qualified name!
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  role_name = snowflake_role.db_role.name
  on_schema {
    schema_name = "\"${snowflake_database.db.name}\".\"my_schema\"" # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

# all schemas in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["MODIFY", "CREATE TABLE"]
  role_name  = snowflake_role.db_role.name
  on_schema {
    all_schemas_in_database = snowflake_database.db.name
  }
}

# future schemas in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["MODIFY", "CREATE TABLE"]
  role_name  = snowflake_role.db_role.name
  on_schema {
    future_schemas_in_database = snowflake_database.db.name
  }
}

##################################
### schema object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["SELECT", "REFERENCES"]
  role_name  = snowflake_role.db_role.name
  on_schema_object {
    object_type = "VIEW"
    object_name = "\"${snowflake_database.db.name}\".\"my_schema\".\"my_view\"" # note this is a fully qualified name!
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  role_name = snowflake_role.db_role.name
  on_schema_object {
    object_type = "VIEW"
    object_name = "\"${snowflake_database.db.name}\".\"my_schema\".\"my_view\"" # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

# all in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["SELECT", "INSERT"]
  role_name  = snowflake_role.db_role.name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.db.name
    }
  }
}

# all in schema
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["SELECT", "INSERT"]
  role_name  = snowflake_role.db_role.name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_schema          = "\"${snowflake_database.db.name}\".\"my_schema\"" # note this is a fully qualified name!
    }
  }
}

# future in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["SELECT", "INSERT"]
  role_name  = snowflake_role.db_role.name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.db.name
    }
  }
}

# future in schema
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges = ["SELECT", "INSERT"]
  role_name  = snowflake_role.db_role.name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_schema          = "\"${snowflake_database.db.name}\".\"my_schema\"" # note this is a fully qualified name!
    }
  }
}
