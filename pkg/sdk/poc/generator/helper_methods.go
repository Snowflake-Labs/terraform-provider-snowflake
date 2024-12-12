package generator

import (
	"fmt"
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

func newIDHelperMethod(structName string, objectIdentifier objectIdentifier) *HelperMethod {
	var args string
	switch objectIdentifier {
	case AccountObjectIdentifier:
		args = "v.Name"
	case DatabaseObjectIdentifier:
		args = "v.DatabaseName, v.Name"
	case SchemaObjectIdentifier:
		args = "v.DatabaseName, v.SchemaName, v.Name"
	default:
		return nil
	}

	returnValue := fmt.Sprintf("New%v(%v)", objectIdentifier, args)
	return newHelperMethod("ID", structName, returnValue, string(objectIdentifier))
}

func newObjectTypeHelperMethod(structName string) *HelperMethod {
	returnValue := fmt.Sprintf("ObjectType%v", structName)
	return newHelperMethod("ObjectType", structName, returnValue, "ObjectType")
}

func containsFieldNames(fields []*Field, names ...string) bool {
	fieldNames := map[string]bool{}
	for _, field := range fields {
		fieldNames[field.Name] = true
	}

	for _, name := range names {
		if _, ok := fieldNames[name]; !ok {
			return false
		}
	}
	return true
}

func checkRequiredFieldsForIDHelperMethod(operations []*Operation, name string, id objectIdentifier) bool {
	for _, op := range operations {
		if op.Name != string(OperationKindShow) {
			continue
		}
		for _, field := range op.HelperStructs {
			if field.Name != name {
				continue
			}
			requiredFields := []string{"Name"}
			switch id {
			case DatabaseObjectIdentifier:
				requiredFields = append(requiredFields, "DatabaseName")
			case SchemaObjectIdentifier:
				requiredFields = append(requiredFields, "DatabaseName", "SchemaName")
			}
			return containsFieldNames(field.Fields, requiredFields...)
		}
	}
	return false
}

// HelperMethodID adds a helper method "ID()" to the interface file that returns the ObjectIdentifier of the object
func (i *Interface) HelperMethodID() *Interface {
	if !checkRequiredFieldsForIDHelperMethod(i.Operations, i.NameSingular, identifierStringToObjectIdentifier(i.IdentifierKind)) {
		fmt.Println("WARNING: Does not contain needed fields for ID helper method. Create the method manually in _ext file or add missing fields.")
		return i
	}
	idKind := identifierStringToObjectIdentifier(i.IdentifierKind)
	i.HelperMethods = append(i.HelperMethods, newIDHelperMethod(i.NameSingular, idKind))
	return i
}

// HelperMethodObjectType adds a helper method "ObjectType()" to the interface file that returns the ObjectType for the struct
func (i *Interface) HelperMethodObjectType() *Interface {
	i.HelperMethods = append(i.HelperMethods, newObjectTypeHelperMethod(i.NameSingular))
	return i
}
