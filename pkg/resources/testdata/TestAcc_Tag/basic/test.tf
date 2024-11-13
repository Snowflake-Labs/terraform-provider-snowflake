resource "snowflake_tag" "t" {
  name     = var.name
  database = var.database
  schema   = var.schema
}
