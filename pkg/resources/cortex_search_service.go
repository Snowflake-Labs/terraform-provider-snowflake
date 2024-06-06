package resources

import (
	"context"
	"log"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

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
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the Cortex search service.",
	},
	"on": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the column to use as the search column for the Cortex search service; must be a text value.",
	},
	"attributes": {
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Specifies the list of columns in the base table to enable filtering on when issuing queries to the service.",
	},
	"warehouse": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The warehouse in which to create the Cortex search service.",
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
		Create: CreateCortexSearchService,
		Read:   ReadCortexSearchService,
		Update: UpdateCortexSearchService,
		Delete: DeleteCortexSearchService,

		Schema: cortexSearchServiceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadCortexSearchServicee implements schema.ReadFunc.
func ReadCortexSearchService(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	cortexSearchService, err := client.CortexSearchServices.ShowByID(context.Background(), id)
	if err != nil {
		log.Printf("[DEBUG] cortex search service (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err := d.Set("name", cortexSearchService.Name); err != nil {
		return err
	}
	if err := d.Set("database", cortexSearchService.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema", cortexSearchService.SchemaName); err != nil {
		return err
	}
	if err := d.Set("comment", cortexSearchService.Comment); err != nil {
		return err
	}
	if err := d.Set("created_on", cortexSearchService.CreatedOn.Format(time.RFC3339)); err != nil {
		return err
	}

	return nil
}

// CreateCortexSearchService implements schema.CreateFunc.
func CreateCortexSearchService(d *schema.ResourceData, meta interface{}) error {
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
		request.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		request.WithOrReplace(true)
	}
	if v, ok := d.GetOk("if_not_exists"); ok && v.(bool) {
		request.WithIfNotExists(true)
	}
	if v, ok := d.GetOk("attributes"); ok && len(v.([]string)) > 0 {
		request.WithAttributes(v.([]string))
	}
	if err := client.CortexSearchServices.Create(context.Background(), request); err != nil {
		return err
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadCortexSearchService(d, meta)
}

// UpdateCortexSearchService implements schema.UpdateFunc.
func UpdateCortexSearchService(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
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

	if runSet {
		request.WithSet(set)
		if err := client.CortexSearchServices.Alter(ctx, request); err != nil {
			return err
		}
	}

	if d.HasChange("comment") {
		err := client.Comments.Set(ctx, &sdk.SetCommentOptions{
			ObjectType: sdk.ObjectTypeCortexSearchService,
			ObjectName: id,
			Value:      sdk.String(d.Get("comment").(string)),
		})
		if err != nil {
			return err
		}
	}
	return ReadCortexSearchService(d, meta)
}

// DeleteCortexSearchService implements schema.DeleteFunc.
func DeleteCortexSearchService(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	request := sdk.NewDropCortexSearchServiceRequest(id)

	if v, ok := d.GetOk("if exists"); ok && v.(bool) {
		request.IfExists = sdk.Bool(v.(bool))
	}

	if err := client.CortexSearchServices.Drop(context.Background(), request); err != nil {
		return err
	}
	d.SetId("")

	return nil
}
