# format is application_role_name (string) | object_type (ACCOUNT_ROLE|APPLICATION) | grantee_name (string)
terraform import snowflake_grant_application_role.example '"my_application"."app_role_1"|ACCOUNT_ROLE|"my_role"'
