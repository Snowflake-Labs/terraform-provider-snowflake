package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SdkObjectDef struct {
	IdType       string
	ObjectType   sdk.ObjectType
	ObjectStruct any
}

var allStructs = []SdkObjectDef{
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeDatabase,
		ObjectStruct: sdk.Database{},
	},
	{
		IdType:       "sdk.DatabaseObjectIdentifier",
		ObjectType:   sdk.ObjectTypeDatabaseRole,
		ObjectStruct: sdk.DatabaseRole{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeRowAccessPolicy,
		ObjectStruct: sdk.RowAccessPolicy{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeUser,
		ObjectStruct: sdk.User{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeView,
		ObjectStruct: sdk.View{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeWarehouse,
		ObjectStruct: sdk.Warehouse{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeResourceMonitor,
		ObjectStruct: sdk.ResourceMonitor{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeMaskingPolicy,
		ObjectStruct: sdk.MaskingPolicy{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeAuthenticationPolicy,
		ObjectStruct: sdk.AuthenticationPolicy{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeTask,
		ObjectStruct: sdk.Task{},
	},
}

func GetSdkObjectDetails() []genhelpers.SdkObjectDetails {
	allSdkObjectsDetails := make([]genhelpers.SdkObjectDetails, len(allStructs))
	for idx, d := range allStructs {
		structDetails := genhelpers.ExtractStructDetails(d.ObjectStruct)
		allSdkObjectsDetails[idx] = genhelpers.SdkObjectDetails{
			IdType:        d.IdType,
			ObjectType:    d.ObjectType,
			StructDetails: structDetails,
		}
	}
	return allSdkObjectsDetails
}
