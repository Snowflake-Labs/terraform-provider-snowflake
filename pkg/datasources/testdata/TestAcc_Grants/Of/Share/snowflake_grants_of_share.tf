resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_share" "test" {
  depends_on = [snowflake_database.test]

  name     = var.share
  accounts = [var.account]
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = snowflake_share.test.name
  privileges  = ["USAGE"]
  on_database = snowflake_database.test.name
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_privileges_to_share.test]

  grants_of {
    share = snowflake_share.test.name
  }
}
