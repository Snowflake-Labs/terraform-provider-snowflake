resource snowflake_integration_grant grant {
  integration_name = "integration"

  privilege = "USAGE"
  roles     = ["role1", "role2"]

  with_grant_option = false
}
