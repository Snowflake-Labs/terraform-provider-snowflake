resource "snowflake_view" "test" {
  name      = var.name
  database  = var.database
  schema    = var.schema
  statement = var.statement

  column {
    column_name = "ID"

    projection_policy {
      policy_name = var.projection_name
    }

    masking_policy {
      policy_name = var.masking_name
      using       = var.masking_using
    }
  }

  column {
    column_name = "FOO"
  }
}
