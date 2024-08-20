resource "snowflake_schema" "test_1" {
  name     = var.name_1
  database = var.database_1
}

resource "snowflake_schema" "test_2" {
  name     = var.name_2
  database = var.database_2
}

resource "snowflake_schema" "test_3" {
  name     = var.name_3
  database = var.database_2
}

data "snowflake_schemas" "test" {
  depends_on = [snowflake_schema.test_1, snowflake_schema.test_2, snowflake_schema.test_3]
  in {
    database = var.in
  }
  starts_with = var.starts_with
}
