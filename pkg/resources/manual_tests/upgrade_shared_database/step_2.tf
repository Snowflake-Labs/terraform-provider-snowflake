# Commands to run
# - terraform init - upgrade
# - terraform plan (should observe upgrader errors similar to: failed to upgrade the state with database created from share, please use snowflake_shared_database or deprecated snowflake_database_old instead)
# - terraform state rm snowflake_database.from_share (remove shared database from the state)

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

resource "snowflake_share" "test" {
  provider = snowflake.second_account
  name = "test_share"
  accounts = ["<primary_account_organization_name>.<primary_account_account_name>"] # TODO: Replace
}

resource "snowflake_database" "test" {
  provider = snowflake.second_account
  name = "test_database"
}

resource "snowflake_grant_privileges_to_share" "test" {
  provider = snowflake.second_account
  privileges = ["USAGE"]
  on_database = snowflake_database.test.name
  to_share = snowflake_share.test.name
}

resource "snowflake_database" "from_share" {
  depends_on = [ snowflake_grant_privileges_to_share.test ]
  name = snowflake_database.test.name
  from_share = {
    provider = "<second_account_account_locator>" # TODO: Replace
    share = snowflake_share.test.name
  }
}
