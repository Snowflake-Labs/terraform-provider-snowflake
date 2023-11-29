resource "snowflake_unsafe_execute" "test" {
  execute = var.execute
  revert  = var.revert
  query   = var.query
}
