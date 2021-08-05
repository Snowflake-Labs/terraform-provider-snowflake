data "snowflake_sequences" "current" {
    database = "MYDB"
    schema   = "MYSCHEMA"
}