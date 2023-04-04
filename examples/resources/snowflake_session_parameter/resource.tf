resource "snowflake_session_parameter" "s" {
  key   = "AUTOCOMMIT"
  value = "false"
  user  = "TEST_USER"
}

resource "snowflake_session_parameter" "s2" {
  key        = "BINARY_OUTPUT_FORMAT"
  value      = "BASE64"
  on_account = true
}
