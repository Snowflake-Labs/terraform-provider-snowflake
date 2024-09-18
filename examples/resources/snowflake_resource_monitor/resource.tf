// Note: Without credit quota and triggers specified in the configuration, the resource monitor is not performing any work.
// More on resource monitor usage: https://docs.snowflake.com/en/user-guide/resource-monitors.
resource "snowflake_resource_monitor" "minimal" {
  name = "resource-monitor-name"
}

// Note: Resource monitors have to be attached to account or warehouse to be able to track credit usage.
resource "snowflake_resource_monitor" "minimal_working" {
  name            = "resource-monitor-name"
  credit_quota    = 100
  suspend_trigger = 100
  notify_users    = ["USERONE", "USERTWO"]
}

resource "snowflake_resource_monitor" "complete" {
  name         = "resource-monitor-name"
  credit_quota = 100

  frequency       = "DAILY"
  start_timestamp = "2030-12-07 00:00"
  end_timestamp   = "2035-12-07 00:00"

  notify_triggers           = [40, 50]
  suspend_trigger           = 50
  suspend_immediate_trigger = 90

  notify_users = ["USERONE", "USERTWO"]
}
