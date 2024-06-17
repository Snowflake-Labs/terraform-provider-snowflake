resource "snowflake_database" "test_1" {
  name = var.name_1
}

resource "snowflake_database" "test_2" {
  name = var.name_2
}

resource "snowflake_database" "test_3" {
  name = var.name_3
}

data "snowflake_databases" "test" {
  depends_on  = [snowflake_database.test_1, snowflake_database.test_2, snowflake_database.test_3]
  starts_with = var.starts_with
}
