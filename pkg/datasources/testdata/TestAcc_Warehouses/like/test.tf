resource "snowflake_warehouse" "test_1" {
  name = var.name_1
}

resource "snowflake_warehouse" "test_2" {
  name = var.name_2
}

resource "snowflake_warehouse" "test_3" {
  name = var.name_3
}

data "snowflake_warehouses" "test" {
  depends_on = [snowflake_warehouse.test_1, snowflake_warehouse.test_2, snowflake_warehouse.test_3]
  like       = var.like
}
