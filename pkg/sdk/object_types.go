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
	ObjectTypeTag              ObjectType = "TAG"
	ObjectTypeUser             ObjectType = "USER"
	ObjectTypeWarehouse        ObjectType = "WAREHOUSE"
)

func (o ObjectType) String() string {
	return string(o)
}

func (o ObjectType) Plural() PluralObjectType {
	switch o {
	case ObjectTypeAccountParameter:
		return PluralObjectTypeAccountParameters
	case ObjectTypeDatabase:
		return PluralObjectTypeDatabases
	case ObjectTypeFailoverGroup:
		return PluralObjectTypeTypeFailoverGroups
	case ObjectTypeIntegration:
		return PluralObjectTypeIntegrations
	case ObjectTypeMaskingPolicy:
		return PluralObjectTypeMaskingPolicies
	case ObjectTypeNetworkPolicy:
		return PluralObjectTypeNetworkPolicies
	case ObjectTypePasswordPolicy:
		return PluralObjectTypePasswordPolicies
	case ObjectTypeResourceMonitor:
		return PluralObjectTypeResourceMonitors
	case ObjectTypeRole:
		return PluralObjectTypeRoles
	case ObjectTypeSchema:
		return PluralObjectTypeSchemas
	case ObjectTypeShare:
		return PluralObjectTypeShares
	case ObjectTypeTag:
		return PluralObjectTypeTags
	case ObjectTypeUser:
		return PluralObjectTypeUsers
	case ObjectTypeWarehouse:
		return PluralObjectTypeWarehouses
	default:
		return PluralObjectType("")
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

type PluralObjectType string

const (
	PluralObjectTypeAccountParameters  PluralObjectType = "ACCOUNT PARAMETERS"
	PluralObjectTypeDatabases          PluralObjectType = "DATABASES"
	PluralObjectTypeTypeFailoverGroups PluralObjectType = "FAILOVER GROUPS"
	PluralObjectTypeIntegrations       PluralObjectType = "INTEGRATIONS"
	PluralObjectTypeMaskingPolicies    PluralObjectType = "MASKING POLICIES"
	PluralObjectTypeNetworkPolicies    PluralObjectType = "NETWORK POLICIES"
	PluralObjectTypePasswordPolicies   PluralObjectType = "PASSWORD POLICIES"
	PluralObjectTypeResourceMonitors   PluralObjectType = "RESOURCE MONITORS"
	PluralObjectTypeRoles              PluralObjectType = "ROLES"
	PluralObjectTypeSchemas            PluralObjectType = "SCHEMAS"
	PluralObjectTypeShares             PluralObjectType = "SHARES"
	PluralObjectTypeTags               PluralObjectType = "TAGS"
	PluralObjectTypeUsers              PluralObjectType = "USERS"
	PluralObjectTypeWarehouses         PluralObjectType = "WAREHOUSES"
)

func (p PluralObjectType) String() string {
	return string(p)
}

func (p PluralObjectType) Singular() ObjectType {
	switch p {
	case PluralObjectTypeAccountParameters:
		return ObjectTypeAccountParameter
	case PluralObjectTypeDatabases:
		return ObjectTypeDatabase
	case PluralObjectTypeTypeFailoverGroups:
		return ObjectTypeFailoverGroup
	case PluralObjectTypeIntegrations:
		return ObjectTypeIntegration
	case PluralObjectTypeMaskingPolicies:
		return ObjectTypeMaskingPolicy
	case PluralObjectTypeNetworkPolicies:
		return ObjectTypeNetworkPolicy
	case PluralObjectTypePasswordPolicies:
		return ObjectTypePasswordPolicy
	case PluralObjectTypeResourceMonitors:
		return ObjectTypeResourceMonitor
	case PluralObjectTypeRoles:
		return ObjectTypeRole
	case PluralObjectTypeSchemas:
		return ObjectTypeSchema
	case PluralObjectTypeShares:
		return ObjectTypeShare
	case PluralObjectTypeUsers:
		return ObjectTypeUser
	case PluralObjectTypeWarehouses:
		return ObjectTypeWarehouse
	default:
		return ObjectType("")
	}
}
