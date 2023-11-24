resource "snowflake_unsafe_execute" "migration" {
  up   = var.up
  down = var.down
}
