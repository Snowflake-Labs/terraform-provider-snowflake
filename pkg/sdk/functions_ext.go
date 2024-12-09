package sdk

func (v *Function) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, v.ArgumentsOld...)
}
