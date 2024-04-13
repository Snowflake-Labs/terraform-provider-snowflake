resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = var.to_share
  privileges  = var.privileges
  on_database = snowflake_database.test.name
}
