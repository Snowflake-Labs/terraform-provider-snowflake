data "snowflake_database" "this" {
    name = "DEMO_DB"
}

resource "snowflake_database" "backup" {
    name  = "BACKUP_${data.snowflake_database.this.name}"
    comment = "Backup of ${data.snowflake_database.this.name} - ${data.snowflake_database.this.comment}"
}
