package sdk

import (
	"strings"

	"golang.org/x/exp/slices"
)

// Object bundles together the object type and name. Its used for DDL statements.
type Object struct {
	ObjectType ObjectType       `ddl:"keyword"`
	Name       ObjectIdentifier `ddl:"identifier"`
}

// ObjectType is the type of object.
type ObjectType string

const (
	ObjectTypeAccount          ObjectType = "ACCOUNT"
	ObjectTypeAccountParameter ObjectType = "ACCOUNT PARAMETER"
	ObjectTypeDatabase         ObjectType = "DATABASE"
	ObjectTypeFailoverGroup    ObjectType = "FAILOVER GROUP"
	ObjectTypeIntegration      ObjectType = "INTEGRATION"
	ObjectTypeMaskingPolicy    ObjectType = "MASKING POLICY"
	ObjectTypeNetworkPolicy    ObjectType = "NETWORK POLICY"
	ObjectTypePasswordPolicy   ObjectType = "PASSWORD POLICY"
	ObjectTypeResourceMonitor  ObjectType = "RESOURCE MONITOR"
	ObjectTypeRole             ObjectType = "ROLE"
	ObjectTypeSchema           ObjectType = "SCHEMA"
	ObjectTypeShare            ObjectType = "SHARE"
	ObjectTypeUser             ObjectType = "USER"
	ObjectTypeWarehouse        ObjectType = "WAREHOUSE"
)

func ObjectTypeFromPluralString(s string) ObjectType {
	// only care about the "ies" endings.
	switch s {
	case "MASKING POLICIES":
		return ObjectTypeMaskingPolicy
	case "NETWORK POLICIES":
		return ObjectTypeNetworkPolicy
	case "PASSWORD POLICIES":
		return ObjectTypePasswordPolicy
	default:
		return ObjectType(s[:len(s)-1])
	}
}

func (o ObjectType) String() string {
	return string(o)
}

func (o ObjectType) Plural() string {
	// only care about the "ies" endings.
	switch o {
	case ObjectTypeMaskingPolicy:
		return "MASKING POLICIES"
	case ObjectTypeNetworkPolicy:
		return "NETWORK POLICIES"
	case ObjectTypePasswordPolicy:
		return "PASSWORD POLICIES"
	default:
		return o.String() + "S"
	}
}

// GetObjectIdentifier returns the ObjectIdentifier for the ObjectType and fully qualified name.
func (o ObjectType) GetObjectIdentifier(fullyQualifiedName string) ObjectIdentifier {
	accountIdentifiers := []ObjectType{
		ObjectTypeAccountParameter,
		ObjectTypeDatabase,
		ObjectTypeFailoverGroup,
		ObjectTypeIntegration,
		ObjectTypeResourceMonitor,
		ObjectTypeRole,
		ObjectTypeShare,
		ObjectTypeUser,
		ObjectTypeWarehouse,
	}
	if slices.Contains(accountIdentifiers, o) {
		return NewAccountObjectIdentifier(fullyQualifiedName)
	}
	parts := strings.Split(fullyQualifiedName, ".")
	dbName := parts[0]
	if o == ObjectTypeSchema {
		schemaName := strings.Join(parts[1:], ".")
		return NewSchemaIdentifier(dbName, schemaName)
	}
	schemaName := parts[1]
	objectName := strings.Join(parts[2:], ".")
	return NewSchemaObjectIdentifier(dbName, schemaName, objectName)
}
