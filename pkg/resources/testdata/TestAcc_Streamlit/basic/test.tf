resource "snowflake_streamlit" "test" {
  database      = var.database
  schema        = var.schema
  name          = var.name
  root_location = var.root_location
  main_file     = var.main_file
}
