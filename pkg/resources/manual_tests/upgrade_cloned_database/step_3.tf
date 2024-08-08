# Commands to run
# - terraform import snowflake__database.cloned '"cloned"' (import cloned database into state)
# - terraform plan (expect empty plan)

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
  data_retention_time_in_days = 0 # to avoid in-place update to -1
}
