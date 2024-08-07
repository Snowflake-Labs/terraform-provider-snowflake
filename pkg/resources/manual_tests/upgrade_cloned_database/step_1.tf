# Commands to run
# - terraform apply

terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "0.92.0"
    }
  }
}

provider "snowflake" {}

resource "snowflake_database" "test" {
  name = "test"
  data_retention_time_in_days = 0 # to avoid in-place update to -1
}

resource "snowflake_database" "cloned" {
  name = "cloned"
  from_database = snowflake_database.test.name
  data_retention_time_in_days = 0 # to avoid in-place update to -1
}
