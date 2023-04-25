package sdk

import (
	"fmt"
	"strings"
)

type ObjectIdentifier interface {
	FullyQualifiedName() string
}

type AccountObjectIdentifier struct {
	Name string
}

func NewAccountObjectIdentifier(name string) *AccountObjectIdentifier {
	return &AccountObjectIdentifier{Name: name}
}

func (v *AccountObjectIdentifier) FullyQualifiedName() string {
	return fmt.Sprintf(`"%v"`, v.Name)
}

type SchemaIdentifier struct {
	DatabaseName string
	SchemaName   string
}

func NewSchemaIdentifier(databaseName, schemaName string) *SchemaIdentifier {
	return &SchemaIdentifier{
		DatabaseName: strings.Trim(databaseName, `"`),
		SchemaName:   strings.Trim(schemaName, `"`),
	}
}

func NewSchemaIdentifierFromFullyQualifiedName(fullyQualifiedName string) *SchemaIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return &SchemaIdentifier{
		DatabaseName: strings.Trim(parts[0], `"`),
		SchemaName:   strings.Trim(parts[1], `"`),
	}
}

func (i *SchemaIdentifier) FullyQualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"`, i.DatabaseName, i.SchemaName)
}

type SchemaObjectIdentifier struct {
	DatabaseName string
	SchemaName   string
	Name         string
}

func NewSchemaObjectIdentifier(databaseName, schemaName, name string) SchemaObjectIdentifier {
	return SchemaObjectIdentifier{
		DatabaseName: strings.Trim(databaseName, `"`),
		SchemaName:   strings.Trim(schemaName, `"`),
		Name:         strings.Trim(name, `"`),
	}
}

func NewSchemaObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string) SchemaObjectIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return SchemaObjectIdentifier{
		DatabaseName: strings.Trim(parts[0], `"`),
		SchemaName:   strings.Trim(parts[1], `"`),
		Name:         strings.Trim(parts[2], `"`),
	}
}

func (i SchemaObjectIdentifier) FullyQualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"."%v"`, i.DatabaseName, i.SchemaName, i.Name)
}

type TableColumnIdentifier struct {
	DatabaseName string
	SchemaName   string
	TableName    string
	ColumnName   string
}

func NewTableColumnIdentifier(databaseName, schemaName, tableName, columnName string) *TableColumnIdentifier {
	return &TableColumnIdentifier{DatabaseName: databaseName, SchemaName: schemaName, TableName: tableName, ColumnName: columnName}
}

func NewTableColumnIdentifierFromFullyQualifiedName(fullyQualifiedName string) *TableColumnIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return &TableColumnIdentifier{
		DatabaseName: strings.Trim(parts[0], `"`),
		SchemaName:   strings.Trim(parts[1], `"`),
		TableName:    strings.Trim(parts[2], `"`),
		ColumnName:   strings.Trim(parts[3], `"`),
	}
}

func (i *TableColumnIdentifier) FullyQualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"."%v"."%v"`, i.DatabaseName, i.SchemaName, i.TableName, i.ColumnName)
}
