resource "snowflake_role" "test1" {
  name    = var.account_role_name_1
  comment = var.comment
}

resource "snowflake_role" "test2" {
  name    = var.account_role_name_2
  comment = var.comment
}

resource "snowflake_role" "test3" {
  name    = var.account_role_name_3
  comment = var.comment
}

data "snowflake_roles" "test" {
  depends_on = [
    snowflake_role.test1,
    snowflake_role.test2,
    snowflake_role.test3,
  ]
  pattern = var.pattern
}
