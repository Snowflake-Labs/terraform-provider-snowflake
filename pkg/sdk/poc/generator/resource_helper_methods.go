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

type ResourceHelperMethodKind uint

const (
	ResourceIDHelperMethod ResourceHelperMethodKind = iota
	ResourceObjectTypeHelperMethod
)

type ResourceHelperMethod struct {
	Name        string
	StructName  string
	ReturnValue string
	ReturnType  string
}

func newResourceHelperMethod(name, structName, returnValue string, returnType string) *ResourceHelperMethod {
	return &ResourceHelperMethod{
		Name:        name,
		StructName:  structName,
		ReturnValue: returnValue,
		ReturnType:  returnType,
	}
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

func newResourceIDHelperMethod(structName string, helperStructs []*Field, identifierString string) *ResourceHelperMethod {
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
	return newResourceHelperMethod("ID", structName, returnValue, string(objectIdentifier))
}

func newResourceObjectTypeHelperMethod(structName string) *ResourceHelperMethod {
	return newResourceHelperMethod("ObjectType", structName, "ObjectType"+structName, "ObjectType")
}

func (s *Operation) withResourceHelperMethods(structName string, helperMethods ...ResourceHelperMethodKind) *Operation {
	for _, helperMethod := range helperMethods {
		switch helperMethod {
		case ResourceIDHelperMethod:
			s.ResourceHelperMethods = append(s.ResourceHelperMethods, newResourceIDHelperMethod(structName, s.HelperStructs, s.ObjectInterface.IdentifierKind))
		case ResourceObjectTypeHelperMethod:
			s.ResourceHelperMethods = append(s.ResourceHelperMethods, newResourceObjectTypeHelperMethod(structName))
		default:
			log.Println("No resourceHelperMethod found for kind:", helperMethod)
		}
	}
	return s
}
