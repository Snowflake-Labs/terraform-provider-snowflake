# format is database_name|schema_name|file_format_name|privilege|with_grant_option|on_future|roles
terraform import snowflake_file_format_grant.example "MY_DATABASE|MY_SCHEMA|MY_FILE_FORMAT|USAGE|false|false|role1,role2'
