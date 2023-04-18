# format is database_name|schema_name|function_name|argument_data_types|privilege|with_grant_option|on_future|roles|shares
terraform import snowflake_function_grant.example "MY_DATABASE|MY_SCHEMA|MY_FUNCTION|ARG1TYPE,ARG2TYPE|USAGE|false|false|role1,role2|share1,share2"
