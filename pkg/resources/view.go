package resources

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var space = regexp.MustCompile(`\s+`)

var viewSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the view; must be unique for the schema in which the view is created. Don't use the | character.",
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
	"or_replace": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Overwrites the View if it exists.",
	},
	// TODO [SNOW-1348118: this is used only during or_replace, we would like to change the behavior before v1
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Retains the access permissions from the original view when a new view is created using the OR REPLACE clause. OR REPLACE must be set when COPY GRANTS is set.",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return oldValue != "" && oldValue != newValue
		},
		RequiredWith: []string{"or_replace"},
	},
	"is_secure": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies that the view is secure. By design, the Snowflake's `SHOW VIEWS` command does not provide information about secure views (consult [view usage notes](https://docs.snowflake.com/en/sql-reference/sql/create-view#usage-notes)) which is essential to manage/import view with Terraform. Use the role owning the view while managing secure views.",
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
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The timestamp at which the view was created.",
	},
	"tag": tagReferenceSchema,
}

// View returns a pointer to the resource representing a view.
func View() *schema.Resource {
	return &schema.Resource{
		Create: CreateView,
		Read:   ReadView,
		Update: UpdateView,
		Delete: DeleteView,

		Schema: viewSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateView implements schema.CreateFunc.
func CreateView(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	s := d.Get("statement").(string)
	createRequest := sdk.NewCreateViewRequest(id, s)

	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		createRequest.WithOrReplace(true)
	}

	if v, ok := d.GetOk("is_secure"); ok && v.(bool) {
		createRequest.WithSecure(true)
	}

	if v, ok := d.GetOk("copy_grants"); ok && v.(bool) {
		createRequest.WithCopyGrants(true)
	}

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(v.(string))
	}

	err := client.Views.Create(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("error creating view %v err = %w", name, err)
	}

	// TODO [SNOW-1348118]: we have to set tags after creation because existing view extractor is not aware of TAG during CREATE
	// Will be discussed with parser topic during resources redesign.
	if _, ok := d.GetOk("tag"); ok {
		err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetTags(getPropertyTags(d, "tag")))
		if err != nil {
			return fmt.Errorf("error setting tags on view %v, err = %w", id, err)
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadView(d, meta)
}

// ReadView implements schema.ReadFunc.
func ReadView(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	view, err := client.Views.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] view (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err = d.Set("name", view.Name); err != nil {
		return err
	}
	if err = d.Set("is_secure", view.IsSecure); err != nil {
		return err
	}
	if err = d.Set("copy_grants", view.HasCopyGrants()); err != nil {
		return err
	}
	if err = d.Set("comment", view.Comment); err != nil {
		return err
	}
	if err = d.Set("schema", view.SchemaName); err != nil {
		return err
	}
	if err = d.Set("database", view.DatabaseName); err != nil {
		return err
	}
	if err = d.Set("created_on", view.CreatedOn); err != nil {
		return err
	}

	if view.Text != "" {
		// Want to only capture the SELECT part of the query because before that is the CREATE part of the view.
		extractor := snowflake.NewViewSelectStatementExtractor(view.Text)
		substringOfQuery, err := extractor.Extract()
		if err != nil {
			return err
		}
		if err = d.Set("statement", substringOfQuery); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("error reading view %v, err = %w, `text` is missing; if the view is secure then the role used by the provider must own the view (consult https://docs.snowflake.com/en/sql-reference/sql/create-view#usage-notes)", d.Id(), err)
	}

	return nil
}

// UpdateView implements schema.UpdateFunc.
func UpdateView(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	// The only way to update the statement field in a view is to perform create or replace with the new statement.
	// In case of any statement change, create or replace will be performed with all the old parameters, except statement
	// and copy grants (which is always set to true to keep the permissions from the previous state).
	if d.HasChange("statement") {
		oldIsSecure, _ := d.GetChange("is_secure")
		oldComment, _ := d.GetChange("comment")
		oldTags, _ := d.GetChange("tag")

		createRequest := sdk.NewCreateViewRequest(id, d.Get("statement").(string)).
			WithOrReplace(true).
			WithCopyGrants(true).
			WithComment(oldComment.(string)).
			WithTag(getTagsFromList(oldTags.([]any)))

		if oldIsSecure.(bool) {
			createRequest.WithSecure(true)
		}

		err := client.Views.Create(ctx, createRequest)
		if err != nil {
			return fmt.Errorf("error when changing property on %v and performing create or replace to update view statements, err = %w", d.Id(), err)
		}
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithRenameTo(newId))
		if err != nil {
			return fmt.Errorf("error renaming view %v err = %w", d.Id(), err)
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	if d.HasChange("comment") {
		if comment := d.Get("comment").(string); comment == "" {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetComment(true))
			if err != nil {
				return fmt.Errorf("error unsetting comment for view %v", d.Id())
			}
		} else {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetComment(comment))
			if err != nil {
				return fmt.Errorf("error updating comment for view %v", d.Id())
			}
		}
	}

	if d.HasChange("is_secure") {
		if d.Get("is_secure").(bool) {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetSecure(true))
			if err != nil {
				return fmt.Errorf("error setting secure for view %v", d.Id())
			}
		} else {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetSecure(true))
			if err != nil {
				return fmt.Errorf("error unsetting secure for view %v", d.Id())
			}
		}
	}

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetTags(unsetTags))
			if err != nil {
				return fmt.Errorf("error unsetting tags on %v, err = %w", d.Id(), err)
			}
		}

		if len(setTags) > 0 {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetTags(setTags))
			if err != nil {
				return fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err)
			}
		}
	}

	return ReadView(d, meta)
}

// DeleteView implements schema.DeleteFunc.
func DeleteView(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
