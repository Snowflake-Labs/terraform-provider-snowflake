package snowflake

func Warehouse(name string) *Builder {
	return &Builder{
		name:       name,
		entityType: WarehouseType,
	}
}
