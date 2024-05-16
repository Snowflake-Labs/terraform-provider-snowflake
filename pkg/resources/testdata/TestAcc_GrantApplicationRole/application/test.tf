locals {
  application_role_identifier = "\"${var.application_name}\".\"app_role_1\""
}

resource "snowflake_grant_application_role" "g" {
  application_role_name = local.application_role_identifier
  application_name      = "\"${var.application_name2}\""
}
