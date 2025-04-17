package resources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
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
		Type:        schema.TypeString,
		Required:    true,
		Description: "An example field being a data type.",
		// TODO: implement
		//DiffSuppressFunc: DiffSuppressDataTypes,
		//ValidateDiagFunc: IsDataTypeValid,
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
	log.Printf("[DEBUG] corretly parsed data type %v", dataType)

	d.SetId(d.Get("env_name").(string))
	return TestResourceReadDataTypeDiffHandling(false)(ctx, d, meta)
}

func TestResourceUpdateDataTypeDiffHandling(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// it seems that it can be a no-op as we don't care about making changes but the changes recognition
	return TestResourceReadDataTypeDiffHandling(false)(ctx, d, meta)
}

func TestResourceReadDataTypeDiffHandling(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		value := oswrapper.Getenv(d.Id())
		if value != "" {
			if err := d.Set("return_data_type", value); err != nil {
				return diag.FromErr(err)
			}
		}

		if withExternalChangesMarking {
			// TODO: show output if needed
		}
		return nil
	}
}

func TestResourceDeleteDataTypeDiffHandling(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	d.SetId("")
	return nil
}
