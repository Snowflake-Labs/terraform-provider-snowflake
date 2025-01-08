# Minimal
resource "snowflake_function_scala" "minimal" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "my_scala_function"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type         = "VARCHAR(100)"
  runtime_version     = "2.12"
  handler             = "TestFunc.echoVarchar"
  function_definition = <<EOT
  class TestFunc {
    def echoVarchar(x : String): String = {
      return x
    }
  }
EOT
}

# Complete
resource "snowflake_function_scala" "complete" {
  database  = snowflake_database.test.name
  schema    = snowflake_schema.test.name
  name      = "my_scala_function"
  is_secure = "false"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  comment                      = "some comment"
  external_access_integrations = ["external_access_integration_name", "external_access_integration_name_2"]
  function_definition          = <<EOT
  class TestFunc {
    def echoVarchar(x : String): String = {
      return x
    }
  }
EOT
  handler                      = "TestFunc.echoVarchar"
  null_input_behavior          = "CALLED ON NULL INPUT"
  return_results_behavior      = "VOLATILE"
  return_type                  = "VARCHAR(100)"
  imports {
    path_on_stage  = "jar_name.jar"
    stage_location = "~"
  }
  imports {
    path_on_stage  = "second_jar_name.jar"
    stage_location = "~"
  }
  packages        = ["com.snowflake:snowpark:1.14.0", "com.snowflake:telemetry:0.1.0"]
  runtime_version = "2.12"
  secrets {
    secret_id            = snowflake_secret_with_authorization_code_grant.one.fully_qualified_name
    secret_variable_name = "abc"
  }
  secrets {
    secret_id            = snowflake_secret_with_authorization_code_grant.two.fully_qualified_name
    secret_variable_name = "def"
  }
  target_path {
    path_on_stage  = "target_jar_name.jar"
    stage_location = snowflake_stage.test.fully_qualified_name
  }
}
