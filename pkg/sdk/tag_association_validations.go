package sdk

var (
	// based on https://docs.snowflake.com/en/user-guide/object-tagging.html#supported-objects
	TagAssociationAllowedObjectTypes = []ObjectType{
		ObjectTypeAccount,
		ObjectTypeApplication,
		ObjectTypeApplicationPackage,
		ObjectTypeDatabase,
		ObjectTypeIntegration,
		ObjectTypeNetworkPolicy,
		ObjectTypeRole,
		ObjectTypeShare,
		ObjectTypeUser,
		ObjectTypeWarehouse,
		ObjectTypeDatabaseRole,
		ObjectTypeSchema,
		ObjectTypeAlert,
		ObjectTypeExternalFunction,
		ObjectTypeExternalTable,
		ObjectTypeGitRepository,
		ObjectTypeIcebergTable,
		ObjectTypeMaterializedView,
		ObjectTypePipe,
		ObjectTypeMaskingPolicy,
		ObjectTypePasswordPolicy,
		ObjectTypeRowAccessPolicy,
		ObjectTypeSessionPolicy,
		ObjectTypeProcedure,
		ObjectTypeStage,
		ObjectTypeStream,
		ObjectTypeTable,
		ObjectTypeTask,
		ObjectTypeView,
		ObjectTypeColumn,
		ObjectTypeEventTable,
	}
	// TODO(SNOW-1229218): Object types should be able tell their id structure and tagAssociationAllowedObjectTypes should be used to filter correct object types.
	TagAssociationTagObjectTypeIsSchemaObjectType = []ObjectType{
		ObjectTypeAlert,
		ObjectTypeExternalFunction,
		ObjectTypeExternalTable,
		ObjectTypeGitRepository,
		ObjectTypeIcebergTable,
		ObjectTypeMaterializedView,
		ObjectTypeColumn, // As a workaround for object_name can be specified as `table_name.column_name`
		ObjectTypePipe,
		ObjectTypeProcedure,
		ObjectTypeStage,
		ObjectTypeStream,
		ObjectTypeTable,
		ObjectTypeTask,
		ObjectTypeView,
		ObjectTypeEventTable,
	}
	TagAssociationAllowedObjectTypesString = make([]string, len(TagAssociationAllowedObjectTypes))
)

func init() {
	for i, v := range TagAssociationAllowedObjectTypes {
		TagAssociationAllowedObjectTypesString[i] = v.String()
	}
}
