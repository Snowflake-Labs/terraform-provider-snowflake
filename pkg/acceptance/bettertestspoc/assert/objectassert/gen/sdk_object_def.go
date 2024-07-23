package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
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
		ObjectType:   sdk.ObjectTypeUser,
		ObjectStruct: sdk.User{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeWarehouse,
		ObjectStruct: sdk.Warehouse{},
	},
}

func GetSdkObjectDetails() []gencommons.SdkObjectDetails {
	allSdkObjectsDetails := make([]gencommons.SdkObjectDetails, len(allStructs))
	for idx, d := range allStructs {
		structDetails := gencommons.ExtractStructDetails(d.ObjectStruct)
		allSdkObjectsDetails[idx] = gencommons.SdkObjectDetails{
			IdType:        d.IdType,
			ObjectType:    d.ObjectType,
			StructDetails: structDetails,
		}
	}
	return allSdkObjectsDetails
}
