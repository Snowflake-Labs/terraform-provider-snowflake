resource "snowflake_warehouse_grant" "grant" {
  warehouse_name = "warehouse"
  privilege      = "MODIFY"

  roles = ["role1", "role2"]

  with_grant_option = false
}
