resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_share" "test" {
  name = var.share_name
}

resource "snowflake_grant_privileges_to_share" "test" {
  share_name    = snowflake_share.test.name
  privileges    = var.privileges
  database_name = snowflake_database.test.name
}
