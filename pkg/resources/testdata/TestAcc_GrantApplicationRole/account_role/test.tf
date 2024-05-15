locals {
  application_role_identifier = "\"${var.application_name}\".\"app_role_1\""
}

resource "snowflake_role" "role" {
  name = var.parent_account_role_name
}

resource "snowflake_grant_application_role" "g" {
  application_role_name    = local.application_role_identifier
  parent_account_role_name = "\"${snowflake_role.role.name}\""
}
