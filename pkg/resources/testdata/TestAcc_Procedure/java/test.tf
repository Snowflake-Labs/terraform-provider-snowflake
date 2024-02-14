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
  language        = "JAVA"
  return_type     = "TABLE (ID NUMBER, NAME VARCHAR, ROLE VARCHAR)"
  runtime_version = "11"
  packages        = ["com.snowflake:snowpark:1.9.0"]
  handler         = "Filter.filterByRole"
  execute_as      = "CALLER"
  comment         = var.comment
  statement       = <<EOT
    import com.snowflake.snowpark_java.*;
    public class Filter {
      public DataFrame filterByRole(Session session, String tableName, String role) {
        DataFrame table = session.table(tableName);
        DataFrame filteredRows = table.filter(Functions.col("role").equal_to(Functions.lit(role)));
        return filteredRows;
      }
    }
  EOT
}
