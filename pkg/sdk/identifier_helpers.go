package sdk

import (
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

func NewObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string) ObjectIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	switch len(parts) {
	case 1:
		return NewAccountObjectIdentifier(fullyQualifiedName)
	case 2:
		return NewSchemaIdentifier(parts[0], parts[1])
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

func (i ExternalObjectIdentifier) Name() string {
	return i.objectIdentifier.Name()
}

func (i ExternalObjectIdentifier) FullyQualifiedName() string {
	return fmt.Sprintf(`%v.%v`, i.accountIdentifier.Name(), i.objectIdentifier.FullyQualifiedName())
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

func (i AccountIdentifier) Name() string {
	if i.organizationName != "" && i.accountName != "" {
		return fmt.Sprintf("%s.%s", i.organizationName, i.accountName)
	}
	return i.accountLocator
}

type AccountObjectIdentifier struct {
	name string
}

func NewAccountObjectIdentifier(name string) AccountObjectIdentifier {
	return AccountObjectIdentifier{name: name}
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

func (i DatabaseObjectIdentifier) Name() string {
	return i.name
}

func (i DatabaseObjectIdentifier) FullyQualifiedName() string {
	if i.databaseName == "" && i.name == "" {
		return ""
	}
	return fmt.Sprintf(`"%v"."%v"`, i.databaseName, i.name)
}

type SchemaObjectIdentifier struct {
	databaseName string
	schemaName   string
	name         string
	arguments    []DataType
}

func NewSchemaObjectIdentifier(databaseName, schemaName, name string) SchemaObjectIdentifier {
	return SchemaObjectIdentifier{
		databaseName: strings.Trim(databaseName, `"`),
		schemaName:   strings.Trim(schemaName, `"`),
		name:         strings.Trim(name, `"`),
	}
}

func NewSchemaObjectIdentifierWithArguments(databaseName, schemaName, name string, arguments []DataType) SchemaObjectIdentifier {
	return SchemaObjectIdentifier{
		databaseName: strings.Trim(databaseName, `"`),
		schemaName:   strings.Trim(schemaName, `"`),
		name:         strings.Trim(name, `"`),
		arguments:    arguments,
	}
}

func NewSchemaObjectIdentifierFromFullyQualifiedName(fullyQualifiedName string) SchemaObjectIdentifier {
	parts := strings.Split(fullyQualifiedName, ".")
	id := SchemaObjectIdentifier{}
	id.databaseName = strings.Trim(parts[0], `"`)
	id.schemaName = strings.Trim(parts[1], `"`)

	// this is either a function or procedure
	if strings.HasSuffix(parts[2], ")") {
		idx := strings.LastIndex(parts[2], "(")
		id.name = strings.Trim(parts[2][:idx], `"`)
		strArgs := strings.Split(strings.Trim(parts[2][idx+1:], `)`), ",")
		id.arguments = make([]DataType, 0)
		for _, arg := range strArgs {
			trimmedArg := strings.TrimSpace(strings.Trim(arg, `"`))
			if trimmedArg == "" {
				continue
			}
			dt, _ := ToDataType(trimmedArg)
			id.arguments = append(id.arguments, dt)
		}
	} else { // this is every other kind of schema object
		id.name = strings.Trim(parts[2], `"`)
	}
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

func (i SchemaObjectIdentifier) Arguments() []DataType {
	return i.arguments
}

func (i SchemaObjectIdentifier) SchemaIdentifier() SchemaIdentifier {
	return NewSchemaIdentifier(i.databaseName, i.schemaName)
}

func (i SchemaObjectIdentifier) FullyQualifiedName() string {
	if i.schemaName == "" && i.databaseName == "" && i.name == "" {
		return ""
	}
	if len(i.arguments) == 0 {
		return fmt.Sprintf(`"%v"."%v"."%v"`, i.databaseName, i.schemaName, i.name)
	}
	// if this is a function or procedure, we need to include the arguments
	args := make([]string, len(i.arguments))
	for i, arg := range i.arguments {
		args[i] = string(arg)
	}
	return fmt.Sprintf(`"%v"."%v"."%v"(%v)`, i.databaseName, i.schemaName, i.name, strings.Join(args, ", "))
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
