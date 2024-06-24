package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var cortexSearchServiceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the Cortex search service. The name must be unique for the schema in which the service is created.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the Cortex search service.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the Cortex search service.",
		ForceNew:    true,
	},
	"on": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the column to use as the search column for the Cortex search service; must be a text value.",
		ForceNew:    true,
	},
	"attributes": {
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Specifies the list of columns in the base table to enable filtering on when issuing queries to the service.",
	},
	"warehouse": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The warehouse in which to create the Cortex search service.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"target_lag": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the maximum target lag time for the Cortex search service.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the Cortex search service.",
	},
	"query": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the query to use to populate the Cortex search service.",
		DiffSuppressFunc: DiffSuppressStatement,
	},
}

// CortexSearchService returns a pointer to the resource representing a Cortex search service.
func CortexSearchService() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateCortexSearchService,
		ReadContext:   ReadCortexSearchService,
		UpdateContext: UpdateCortexSearchService,
		DeleteContext: DeleteCortexSearchService,

		Schema: cortexSearchServiceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadCortexSearchServicee implements schema.ReadFunc.
func ReadCortexSearchService(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	cortexSearchService, err := client.CortexSearchServices.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			log.Printf("[DEBUG] cortex search service (%s) not found", d.Id())
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query cortex search service. Marking the resource as removed.",
					Detail:   fmt.Sprintf("CortexSearchServiceName: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set("name", cortexSearchService.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database", cortexSearchService.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", cortexSearchService.SchemaName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", cortexSearchService.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_on", cortexSearchService.CreatedOn.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// CreateCortexSearchService implements schema.CreateFunc.
func CreateCortexSearchService(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	on := d.Get("on").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	warehouse := sdk.NewAccountObjectIdentifier(d.Get("warehouse").(string))
	target_lag := d.Get("target_lag").(string)
	query := d.Get("query").(string)

	request := sdk.NewCreateCortexSearchServiceRequest(id, on, warehouse, target_lag, query)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}
	if v, ok := d.GetOk("attributes"); ok && len(v.([]string)) > 0 {
		attributes := sdk.AttributesRequest{
			Columns: v.([]string),
		}
		request.WithAttributes(attributes)
	}
	var diags diag.Diagnostics
	if err := client.CortexSearchServices.Create(ctx, request); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return append(diags, ReadCortexSearchService(ctx, d, meta)...)
}

// UpdateCortexSearchService implements schema.UpdateFunc.
func UpdateCortexSearchService(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	request := sdk.NewAlterCortexSearchServiceRequest(id)

	runSet := false
	set := sdk.NewCortexSearchServiceSetRequest()
	if d.HasChange("target_lag") {
		tl := d.Get("target_lag").(string)
		set.WithTargetLag(tl)
		runSet = true
	}

	if d.HasChange("warehouse") {
		warehouseName := d.Get("warehouse").(string)
		set.WithWarehouse(sdk.NewAccountObjectIdentifier(warehouseName))
		runSet = true
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		set.WithComment(comment)
		runSet = true
	}

	var diags diag.Diagnostics
	if runSet {
		request.WithSet(*set)
		if err := client.CortexSearchServices.Alter(ctx, request); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return append(diags, ReadCortexSearchService(ctx, d, meta)...)
}

// DeleteCortexSearchService implements schema.DeleteFunc.
func DeleteCortexSearchService(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	request := sdk.NewDropCortexSearchServiceRequest(id).WithIfExists(true)

	if err := client.CortexSearchServices.Drop(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}
