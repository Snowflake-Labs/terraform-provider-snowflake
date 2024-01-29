package resources

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"

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

// View returns a pointer to the resource representing a view.
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

type materializedViewID struct {
	DatabaseName string
	SchemaName   string
	ViewName     string
}

const (
	materializedViewDelimiter = '|'
)

// String() takes in a materializedViewID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|ExternalTableName.
func (si *materializedViewID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = materializedViewDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.ViewName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strMeterilizedViewID := strings.TrimSpace(buf.String())
	return strMeterilizedViewID, nil
}

// materializedViewIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|MaterializedViewName
// and returns a externalTableID object.
func materializedViewIDFromString(stringID string) (*materializedViewID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = materializedViewDelimiter
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

	materializedViewResult := &materializedViewID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		ViewName:     lines[0][2],
	}
	return materializedViewResult, nil
}

// CreateMaterializedView implements schema.CreateFunc.
func CreateMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	warehouse := d.Get("warehouse").(string)
	s := d.Get("statement").(string)

	builder := snowflake.NewMaterializedViewBuilder(name).WithDB(database).WithSchema(schema).WithWarehouse(warehouse).WithStatement(s)

	// Set optionals
	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		builder.WithReplace()
	}

	if v, ok := d.GetOk("is_secure"); ok && v.(bool) {
		builder.WithSecure()
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	q := builder.Create()
	log.Print("[DEBUG] xxx ", q)
	if err := snowflake.ExecMulti(db, q); err != nil {
		return fmt.Errorf("error creating view %v err = %w", name, err)
	}

	materializedViewID := &materializedViewID{
		DatabaseName: database,
		SchemaName:   schema,
		ViewName:     name,
	}
	dataIDInput, err := materializedViewID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadMaterializedView(d, meta)
}

// ReadMaterializedView implements schema.ReadFunc.
func ReadMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	materializedViewID, err := materializedViewIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := materializedViewID.DatabaseName
	schema := materializedViewID.SchemaName
	view := materializedViewID.ViewName

	q := snowflake.NewMaterializedViewBuilder(view).WithDB(dbName).WithSchema(schema).Show()
	row := snowflake.QueryRow(db, q)
	v, err := snowflake.ScanMaterializedView(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] view (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("name", v.Name.String); err != nil {
		return err
	}

	if err := d.Set("is_secure", v.IsSecure); err != nil {
		return err
	}

	if err := d.Set("comment", v.Comment.String); err != nil {
		return err
	}

	if err := d.Set("schema", v.SchemaName.String); err != nil {
		return err
	}

	// Want to only capture the Select part of the query because before that is the Create part of the view which we no longer care about

	extractor := snowflake.NewViewSelectStatementExtractor(v.Text.String)
	substringOfQuery, err := extractor.ExtractMaterializedView()
	if err != nil {
		return err
	}

	if err := d.Set("statement", substringOfQuery); err != nil {
		return err
	}

	return d.Set("database", v.DatabaseName.String)
}

// UpdateMaterializedView implements schema.UpdateFunc.
func UpdateMaterializedView(d *schema.ResourceData, meta interface{}) error {
	mvid, err := materializedViewIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := mvid.DatabaseName
	schema := mvid.SchemaName
	view := mvid.ViewName

	builder := snowflake.NewMaterializedViewBuilder(view).WithDB(dbName).WithSchema(schema)

	db := meta.(*sql.DB)
	if d.HasChange("name") {
		name := d.Get("name")

		q := builder.Rename(name.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error renaming view %v err = %w", d.Id(), err)
		}
		materializedViewID := &materializedViewID{
			DatabaseName: dbName,
			SchemaName:   schema,
			ViewName:     name.(string),
		}
		dataIDInput, err := materializedViewID.String()
		if err != nil {
			return err
		}
		d.SetId(dataIDInput)
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")

		if c := comment.(string); c == "" {
			q := builder.RemoveComment()
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting comment for view %v err = %w", d.Id(), err)
			}
		} else {
			q := builder.ChangeComment(c)
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error updating comment for view %v err = %w", d.Id(), err)
			}
		}
	}
	if d.HasChange("is_secure") {
		secure := d.Get("is_secure")

		if secure.(bool) {
			q := builder.Secure()
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error setting secure for view %v err = %w", d.Id(), err)
			}
		} else {
			q := builder.Unsecure()
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting secure for materialized view %v err = %w", d.Id(), err)
			}
		}
	}

	handleErr := handleTagChanges(db, d, builder)
	if handleErr != nil {
		return handleErr
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
