package generator

import (
	"log"
	"slices"
)

type ShowObjectIdMethod struct {
	StructName     string
	IdentifierKind objectIdentifierKind
	Args           []string
}

func newShowObjectIDMethod(structName string, idType objectIdentifierKind) *ShowObjectIdMethod {
	return &ShowObjectIdMethod{
		StructName:     structName,
		IdentifierKind: idType,
		Args:           idTypeParts[idType],
	}
}

var idTypeParts map[objectIdentifierKind][]string = map[objectIdentifierKind][]string{
	AccountObjectIdentifier:  {"Name"},
	DatabaseObjectIdentifier: {"DatabaseName", "Name"},
	SchemaObjectIdentifier:   {"DatabaseName", "SchemaName", "Name"},
}

func checkRequiredFieldsForIdMethod(structName string, helperStructs []*Field, idKind objectIdentifierKind) bool {
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

type ShowObjectTypeMethod struct {
	StructName string
}

func newShowObjectTypeMethod(structName string) *ShowObjectTypeMethod {
	return &ShowObjectTypeMethod{StructName: structName}
}
