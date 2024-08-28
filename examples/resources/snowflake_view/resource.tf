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
# resource with attached policies
resource "snowflake_view" "test" {
  database        = "database"
  schema          = "schema"
  name            = "view"
  comment         = "comment"
  is_secure       = "true"
  change_tracking = "true"
  is_temporary    = "true"
  row_access_policy {
    policy_name = "row_access_policy"
    on          = ["id"]
  }
  aggregation_policy {
    policy_name = "aggregation_policy"
    entity_key  = ["id"]
  }
  statement = <<-SQL
    select id from foo;
SQL
}
