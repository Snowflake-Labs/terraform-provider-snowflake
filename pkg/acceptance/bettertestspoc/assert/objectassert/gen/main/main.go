//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func main() {
	genhelpers.NewGenerator(
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

func getSdkObjectDetails() []genhelpers.SdkObjectDetails {
	allSdkObjectsDetails := make([]genhelpers.SdkObjectDetails, len(allStructs))
	for idx, d := range allStructs {
		structDetails := genhelpers.ExtractStructDetails(d.objectStruct)
		allSdkObjectsDetails[idx] = genhelpers.SdkObjectDetails{
			IdType:        d.idType,
			ObjectType:    d.objectType,
			StructDetails: structDetails,
		}
	}
	return allSdkObjectsDetails
}

func getFilename(_ genhelpers.SdkObjectDetails, model gen.SnowflakeObjectAssertionsModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_snowflake" + "_gen.go"
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
