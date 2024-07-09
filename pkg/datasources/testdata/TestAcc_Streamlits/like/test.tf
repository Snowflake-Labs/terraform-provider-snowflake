resource "snowflake_streamlit" "test_1" {
  name      = var.name_1
  database  = var.database
  schema    = var.schema
  stage     = var.stage
  main_file = var.main_file
}

resource "snowflake_streamlit" "test_2" {
  name      = var.name_2
  database  = var.database
  schema    = var.schema
  stage     = var.stage
  main_file = var.main_file
}

resource "snowflake_streamlit" "test_3" {
  name      = var.name_3
  database  = var.database
  schema    = var.schema
  stage     = var.stage
  main_file = var.main_file
}

data "snowflake_streamlits" "test" {
  depends_on = [snowflake_streamlit.test_1, snowflake_streamlit.test_2, snowflake_streamlit.test_3]

  like = var.like
}
