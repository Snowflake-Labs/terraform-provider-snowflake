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

func (s *Operation) withShowObjectMethods(structName string, showObjectMethodsKind ...ShowObjectMethodType) *Operation {
	for _, methodKind := range showObjectMethodsKind {
		switch methodKind {
		case ShowObjectIdMethod:
			id, err := identifierStringToObjectIdentifier(s.ObjectInterface.IdentifierKind)
			if err != nil {
				log.Printf("[WARN] for showObjectIdMethod: %v", err)
				continue
			}
			if !hasRequiredFieldsForIDMethod(structName, s.HelperStructs, id) {
				log.Printf("[WARN] struct '%s' does not contain needed fields to build ID() helper method. Create the method manually in _ext file or add missing fields: %v.\n", structName, idTypeParts[id])
				continue
			}
			s.ShowObjectMethods = append(s.ShowObjectMethods, newShowObjectIDMethod(structName, s.HelperStructs, id))
		case ShowObjectTypeMethod:
			s.ShowObjectMethods = append(s.ShowObjectMethods, newShowObjectTypeMethod(structName))
		default:
			log.Println("[WARN] no showObjectMethod found for kind:", methodKind)
		}
	}
	return s
}

func hasRequiredFieldsForIDMethod(structName string, helperStructs []*Field, idType objectIdentifierKind) bool {
	if requiredFields, ok := idTypeParts[idType]; ok {
		for _, field := range helperStructs {
			if field.Name == structName {
				return containsFieldNames(field.Fields, requiredFields...)
			}
		}
	}
	log.Printf("[WARN] no required fields mapping defined for identifier %s", idType)
	return false
}

var idTypeParts map[objectIdentifierKind][]string = map[objectIdentifierKind][]string{
	AccountObjectIdentifier:  {"Name"},
	DatabaseObjectIdentifier: {"DatabaseName", "Name"},
	SchemaObjectIdentifier:   {"DatabaseName", "SchemaName", "Name"},
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

func newShowObjectIDMethod(structName string, helperStructs []*Field, idType objectIdentifierKind) *ShowObjectMethod {
	requiredFields := idTypeParts[idType]
	var args string
	for _, field := range requiredFields {
		args += fmt.Sprintf("v.%v, ", field)
	}

	returnValue := fmt.Sprintf("New%v(%v)", idType, args)
	return newShowObjectMethod("ID", structName, returnValue, string(idType))
}

func newShowObjectTypeMethod(structName string) *ShowObjectMethod {
	return newShowObjectMethod("ObjectType", structName, "ObjectType"+structName, "ObjectType")
}
