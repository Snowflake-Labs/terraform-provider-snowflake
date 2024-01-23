resource "snowflake_storage_integration" "test" {
  name                      = var.name
  storage_allowed_locations = var.allowed_locations
  storage_provider          = "S3"
  storage_aws_role_arn      = "arn:aws:iam::000000000001:/role/test"
}
