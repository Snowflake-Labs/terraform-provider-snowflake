resource "snowflake_row_access_policy" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.arguments
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body = "case when current_role() in ('ANALYST') then true else false end"
}
