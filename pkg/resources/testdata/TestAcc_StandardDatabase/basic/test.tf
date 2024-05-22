resource "snowflake_standard_database" "test" {
  name          = var.name
  comment       = var.comment
}
