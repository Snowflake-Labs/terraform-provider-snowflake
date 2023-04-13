# format is database|schema|tag|privilege|with_grant_option|roles
terraform import snowflake_tag_grant.example "MY_DATABASE|MY_SCHEMA|MY_TAG|USAGE|false|role1,role2"
