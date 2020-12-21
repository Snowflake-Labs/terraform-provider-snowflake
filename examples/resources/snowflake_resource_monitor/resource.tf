resource snowflake_resource_monitor monitor {
  name         = "monitor"
  credit_quota = 100

  frequency       = "DAILY"
  start_timestamp = "2020-12-07 00:00"
  end_timestamp   = "2021-12-07 00:00"

  notify_triggers            = [40]
  suspend_triggers           = [50]
  suspend_immediate_triggers = [90]
}
