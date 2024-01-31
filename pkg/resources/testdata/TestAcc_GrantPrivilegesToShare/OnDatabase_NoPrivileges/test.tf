resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_share" "test" {
  name = var.to_share
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = snowflake_share.test.name
  on_database = snowflake_database.test.name
}
