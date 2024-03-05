resource "snowflake_share" "test" {
  depends_on = [snowflake_database.test, snowflake_schema.test]
  name       = var.to_share
}

resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

