resource "snowflake_stream_on_table" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema

  copy_grants       = true
  table             = var.table
  append_only       = "true"
  show_initial_rows = "true"

  before {
    timestamp = try(var.before["timestamp"], null)
    offset    = try(var.before["offset"], null)
    stream    = try(var.before["stream"], null)
    statement = try(var.before["statement"], null)
  }

  comment = var.comment
}
