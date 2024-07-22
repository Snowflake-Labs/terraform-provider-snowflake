resource "snowflake_streamlit" "test" {
  database  = var.database
  schema    = var.schema
  stage     = var.stage
  name      = var.name
  main_file = var.main_file
}
