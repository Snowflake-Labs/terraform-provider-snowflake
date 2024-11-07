resource "snowflake_task" "test" {
  count = length(var.tasks)
  name          = var.tasks[count.index].name
  database      = var.tasks[count.index].database
  schema        = var.tasks[count.index].schema
  started       = var.tasks[count.index].started
  sql_statement = var.tasks[count.index].sql_statement

  # Optionals
  dynamic "schedule" {
    for_each = [for element in [var.tasks[count.index].schedule] : element if element != null]
    content {
      minutes    = lookup(var.tasks[count.index].schedule, "minutes", null)
      using_cron = lookup(var.tasks[count.index].schedule, "cron", null)
    }
  }

  comment = var.tasks[count.index].comment
  after = var.tasks[count.index].after
  finalize = var.tasks[count.index].finalize

  # Parameters
  suspend_task_after_num_failures               = var.tasks[count.index].suspend_task_after_num_failures
}
