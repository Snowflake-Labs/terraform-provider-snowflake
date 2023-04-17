# format is database_name|schema_name|procedure_name|argument_data_types|privilege|with_grant_option|on_future|roles|shares
terraform import snowflake_procedure_grant.example "MY_DATABASE|MY_SCHEMA|MY_PROCEDURE|ARG1TYPE,ARG2TYPE|USAGE|false|false|role1,role2|share1,share2"
