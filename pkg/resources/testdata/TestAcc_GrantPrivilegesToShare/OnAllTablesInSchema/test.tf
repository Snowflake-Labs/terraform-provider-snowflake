resource "snowflake_share" "test" {
  name = var.to_share
}

resource "snowflake_database" "test" {
  depends_on = [snowflake_share.test]
  name       = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = snowflake_share.test.name
  privileges  = ["USAGE"]
  on_database = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share                = snowflake_share.test.name
  privileges              = var.privileges
  on_all_tables_in_schema = "\"${snowflake_schema.test.database}\".\"${snowflake_schema.test.name}\""
}
