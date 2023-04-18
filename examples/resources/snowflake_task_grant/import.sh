# format is database_name|schema_name|task_name|privilege|with_grant_option|on_future|roles"
terraform import snowflake_task_grant.example "MY_DATABASE|MY_SCHEMA|MY_TASK|OPERATE|false|false|role1,role2"
