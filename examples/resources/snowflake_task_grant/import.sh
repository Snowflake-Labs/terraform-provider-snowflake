# format is database_name|schema_name|task_name|privilege|with_grant_option|on_future|on_all|roles"
terraform import snowflake_task_grant.example "MY_DATABASE|MY_SCHEMA|MY_TASK|OPERATE|false|false|false|role1,role2"
