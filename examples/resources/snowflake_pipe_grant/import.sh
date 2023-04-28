# format is database_name|schema_name|pipe_name|privilege|with_grant_option|on_future|roles
terraform import snowflake_pipe_grant.example "MY_DATABASE|MY_SCHEMA|MY_PIPE_NAME|OPERATE|false|false|role1,role2'
