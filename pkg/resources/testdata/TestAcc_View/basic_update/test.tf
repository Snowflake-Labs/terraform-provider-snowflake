resource "snowflake_view" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema

  dynamic "column" {
    for_each = var.columns
    content {
      column_name = column.value["column_name"]
    }
  }

  row_access_policy {
    policy_name = var.row_access_policy
    on          = var.row_access_policy_on
  }
  aggregation_policy {
    policy_name = var.aggregation_policy
    entity_key  = var.aggregation_policy_entity_key
  }
  data_metric_function {
    function_name   = var.data_metric_function
    on              = var.data_metric_function_on
    schedule_status = var.schedule_status
  }
  data_metric_schedule {
    using_cron = var.data_metric_schedule_using_cron
  }
  statement = var.statement
  comment   = var.comment
}
