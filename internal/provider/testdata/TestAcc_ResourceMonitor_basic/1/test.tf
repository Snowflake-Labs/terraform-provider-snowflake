resource "snowflake_resource_monitor" "test" {
	name            = var.name
	credit_quota    = var.credit_quota
	set_for_account = var.set_for_account
 	notify_triggers = [40]
	suspend_trigger = 80
	suspend_immediate_trigger = 90
	warehouses      = [snowflake_warehouse.warehouse.id]
}
