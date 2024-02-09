resource "snowflake_procedure" "p" {
  database = var.database
  schema   = var.schema
  name     = var.name
  arguments {
    name = "TABLE_NAME"
    type = "VARCHAR"
  }
  arguments {
    name = "ROLE"
    type = "VARCHAR"
  }
  language        = "PYTHON"
  return_type     = "TABLE (ID NUMBER, NAME VARCHAR, ROLE VARCHAR)"
  runtime_version = "3.8"
  packages        = ["snowflake-snowpark-python"]
  handler         = "filter_by_role"
  execute_as      = "CALLER"
  comment         = var.comment
  statement       = <<EOT
from snowflake.snowpark.functions import col
def filter_by_role(session, table_name, role):
  df = session.table(table_name)
  return df.filter(col("role") == role)
  EOT
}
