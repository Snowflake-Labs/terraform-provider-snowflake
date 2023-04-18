package snowflake

import "fmt"

type Identifier interface {
	QualifiedName() string
}

type TopLevelIdentifier struct {
	Name string
}

func (i *TopLevelIdentifier) QualifiedName() string {
	return i.Name
}

type SchemaIdentifier struct {
	Database string
	Schema   string
}

func (i *SchemaIdentifier) QualifiedName() string {
	return fmt.Sprintf("%v.%v", i.Database, i.Schema)
}

type SchemaObjectIdentifier struct {
	Database   string
	Schema     string
	ObjectName string `db:"NAME"`
}

func (i *SchemaObjectIdentifier) QualifiedName() string {
	return fmt.Sprintf("%v.%v.%v", i.Database, i.Schema, i.ObjectName)
}
