resource "snowflake_alert" "alert" {
  database  = "database"
  schema    = "schema"
  name      = "alert"
  warehouse = "warehouse"
  alert_schedule {
    interval = 10
  }
  condition = "select 1 as c"
  action    = "select 1 as c"
  enabled   = true
  comment   = "my alert"
}
