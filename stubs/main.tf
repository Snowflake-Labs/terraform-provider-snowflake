terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "~> 0.61"
    }
  }
}

provider "snowflake" {
  role = "SYSADMIN"
}

resource "snowflake_database" "db" {
  name = "TEST_DB"
}

resource "snowflake_schema" "schema" {
  database = snowflake_database.db.name
  name     = "TEST_SCHEMA"
}
