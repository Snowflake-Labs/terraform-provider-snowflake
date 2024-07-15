##################################
### grant database role to account role
##################################

resource "snowflake_database_role" "database_role" {
  database = var.database
  name     = var.database_role_name
}

resource "snowflake_account_role" "parent_role" {
  name = var.parent_role_name
}

resource "snowflake_grant_database_role" "g" {
  database_role_name = "\"${var.database}\".\"${snowflake_database_role.database_role.name}\""
  parent_role_name   = snowflake_account_role.parent_role.name
}

##################################
### grant database role to database role
##################################

resource "snowflake_database_role" "database_role" {
  database = var.database
  name     = var.database_role_name
}

resource "snowflake_database_role" "parent_database_role" {
  database = var.database
  name     = var.parent_database_role_name
}

resource "snowflake_grant_database_role" "g" {
  database_role_name        = "\"${var.database}\".\"${snowflake_database_role.database_role.name}\""
  parent_database_role_name = "\"${var.database}\".\"${snowflake_database_role.parent_database_role.name}\""
}

##################################
### grant database role to share
##################################

resource "snowflake_grant_database_role" "g" {
  database_role_name = "\"${var.database}\".\"${snowflake_database_role.database_role.name}\""
  share_name         = snowflake_share.share.name
}
