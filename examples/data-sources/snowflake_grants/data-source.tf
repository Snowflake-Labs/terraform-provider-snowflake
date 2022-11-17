# list all grants on account
data "snowflake_grants" "grants" {
  grants_on {
    account = true
  }
}

# list all grants in database with name "tst"
data "snowflake_grants" "grants2" {
  grants_on {
    object_name = "\"tst\""
    object_type = "DATABASE"
  }
}

# list all grants to role with name "ACCOUNTADMIN"
data "snowflake_grants" "grants3" {
  grants_to {
    role = "ACCOUNTADMIN"
  }
}

# list all grants of role with name "ACCOUNTADMIN"
data "snowflake_grants" "grants4" {
  grants_of {
    role = "ACCOUNTADMIN"
  }
}

# list all grants in database with name "tst"
data "snowflake_grants" "grants5" {
  future_grants_in {
    database = "\"tst\""
  }
}

# list all future grants in schema with name "mydatabase" and database with name "myschema"
data "snowflake_grants" "grants6" {
  future_grants_in {
    schema {
      database_name = "\"mydatabase\""
      schema_name   = "\"myschema\""
    }
  }
}

# list all future grants to role with name "ACCOUNTADMIN"
data "snowflake_grants" "grants7" {
  future_grants_to {
    role = "ACCOUNTADMIN"
  }
}
