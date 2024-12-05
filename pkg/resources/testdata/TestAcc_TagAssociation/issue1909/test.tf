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
    name = var.column_name
    type = "VARIANT"
  }
}
resource "snowflake_table" "test2" {
  name     = var.table_name2
  database = var.database
  schema   = var.schema
  column {
    name = var.column_name
    type = "VARIANT"
  }
}

resource "snowflake_tag_association" "test" {
  object_identifiers = [var.column_fully_qualified_name, var.column2_fully_qualified_name]
  object_type        = "COLUMN"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = "v1"
  depends_on         = [snowflake_table.test, snowflake_table.test2]
}
