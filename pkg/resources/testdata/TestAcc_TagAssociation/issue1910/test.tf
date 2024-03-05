resource "snowflake_tag" "test" {
  name     = var.tag_name
  database = var.database
  schema   = var.schema
}

resource "snowflake_account" "test" {
  name                 = var.account_name
  admin_name           = "someadmin"
  admin_password       = "123456"
  first_name           = "Ad"
  last_name            = "Min"
  email                = "admin@example.com"
  must_change_password = false
  edition              = "BUSINESS_CRITICAL"
  grace_period_in_days = 4
}

resource "snowflake_tag_association" "test" {
  object_identifier {
    name = snowflake_account.test.name
  }
  object_type = "ACCOUNT"
  tag_id      = snowflake_tag.test.id
  tag_value   = "v1"
}
