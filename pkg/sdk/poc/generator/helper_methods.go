package generator

import (
	"fmt"
	"log"
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

var requiredFieldsForIDMethodMapping map[objectIdentifier][]string = map[objectIdentifier][]string{
	AccountObjectIdentifier:  {"Name"},
	DatabaseObjectIdentifier: {"Name", "DatabaseName"},
	SchemaObjectIdentifier:   {"Name", "DatabaseName", "SchemaName"},
}

func newIDHelperMethod(structName string, objectIdentifier objectIdentifier) *HelperMethod {
	var args string
	fields := requiredFieldsForIDMethodMapping[objectIdentifier]
	for _, field := range fields {
		args += fmt.Sprintf("v.%v, ", field)
	}
	returnValue := fmt.Sprintf("New%v(%v)", objectIdentifier, args)
	return newHelperMethod("ID", structName, returnValue, string(objectIdentifier))
}

func newObjectTypeHelperMethod(structName string) *HelperMethod {
	returnValue := fmt.Sprintf("ObjectType%v", structName)
	return newHelperMethod("ObjectType", structName, returnValue, "ObjectType")
}

func containsFieldNames(fields []*Field, names ...string) bool {
	fieldNames := map[string]any{}
	for _, field := range fields {
		fieldNames[field.Name] = nil
	}

	for _, name := range names {
		if _, ok := fieldNames[name]; !ok {
			return false
		}
	}
	return true
}

func hasRequiredFieldsForIDMethod(operations []*Operation, structName string, requiredFields ...string) bool {
	for _, op := range operations {
		if op.Name != string(OperationKindShow) {
			continue
		}
		for _, field := range op.HelperStructs {
			if field.Name == structName {
				return containsFieldNames(field.Fields, requiredFields...)
			}
		}
		log.Printf("WARNING: Struct: '%s' not found in '%s' operation. Couldn't generate ID() helper method.", structName, OperationKindShow)
	}
	log.Printf("WARNING: Operation: '%s' not found. Couldn't generate ID() helper method.", OperationKindShow)
	return false
}

// HelperMethodID adds a helper method "ID()" to the interface file that returns the ObjectIdentifier of the object
func (i *Interface) HelperMethodID() *Interface {
	identifierKind := identifierStringToObjectIdentifier(i.IdentifierKind)
	requiredFields := requiredFieldsForIDMethodMapping[identifierKind]
	if !hasRequiredFieldsForIDMethod(i.Operations, i.NameSingular, requiredFields...) {
		log.Printf("WARNING: Struct '%s' does not contain needed fields to build ID() helper method. Create the method manually in _ext file or add missing one of required fields: %v.\n", i.NameSingular, requiredFields)
		return i
	}
	i.HelperMethods = append(i.HelperMethods, newIDHelperMethod(i.NameSingular, identifierKind))
	return i
}

// HelperMethodObjectType adds a helper method "ObjectType()" to the interface file that returns the ObjectType for the struct
func (i *Interface) HelperMethodObjectType() *Interface {
	i.HelperMethods = append(i.HelperMethods, newObjectTypeHelperMethod(i.NameSingular))
	return i
}
