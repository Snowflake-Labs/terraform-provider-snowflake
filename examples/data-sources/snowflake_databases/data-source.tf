data "snowflake_databases" "this" {}

resource "snowflake_database" "backups" {
    for_each = { for x in data.snowflake_databases.this.databases: x.name => x }

    name  = "BACKUP_${each.key}"
    comment = "Backup of ${each.key} - ${each.value.comment}"
}
