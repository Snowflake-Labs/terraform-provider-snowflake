resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_tag" "test" {
  name = var.tag_name
  database = snowflake_database.test.name
  schema = snowflake_schema.test.name
}

resource "snowflake_share" "test" {
  name = var.share_name
}
