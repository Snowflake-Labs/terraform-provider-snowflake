resource "snowflake_tag" "test" {
  name     = var.tag_name
  database = var.database
  schema   = var.schema
}

resource "snowflake_table" "test" {
  name     = var.table_name
  database = var.database
  schema   = var.schema
  column {
    name = "test_column"
    type = "VARIANT"
  }
}

resource "snowflake_tag_association" "test" {
  object_name = "\"${var.database}\".\"${var.schema}\".\"${var.table_name}\""
  object_type = "TABLE"
  tag_id      = snowflake_tag.test.id
  tag_value   = "v1"
}
