package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var materializedViewSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the view; must be unique for the schema in which the view is created.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the view. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the view. Don't use the | character.",
		ForceNew:    true,
	},
	"warehouse": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The warehouse name.",
		ForceNew:    true,
	},
	"or_replace": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Overwrites the View if it exists.",
	},
	"is_secure": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies that the view is secure.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the view.",
	},
	"statement": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the query used to create the view.",
		ForceNew:         true,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"tag": tagReferenceSchema,
}

// MaterializedView returns a pointer to the resource representing a view.
func MaterializedView() *schema.Resource {
	return &schema.Resource{
		Create: CreateMaterializedView,
		Read:   ReadMaterializedView,
		Update: UpdateMaterializedView,
		Delete: DeleteMaterializedView,

		Schema: materializedViewSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateMaterializedView implements schema.CreateFunc.
func CreateMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	s := d.Get("statement").(string)
	createRequest := sdk.NewCreateMaterializedViewRequest(id, s)

	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		createRequest.WithOrReplace(sdk.Bool(true))
	}

	if v, ok := d.GetOk("is_secure"); ok && v.(bool) {
		createRequest.WithSecure(sdk.Bool(true))
	}

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	warehouseName := d.Get("warehouse").(string)
	// TODO [SNOW-867235]: this was the old implementation, it's left for now, we will address this with resources rework discussions
	err := client.Sessions.UseWarehouse(ctx, sdk.NewAccountObjectIdentifier(warehouseName))
	if err != nil {
		return fmt.Errorf("error setting warehouse %s while creating materialized view %v err = %w", warehouseName, name, err)
	}

	err = client.MaterializedViews.Create(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("error creating materialized view %v err = %w", name, err)
	}

	// TODO [SNOW-867235]: we have to set tags after creation because existing materialized view extractor is not aware of TAG during CREATE
	// Will be discussed with parser topic during resources redesign.
	if _, ok := d.GetOk("tag"); ok {
		err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetTags(getPropertyTags(d, "tag")))
		if err != nil {
			return fmt.Errorf("error setting tags on materialized view %v, err = %w", id, err)
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadMaterializedView(d, meta)
}

// ReadMaterializedView implements schema.ReadFunc.
func ReadMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	materializedView, err := client.MaterializedViews.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] materialized view (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("name", materializedView.Name); err != nil {
		return err
	}

	if err := d.Set("is_secure", materializedView.IsSecure); err != nil {
		return err
	}

	if err := d.Set("comment", materializedView.Comment); err != nil {
		return err
	}

	if err := d.Set("schema", materializedView.SchemaName); err != nil {
		return err
	}

	if err := d.Set("database", materializedView.DatabaseName); err != nil {
		return err
	}

	// TODO [SNOW-867235]: what do we do with these extractors (added as discussion topic)?
	// Want to only capture the SELECT part of the query because before that is the CREATE part of the view.
	extractor := snowflake.NewViewSelectStatementExtractor(materializedView.Text)
	substringOfQuery, err := extractor.ExtractMaterializedView()
	if err != nil {
		return err
	}

	if err := d.Set("statement", substringOfQuery); err != nil {
		return err
	}

	return nil
}

// UpdateMaterializedView implements schema.UpdateFunc.
func UpdateMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("name") {
		newName := d.Get("name").(string)

		newId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), newName)

		err := client.MaterializedViews.Alter(ctx, sdk.NewAlterMaterializedViewRequest(id).WithRenameTo(&newId))
		if err != nil {
			return fmt.Errorf("error renaming materialized view %v err = %w", d.Id(), err)
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	var runSetStatement bool
	var runUnsetStatement bool
	setRequest := sdk.NewMaterializedViewSetRequest()
	unsetRequest := sdk.NewMaterializedViewUnsetRequest()

	if d.HasChange("comment") {
		comment := d.Get("comment")
		if c := comment.(string); c == "" {
			runUnsetStatement = true
			unsetRequest.WithComment(sdk.Bool(true))
		} else {
			runSetStatement = true
			setRequest.WithComment(sdk.String(d.Get("comment").(string)))
		}
	}
	if d.HasChange("is_secure") {
		if d.Get("is_secure").(bool) {
			runSetStatement = true
			setRequest.WithSecure(sdk.Bool(true))
		} else {
			runUnsetStatement = true
			unsetRequest.WithSecure(sdk.Bool(true))
		}
	}

	if runSetStatement {
		err := client.MaterializedViews.Alter(ctx, sdk.NewAlterMaterializedViewRequest(id).WithSet(setRequest))
		if err != nil {
			return fmt.Errorf("error updating materialized view: %w", err)
		}
	}

	if runUnsetStatement {
		err := client.MaterializedViews.Alter(ctx, sdk.NewAlterMaterializedViewRequest(id).WithUnset(unsetRequest))
		if err != nil {
			return fmt.Errorf("error updating materialized view: %w", err)
		}
	}

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			// TODO [SNOW-1022645]: view is used on purpose here; change after we have an agreement on situations like this in the SDK
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetTags(unsetTags))
			if err != nil {
				return fmt.Errorf("error unsetting tags on %v, err = %w", d.Id(), err)
			}
		}

		if len(setTags) > 0 {
			// TODO [SNOW-1022645]: view is used on purpose here; change after we have an agreement on situations like this in the SDK
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetTags(setTags))
			if err != nil {
				return fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err)
			}
		}
	}

	return ReadMaterializedView(d, meta)
}

// DeleteMaterializedView implements schema.DeleteFunc.
func DeleteMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.MaterializedViews.Drop(ctx, sdk.NewDropMaterializedViewRequest(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
