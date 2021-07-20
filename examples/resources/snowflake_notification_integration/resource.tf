resource snowflake_notification_integration integration {
  name    = "notification"
  comment = "A notification integration."
  
  enabled   = true
  type      = "QUEUE"
  direction = "OUTBOUND"

  # AZURE_STORAGE_QUEUE
  notification_provider           = "AZURE_STORAGE_QUEUE"
  azure_storage_queue_primary_uri = "..."
  azure_tenant_id                 = "..."

  # AWS_SQS
  notification_provider = "AWS_SQS"
  aws_sqs_arn           = "..." 
  aws_sqs_role_arn      = "..."
}
