resource "snowflake_procedure" "p" {
  database            = var.database
  schema              = var.schema
  name                = var.name
  arguments {
    name = "TABLE_NAME"
    type = "VARCHAR"
  }
  arguments {
    name = "ROLE"
    type = "VARCHAR"
  }
  language            = "SCALA"
  return_type         = "TABLE (ID NUMBER, NAME VARCHAR, ROLE VARCHAR)"
  runtime_version     = "2.12"
  packages            = ["com.snowflake:snowpark:1.9.0"]
  handler             = "Filter.filterByRole"
  execute_as          = "CALLER"
  comment             = var.comment
  statement           = <<EOT
		import com.snowflake.snowpark.functions._
		import com.snowflake.snowpark._
		object Filter {
			def filterByRole(session: Session, tableName: String, role: String): DataFrame = {
				val table = session.table(tableName)
				val filteredRows = table.filter(col("role") === role)
				return filteredRows
			}
		}
  EOT
}
