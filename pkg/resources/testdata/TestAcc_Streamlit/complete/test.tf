resource "snowflake_streamlit" "test" {
  schema                       = var.schema
  database                     = var.database
  name                         = var.name
  stage                        = var.stage
  directory_location           = var.directory_location
  main_file                    = var.main_file
  query_warehouse              = var.query_warehouse
  external_access_integrations = var.external_access_integrations
  title                        = var.title
  comment                      = var.comment
}
