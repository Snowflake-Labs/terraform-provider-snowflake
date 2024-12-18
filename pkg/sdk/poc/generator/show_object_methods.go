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

type ShowObjectMethodKind uint

const (
	ShowObjectIdMethod ShowObjectMethodKind = iota
	ShowObjectTypeMethod
)

type ShowObjectMethod struct {
	Name        string
	StructName  string
	ReturnValue string
	ReturnType  string
}

func newShowObjectMethod(name, structName, returnValue string, returnType string) *ShowObjectMethod {
	return &ShowObjectMethod{
		Name:        name,
		StructName:  structName,
		ReturnValue: returnValue,
		ReturnType:  returnType,
	}
}

var idTypeParts map[objectIdentifier][]string = map[objectIdentifier][]string{
	AccountObjectIdentifier:  {"Name"},
	DatabaseObjectIdentifier: {"DatabaseName", "Name"},
	SchemaObjectIdentifier:   {"DatabaseName", "SchemaName", "Name"},
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
		if !slices.Contains(fieldNames, name) {
			return false
		}
	}
	return true
}

func (s *Operation) withShowObjectMethods(structName string, showObjectMethodsKind ...ShowObjectMethodKind) *Operation {
	for _, methodKind := range showObjectMethodsKind {
		switch methodKind {
		case ShowObjectIdMethod:
			s.ShowObjectMethods = append(s.ShowObjectMethods, newShowObjectIDMethod(structName, s.HelperStructs, s.ObjectInterface.IdentifierKind))
		case ShowObjectTypeMethod:
			s.ShowObjectMethods = append(s.ShowObjectMethods, newShowObjectTypeMethod(structName))
		default:
			log.Println("No showObjectMethod found for kind:", methodKind)
		}
	}
	return s
}

func newShowObjectIDMethod(structName string, helperStructs []*Field, identifierString string) *ShowObjectMethod {
	objectIdentifier := identifierStringToObjectIdentifier(identifierString)
	requiredFields, ok := idTypeParts[objectIdentifier]
	if !ok {
		log.Printf("[WARN]: No required fields mapping defined for identifier %s", objectIdentifier)
		return nil
	}
	if !hasRequiredFieldsForIDMethod(structName, helperStructs, requiredFields...) {
		log.Printf("[WARN]: Struct '%s' does not contain needed fields to build ID() helper method. Create the method manually in _ext file or add missing one of required fields: %v.\n", structName, requiredFields)
		return nil
	}

	var args string
	for _, field := range requiredFields {
		args += fmt.Sprintf("v.%v, ", field)
	}

	returnValue := fmt.Sprintf("New%v(%v)", objectIdentifier, args)
	return newShowObjectMethod("ID", structName, returnValue, string(objectIdentifier))
}

func newShowObjectTypeMethod(structName string) *ShowObjectMethod {
	return newShowObjectMethod("ObjectType", structName, "ObjectType"+structName, "ObjectType")
}
