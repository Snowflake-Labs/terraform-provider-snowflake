resource "snowflake_external_volume" "complete" {
  name         = var.name
  comment      = var.comment
  allow_writes = var.allow_writes
  dynamic "storage_location" {
    for_each = var.s3_storage_locations
    content {
      storage_location_name = storage_location.value["storage_location_name"]
      storage_provider      = storage_location.value["storage_provider"]
      storage_base_url      = storage_location.value["storage_base_url"]
      storage_aws_role_arn  = storage_location.value["storage_aws_role_arn"]
    }
  }
  dynamic "storage_location" {
    for_each = var.gcs_storage_locations
    content {
      storage_location_name = storage_location.value["storage_location_name"]
      storage_provider      = storage_location.value["storage_provider"]
      storage_base_url      = storage_location.value["storage_base_url"]
    }
  }
  dynamic "storage_location" {
    for_each = var.azure_storage_locations
    content {
      storage_location_name = storage_location.value["storage_location_name"]
      storage_provider      = storage_location.value["storage_provider"]
      storage_base_url      = storage_location.value["storage_base_url"]
      azure_tenant_id       = storage_location.value["azure_tenant_id"]
    }
  }
}
