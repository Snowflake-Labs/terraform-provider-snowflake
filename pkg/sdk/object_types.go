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
	ObjectTypeAlert            ObjectType = "ALERT"
	ObjectTypeAccountParameter ObjectType = "ACCOUNT PARAMETER"
	ObjectTypeDatabase         ObjectType = "DATABASE"
	ObjectTypeExternalTable    ObjectType = "EXTERNAL TABLE"
	ObjectTypeFailoverGroup    ObjectType = "FAILOVER GROUP"
	ObjectTypeFileFormat       ObjectType = "FILE FORMAT"
	ObjectTypeIntegration      ObjectType = "INTEGRATION"
	ObjectTypeMaskingPolicy    ObjectType = "MASKING POLICY"
	ObjectTypeNetworkPolicy    ObjectType = "NETWORK POLICY"
	ObjectTypePasswordPolicy   ObjectType = "PASSWORD POLICY"
	ObjectTypeReplicationGroup ObjectType = "REPLICATION GROUP"
	ObjectTypeResourceMonitor  ObjectType = "RESOURCE MONITOR"
	ObjectTypeRole             ObjectType = "ROLE"
	ObjectTypeSchema           ObjectType = "SCHEMA"
	ObjectTypeSessionPolicy    ObjectType = "SESSION POLICY"
	ObjectTypeShare            ObjectType = "SHARE"
	ObjectTypeTable            ObjectType = "TABLE"
	ObjectTypeTag              ObjectType = "TAG"
	ObjectTypeTask             ObjectType = "TASK"
	ObjectTypeUser             ObjectType = "USER"
	ObjectTypeWarehouse        ObjectType = "WAREHOUSE"
)

func (o ObjectType) String() string {
	return string(o)
}

func objectTypeSingularToPluralMap() map[ObjectType]PluralObjectType {
	return map[ObjectType]PluralObjectType{
		ObjectTypeAccountParameter: PluralObjectTypeAccountParameters,
		ObjectTypeDatabase:         PluralObjectTypeDatabases,
		ObjectTypeFailoverGroup:    PluralObjectTypeTypeFailoverGroups,
		ObjectTypeIntegration:      PluralObjectTypeIntegrations,
		ObjectTypeMaskingPolicy:    PluralObjectTypeMaskingPolicies,
		ObjectTypeNetworkPolicy:    PluralObjectTypeNetworkPolicies,
		ObjectTypePasswordPolicy:   PluralObjectTypePasswordPolicies,
		ObjectTypeReplicationGroup: PluralObjectTypeReplicationGroups,
		ObjectTypeResourceMonitor:  PluralObjectTypeResourceMonitors,
		ObjectTypeRole:             PluralObjectTypeRoles,
		ObjectTypeSchema:           PluralObjectTypeSchemas,
		ObjectTypeSessionPolicy:    PluralObjectTypeSessionPolicies,
		ObjectTypeShare:            PluralObjectTypeShares,
		ObjectTypeTable:            PluralObjectTypeTables,
		ObjectTypeTag:              PluralObjectTypeTags,
		ObjectTypeTask:             PluralObjectTypeTasks,
		ObjectTypeUser:             PluralObjectTypeUsers,
		ObjectTypeWarehouse:        PluralObjectTypeWarehouses,
	}
}

func pluralObjectTypeToSingularMap() map[PluralObjectType]ObjectType {
	m := make(map[PluralObjectType]ObjectType)
	for k, v := range objectTypeSingularToPluralMap() {
		m[v] = k
	}
	return m
}

func (o ObjectType) Plural() PluralObjectType {
	if plural, ok := objectTypeSingularToPluralMap()[o]; ok {
		return plural
	}
	return PluralObjectType(o + "S")
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
	PluralObjectTypeReplicationGroups  PluralObjectType = "REPLICATION GROUPS"
	PluralObjectTypeResourceMonitors   PluralObjectType = "RESOURCE MONITORS"
	PluralObjectTypeRoles              PluralObjectType = "ROLES"
	PluralObjectTypeSchemas            PluralObjectType = "SCHEMAS"
	PluralObjectTypeSessionPolicies    PluralObjectType = "SESSION POLICIES"
	PluralObjectTypeShares             PluralObjectType = "SHARES"
	PluralObjectTypeTables             PluralObjectType = "TABLES"
	PluralObjectTypeTags               PluralObjectType = "TAGS"
	PluralObjectTypeTasks              PluralObjectType = "TASKS"
	PluralObjectTypeUsers              PluralObjectType = "USERS"
	PluralObjectTypeWarehouses         PluralObjectType = "WAREHOUSES"
)

func (p PluralObjectType) String() string {
	return string(p)
}

func (p PluralObjectType) Singular() ObjectType {
	if singular, ok := pluralObjectTypeToSingularMap()[p]; ok {
		return singular
	}
	return ObjectType(strings.TrimSuffix(string(p), "S"))
}
