##################################
### grant account role to account role
##################################

resource "snowflake_account_role" "role" {
  name = "ROLE"
}

resource "snowflake_account_role" "parent_role" {
  name = "PARENT_ROLE"
}

resource "snowflake_grant_account_role" "g" {
  role_name        = snowflake_account_role.role.name
  parent_role_name = snowflake_account_role.parent_role.name
}


##################################
### grant account role to user
##################################

resource "snowflake_account_role" "role" {
  name = "ROLE"
}

resource "snowflake_user" "user" {
  name = "USER"
}

resource "snowflake_grant_account_role" "g" {
  role_name = snowflake_account_role.role.name
  user_name = snowflake_user.user.name
}
