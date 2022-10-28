resource "snowflake_user_grant" "grant" {
  user_name = "user"
  privilege = "MONITOR"

  roles = ["role1", "role2"]

  with_grant_option = false
}
