variable "database_role_name" {
  type = string
}

variable "share_name" {
  type = string
}

variable "database" {
  type = string
}

resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_database_role" "test" {
  database = snowflake_database.test.name
  name     = var.database_role_name
}

resource "snowflake_share" "test" {
  name = var.share_name
}

resource "snowflake_grant_privileges_to_share" "test" {
  privileges  = ["USAGE"]
  on_database = snowflake_database.test.name
  to_share    = snowflake_share.test.name
}

resource "snowflake_grant_database_role" "test" {
  database_role_name = "\"${snowflake_database.test.name}\".\"${snowflake_database_role.test.name}\""
  share_name         = snowflake_share.test.name
  depends_on         = [snowflake_grant_privileges_to_share.test]
}
