resource "snowflake_role" "test" {
  name    = var.name
  comment = var.comment
}
