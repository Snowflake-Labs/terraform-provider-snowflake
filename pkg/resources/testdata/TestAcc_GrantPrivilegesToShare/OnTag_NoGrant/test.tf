resource "snowflake_share" "test" {
  name = var.to_share
}

resource "snowflake_database" "test" {
  depends_on = [snowflake_share.test]
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_tag" "test" {
  name     = var.on_tag
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
}
