resource snowflake_schema schema {
  database = "db"
  name     = "schema"
  comment  = "A schema."

  is_transient        = false
  is_managed          = false
  data_retention_days = 1
}
