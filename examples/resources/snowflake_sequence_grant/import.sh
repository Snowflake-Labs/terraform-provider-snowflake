# format is database_name|schema_name|sequence_name|privilege|with_grant_option|on_future|roles
terraform import snowflake_sequence_grant.example "MY_DATABASE|MY_SCHEMA|MY_SEQUENCE|USAGE|false|false|role1,role2"
