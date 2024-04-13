locals {
  application_role_identifier = "\"my_appplication\".\"app_role_1\""
}

##################################
### grant application role to account role
##################################


resource "snowflake_role" "role" {
  name = "my_role"
}

resource "snowflake_grant_application_role" "g" {
  application_role_name    = local.application_role_identifier
  parent_account_role_name = snowflake_role.role.name
}

##################################
### grant application role to application
##################################

resource "snowflake_grant_application_role" "g" {
  application_role_name = local.application_role_identifier
  application_name      = "my_second_application"
}
