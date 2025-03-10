terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "1.0.3"
    }
  }
}

provider "snowflake" {
}

resource "snowflake_account" "test_account" {
  grace_period_in_days = 3
  name                 = "<name>" # TODO: Replace
  admin_name           = "<admin_name>" # TODO: Replace
  admin_password       = "<admin_password>" # TODO: Replace
  admin_user_type      = "SERVICE"
  email                = "<email>" # TODO: Replace
  edition              = "STANDARD"
  region               = "<region>" # TODO: Replace (if needed; can be filled after the import)
  comment              = ""
}
