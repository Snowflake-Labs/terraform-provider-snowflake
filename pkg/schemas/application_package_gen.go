// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowApplicationPackageSchema represents output of SHOW query for the single ApplicationPackage.
var ShowApplicationPackageSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_default": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_current": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"distribution": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"retention_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"options": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"dropped_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"application_class": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowApplicationPackageSchema

func ApplicationPackageToSchema(applicationPackage *sdk.ApplicationPackage) map[string]any {
	applicationPackageSchema := make(map[string]any)
	applicationPackageSchema["created_on"] = applicationPackage.CreatedOn
	applicationPackageSchema["name"] = applicationPackage.Name
	applicationPackageSchema["is_default"] = applicationPackage.IsDefault
	applicationPackageSchema["is_current"] = applicationPackage.IsCurrent
	applicationPackageSchema["distribution"] = applicationPackage.Distribution
	applicationPackageSchema["owner"] = applicationPackage.Owner
	applicationPackageSchema["comment"] = applicationPackage.Comment
	applicationPackageSchema["retention_time"] = applicationPackage.RetentionTime
	applicationPackageSchema["options"] = applicationPackage.Options
	applicationPackageSchema["dropped_on"] = applicationPackage.DroppedOn
	applicationPackageSchema["application_class"] = applicationPackage.ApplicationClass
	return applicationPackageSchema
}

var _ = ApplicationPackageToSchema
