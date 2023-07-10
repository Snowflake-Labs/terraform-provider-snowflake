resource "snowflake_email_notification_integration" "email_int" {
  name    = "notification"
  comment = "A notification integration."

  enabled            = true
  allowed_recipients = ["john.doe@gmail.com"]
}
