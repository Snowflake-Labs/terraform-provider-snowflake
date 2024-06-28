resource "snowflake_warehouse" "test" {
  name    = var.name
  comment = var.comment
}

data "snowflake_warehouses" "test" {
  depends_on = [snowflake_warehouse.test]
  like       = var.name
}
