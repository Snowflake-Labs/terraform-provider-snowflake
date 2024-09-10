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

provider "snowflake" {
  profile = "secondary_test_account"
  alias = second_account
}

resource "snowflake_database" "primary" {
  provider = snowflake.second_account
  name = "test"
  data_retention_time_in_days = 0 # to avoid in-place update to -1
  replication_configuration {
    accounts             = ["<second_account_account_locator>"] # TODO: Replace
    ignore_edition_check = true
  }
}

resource "snowflake_database" "secondary" {
  name = "test"
  data_retention_time_in_days = 0 # to avoid in-place update to -1
  from_replica = "<second_account_account_locator>.\"${snowflake_database.primary.name}\"" # TODO: Replace
}
