resource "snowflake_shared_database" "test" {
  name       = var.name
  from_share = var.from_share
  comment    = var.comment
}
