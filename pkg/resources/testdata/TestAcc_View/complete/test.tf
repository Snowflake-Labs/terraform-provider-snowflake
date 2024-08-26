resource "snowflake_view" "test" {
  name            = var.name
  comment         = var.comment
  database        = var.database
  schema          = var.schema
  is_secure       = var.is_secure
  or_replace      = var.or_replace
  copy_grants     = var.copy_grants
  change_tracking = var.change_tracking
  is_temporary    = var.is_temporary
  row_access_policy {
    policy_name = var.row_access_policy
    on          = var.row_access_policy_on

  }
  aggregation_policy {
    policy_name = var.aggregation_policy
    entity_key  = var.aggregation_policy_entity_key
  }
  statement  = var.statement
  depends_on = [snowflake_unsafe_execute.use_warehouse]
}

resource "snowflake_unsafe_execute" "use_warehouse" {
  execute = "USE WAREHOUSE \"${var.warehouse}\""
  revert  = "SELECT 1"
}
