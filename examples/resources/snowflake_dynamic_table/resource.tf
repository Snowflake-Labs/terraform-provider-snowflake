# https://docs.snowflake.com/en/sql-reference/sql/create-dynamic-table#examples
resource "snowflake_dynamic_table" "dt" {
  name     = "product"
  database = "mydb"
  schema   = "myschema"
  target_lag {
    maximum_duration = "20 minutes"
  }
  warehouse = "mywh"
  query     = "SELECT product_id, product_name FROM \"mydb\".\"myschema\".\"staging_table\""
  comment   = "example comment"
}
