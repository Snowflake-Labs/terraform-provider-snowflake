resource "snowflake_view" "test" {
  name            = var.name
  comment         = var.comment
  database        = var.database
  schema          = var.schema
  is_secure       = var.is_secure
  copy_grants     = var.copy_grants
  change_tracking = var.change_tracking
  is_temporary    = var.is_temporary
  data_metric_functions {
    function_name = var.data_metric_function
    on            = var.data_metric_function_on
  }
  data_metric_schedule {
    using_cron = var.data_metric_schedule_using_cron
  }
  row_access_policy {
    policy_name = var.row_access_policy
    on          = var.row_access_policy_on
  }
  aggregation_policy {
    policy_name = var.aggregation_policy
    entity_key  = var.aggregation_policy_entity_key
  }
  statement = var.statement
}
