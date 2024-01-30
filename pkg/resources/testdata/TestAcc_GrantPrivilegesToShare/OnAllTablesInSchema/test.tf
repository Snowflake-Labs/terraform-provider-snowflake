resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_share" "test" {
  name = var.share_name
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  share_name    = snowflake_share.test.name
  privileges    = ["USAGE"]
  database_name = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_share" "test" {
  share_name           = snowflake_share.test.name
  privileges           = var.privileges
  all_tables_in_schema = "\"${snowflake_schema.test.database}\".\"${snowflake_schema.test.name}\""
}
