# basic resource
resource "snowflake_tag" "tag" {
  name     = "tag"
  database = "database"
  schema   = "schema"
}

# complete resource
resource "snowflake_tag" "tag" {
  name             = "tag"
  database         = "database"
  schema           = "schema"
  comment          = "comment"
  allowed_values   = ["finance", "engineering", ""]
  masking_policies = [snowfalke_masking_policy.example.fully_qualified_name]
}
