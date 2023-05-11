package sdk

// Object bundles together the object type and name. Its used for DDL statements.
type Object struct {
	ObjectType ObjectType       `ddl:"keyword"`
	Name       ObjectIdentifier `ddl:"identifier"`
}

// ObjectType is the type of object.
type ObjectType string

const (
	ObjectTypeAccountParameter ObjectType = "ACCOUNT PARAMETER"
	ObjectTypeDatabase         ObjectType = "DATABASE"
	ObjectTypeFailoverGroup    ObjectType = "FAILOVER GROUP"
	ObjectTypeIntegration      ObjectType = "INTEGRATION"
	ObjectTypeMaskingPolicy    ObjectType = "MASKING POLICY"
	ObjectTypeNetworkPolicy     ObjectType = "NETWORK POLICY"
	ObjectTypePasswordPolicy   ObjectType = "PASSWORD POLICY"
	ObjectTypeResourceMonitor  ObjectType = "RESOURCE MONITOR"
	ObjectTypeRole             ObjectType = "ROLE"
	ObjectTypeShare            ObjectType = "SHARE"
	ObjectTypeUser             ObjectType = "USER"
	ObjectTypeWarehouse        ObjectType = "WAREHOUSE"
)

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
