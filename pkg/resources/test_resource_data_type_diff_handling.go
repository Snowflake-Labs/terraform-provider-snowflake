package resources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testResourceDataTypeDiffHandlingSchema = map[string]*schema.Schema{
	"env_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Used to make the tests faster (instead of communicating with SF, we read from environment variable).",
	},
	"return_data_type": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "An example field being a data type.",
		DiffSuppressFunc: DiffSuppressDataTypes,
		ValidateDiagFunc: IsDataTypeValid,
	},
}

func TestResourceDataTypeDiffHandling() *schema.Resource {
	return &schema.Resource{
		CreateContext: TestResourceCreateDataTypeDiffHandling,
		UpdateContext: TestResourceUpdateDataTypeDiffHandling,
		ReadContext:   TestResourceReadDataTypeDiffHandling(true),
		DeleteContext: TestResourceDeleteDataTypeDiffHandling,

		Schema: testResourceDataTypeDiffHandlingSchema,
	}
}

func TestResourceCreateDataTypeDiffHandling(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	dataType, err := readDatatypeCommon(d, "return_data_type")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] correctly parsed data type %v, new sql: %s", dataType, sqlNew(dataType))

	d.SetId(d.Get("env_name").(string))
	return TestResourceReadDataTypeDiffHandling(false)(ctx, d, meta)
}

func TestResourceUpdateDataTypeDiffHandling(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// it seems that it can be a no-op as we don't care about making changes but the changes recognition
	log.Printf("[DEBUG] handling update")
	return TestResourceReadDataTypeDiffHandling(false)(ctx, d, meta)
}

func TestResourceReadDataTypeDiffHandling(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if withExternalChangesMarking {
			value := oswrapper.Getenv(d.Id())
			if value != "" {
				externalDataType, err := datatypes.ParseDataType(value)
				if err != nil {
					return diag.FromErr(err)
				}
				currentConfigDataType, err := readDatatypeCommon(d, "return_data_type")
				if err != nil {
					return diag.FromErr(err)
				}
				if datatypes.AreDefinitelyDifferent(currentConfigDataType, externalDataType) {
					if err := d.Set("return_data_type", value); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
		return nil
	}
}

func TestResourceDeleteDataTypeDiffHandling(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	d.SetId("")
	return nil
}
