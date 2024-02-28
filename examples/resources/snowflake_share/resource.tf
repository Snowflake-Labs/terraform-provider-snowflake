resource "snowflake_share" "test" {
  name     = "share_name"
  comment  = "cool comment"
  accounts = ["organizationName.accountName"]
}

resource "snowflake_database" "example" {
  # remember to define dependency between objects on a share, because shared objects have to be dropped before dropping share
  depends_on = [snowflake_share.test]
  name       = "test"
}
