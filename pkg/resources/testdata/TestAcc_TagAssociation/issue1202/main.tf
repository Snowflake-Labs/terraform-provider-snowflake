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
  object_identifiers = [snowflake_table.test.fully_qualified_name]
  object_type        = "TABLE"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = "v1"
}
