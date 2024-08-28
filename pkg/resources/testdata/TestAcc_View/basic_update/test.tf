resource "snowflake_view" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  row_access_policy {
    policy_name = var.row_access_policy
    on          = var.row_access_policy_on

  }
  aggregation_policy {
    policy_name = var.aggregation_policy
    entity_key  = var.aggregation_policy_entity_key
  }
  data_metric_functions {
    function_name = var.data_metric_function
    on            = var.data_metric_function_on
  }
  data_metric_schedule {
    using_cron = var.data_metric_schedule_using_cron
  }
  statement = var.statement
  comment   = var.comment
}
