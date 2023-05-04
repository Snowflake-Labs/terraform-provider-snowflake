package sdk

// placeholder for the real implementation.
type DatabaseCreateOptions struct{}

type Database struct {
	Name string
}

func (v *Database) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}
