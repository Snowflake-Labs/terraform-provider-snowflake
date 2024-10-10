resource "snowflake_external_volume" "complete" {
  name         = var.name
  comment      = var.comment
  allow_writes = var.allow_writes
  dynamic "storage_location" {
    for_each = var.storage_location
    content {
      storage_location_name = storage_location.value["storage_location_name"]
      storage_provider      = storage_location.value["storage_provider"]
      storage_base_url      = storage_location.value["storage_base_url"]
      storage_aws_role_arn  = storage_location.value["storage_aws_role_arn"]
    }
  }
}
