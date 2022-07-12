resource snowflake_monitor_grant grant {
  monitor_name      = "monitor"
  privilege         = "MODIFY"
  roles             = ["role1"]
  with_grant_option = false
}
