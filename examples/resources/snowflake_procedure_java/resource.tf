# basic example
resource "snowflake_procedure_java" "basic" {
  database = "Database"
  schema   = "Schema"
  name     = "ProcedureName"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  handler              = "TestFunc.echoVarchar"
  procedure_definition = <<EOT
  import com.snowflake.snowpark_java.*;
  class TestFunc {
    public static String echoVarchar(Session session, String x) {
      return x;
    }
  }
EOT
  runtime_version      = "11"
  snowpark_package     = "1.14.0"
}

# full example
resource "snowflake_procedure_java" "full" {
  database = "Database"
  schema   = "Schema"
  name     = "ProcedureName"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  handler              = "TestFunc.echoVarchar"
  procedure_definition = <<EOT
    import com.snowflake.snowpark_java.*;
  class TestFunc {
    public static String echoVarchar(Session session, String x) {
      return x;
    }
  }
EOT
  runtime_version      = "11"
  snowpark_package     = "1.14.0"

  comment    = "some comment"
  execute_as = "CALLER"
  target_path {
    path_on_stage  = "tf-1734028493-OkoTf.jar"
    stage_location = snowflake_stage.example.fully_qualified_name
  }
  packages = ["com.snowflake:telemetry:0.1.0"]
  imports {
    path_on_stage  = "tf-1734028486-OLJpF.jar"
    stage_location = "~"
  }
  imports {
    path_on_stage  = "tf-1734028491-EMoDC.jar"
    stage_location = "~"
  }
  is_secure           = "false"
  null_input_behavior = "CALLED ON NULL INPUT"
  external_access_integrations = [
    "INTEGRATION_1", "INTEGRATION_2"
  ]
  secrets {
    secret_id            = snowflake_secret_with_generic_string.example1.fully_qualified_name
    secret_variable_name = "abc"
  }
  secrets {
    secret_id            = snowflake_secret_with_generic_string.example2.fully_qualified_name
    secret_variable_name = "def"
  }
}
