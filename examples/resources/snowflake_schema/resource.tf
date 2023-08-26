resource "snowflake_schema" "schema" {
  database = "database"
  name     = "schema"
  comment  = "A schema."

  is_transient        = false
  is_managed          = false
  data_retention_days = 1
}
