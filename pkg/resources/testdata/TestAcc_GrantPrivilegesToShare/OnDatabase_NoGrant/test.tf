resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_share" "test" {
  name = var.to_share
}
