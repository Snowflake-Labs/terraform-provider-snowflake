resource "snowflake_unsafe_execute" "migration" {
  execute = var.execute
  revert  = var.revert
}
