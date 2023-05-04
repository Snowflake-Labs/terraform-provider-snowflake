package sdk

import (
	"fmt"
	"strings"
)

type ObjectIdentifier interface {
	Literal() string
	Name() string
	FullyQualifiedName() string
}

type AccountObjectIdentifier struct {
	name string
}

func NewAccountObjectIdentifier(name string) AccountObjectIdentifier {
	return AccountObjectIdentifier{name: name}
}

func (i AccountObjectIdentifier) Literal() string {
	return i.name
}

func (i AccountObjectIdentifier) Name() string {
	return i.name
}

func (i AccountObjectIdentifier) FullyQualifiedName() string {
	if i.name == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"`, i.name)
}

type SchemaIdentifier struct {
	databaseName string
	schemaName   string
}

func NewSchemaIdentifier(databaseName, schemaName string) SchemaIdentifier {
	return SchemaIdentifier{
		databaseName: strings.Trim(databaseName, `"`),
		schemaName:   strings.Trim(schemaName, `"`),
	}
}

func NewSchemaIdentifierFromFullyQualifiedName(fullyQualifiedName string) SchemaIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return SchemaIdentifier{
		databaseName: strings.Trim(parts[0], `"`),
		schemaName:   strings.Trim(parts[1], `"`),
	}
}

func (i SchemaIdentifier) Literal() string {
	return fmt.Sprintf(`%v.%v`, i.databaseName, i.schemaName)
}

func (i SchemaIdentifier) DatabaseName() string {
	return i.databaseName
}

func (i SchemaIdentifier) Name() string {
	return i.schemaName
}

func (i SchemaIdentifier) FullyQualifiedName() string {
	if i.schemaName == "" && i.databaseName == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"`, i.databaseName, i.schemaName)
}

type SchemaObjectIdentifier struct {
	databaseName string
	schemaName   string
	name         string
}

func NewSchemaObjectIdentifier(databaseName, schemaName, name string) SchemaObjectIdentifier {
	return SchemaObjectIdentifier{
		databaseName: strings.Trim(databaseName, `"`),
		schemaName:   strings.Trim(schemaName, `"`),
		name:         strings.Trim(name, `"`),
	}
}

func NewSchemaObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string) SchemaObjectIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return SchemaObjectIdentifier{
		databaseName: strings.Trim(parts[0], `"`),
		schemaName:   strings.Trim(parts[1], `"`),
		name:         strings.Trim(parts[2], `"`),
	}
}

func (i SchemaObjectIdentifier) Literal() string {
	return fmt.Sprintf(`%v.%v.%v`, i.databaseName, i.schemaName, i.name)
}

func (i SchemaObjectIdentifier) DatabaseName() string {
	return i.databaseName
}

func (i SchemaObjectIdentifier) SchemaName() string {
	return i.schemaName
}

func (i SchemaObjectIdentifier) Name() string {
	return i.name
}

func (i SchemaObjectIdentifier) FullyQualifiedName() string {
	if i.schemaName == "" && i.databaseName == "" && i.name == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"."%v"`, i.databaseName, i.schemaName, i.name)
}

type TableColumnIdentifier struct {
	databaseName string
	schemaName   string
	tableName    string
	columnName   string
}

func NewTableColumnIdentifier(databaseName, schemaName, tableName, columnName string) TableColumnIdentifier {
	return TableColumnIdentifier{databaseName: databaseName, schemaName: schemaName, tableName: tableName, columnName: columnName}
}

func NewTableColumnIdentifierFromFullyQualifiedName(fullyQualifiedName string) TableColumnIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return TableColumnIdentifier{
		databaseName: strings.Trim(parts[0], `"`),
		schemaName:   strings.Trim(parts[1], `"`),
		tableName:    strings.Trim(parts[2], `"`),
		columnName:   strings.Trim(parts[3], `"`),
	}
}

func (i TableColumnIdentifier) Literal() string {
	return fmt.Sprintf(`%v.%v.%v.%v`, i.databaseName, i.schemaName, i.tableName, i.columnName)
}

func (i TableColumnIdentifier) DatabaseName() string {
	return i.databaseName
}

func (i TableColumnIdentifier) SchemaName() string {
	return i.schemaName
}

func (i TableColumnIdentifier) TableName() string {
	return i.tableName
}

func (i TableColumnIdentifier) Name() string {
	return i.columnName
}

func (i TableColumnIdentifier) FullyQualifiedName() string {
	if i.schemaName == "" && i.databaseName == "" && i.tableName == "" && i.columnName == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"."%v"."%v"`, i.databaseName, i.schemaName, i.tableName, i.columnName)
}
