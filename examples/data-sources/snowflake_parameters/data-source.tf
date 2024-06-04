resource "snowflake_standard_database" "d" {
  name = "TEST_DB"
}

// read all object parameters in database TEST_DB
data "snowflake_parameters" "p" {
  parameter_type = "OBJECT"
  object_type    = "DATABASE"
  object_name    = snowflake_standard_database.d.name
}

// read all account parameters with the pattern '%TIMESTAMP%'
data "snowflake_parameters" "p2" {
  parameter_type = "ACCOUNT"
  pattern        = "%TIMESTAMP%"
}

// read the exact session parameter ROWS_PER_RESULTSET
data "snowflake_parameters" "p3" {
  parameter_type = "SESSION"
  pattern        = "ROWS_PER_RESULTSET"
  user           = "TEST_USER"
}
