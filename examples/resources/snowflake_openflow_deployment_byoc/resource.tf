# basic resource
resource "snowflake_openflow_deployment_byoc" "basic" {
  name     = "my_byoc_deployment"
  vpc_type = "MANAGED"
}

# in a VPC you provide yourself
resource "snowflake_openflow_deployment_byoc" "provided_vpc" {
  name     = "my_byoc_deployment"
  vpc_type = "PROVIDED"
}

# complete resource
resource "snowflake_openflow_deployment_byoc" "complete" {
  name         = "my_byoc_deployment_complete"
  vpc_type     = "MANAGED"
  display_name = "My BYOC deployment"
  # an existing event table, created with CREATE EVENT TABLE outside Terraform: the provider has no
  # resource for one, and a regular table will not do because Snowflake requires the event schema
  event_table                    = "\"<database_name>\".\"<schema_name>\".\"<event_table_name>\""
  custom_ingress_hostname        = "openflow.example.com"
  use_private_link               = "true"
  use_user_auth_over_privatelink = "true"
  comment                        = "Managed by Terraform."

  timeouts {
    create = "60m"
    delete = "60m"
  }
}
