# format is database_name|privilege|with_grant_option|roles|shares
terraform import snowflake_database_grant.example "MY_DATABASE|USAGE|false|role1,role2|share1,share2"
