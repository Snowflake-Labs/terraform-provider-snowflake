resource "snowflake_resource_monitor" "test" {
  name     = var.name
  notify_users = var.notify_users
  credit_quota = var.credit_quota
  frequency = var.frequency
  start_timestamp = var.start_timestamp
  end_timestamp = var.end_timestamp

  dynamic "trigger" {
    for_each = var.trigger

    content {
      threshold = trigger.value.threshold
      on_threshold_reached = trigger.value.on_threshold_reached
    }
  }
}
