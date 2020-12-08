resource snowflake_integratin_grant grant {
  integration_name = "integration"

  privilege = "USAGE"
  roles     = ["role1", "role2"]

  with_grant_option = false
}
