package snowflake

type UserBuilder struct {
	name string
}

func User(name string) *Builder {
	return &Builder{
		entityType: UserType,
		name:       name,
	}
}
