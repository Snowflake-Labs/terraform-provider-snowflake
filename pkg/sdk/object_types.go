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
	ObjectTypeAccount            ObjectType = "ACCOUNT"
	ObjectTypeManagedAccount     ObjectType = "MANAGED ACCOUNT"
	ObjectTypeUser               ObjectType = "USER"
	ObjectTypeDatabaseRole       ObjectType = "DATABASE ROLE"
	ObjectTypeRole               ObjectType = "ROLE"
	ObjectTypeIntegration        ObjectType = "INTEGRATION"
	ObjectTypeNetworkPolicy      ObjectType = "NETWORK POLICY"
	ObjectTypePasswordPolicy     ObjectType = "PASSWORD POLICY"
	ObjectTypeSessionPolicy      ObjectType = "SESSION POLICY"
	ObjectTypeReplicationGroup   ObjectType = "REPLICATION GROUP"
	ObjectTypeFailoverGroup      ObjectType = "FAILOVER GROUP"
	ObjectTypeConnection         ObjectType = "CONNECTION"
	ObjectTypeParameter          ObjectType = "PARAMETER"
	ObjectTypeWarehouse          ObjectType = "WAREHOUSE"
	ObjectTypeResourceMonitor    ObjectType = "RESOURCE MONITOR"
	ObjectTypeDatabase           ObjectType = "DATABASE"
	ObjectTypeSchema             ObjectType = "SCHEMA"
	ObjectTypeShare              ObjectType = "SHARE"
	ObjectTypeTable              ObjectType = "TABLE"
	ObjectTypeDynamicTable       ObjectType = "DYNAMIC TABLE"
	ObjectTypeExternalTable      ObjectType = "EXTERNAL TABLE"
	ObjectTypeEventTable         ObjectType = "EVENT TABLE"
	ObjectTypeView               ObjectType = "VIEW"
	ObjectTypeMaterializedView   ObjectType = "MATERIALIZED VIEW"
	ObjectTypeSequence           ObjectType = "SEQUENCE"
	ObjectTypeFunction           ObjectType = "FUNCTION"
	ObjectTypeExternalFunction   ObjectType = "EXTERNAL FUNCTION"
	ObjectTypeProcedure          ObjectType = "PROCEDURE"
	ObjectTypeStream             ObjectType = "STREAM"
	ObjectTypeTask               ObjectType = "TASK"
	ObjectTypeMaskingPolicy      ObjectType = "MASKING POLICY"
	ObjectTypeRowAccessPolicy    ObjectType = "ROW ACCESS POLICY"
	ObjectTypeTag                ObjectType = "TAG"
	ObjectTypeSecret             ObjectType = "SECRET"
	ObjectTypeStage              ObjectType = "STAGE"
	ObjectTypeFileFormat         ObjectType = "FILE FORMAT"
	ObjectTypePipe               ObjectType = "PIPE"
	ObjectTypeAlert              ObjectType = "ALERT"
	ObjectTypeApplication        ObjectType = "APPLICATION"
	ObjectTypeApplicationPackage ObjectType = "APPLICATION PACKAGE"
	ObjectTypeApplicationRole    ObjectType = "APPLICATION ROLE"
	ObjectTypeStreamlit          ObjectType = "STREAMLIT"
)

func (o ObjectType) String() string {
	return string(o)
}

func objectTypeSingularToPluralMap() map[ObjectType]PluralObjectType {
	return map[ObjectType]PluralObjectType{
		ObjectTypeAccount:            PluralObjectTypeAccounts,
		ObjectTypeManagedAccount:     PluralObjectTypeManagedAccounts,
		ObjectTypeUser:               PluralObjectTypeUsers,
		ObjectTypeDatabaseRole:       PluralObjectTypeDatabaseRoles,
		ObjectTypeRole:               PluralObjectTypeRoles,
		ObjectTypeIntegration:        PluralObjectTypeIntegrations,
		ObjectTypeNetworkPolicy:      PluralObjectTypeNetworkPolicies,
		ObjectTypePasswordPolicy:     PluralObjectTypePasswordPolicies,
		ObjectTypeSessionPolicy:      PluralObjectTypeSessionPolicies,
		ObjectTypeReplicationGroup:   PluralObjectTypeReplicationGroups,
		ObjectTypeFailoverGroup:      PluralObjectTypeFailoverGroups,
		ObjectTypeConnection:         PluralObjectTypeConnections,
		ObjectTypeParameter:          PluralObjectTypeParameters,
		ObjectTypeWarehouse:          PluralObjectTypeWarehouses,
		ObjectTypeResourceMonitor:    PluralObjectTypeResourceMonitors,
		ObjectTypeDatabase:           PluralObjectTypeDatabases,
		ObjectTypeSchema:             PluralObjectTypeSchemas,
		ObjectTypeShare:              PluralObjectTypeShares,
		ObjectTypeTable:              PluralObjectTypeTables,
		ObjectTypeDynamicTable:       PluralObjectTypeDynamicTables,
		ObjectTypeExternalTable:      PluralObjectTypeExternalTables,
		ObjectTypeEventTable:         PluralObjectTypeEventTables,
		ObjectTypeView:               PluralObjectTypeViews,
		ObjectTypeMaterializedView:   PluralObjectTypeMaterializedViews,
		ObjectTypeSequence:           PluralObjectTypeSequences,
		ObjectTypeFunction:           PluralObjectTypeFunctions,
		ObjectTypeExternalFunction:   PluralObjectTypeExternalFunctions,
		ObjectTypeProcedure:          PluralObjectTypeProcedures,
		ObjectTypeStream:             PluralObjectTypeStreams,
		ObjectTypeTask:               PluralObjectTypeTasks,
		ObjectTypeMaskingPolicy:      PluralObjectTypeMaskingPolicies,
		ObjectTypeRowAccessPolicy:    PluralObjectTypeRowAccessPolicies,
		ObjectTypeTag:                PluralObjectTypeTags,
		ObjectTypeSecret:             PluralObjectTypeSecrets,
		ObjectTypeStage:              PluralObjectTypeStages,
		ObjectTypeFileFormat:         PluralObjectTypeFileFormats,
		ObjectTypePipe:               PluralObjectTypePipes,
		ObjectTypeAlert:              PluralObjectTypeAlerts,
		ObjectTypeApplication:        PluralObjectTypeApplications,
		ObjectTypeApplicationPackage: PluralObjectTypeApplicationPackages,
		ObjectTypeApplicationRole:    PluralObjectTypeApplicationRoles,
		ObjectTypeStreamlit:          PluralObjectTypeStreamlits,
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
	accountObjectIdentifiers := []ObjectType{
		ObjectTypeParameter,
		ObjectTypeDatabase,
		ObjectTypeFailoverGroup,
		ObjectTypeIntegration,
		ObjectTypeResourceMonitor,
		ObjectTypeRole,
		ObjectTypeShare,
		ObjectTypeUser,
		ObjectTypeWarehouse,
	}
	if slices.Contains(accountObjectIdentifiers, o) {
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
	PluralObjectTypeAccounts            = "ACCOUNTS"
	PluralObjectTypeManagedAccounts     = "MANAGED ACCOUNTS"
	PluralObjectTypeUsers               = "USERS"
	PluralObjectTypeDatabaseRoles       = "DATABASE ROLES"
	PluralObjectTypeRoles               = "ROLES"
	PluralObjectTypeIntegrations        = "INTEGRATIONS"
	PluralObjectTypeNetworkPolicies     = "NETWORK POLICIES"
	PluralObjectTypePasswordPolicies    = "PASSWORD POLICIES"
	PluralObjectTypeSessionPolicies     = "SESSION POLICIES"
	PluralObjectTypeReplicationGroups   = "REPLICATION GROUPS"
	PluralObjectTypeFailoverGroups      = "FAILOVER GROUPS"
	PluralObjectTypeConnections         = "CONNECTIONS"
	PluralObjectTypeParameters          = "PARAMETERS"
	PluralObjectTypeWarehouses          = "WAREHOUSES"
	PluralObjectTypeResourceMonitors    = "RESOURCE MONITORS"
	PluralObjectTypeDatabases           = "DATABASES"
	PluralObjectTypeSchemas             = "SCHEMAS"
	PluralObjectTypeShares              = "SHARES"
	PluralObjectTypeTables              = "TABLES"
	PluralObjectTypeDynamicTables       = "DYNAMIC TABLES"
	PluralObjectTypeExternalTables      = "EXTERNAL TABLES"
	PluralObjectTypeEventTables         = "EVENT TABLES"
	PluralObjectTypeViews               = "VIEWS"
	PluralObjectTypeMaterializedViews   = "MATERIALIZED VIEWS"
	PluralObjectTypeSequences           = "SEQUENCES"
	PluralObjectTypeFunctions           = "FUNCTIONS"
	PluralObjectTypeExternalFunctions   = "EXTERNAL FUNCTIONS"
	PluralObjectTypeProcedures          = "PROCEDURES"
	PluralObjectTypeStreams             = "STREAMS"
	PluralObjectTypeTasks               = "TASKS"
	PluralObjectTypeMaskingPolicies     = "MASKING POLICIES"
	PluralObjectTypeRowAccessPolicies   = "ROW ACCESS POLICIES"
	PluralObjectTypeTags                = "TAGS"
	PluralObjectTypeSecrets             = "SECRETS"
	PluralObjectTypeStages              = "STAGES"
	PluralObjectTypeFileFormats         = "FILE FORMATS"
	PluralObjectTypePipes               = "PIPES"
	PluralObjectTypeAlerts              = "ALERTS"
	PluralObjectTypeApplications        = "APPLICATIONS"
	PluralObjectTypeApplicationPackages = "APPLICATION PACKAGES"
	PluralObjectTypeApplicationRoles    = "APPLICATION ROLES"
	PluralObjectTypeStreamlits          = "STREAMLITS"
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
