resource "snowflake_share" "test" {
  name       = var.to_share
  depends_on = [snowflake_database.test]
}

resource "snowflake_database" "test" {
  name = var.database
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
  on_all_tables_in_schema = snowflake_schema.test.fully_qualified_name
  depends_on              = [snowflake_grant_privileges_to_share.test_setup]
}
