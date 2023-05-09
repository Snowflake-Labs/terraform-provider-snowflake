package sdk

import (
	"fmt"
	"strings"
)

type Identifier interface {
	Name() string
	FullyQualifiedName() string
}

type ObjectIdentifier interface {
	Identifier
	ObjectType() ObjectType
}

type ExternalObjectIdentifier struct {
	objectIdentifier  ObjectIdentifier
	accountIdentifier AccountIdentifier
}

func NewExternalObjectIdentifier(objectIdentifier ObjectIdentifier, accountIdentifier AccountIdentifier) ExternalObjectIdentifier {
	return ExternalObjectIdentifier{objectIdentifier: objectIdentifier, accountIdentifier: accountIdentifier}
}

func (i ExternalObjectIdentifier) Name() string {
	return i.objectIdentifier.Name()
}

func (i ExternalObjectIdentifier) FullyQualifiedName() string {
	return fmt.Sprintf(`%v.%v`, i.accountIdentifier.FullyQualifiedName(), i.objectIdentifier.FullyQualifiedName())
}

func (i ExternalObjectIdentifier) ObjectType() ObjectType {
	return i.objectIdentifier.ObjectType()
}

type AccountIdentifier struct {
	accountName      string
	organizationName string
}

func NewAccountIdentifier(accountName, organizationName string) AccountIdentifier {
	return AccountIdentifier{accountName: accountName, organizationName: organizationName}
}

func NewAccountIdentifierFromFullyQualifiedName(fullyQualifiedName string) AccountIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return AccountIdentifier{
		accountName:      strings.Trim(parts[0], `"`),
		organizationName: strings.Trim(parts[1], `"`),
	}
}

func (i AccountIdentifier) Name() string {
	return i.accountName
}

func (i AccountIdentifier) Organization() string {
	return i.organizationName
}

func (i AccountIdentifier) FullyQualifiedName() string {
	if i.accountName == "" && i.organizationName == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"`, i.accountName, i.organizationName)
}

type AccountLevelIdentifier struct {
	name       string
	objectType ObjectType
}

func NewAccountLevelIdentifier(name string, objectType ObjectType) AccountLevelIdentifier {
	return AccountLevelIdentifier{name: name, objectType: objectType}
}

func (i AccountLevelIdentifier) Name() string {
	return i.name
}

func (i AccountLevelIdentifier) FullyQualifiedName() string {
	if i.name == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"`, i.name)
}

func (i AccountLevelIdentifier) ObjectType() ObjectType {
	return i.objectType
}

type SchemaIdentifier struct {
	databaseName string
	schemaName   string
	objectType   ObjectType
}

func NewSchemaIdentifier(databaseName, schemaName string) SchemaIdentifier {
	return SchemaIdentifier{
		databaseName: strings.Trim(databaseName, `"`),
		schemaName:   strings.Trim(schemaName, `"`),
		objectType:   ObjectTypeSchema,
	}
}

func NewSchemaIdentifierFromFullyQualifiedName(fullyQualifiedName string) SchemaIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return SchemaIdentifier{
		databaseName: strings.Trim(parts[0], `"`),
		schemaName:   strings.Trim(parts[1], `"`),
		objectType:   ObjectTypeSchema,
	}
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

func (i SchemaIdentifier) ObjectType() ObjectType {
	return i.objectType
}

type SchemaObjectIdentifier struct {
	databaseName string
	schemaName   string
	name         string
	objectType   ObjectType
}

func NewSchemaObjectIdentifier(databaseName, schemaName, name string, objectType ObjectType) SchemaObjectIdentifier {
	return SchemaObjectIdentifier{
		databaseName: strings.Trim(databaseName, `"`),
		schemaName:   strings.Trim(schemaName, `"`),
		name:         strings.Trim(name, `"`),
		objectType:   objectType,
	}
}

func NewSchemaObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string, objectType ObjectType) SchemaObjectIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return SchemaObjectIdentifier{
		databaseName: strings.Trim(parts[0], `"`),
		schemaName:   strings.Trim(parts[1], `"`),
		name:         strings.Trim(parts[2], `"`),
		objectType:   objectType,
	}
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

func (i SchemaObjectIdentifier) ObjectType() ObjectType {
	return i.objectType
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

func (i TableColumnIdentifier) ObjectType() ObjectType {
	return ObjectTypeTableColumn
}

type InboundShareIdentifier struct {
	providerAccount string
	shareName       string
}

func NewInboundShareIdentifier(providerAccount, shareName string) InboundShareIdentifier {
	return InboundShareIdentifier{providerAccount: providerAccount, shareName: shareName}
}

func NewInboundShareIdentifierFromFullyQualifiedName(fullyQualifiedName string) InboundShareIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return InboundShareIdentifier{
		providerAccount: strings.Trim(parts[0], `"`),
		shareName:       strings.Trim(parts[1], `"`),
	}
}

func (i InboundShareIdentifier) Account() string {
	return i.providerAccount
}

func (i InboundShareIdentifier) Name() string {
	return i.shareName
}

func (i InboundShareIdentifier) FullyQualifiedName() string {
	if i.providerAccount == "" && i.shareName == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"`, i.providerAccount, i.shareName)
}
