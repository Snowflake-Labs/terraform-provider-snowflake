package generator

import (
	"fmt"
	"log"
	"slices"
)

type ShowObjectMethodType uint

const (
	ShowObjectIdMethod ShowObjectMethodType = iota
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

func hasRequiredFieldsForIDMethod(structName string, helperStructs []*Field, idType objectIdentifier) bool {
	if requiredFields, ok := idTypeParts[idType]; ok {
		for _, field := range helperStructs {
			if field.Name == structName {
				return containsFieldNames(field.Fields, requiredFields...)
			}
		}
	}
	log.Printf("[WARN]: No required fields mapping defined for identifier %s", idType)
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

func (s *Operation) withShowObjectMethods(structName string, showObjectMethodsKind ...ShowObjectMethodType) *Operation {
	for _, methodKind := range showObjectMethodsKind {
		switch methodKind {
		case ShowObjectIdMethod:
			id, err := identifierStringToObjectIdentifier(s.ObjectInterface.IdentifierKind)
			if err != nil {
				log.Printf("[WARN]: %v, for showObjectIdMethod", err)
				continue
			}
			s.ShowObjectMethods = append(s.ShowObjectMethods, newShowObjectIDMethod(structName, s.HelperStructs, id))
		case ShowObjectTypeMethod:
			s.ShowObjectMethods = append(s.ShowObjectMethods, newShowObjectTypeMethod(structName))
		default:
			log.Println("No showObjectMethod found for kind:", methodKind)
		}
	}
	return s
}

func newShowObjectIDMethod(structName string, helperStructs []*Field, idType objectIdentifier) *ShowObjectMethod {
	if !hasRequiredFieldsForIDMethod(structName, helperStructs, idType) {
		log.Printf("[WARN]: Struct '%s' does not contain needed fields to build ID() helper method. Create the method manually in _ext file or add missing fields: %v.\n", structName, idTypeParts[idType])
		return nil
	}
	fields := idTypeParts[idType]
	var args string
	for _, field := range fields {
		args += fmt.Sprintf("v.%v, ", field)
	}

	returnValue := fmt.Sprintf("New%v(%v)", idType, args)
	return newShowObjectMethod("ID", structName, returnValue, string(idType))
}

func newShowObjectTypeMethod(structName string) *ShowObjectMethod {
	return newShowObjectMethod("ObjectType", structName, "ObjectType"+structName, "ObjectType")
}
