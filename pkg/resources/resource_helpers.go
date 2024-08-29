package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringAttributeCreate(d *schema.ResourceData, key string, createField **string) {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.String(v.(string))
	}
}

func boolAttributeCreate(d *schema.ResourceData, key string, createField **bool) {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Bool(v.(bool))
	}
}

func intAttributeCreate(d *schema.ResourceData, key string, createField **int) {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Int(v.(int))
	}
}

// TODO: NewAccountObjectIdentifierFromFullyQualifiedName versus one of the new functions
func accountObjectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.AccountObjectIdentifier) {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(v.(string)))
	}
}
