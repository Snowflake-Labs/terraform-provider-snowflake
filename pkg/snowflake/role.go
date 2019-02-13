package snowflake

func Role(name string) *Builder {
	return &Builder{
		entityType: RoleType,
		name:       name,
	}
}
