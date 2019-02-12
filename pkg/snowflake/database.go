package snowflake

type DatabaseBuilder struct {
	name string
}

func Database(name string) *Builder {
	return &Builder{
		name:       name,
		entityType: DatabaseType,
	}
}
