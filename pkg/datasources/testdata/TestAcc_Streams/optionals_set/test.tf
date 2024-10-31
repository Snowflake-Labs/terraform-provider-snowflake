resource "snowflake_stream_on_table" "test" {
  name     = var.name
  schema   = var.schema
  database = var.database

  table       = var.table
  append_only = var.append_only

  comment = var.comment
}

data "snowflake_streams" "test" {
  depends_on = [snowflake_stream_on_table.test]

  like = var.name
}
