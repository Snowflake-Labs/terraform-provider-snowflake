variable "database_role_name" {
  type = string
}

variable "share_name" {
  type = string
}

variable "database" {
  type = string
}

resource "snowflake_database_role" "database_role" {
  database = var.database
  name     = var.database_role_name
}

resource "snowflake_share" "share" {
  name = var.share_name
}

// todo: add grant_privileges_to_share resource

resource "snowflake_grant_database_role" "g" {
  database_role_name = snowflake_database_role.database_role.name
  share_name         = snowflake_share.name
}
