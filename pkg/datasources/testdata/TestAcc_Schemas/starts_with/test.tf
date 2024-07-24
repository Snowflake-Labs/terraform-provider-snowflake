resource "snowflake_schema" "test_1" {
  name     = var.name_1
  database = var.database
}

resource "snowflake_schema" "test_2" {
  name     = var.name_2
  database = var.database
}

resource "snowflake_schema" "test_3" {
  name     = var.name_3
  database = var.database
}

data "snowflake_schemas" "test" {
  depends_on  = [snowflake_schema.test_1, snowflake_schema.test_2, snowflake_schema.test_3]
  starts_with = var.starts_with
}
