package snowflake

import (
	"fmt"
	"strings"
)

type Identifier interface {
	QualifiedName() string
}

type TopLevelIdentifier struct {
	Name string
}

func (i *TopLevelIdentifier) QualifiedName() string {
	return i.Name
}

func TopLevelIdentifierFromQualifiedName(name string) *TopLevelIdentifier {
	return &TopLevelIdentifier{
		Name: name,
	}
}

type SchemaIdentifier struct {
	Database string
	Schema   string
}

func (i *SchemaIdentifier) QualifiedName() string {
	return fmt.Sprintf("%v.%v", i.Database, i.Schema)
}

func SchemaIdentifierFromQualifiedName(name string) *SchemaIdentifier {
	parts := strings.Split(name, ".")
	return &SchemaIdentifier{
		Database: parts[0],
		Schema:   parts[1],
	}
}

type SchemaObjectIdentifier struct {
	Database   string
	Schema     string
	ObjectName string `db:"NAME"`
}

func (i *SchemaObjectIdentifier) QualifiedName() string {
	return fmt.Sprintf("%v.%v.%v", i.Database, i.Schema, i.ObjectName)
}

func SchemaObjectIdentifierFromQualifiedName(name string) *SchemaObjectIdentifier {
	parts := strings.Split(name, ".")
	return &SchemaObjectIdentifier{
		Database:   parts[0],
		Schema:     parts[1],
		ObjectName: parts[2],
	}
}
