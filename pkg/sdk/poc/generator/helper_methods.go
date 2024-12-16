package generator

import (
	"fmt"
	"log"
	"slices"
)

type objectIdentifier string

const (
	AccountObjectIdentifier             objectIdentifier = "AccountObjectIdentifier"
	DatabaseObjectIdentifier            objectIdentifier = "DatabaseObjectIdentifier"
	SchemaObjectIdentifier              objectIdentifier = "SchemaObjectIdentifier"
	SchemaObjectIdentifierWithArguments objectIdentifier = "SchemaObjectIdentifierWithArguments"
)

func identifierStringToObjectIdentifier(s string) objectIdentifier {
	switch s {
	case "AccountObjectIdentifier":
		return AccountObjectIdentifier
	case "DatabaseObjectIdentifier":
		return DatabaseObjectIdentifier
	case "SchemaObjectIdentifier":
		return SchemaObjectIdentifier
	case "SchemaObjectIdentifierWithArguments":
		return SchemaObjectIdentifierWithArguments
	default:
		return ""
	}
}

type ObjectHelperMethodKind uint

const (
	ObjectHelperMethodID ObjectHelperMethodKind = iota
	ObjectHelperMethodObjectType
)

type HelperMethod struct {
	Name        string
	StructName  string
	ReturnValue string
	ReturnType  string
}

func newHelperMethod(name, structName, returnValue string, returnType string) *HelperMethod {
	return &HelperMethod{
		Name:        name,
		StructName:  structName,
		ReturnValue: returnValue,
		ReturnType:  returnType,
	}
}

func newObjectHelperMethodID(structName string, helperStructs []*Field, identifierString string) *HelperMethod {
	objectIdentifier := identifierStringToObjectIdentifier(identifierString)
	requiredFields, ok := requiredFieldsForIDMethodMapping[objectIdentifier]
	if !ok {
		log.Printf("WARNING: No required fields mapping defined for identifier %s", objectIdentifier)
		return nil
	}
	if !hasRequiredFieldsForIDMethod(structName, helperStructs, requiredFields...) {
		log.Printf("WARNING: Struct '%s' does not contain needed fields to build ID() helper method. Create the method manually in _ext file or add missing one of required fields: %v.\n", structName, requiredFields)
		return nil
	}

	var args string
	for _, field := range requiredFields {
		args += fmt.Sprintf("v.%v, ", field)
	}

	returnValue := fmt.Sprintf("New%v(%v)", objectIdentifier, args)
	return newHelperMethod("ID", structName, returnValue, string(objectIdentifier))
}

func newObjectHelperMethodObjectType(structName string) *HelperMethod {
	returnValue := fmt.Sprintf("ObjectType%v", structName)
	return newHelperMethod("ObjectType", structName, returnValue, "ObjectType")
}

var requiredFieldsForIDMethodMapping map[objectIdentifier][]string = map[objectIdentifier][]string{
	AccountObjectIdentifier:  {"Name"},
	DatabaseObjectIdentifier: {"Name", "DatabaseName"},
	SchemaObjectIdentifier:   {"Name", "DatabaseName", "SchemaName"},
}

func hasRequiredFieldsForIDMethod(structName string, helperStructs []*Field, requiredFields ...string) bool {
	for _, field := range helperStructs {
		if field.Name == structName {
			return containsFieldNames(field.Fields, requiredFields...)
		}
	}
	return false
}

func containsFieldNames(fields []*Field, names ...string) bool {
	fieldNames := []string{}
	for _, field := range fields {
		fieldNames = append(fieldNames, field.Name)
	}

	for _, name := range names {
		if !slices.Contains(names, name) {
			return false
		}
	}

	return true
}

func (s *Operation) withObjectHelperMethods(structName string, helperMethods ...ObjectHelperMethodKind) *Operation {
	for _, helperMethod := range helperMethods {
		switch helperMethod {
		case ObjectHelperMethodID:
			s.HelperMethods = append(s.HelperMethods, newObjectHelperMethodID(structName, s.HelperStructs, s.ObjectInterface.IdentifierKind))
		case ObjectHelperMethodObjectType:
			s.HelperMethods = append(s.HelperMethods, newObjectHelperMethodObjectType(structName))
		default:
			log.Println("No object helper method found for kind:", helperMethod)
		}
	}
	return s
}
