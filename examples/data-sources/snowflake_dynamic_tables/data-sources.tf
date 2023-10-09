data "snowflake_dynamic_tables" "dts" {
    like {
        pattern = "product"
    }
    in {
        database = "mydb"
    }
}

output "dt" {
    value = data.snowflake_dynamic_tables.dts.records[0]
}
