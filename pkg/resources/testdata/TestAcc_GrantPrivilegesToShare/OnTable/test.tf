resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_table" "test" {
  name = var.table_name
  database = snowflake_database.test.name
  schema = snowflake_schema.test.name
  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
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
  table_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_table.test.name}\""
}
