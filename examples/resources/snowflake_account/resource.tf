## Minimal
resource "snowflake_account" "minimal" {
  name                 = "ACCOUNT_NAME"
  admin_name           = "ADMIN_NAME"
  admin_password       = "ADMIN_PASSWORD"
  email                = "admin@email.com"
  edition              = "STANDARD"
  grace_period_in_days = 3
}

## Complete (with SERVICE user type)
resource "snowflake_account" "complete" {
  name                 = "ACCOUNT_NAME"
  admin_name           = "ADMIN_NAME"
  admin_rsa_public_key = "<public_key>"
  admin_user_type      = "SERVICE"
  email                = "admin@email.com"
  edition              = "STANDARD"
  region_group         = "PUBLIC"
  region               = "AWS_US_WEST_2"
  comment              = "some comment"
  is_org_admin         = "true"
  grace_period_in_days = 3
}

## Complete (with PERSON user type)
resource "snowflake_account" "complete" {
  name                 = "ACCOUNT_NAME"
  admin_name           = "ADMIN_NAME"
  admin_password       = "ADMIN_PASSWORD"
  admin_user_type      = "PERSON"
  first_name           = "first_name"
  last_name            = "last_name"
  email                = "admin@email.com"
  must_change_password = "false"
  edition              = "STANDARD"
  region_group         = "PUBLIC"
  region               = "AWS_US_WEST_2"
  comment              = "some comment"
  is_org_admin         = "true"
  grace_period_in_days = 3
}
