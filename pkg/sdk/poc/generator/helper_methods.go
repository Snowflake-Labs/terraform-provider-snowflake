package generator

import "fmt"

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

func newHelperMethod(name, structName, returnValue, returnType string) *HelperMethod {
	return &HelperMethod{
		Name:        name,
		StructName:  structName,
		ReturnValue: returnValue,
		ReturnType:  returnType,
	}
}

func newIDHelperMethod(structName, objectIdentifier string) *HelperMethod {
	var args string
	switch objectIdentifier {
	case "AccountObjectIdentifier":
		args = "v.Name"
	case "DatabaseObjectIdentifier":
		args = "v.DatabaseName, v.Name"
	case "SchemaObjectIdentifier":
		args = "v.DatabaseName, v.SchemaName, v.Name"
	default:
		return nil
	}
	returnValue := fmt.Sprintf("New%v(%v)", objectIdentifier, args)
	return newHelperMethod("ID", structName, returnValue, objectIdentifier)
}

func newObjectTypeHelperMethod(structName string) *HelperMethod {
	return newHelperMethod("ObjectType", structName, fmt.Sprintf("ObjectType%v", structName), "ObjectType")
}

func (i *Interface) ID() *Interface {
	i.HelperMethods = append(i.HelperMethods, newIDHelperMethod(i.NameSingular, i.IdentifierKind))
	return i
}

func (i *Interface) ObjectType() *Interface {
	i.HelperMethods = append(i.HelperMethods, newObjectTypeHelperMethod(i.NameSingular))
	return i
}
