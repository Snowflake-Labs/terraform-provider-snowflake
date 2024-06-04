// Provider configuration
provider "snowflake" {
  region    = "REGION" // Default is "us-west-2"
  username  = "USERNAME"
  account   = "ACCOUNT"
  password  = "PASSWORD"
  role      = "MY_ROLE"
  warehouse = "MY_WH" // Optional attribute, some resources (e.g. Python UDFs)' require a warehouse to create and can also be set optionally from the `SNOWFLAKE_WAREHOUSE` environment variable
}

// Create database
resource "snowflake_standard_database" "db" {
  name                = "MY_DB"
  data_retention_days = 1
}

// Create schema
resource "snowflake_schema" "schema" {
  database            = snowflake_standard_database.db.name
  name                = "MY_SCHEMA"
  data_retention_days = 1
}

// Example for Java language
resource "snowflake_function" "test_funct_java" {
  name     = "my_java_func"
  database = "MY_DB"
  schema   = "MY_SCHEMA"
  arguments {
    name = "arg1"
    type = "number"
  }
  comment     = "Example for java language"
  return_type = "varchar"
  language    = "java"
  handler     = "CoolFunc.test"
  statement   = "class CoolFunc {public static String test(int n) {return \"hello!\";}}"
}

// Example for Python language
resource "snowflake_function" "python_test" {
  name     = "MY_PYTHON_FUNC"
  database = "MY_DB"
  schema   = "MY_SCHEMA"
  arguments {
    name = "arg1"
    type = "number"
  }
  comment             = "Example for Python language"
  return_type         = "NUMBER(38,0)"
  null_input_behavior = "CALLED ON NULL INPUT"
  return_behavior     = "VOLATILE"
  language            = "python"
  runtime_version     = "3.8"
  handler             = "add_py"
  statement           = "def add_py(i): return i+1"
}

// Example SQL language
resource "snowflake_function" "sql_test" {
  name     = "MY_SQL_FUNC"
  database = "MY_DB"
  schema   = "MY_SCHEMA"
  arguments {
    name = "arg1"
    type = "number"
  }
  comment             = "Example for SQL language"
  return_type         = "NUMBER(38,0)"
  null_input_behavior = "CALLED ON NULL INPUT"
  return_behavior     = "VOLATILE"
  statement           = "select arg1 + 1"
}
