package sdk

// Stage is a placeholder for now, will be implemented later.
type Stage struct {
	DatabaseName string
	SchemaName   string
	Name         string
}

func (v *Stage) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Stage) ObjectType() ObjectType {
	return ObjectTypeStage
}
