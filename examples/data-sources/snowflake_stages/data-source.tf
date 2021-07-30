data "snowflake_stages" "current" {
    database = "MYDB"
    schema   = "MYSCHEMA"
}