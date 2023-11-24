resource "snowflake_unsafe_migration" "migration" {
  up   = var.up
  down = var.down
}
