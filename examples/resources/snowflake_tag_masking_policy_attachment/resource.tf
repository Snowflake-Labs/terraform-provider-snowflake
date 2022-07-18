resource "snowflake_database" "test" {
  name                        = "TEST_DB1"
  data_retention_time_in_days = 1
}

resource "snowflake_database" "test2" {
  name                        = "TEST_DB2"
  data_retention_time_in_days = 1
}


resource "snowflake_schema" "test2" {
  database            = snowflake_database.test2.name
  name                = "FOOBAR2"
  data_retention_days = snowflake_database.test2.data_retention_time_in_days
}

resource "snowflake_schema" "test" {
  database            = snowflake_database.test.name
  name                = "FOOBAR"
  data_retention_days = snowflake_database.test.data_retention_time_in_days
}

resource "snowflake_tag" "this" {
  name     = upper("test_tag")
  database = snowflake_database.test2.name
  schema   = snowflake_schema.test2.name
}

resource "snowflake_masking_policy" "example_masking_policy" {
  name               = "EXAMPLE_MASKING_POLICY"
  database           = snowflake_database.test.name
  schema             = snowflake_schema.test.name
  value_data_type    = "string"
  masking_expression = "case when current_role() in ('ACCOUNTADMIN') then val else sha2(val, 512) end"
  return_data_type   = "string"
}

resource "snowflake_tag_masking_policy_attachment" "name" {
  tag_database            = snowflake_database.test2.name
  tag_schema              = snowflake_schema.test2.name
  tag_name                = snowflake_tag.this.name
  masking_policy_database = snowflake_database.test.name
  masking_policy_schema   = snowflake_schema.test.name
  masking_policy_name     = snowflake_masking_policy.example_masking_policy.name

}