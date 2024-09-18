resource "snowflake_tag" "test" {
  name           = var.name
  database       = var.database
  schema         = var.schema
  comment        = var.comment
  allowed_values = ["alv1", "alv2"]
}

resource "snowflake_masking_policy" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  argument {
    name = "val"
    type = "VARCHAR"
  }
  body             = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
  return_data_type = "VARCHAR(16777216)"
  comment          = "Terraform acceptance test"
}

resource "snowflake_tag_masking_policy_association" "test" {
  tag_id            = "${snowflake_tag.test.database}.${snowflake_tag.test.schema}.${snowflake_tag.test.name}"
  masking_policy_id = "${snowflake_masking_policy.test.database}.${snowflake_masking_policy.test.schema}.${snowflake_masking_policy.test.name}"
}
