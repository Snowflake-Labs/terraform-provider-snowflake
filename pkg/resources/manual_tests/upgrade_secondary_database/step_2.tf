# Commands to run
# - terraform init - upgrade
# - terraform plan (should observe upgrader errors similar to: failed to upgrade the state with database created from replica, please use snowflake_secondary_database or deprecated snowflake_database_old instead)
# - terraform state rm snowflake_database.secondary (remove secondary database from the state)

terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = ">= 0.92.0" # latest
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
