resource "snowflake_tag" "test" {
  name             = var.name
  database         = var.database
  schema           = var.schema
  comment          = var.comment
  allowed_values   = var.allowed_values
  masking_policies = var.masking_policies
}
