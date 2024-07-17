# basic resource
resource "snowflake_streamlit" "streamlit" {
  database  = "database"
  schema    = "schema"
  name      = "streamlit"
  stage     = "streamlit_db.streamlit_schema.streamlit_stage"
  main_file = "/streamlit_main.py"
}
# resource with all fields set
resource "snowflake_streamlit" "streamlit" {
  database                     = "database"
  schema                       = "schema"
  name                         = "streamlit"
  stage                        = "streamlit_db.streamlit_schema.streamlit_stage"
  directory_location           = "src"
  main_file                    = "streamlit_main.py"
  query_warehouse              = "warehouse"
  external_access_integrations = ["integration_id"]
  title                        = "title"
  comment                      = "comment"
}
