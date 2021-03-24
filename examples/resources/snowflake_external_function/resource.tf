resource "snowflake_external_function" "test_ext_func" {
  name = "my_function"
  database = "my_test_db"
  schema   = "my_test_schema"
  arg {
    name = "arg1"
    type = "varchar"
  }
  arg {
    name = "arg2"
    type = "varchar"
  }
  return_type = "varchar"
  return_behavior = "IMMUTABLE"
  api_integration = "api_integration_name"
  url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}