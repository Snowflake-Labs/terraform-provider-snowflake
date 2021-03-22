resource snowflake_account_grant grant {
  roles             = ["role1", "role2"]
  privilege         = "CREATE ROLE"
  with_grant_option = false
}
