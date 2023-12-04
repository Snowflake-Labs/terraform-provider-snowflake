##################################
### global privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_role" "g1" {
  privileges = ["MODIFY", "USAGE"]
  role_name  = snowflake_role.r.name
  on_account = true
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_role" "g2" {
  role_name         = snowflake_role.r.name
  on_account        = true
  all_privileges    = true
  with_grant_option = true
}

##################################
### account object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_role" "g3" {
  privileges = ["CREATE", "MONITOR"]
  role_name  = snowflake_role.r.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.d.name
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_role" "g4" {
  role_name = snowflake_role.r.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.d.name
  }
  all_privileges    = true
  with_grant_option = true
}

##################################
### schema privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_role" "g5" {
  privileges = ["MODIFY", "CREATE TABLE"]
  role_name  = snowflake_role.r.name
  on_schema {
    schema_name = "\"my_db\".\"my_schema\"" # note this is a fully qualified name!
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_role" "g6" {
  role_name = snowflake_role.r.name
  on_schema {
    schema_name = "\"my_db\".\"my_schema\"" # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

# all schemas in database
resource "snowflake_grant_privileges_to_role" "g7" {
  privileges = ["MODIFY", "CREATE TABLE"]
  role_name  = snowflake_role.r.name
  on_schema {
    all_schemas_in_database = snowflake_database.d.name
  }
}

# future schemas in database
resource "snowflake_grant_privileges_to_role" "g8" {
  privileges = ["MODIFY", "CREATE TABLE"]
  role_name  = snowflake_role.r.name
  on_schema {
    future_schemas_in_database = snowflake_database.d.name
  }
}

##################################
### schema object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_role" "g9" {
  privileges = ["SELECT", "REFERENCES"]
  role_name  = snowflake_role.r.name
  on_schema_object {
    object_type = "VIEW"
    object_name = "\"my_db\".\"my_schema\".\"my_view\"" # note this is a fully qualified name!
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_role" "g10" {
  role_name = snowflake_role.r.name
  on_schema_object {
    object_type = "VIEW"
    object_name = "\"my_db\".\"my_schema\".\"my_view\"" # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

# all in database
resource "snowflake_grant_privileges_to_role" "g11" {
  privileges = ["SELECT", "INSERT"]
  role_name  = snowflake_role.r.name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.d.name
    }
  }
}

# all in schema
resource "snowflake_grant_privileges_to_role" "g12" {
  privileges = ["SELECT", "INSERT"]
  role_name  = snowflake_role.r.name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_schema          = "\"my_db\".\"my_schema\"" # note this is a fully qualified name!
    }
  }
}

# future in database
resource "snowflake_grant_privileges_to_role" "g13" {
  privileges = ["SELECT", "INSERT"]
  role_name  = snowflake_role.r.name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.d.name
    }
  }
}

# future in schema
resource "snowflake_grant_privileges_to_role" "g14" {
  privileges = ["SELECT", "INSERT"]
  role_name  = snowflake_role.r.name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_schema          = "\"my_db\".\"my_schema\"" # note this is a fully qualified name!
    }
  }
}
