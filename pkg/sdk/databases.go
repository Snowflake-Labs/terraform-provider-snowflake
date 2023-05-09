package sdk

// placeholder for the real implementation.
type DatabaseCreateOptions struct{}

type Database struct {
	Name string
}

func (v *Database) ID() AccountLevelIdentifier {
	return NewAccountLevelIdentifier(v.Name, ObjectTypeDatabase)
}
