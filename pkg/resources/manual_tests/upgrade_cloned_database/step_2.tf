# Commands to run
# - terraform init - upgrade
# - terraform plan (should observe upgrader errors similar to: failed to upgrade the state with database created from database, please use snowflake_database or deprecated snowflake_database_old instead...)
# - terraform state rm snowflake_database.cloned (remove cloned database from the state)

terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = ">= 0.92.0" # latest
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
