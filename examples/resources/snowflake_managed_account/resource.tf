resource snowflake_managed_account account {
  name           = "managed account"
  admin_name     = "admin"
  admin_password = "secret"
  type           = "READER"
  comment        = "A managed account."
  cloud          = "aws"
  region         = "us-west-2"
  locator        = "managed-account"
}
