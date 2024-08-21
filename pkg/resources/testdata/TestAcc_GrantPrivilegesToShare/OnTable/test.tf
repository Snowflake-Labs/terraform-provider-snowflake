resource "snowflake_share" "test" {
  depends_on = [snowflake_database.test]
  name       = var.to_share
}

resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_table" "test" {
  name     = var.on_table
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = snowflake_share.test.fully_qualified_name
  privileges  = ["USAGE"]
  on_database = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = snowflake_share.test.fully_qualified_name
  privileges = var.privileges
  on_table   = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_table.test.name}\""
  depends_on = [snowflake_grant_privileges_to_share.test_setup]
}
