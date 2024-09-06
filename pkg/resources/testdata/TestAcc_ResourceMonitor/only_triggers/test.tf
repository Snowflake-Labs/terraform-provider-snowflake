resource "snowflake_resource_monitor" "test" {
  name     = var.name

  dynamic "trigger" {
    for_each = var.trigger

    content {
      threshold = trigger.value.threshold
      on_threshold_reached = trigger.value.on_threshold_reached
    }
  }
}
