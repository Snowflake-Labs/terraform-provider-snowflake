resource "snowflake_streamlit" "test" {
  # schema                       = "\"${var.database}\".\"${var.schema}\""
  schema                       = var.schema
  database                     = var.database
  name                         = var.name
  root_location                = "@\"${var.database}\".\"${var.schema}\".\"${var.stage}\""
  main_file                    = var.main_file
  query_warehouse              = var.query_warehouse
  external_access_integrations = var.external_access_integrations
  title                        = var.title
  comment                      = var.comment
}
