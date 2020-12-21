resource snowflake_warehouse_grant grant {
  warehouse_name = "wh"
  privilege      = "MODIFY"

  roles = [
    "role1",
  ]

  with_grant_option = false
}
