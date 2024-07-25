package sdk

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type Identifier interface {
	Name() string
}

type ObjectIdentifier interface {
	Identifier
	FullyQualifiedName() string
}

// TODO(SNOW-999049): This function will be tested/improved/used more wiedely during the identifiers rework.
// Right now, the implementation is just a copy of DecodeSnowflakeParameterID used in resources.
func ParseObjectIdentifier(identifier string) (ObjectIdentifier, error) {
	reader := csv.NewReader(strings.NewReader(identifier))
	reader.Comma = '.'
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read identifier: %s, err = %w", identifier, err)
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("incompatible identifier: %s", identifier)
	}
	parts := lines[0]
	switch len(parts) {
	case 1:
		return NewAccountObjectIdentifier(parts[0]), nil
	case 2:
		return NewDatabaseObjectIdentifier(parts[0], parts[1]), nil
	case 3:
		return NewSchemaObjectIdentifier(parts[0], parts[1], parts[2]), nil
	case 4:
		return NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3]), nil
	default:
		return nil, fmt.Errorf("unable to classify identifier: %s", identifier)
	}
}

func NewObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string) ObjectIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	switch len(parts) {
	case 1:
		return NewAccountObjectIdentifier(fullyQualifiedName)
	case 2:
		return NewDatabaseObjectIdentifier(parts[0], parts[1])
	case 3:
		return NewSchemaObjectIdentifier(parts[0], parts[1], parts[2])
	case 4:
		return NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3])
	}
	return NewAccountObjectIdentifier(fullyQualifiedName)
}

// for objects that live in other accounts
type ExternalObjectIdentifier struct {
	objectIdentifier  ObjectIdentifier
	accountIdentifier AccountIdentifier
}

func NewExternalObjectIdentifier(accountIdentifier AccountIdentifier, objectIdentifier ObjectIdentifier) ExternalObjectIdentifier {
	return ExternalObjectIdentifier{
		objectIdentifier:  objectIdentifier,
		accountIdentifier: accountIdentifier,
	}
}

func NewExternalObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string) ExternalObjectIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")

	if len(parts) == 1 {
		return ExternalObjectIdentifier{
			objectIdentifier:  NewAccountObjectIdentifier(fullyQualifiedName),
			accountIdentifier: NewAccountIdentifier("", ""),
		}
	}

	if len(parts) == 2 {
		accountLocator := parts[0]
		objectName := parts[1]

		return ExternalObjectIdentifier{
			objectIdentifier:  NewAccountObjectIdentifier(objectName),
			accountIdentifier: NewAccountIdentifierFromAccountLocator(accountLocator),
		}
	}

	orgName := parts[0]
	accountName := parts[1]
	objectName := strings.Join(parts[2:], ".")

	return ExternalObjectIdentifier{
		objectIdentifier:  NewAccountObjectIdentifier(objectName),
		accountIdentifier: NewAccountIdentifier(orgName, accountName),
	}
}

func (i ExternalObjectIdentifier) AccountIdentifier() AccountIdentifier {
	return i.accountIdentifier
}

func (i ExternalObjectIdentifier) Name() string {
	return i.objectIdentifier.Name()
}

func (i ExternalObjectIdentifier) FullyQualifiedName() string {
	return fmt.Sprintf(`%v.%v`, i.accountIdentifier.FullyQualifiedName(), i.objectIdentifier.FullyQualifiedName())
}

type AccountIdentifier struct {
	organizationName string
	accountName      string
	accountLocator   string
}

func NewAccountIdentifier(organizationName, accountName string) AccountIdentifier {
	return AccountIdentifier{
		organizationName: strings.Trim(organizationName, `"`),
		accountName:      strings.Trim(accountName, `"`),
	}
}

func NewAccountIdentifierFromAccountLocator(accountLocator string) AccountIdentifier {
	return AccountIdentifier{
		accountLocator: accountLocator,
	}
}

func NewAccountIdentifierFromFullyQualifiedName(fullyQualifiedName string) AccountIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	if len(parts) == 1 {
		return NewAccountIdentifierFromAccountLocator(fullyQualifiedName)
	}
	organizationName := strings.Trim(parts[0], `"`)
	accountName := strings.Trim(parts[1], `"`)
	return NewAccountIdentifier(organizationName, accountName)
}

func (i AccountIdentifier) OrganizationName() string {
	return i.organizationName
}

func (i AccountIdentifier) AccountName() string {
	return i.accountName
}

func (i AccountIdentifier) Name() string {
	if i.organizationName != "" && i.accountName != "" {
		return fmt.Sprintf("%s.%s", i.organizationName, i.accountName)
	}
	return i.accountLocator
}

func (i AccountIdentifier) FullyQualifiedName() string {
	if i.organizationName != "" && i.accountName != "" {
		return fmt.Sprintf(`"%s"."%s"`, i.organizationName, i.accountName)
	}
	return fmt.Sprintf(`"%s"`, i.accountLocator)
}

type AccountObjectIdentifier struct {
	name string
}

func NewAccountObjectIdentifier(name string) AccountObjectIdentifier {
	return AccountObjectIdentifier{
		name: strings.Trim(name, `"`),
	}
}

func NewAccountObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string) AccountObjectIdentifier {
	name := strings.Trim(fullyQualifiedName, `"`)
	return AccountObjectIdentifier{name: name}
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

type DatabaseObjectIdentifier struct {
	databaseName string
	name         string
}

func NewDatabaseObjectIdentifier(databaseName, name string) DatabaseObjectIdentifier {
	return DatabaseObjectIdentifier{
		databaseName: strings.Trim(databaseName, `"`),
		name:         strings.Trim(name, `"`),
	}
}

func NewDatabaseObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string) DatabaseObjectIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	return DatabaseObjectIdentifier{
		databaseName: strings.Trim(parts[0], `"`),
		name:         strings.Trim(parts[1], `"`),
	}
}

func (i DatabaseObjectIdentifier) DatabaseName() string {
	return i.databaseName
}

func (i DatabaseObjectIdentifier) DatabaseId() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(i.databaseName)
}

func (i DatabaseObjectIdentifier) Name() string {
	return i.name
}

func (i DatabaseObjectIdentifier) FullyQualifiedName() string {
	if i.name == "" && i.databaseName == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"`, i.databaseName, i.name)
}

type SchemaObjectIdentifier struct {
	databaseName string
	schemaName   string
	name         string
}

func NewSchemaObjectIdentifierInSchema(schemaId DatabaseObjectIdentifier, name string) SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(schemaId.DatabaseName(), schemaId.Name(), name)
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
	id := SchemaObjectIdentifier{}
	id.databaseName = strings.Trim(parts[0], `"`)
	id.schemaName = strings.Trim(parts[1], `"`)
	id.name = strings.Trim(parts[2], `"`)
	return id
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

func (i SchemaObjectIdentifier) SchemaId() DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(i.databaseName, i.schemaName)
}

func (i SchemaObjectIdentifier) DatabaseId() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(i.databaseName)
}

func (i SchemaObjectIdentifier) FullyQualifiedName() string {
	if i.schemaName == "" && i.databaseName == "" && i.name == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"."%v"`, i.databaseName, i.schemaName, i.name)
}

// TODO:
// - Add parser
// - Add to IsValidIdentifier
// - Handle in the sql_builder
// - Use in function,procedure,external_function
// - Function (Test, Impl)
// - Fix after arguments removed from SchemaObjectIdentifier
// - Look for todos on SNOW-999049

type SchemaObjectIdentifierWithArguments struct {
	databaseName string
	schemaName   string
	name         string
	arguments    []string
}

func NewSchemaObjectIdentifierWithArguments(databaseName, schemaName, name string, arguments ...string) SchemaObjectIdentifierWithArguments {
	return SchemaObjectIdentifierWithArguments{
		databaseName: strings.Trim(databaseName, `"`),
		schemaName:   strings.Trim(schemaName, `"`),
		name:         strings.Trim(name, `"`),
		arguments:    arguments,
	}
}

// TODO:
//func NewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName(fullyQualifiedName string) SchemaObjectIdentifierWithArguments {
//	parts := strings.Split(fullyQualifiedName, ".")
//	id := SchemaObjectIdentifierWithArguments{
//		databaseName: strings.Trim(parts[0], `"`),
//		schemaName:   strings.Trim(parts[1], `"`),
//		name:         strings.Trim(parts[2], `"`),
//		// TODO: Arguments
//	}
//	return id
//}

func (i SchemaObjectIdentifierWithArguments) DatabaseName() string {
	return i.databaseName
}

func (i SchemaObjectIdentifierWithArguments) SchemaName() string {
	return i.schemaName
}

func (i SchemaObjectIdentifierWithArguments) Name() string {
	return i.name
}

func (i SchemaObjectIdentifierWithArguments) Arguments() []string {
	return i.arguments
}

func (i SchemaObjectIdentifierWithArguments) SchemaId() DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(i.databaseName, i.schemaName)
}

func (i SchemaObjectIdentifierWithArguments) DatabaseId() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(i.databaseName)
}

func (i SchemaObjectIdentifierWithArguments) FullyQualifiedName() string {
	if i.schemaName == "" && i.databaseName == "" && i.name == "" && len(i.arguments) == 0 {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"."%v"(%v)`, i.databaseName, i.schemaName, i.name, strings.Join(i.arguments, ", "))
}

type TableColumnIdentifier struct {
	databaseName string
	schemaName   string
	tableName    string
	columnName   string
}

func NewTableColumnIdentifier(databaseName, schemaName, tableName, columnName string) TableColumnIdentifier {
	return TableColumnIdentifier{
		databaseName: strings.Trim(databaseName, `"`),
		schemaName:   strings.Trim(schemaName, `"`),
		tableName:    strings.Trim(tableName, `"`),
		columnName:   strings.Trim(columnName, `"`),
	}
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

func (i TableColumnIdentifier) SchemaObjectId() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(i.databaseName, i.schemaName, i.tableName)
}
