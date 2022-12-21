package snowflake

// ObjectType is the type of object.
type ObjectType string

const (
	ObjectTypeDatabase         ObjectType = "DATABASE"
	ObjectTypeSchema           ObjectType = "SCHEMA"
	ObjectTypeTable            ObjectType = "TABLE"
	ObjectTypeReplicationGroup ObjectType = "REPLICATION GROUP"
	ObjectTypeFailoverGroup    ObjectType = "FAILOVER GROUP"
	ObjectTypeWarehouse        ObjectType = "WAREHOUSE"
	ObjectTypePipe             ObjectType = "PIPE"
	ObjectTypeUser             ObjectType = "USER"
	ObjectTypeShare            ObjectType = "SHARE"
	ObjectTypeTask             ObjectType = "TASK"
)

func (o ObjectType) String() string {
	return string(o)
}
