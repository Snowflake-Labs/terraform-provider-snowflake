##################################
### grant account role to account role
##################################

resource "snowflake_account_role" "role" {
  name = var.role_name
}

resource "snowflake_account_role" "parent_role" {
  name = var.parent_role_name
}

resource "snowflake_grant_account_role" "g" {
  role_name        = snowflake_account_role.role.name
  parent_role_name = snowflake_account_role.parent_role.name
}


##################################
### grant account role to user
##################################

resource "snowflake_account_role" "role" {
  name = var.role_name
}

resource "snowflake_user" "user" {
  name = var.user_name
}

resource "snowflake_grant_account_role" "g" {
  role_name = snowflake_role.role.name
  user_name = snowflake_user.user.name
}
