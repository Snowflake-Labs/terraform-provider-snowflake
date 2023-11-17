resource "snowflake_schema" "other_schema" {
    database = var.database
	name     = var.other_schema
	comment  = "Other schema"
}

resource "snowflake_table" "t" {
	database = var.database
	schema   = var.schema
	name     = var.table_name
    change_tracking = true
	column {
		name = "id"
		type = "NUMBER(38,0)"
	}
}
