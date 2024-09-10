resource "snowflake_database" "test" {
  name = var.database_name
}

resource "snowflake_schema" "test" {
  database = snowflake_database.test.fully_qualified_name
  name     = var.schema_name
}
