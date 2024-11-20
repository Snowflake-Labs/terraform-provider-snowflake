resource "snowflake_tag" "test" {
  name     = var.tag_name
  database = var.database
  schema   = var.schema
}
resource "snowflake_table" "test" {
  name     = var.table_name
  database = var.database
  schema   = var.schema
  // TODO(SNOW-1348114): use only one column, if possible.
  // We need a dummy column here because a table must have at least one column, and when we rename the second one in the config, it gets dropped for a moment.
  column {
    name = "DUMMY"
    type = "VARIANT"
  }
  column {
    name = var.column_name
    type = "VARIANT"
  }
}
resource "snowflake_tag_association" "test" {
  object_identifiers = [var.column_fully_qualified_name]
  object_type        = "COLUMN"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = "TAG_VALUE"
  depends_on         = [snowflake_table.test]
}
