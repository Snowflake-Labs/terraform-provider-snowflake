resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_share" "test" {
  name = var.share_name
}

resource "snowflake_grant_privileges_to_share" "test" {
  depends_on = [snowflake_share.test]
  share_name = var.share_account_name
  privileges = var.privileges
  database_name = snowflake_database.test.name
}
