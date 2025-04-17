package resources

import (
	"context"

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
	// TODO: implement

	d.SetId(d.Get("env_name").(string))
	return ReadObjectRenamingListsAndSets(false)(ctx, d, meta)
}

func TestResourceUpdateDataTypeDiffHandling(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// TODO: implement

	return ReadObjectRenamingListsAndSets(false)(ctx, d, meta)
}

func TestResourceReadDataTypeDiffHandling(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		// TODO: implement

		return nil
	}
}

func TestResourceDeleteDataTypeDiffHandling(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	d.SetId("")
	return nil
}
