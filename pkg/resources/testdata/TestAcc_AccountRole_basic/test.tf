resource "snowflake_account_role" "test" {
  name    = var.name
  comment = var.comment
}
