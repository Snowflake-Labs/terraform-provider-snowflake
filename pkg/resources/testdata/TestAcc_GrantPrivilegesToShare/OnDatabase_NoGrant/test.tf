resource "snowflake_share" "test" {
  name = var.to_share
}

resource "snowflake_database" "test" {
  depends_on = [snowflake_share.test]
  name       = var.database
}
