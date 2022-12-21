resource "snowflake_account_parameter" "p" {
  key   = "ALLOW_ID_TOKEN"
  value = "true"
}

resource "snowflake_account_parameter" "p2" {
  key   = "CLIENT_ENCRYPTION_KEY_SIZE"
  value = "256"
}
