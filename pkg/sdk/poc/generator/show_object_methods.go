package generator

import (
	"fmt"
	"log"
	"slices"
)


type ShowObjectMethod struct {
	Name        string
	StructName  string
	ReturnValue string
	ReturnType  string
}

type ShowObjectIdMethod struct {

}

func newShowObjectMethod(name, structName, returnValue string, returnType string) *ShowObjectMethod {
	return &ShowObjectMethod{
		Name:        name,
		StructName:  structName,
		ReturnValue: returnValue,
		ReturnType:  returnType,
	}
}

func checkRequiredFieldsForIDMethod(structName string, helperStructs []*Field, idKind objectIdentifierKind) bool {
	if requiredFields, ok := idTypeParts[idKind]; ok {
		for _, field := range helperStructs {
			if field.Name == structName {
				return containsFieldNames(field.Fields, requiredFields...)
			}
		}
	}
	log.Printf("[WARN] no required fields mapping defined for identifier %s", idKind)
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

func newShowObjectIDMethod(structName string, idType objectIdentifierKind) *ShowObjectMethod {
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
