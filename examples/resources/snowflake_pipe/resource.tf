resource snowflake_pipe pipe {
  database = "db"
  schema   = "schema"
  name     = "pipe"

  comment = "A pipe."

  copy_statement = "copy into mytable from @mystage"
  auto_ingest    = false

  aws_sns_topic_arn    = "..."
  notification_channel = "..."
  owner                = "role1"
}
