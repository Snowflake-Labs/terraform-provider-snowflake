resource "snowflake_share" "test" {
  name     = "share_name"
  comment  = "cool comment"
  accounts = ["organizationName.accountName"]
}
