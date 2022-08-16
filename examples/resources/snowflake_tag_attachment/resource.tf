resource "snowflake_tag_attachment" "test_tag_attachment" {
  // Required
  resourceId = "test_user_name"
  objectType = "USER"
  tagName    = "test_tag_name"
  tagValue   = "test_user_value"
}
