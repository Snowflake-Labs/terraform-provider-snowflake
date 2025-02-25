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
  name                 = "<name>"
  admin_name           = "<admin_name>"
  admin_password       = "<admin_password>"
  admin_user_type      = "SERVICE"
  email                = "<email>"
  edition              = "STANDARD"
  region               = "<region>" # if needed
  comment              = ""
}
