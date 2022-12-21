resource "snowflake_session_parameter" "s" {
  key   = "AUTOCOMMIT"
  value = "false"
}
