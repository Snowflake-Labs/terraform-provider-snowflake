resource "snowflake_unsafe_execute" "test" {
  execute = var.execute
  revert  = var.revert
  read  = var.read
}
