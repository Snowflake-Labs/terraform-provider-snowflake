package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
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

// View returns a pointer to the resource representing a view
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

//String() takes in a materializedViewID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|ExternalTableName
func (si *materializedViewID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = materializedViewDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.ViewName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strMeterilizedViewID := strings.TrimSpace(buf.String())
	return strMeterilizedViewID, nil
}

// materializedViewIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|MaterializedViewName
// and returns a externalTableID object
func materializedViewIDFromString(stringID string) (*materializedViewID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = materializedViewDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
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

// CreateMaterializedView implements schema.CreateFunc
func CreateMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	warehouse := d.Get("warehouse").(string)
	s := d.Get("statement").(string)

	builder := snowflake.MaterializedView(name).WithDB(database).WithSchema(schema).WithWarehouse(warehouse).WithStatement(s)

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
	err := snowflake.ExecMulti(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating view %v", name)
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

// ReadMaterializedView implements schema.ReadFunc
func ReadMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	materializedViewID, err := materializedViewIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := materializedViewID.DatabaseName
	schema := materializedViewID.SchemaName
	view := materializedViewID.ViewName

	q := snowflake.MaterializedView(view).WithDB(dbName).WithSchema(schema).Show()
	row := snowflake.QueryRow(db, q)
	v, err := snowflake.ScanMaterializedView(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] view (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("name", v.Name.String)
	if err != nil {
		return err
	}

	err = d.Set("is_secure", v.IsSecure)
	if err != nil {
		return err
	}

	err = d.Set("comment", v.Comment.String)
	if err != nil {
		return err
	}

	err = d.Set("schema", v.SchemaName.String)
	if err != nil {
		return err
	}

	// Want to only capture the Select part of the query because before that is the Create part of the view which we no longer care about

	extractor := snowflake.NewViewSelectStatementExtractor(v.Text.String)
	substringOfQuery, err := extractor.ExtractMaterializedView()
	if err != nil {
		return err
	}

	err = d.Set("statement", substringOfQuery)
	if err != nil {
		return err
	}

	return d.Set("database", v.DatabaseName.String)
}

// UpdateMaterializedView implements schema.UpdateFunc
func UpdateMaterializedView(d *schema.ResourceData, meta interface{}) error {
	materializedViewID, err := materializedViewIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := materializedViewID.DatabaseName
	schema := materializedViewID.SchemaName
	view := materializedViewID.ViewName

	builder := snowflake.MaterializedView(view).WithDB(dbName).WithSchema(schema)

	db := meta.(*sql.DB)
	if d.HasChange("name") {
		name := d.Get("name")

		q := builder.Rename(name.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error renaming view %v", d.Id())
		}

		d.SetId(fmt.Sprintf("%v|%v|%v", dbName, schema, name.(string)))
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")

		if c := comment.(string); c == "" {
			q := builder.RemoveComment()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting comment for view %v", d.Id())
			}
		} else {
			q := builder.ChangeComment(c)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error updating comment for view %v", d.Id())
			}
		}
	}
	if d.HasChange("is_secure") {
		secure := d.Get("is_secure")

		if secure.(bool) {
			q := builder.Secure()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error setting secure for view %v", d.Id())
			}
		} else {
			q := builder.Unsecure()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting secure for materialized view %v", d.Id())
			}
		}
	}

	handleTagChanges(db, d, builder)

	return ReadMaterializedView(d, meta)
}

// DeleteMaterializedView implements schema.DeleteFunc
func DeleteMaterializedView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	materializedViewID, err := materializedViewIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := materializedViewID.DatabaseName
	schema := materializedViewID.SchemaName
	view := materializedViewID.ViewName

	q := snowflake.MaterializedView(view).WithDB(dbName).WithSchema(schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting materialized view %v", d.Id())
	}

	d.SetId("")

	return nil
}
