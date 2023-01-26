provider "snowflake" {
  role  = "ORGADMIN"
  alias = "orgadmin"
}

resource "snowflake_account" "ac1" {
  provider             = snowflake.orgadmin
  name                 = "SNOWFLAKE_TEST_ACCOUNT"
  admin_name           = "John Doe"
  admin_password       = "Abcd1234!"
  email                = "john.doe@snowflake.com"
  first_name           = "John"
  last_name            = "Doe"
  must_change_password = true
  edition              = "STANDARD"
  comment              = "Snowflake Test Account"
  region               = "AWS_US_WEST_2"
}
