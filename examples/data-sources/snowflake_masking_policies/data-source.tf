data "snowflake_masking_policies" "current" {
    database = "MYDB"
    schema   = "MYSCHEMA"
}