package generator

import "fmt"

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

func newIDHelperMethod(idKind idPrefix, structName, returnType string) *HelperMethod {
	var returnValue string
	switch idKind {
	case AccountIdentifierPrefix:
		returnValue = "NewAccountObjectIdentifier(v.Name)"
	case DatabaseIdentifierPrefix:
		returnValue = "NewDatabaseObjectIdentifier(v.DatabaseName, v.Name)"
	case SchemaIdentifierPrefix:
		returnValue = "NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)"
	default:
		return nil
	}
	return newHelperMethod("ID", structName, returnValue, returnType)
}

func newObjectTypeHelperMethod(structName string) *HelperMethod {
	return newHelperMethod("ObjectType", structName, fmt.Sprintf("ObjectType%v", structName), "ObjectType")
}

func (i *Interface) ID() *Interface {
	i.HelperMethods = append(i.HelperMethods, newIDHelperMethod(i.ObjectIdentifierPrefix(), i.NameSingular, i.IdentifierKind))
	return i
}

func (i *Interface) ObjectType() *Interface {
	i.HelperMethods = append(i.HelperMethods, newObjectTypeHelperMethod(i.NameSingular))
	return i
}

func (i *Interface) ObjectHelperMethods() *Interface {
	return i.ID().ObjectType()
}
