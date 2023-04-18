# format is database_name|schema_name|stage_name|privilege|with_grant_option|on_future|on_all|roles
terraform import snowflake_stage_grant.example "MY_DATABASE|MY_SCHEMA|MY_STAGE|USAGE|false|false|false|role1,role2"
