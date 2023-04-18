# format is database_name|schema_name|stream_name|privilege|with_grant_option|on_future|roles"
terraform import snowflake_stream_grant.example "MY_DATABASE|MY_SCHEMA|MY_STREAM|SELECT|false|false|role1,role2"
