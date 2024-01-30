resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_share" "test" {
  name = var.share_name
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  depends_on = [snowflake_share.test]
  share_name = var.share_account_name
  privileges = ["USAGE"]
  database_name = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_share" "test" {
  depends_on = [snowflake_share.test]
  share_name = var.share_account_name
  privileges = var.privileges
  schema_name = "\"${snowflake_schema.test.database}\".\"${snowflake_schema.test.name}\""
}
