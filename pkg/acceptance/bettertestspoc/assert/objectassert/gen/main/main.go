//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func main() {
	gencommons.NewGenerator(
		getSdkObjectDetails,
		gen.ModelFromSdkObjectDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

type SdkObjectDef struct {
	idType       string
	objectType   sdk.ObjectType
	objectStruct any
}

func getSdkObjectDetails() []gencommons.SdkObjectDetails {
	allSdkObjectsDetails := make([]gencommons.SdkObjectDetails, len(allStructs))
	for idx, d := range allStructs {
		structDetails := gencommons.ExtractStructDetails(d.objectStruct)
		allSdkObjectsDetails[idx] = gencommons.SdkObjectDetails{
			IdType:        d.idType,
			ObjectType:    d.objectType,
			StructDetails: structDetails,
		}
	}
	return allSdkObjectsDetails
}

func getFilename(_ gencommons.SdkObjectDetails, model gen.SnowflakeObjectAssertionsModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_snowflake" + "_gen.go"
}

var allStructs = []SdkObjectDef{
	{
		idType:       "sdk.AccountObjectIdentifier",
		objectType:   sdk.ObjectTypeUser,
		objectStruct: sdk.User{},
	},
	{
		idType:       "sdk.AccountObjectIdentifier",
		objectType:   sdk.ObjectTypeWarehouse,
		objectStruct: sdk.Warehouse{},
	},
}
