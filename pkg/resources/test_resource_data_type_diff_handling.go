package resources

import (
	"context"
	"log"
	"os"

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
		CreateContext: TestResourceDataTypeDiffHandlingCreate,
		UpdateContext: TestResourceDataTypeDiffHandlingUpdate,
		ReadContext:   TestResourceDataTypeDiffHandlingRead(true),
		DeleteContext: TestResourceDataTypeDiffHandlingDelete,

		Schema: testResourceDataTypeDiffHandlingSchema,
	}
}

func TestResourceDataTypeDiffHandlingCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Get("env_name").(string)
	log.Printf("[DEBUG] handling create for %s", envName)

	dataType, err := readDatatypeCommon(d, "return_data_type")
	if err != nil {
		return diag.FromErr(err)
	}
	if err := testResourceDataTypeDiffHandlingSet(envName, dataType); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(envName)
	return TestResourceDataTypeDiffHandlingRead(false)(ctx, d, meta)
}

func TestResourceDataTypeDiffHandlingUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling update for %s", envName)

	if d.HasChange("return_data_type") {
		dataType, err := readChangedDatatypeCommon(d, "return_data_type")
		if err != nil {
			return diag.FromErr(err)
		}
		if err := testResourceDataTypeDiffHandlingSet(envName, dataType); err != nil {
			return diag.FromErr(err)
		}
	}

	return TestResourceDataTypeDiffHandlingRead(false)(ctx, d, meta)
}

func TestResourceDataTypeDiffHandlingRead(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		envName := d.Id()
		log.Printf("[DEBUG] handling read for %s, with marking external changes: %t", envName, withExternalChangesMarking)

		value := oswrapper.Getenv(envName)
		log.Printf("[DEBUG] env %s value is `%s`", envName, value)
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

		if withExternalChangesMarking {
			// TODO: implement if needed
		}
		return nil
	}
}

func TestResourceDataTypeDiffHandlingDelete(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling delete for %s", envName)

	if err := testResourceDataTypeDiffHandlingUnset(envName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func testResourceDataTypeDiffHandlingSet(envName string, dataType datatypes.DataType) error {
	log.Printf("[DEBUG] setting env %s to value `%s`", envName, dataType.ToSql())
	return os.Setenv(envName, dataType.ToSql())
}

func testResourceDataTypeDiffHandlingUnset(envName string) error {
	log.Printf("[DEBUG] unsetting env %s", envName)
	return os.Setenv(envName, "")
}
