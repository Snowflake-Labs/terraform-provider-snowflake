# Commands to run
# - terraform import snowflake_secondary_database.secondary '"test"' (import secondary database into state)
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

provider "snowflake" {
  profile = "secondary_test_account"
  alias = second_account
}

resource "snowflake_database" "primary" {
  provider = snowflake.second_account
  name = "test"
  data_retention_time_in_days = 0 # to avoid in-place update to -1
  replication {
    enable_to_account {
      account_identifier = "<second_account_organization_name>.<second_account_account_name>" # TODO: Replace
      with_failover      = true
    }
    ignore_edition_check = true
  }
}

resource "snowflake_secondary_database" "secondary" {
  name = "test"
  data_retention_time_in_days = 0 # to avoid in-place update to -1
  from_replica = "\"<second_account_organization_name>\".\"<second_account_account_name>\".\"${snowflake_database.primary.name}\"" # TODO: Replace
}
