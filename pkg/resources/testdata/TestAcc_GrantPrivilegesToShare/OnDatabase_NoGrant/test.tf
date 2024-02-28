resource "snowflake_share" "test" {
  depends_on = [snowflake_database.test]
  name       = var.to_share
}

resource "snowflake_database" "test" {
  name = var.database
}
