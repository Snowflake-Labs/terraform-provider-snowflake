package resources

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"regexp"
	"strings"

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
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Retains the access permissions from the original view when a new view is created using the OR REPLACE clause.",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return oldValue != "" && oldValue != newValue
		},
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

type ViewID struct {
	DatabaseName string
	SchemaName   string
	ViewName     string
}

const (
	viewDelimiter = '|'
)

// String() takes in a viewID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|viewName.
func (si *ViewID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = viewDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.ViewName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strViewID := strings.TrimSpace(buf.String())
	return strViewID, nil
}

// viewIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|viewName
// and returns a externalTableID object.
func viewIDFromString(stringID string) (*ViewID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = viewDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	viewResult := &ViewID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		ViewName:     lines[0][2],
	}
	return viewResult, nil
}

// CreateView implements schema.CreateFunc.
func CreateView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	s := d.Get("statement").(string)
	createRequest := sdk.NewCreateViewRequest(id, s)

	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		createRequest.WithOrReplace(sdk.Bool(true))
	}

	if v, ok := d.GetOk("is_secure"); ok && v.(bool) {
		createRequest.WithSecure(sdk.Bool(true))
	}

	if v, ok := d.GetOk("copy_grants"); ok && v.(bool) {
		createRequest.WithCopyGrants(sdk.Bool(true))
	}

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	if _, ok := d.GetOk("tag"); ok {
		createRequest.WithTag(getPropertyTags(d, "tag"))
	}

	err := client.Views.Create(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("error creating view %v err = %w", name, err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadView(d, meta)
}

// ReadView implements schema.ReadFunc.
func ReadView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
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

	// TODO [TODO]: what do we do with these extractors?
	// Want to only capture the SELECT part of the query because before that is the CREATE part of the view.
	extractor := snowflake.NewViewSelectStatementExtractor(view.Text)
	substringOfQuery, err := extractor.Extract()
	if err != nil {
		return err
	}

	if err = d.Set("statement", substringOfQuery); err != nil {
		return err
	}

	return nil
}

// UpdateView implements schema.UpdateFunc.
func UpdateView(d *schema.ResourceData, meta interface{}) error {
	viewID, err := viewIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := viewID.DatabaseName
	schema := viewID.SchemaName
	view := viewID.ViewName
	builder := snowflake.NewViewBuilder(view).WithDB(dbName).WithSchema(schema)
	db := meta.(*sql.DB)

	// The only way to update the statement field in a view is to perform create or replace with the new statement.
	// In case of any statement change, create or replace will be performed with all the old parameters, except statement
	// and copy grants (which is always set to true to keep the permissions from the previous state).
	if d.HasChange("statement") {
		isSecureOld, _ := d.GetChange("is_secure")
		commentOld, _ := d.GetChange("comment")
		tagsOld, _ := d.GetChange("tag")

		if isSecureOld.(bool) {
			builder.WithSecure()
		}

		query, err := builder.
			WithReplace().
			WithStatement(d.Get("statement").(string)).
			WithCopyGrants().
			WithComment(commentOld.(string)).
			WithTags(getTags(tagsOld).toSnowflakeTagValues()).
			Create()
		if err != nil {
			return fmt.Errorf("error when building sql query on %v, err = %w", d.Id(), err)
		}

		if err := snowflake.Exec(db, query); err != nil {
			return fmt.Errorf("error when changing property on %v and performing create or replace to update view statements, err = %w", d.Id(), err)
		}
	}

	if d.HasChange("name") {
		name := d.Get("name")
		q, err := builder.Rename(name.(string))
		if err != nil {
			return err
		}
		if err = snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error renaming view %v", d.Id())
		}
		viewID := &ViewID{
			DatabaseName: dbName,
			SchemaName:   schema,
			ViewName:     name.(string),
		}
		dataIDInput, err := viewID.String()
		if err != nil {
			return err
		}
		d.SetId(dataIDInput)
	}

	if d.HasChange("comment") {
		if comment := d.Get("comment").(string); comment == "" {
			q, err := builder.RemoveComment()
			if err != nil {
				return err
			}
			if err = snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting comment for view %v", d.Id())
			}
		} else {
			q, err := builder.ChangeComment(comment)
			if err != nil {
				return err
			}
			if err = snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error updating comment for view %v", d.Id())
			}
		}
	}

	if d.HasChange("is_secure") {
		if d.Get("is_secure").(bool) {
			q, err := builder.Secure()
			if err != nil {
				return err
			}
			if err = snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error setting secure for view %v", d.Id())
			}
		} else {
			q, err := builder.Unsecure()
			if err != nil {
				return err
			}
			if err = snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting secure for view %v", d.Id())
			}
		}
	}

	tagChangeErr := handleTagChanges(db, d, builder)
	if tagChangeErr != nil {
		return tagChangeErr
	}
	if d.HasChange("tag") {
		o, n := d.GetChange("tag")
		removed, added, changed := getTags(o).diffs(getTags(n))
		for _, tA := range removed {
			q := builder.UnsetTag(tA.toSnowflakeTagValue())
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error dropping tag on %v", d.Id())
			}
		}
		for _, tA := range added {
			q := builder.AddTag(tA.toSnowflakeTagValue())
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error adding column on %v", d.Id())
			}
		}
		for _, tA := range changed {
			q := builder.ChangeTag(tA.toSnowflakeTagValue())
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error changing property on %v", d.Id())
			}
		}
	}

	return ReadView(d, meta)
}

// DeleteView implements schema.DeleteFunc.
func DeleteView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
