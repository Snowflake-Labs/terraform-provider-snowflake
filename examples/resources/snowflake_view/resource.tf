# basic resource
resource "snowflake_view" "view" {
  database  = "database"
  schema    = "schema"
  name      = "view"
  statement = <<-SQL
    select * from foo;
SQL
}

# recursive view
resource "snowflake_view" "view" {
  database     = "database"
  schema       = "schema"
  name         = "view"
  is_recursive = "true"
  statement    = <<-SQL
    select * from foo;
SQL
}
# resource with attached policies, columns and data metric functions
resource "snowflake_view" "test" {
  database        = "database"
  schema          = "schema"
  name            = "view"
  comment         = "comment"
  is_secure       = "true"
  change_tracking = "true"
  is_temporary    = "true"
  column {
    column_name = "id"
    comment     = "column comment"

  }
  column {
    column_name = "address"
    projection_policy {
      policy_name = "projection_policy"
    }

    masking_policy {
      policy_name = "masking_policy"
      using       = ["address"]
    }

  }
  row_access_policy {
    policy_name = "row_access_policy"
    on          = ["id"]
  }
  aggregation_policy {
    policy_name = "aggregation_policy"
    entity_key  = ["id"]
  }
  data_metric_function {
    function_name = "data_metric_function"
    on            = ["id"]
  }
  data_metric_schedule {
    using_cron = "15 * * * * UTC"
  }
  statement = <<-SQL
    SELECT id, address FROM TABLE;
SQL
}
