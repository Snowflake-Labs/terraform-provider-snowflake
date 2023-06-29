package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

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
		ForceNew:    true,
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
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The timestamp at which the view was created.",
	},
	"tag": tagReferenceSchema,
}

func normalizeQuery(str string) string {
	return strings.TrimSpace(space.ReplaceAllString(str, " "))
}

// DiffSuppressStatement will suppress diffs between statements if they differ in only case or in
// runs of whitespace (\s+ = \s). This is needed because the snowflake api does not faithfully
// round-trip queries so we cannot do a simple character-wise comparison to detect changes.
//
// Warnings: We will have false positives in cases where a change in case or run of whitespace is
// semantically significant.
//
// If we can find a sql parser that can handle the snowflake dialect then we should switch to parsing
// queries and either comparing ASTs or emitting a canonical serialization for comparison. I couldn't
// find such a library.
func DiffSuppressStatement(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(normalizeQuery(old), normalizeQuery(new))
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
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	s := d.Get("statement").(string)

	builder := snowflake.NewViewBuilder(name).WithDB(database).WithSchema(schema).WithStatement(s)

	// Set optionals
	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		builder.WithReplace()
	}

	if v, ok := d.GetOk("is_secure"); ok && v.(bool) {
		builder.WithSecure()
	}

	if v, ok := d.GetOk("copy_grants"); ok && v.(bool) {
		builder.WithCopyGrants()
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	q, err := builder.Create()
	if err != nil {
		return err
	}
	err = snowflake.Exec(db, q)
	if err != nil {
		return fmt.Errorf("error creating view %v", name)
	}

	viewID := &ViewID{
		DatabaseName: database,
		SchemaName:   schema,
		ViewName:     name,
	}
	dataIDInput, err := viewID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadView(d, meta)
}

// ReadView implements schema.ReadFunc.
func ReadView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	viewID, err := viewIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := viewID.DatabaseName
	schema := viewID.SchemaName
	view := viewID.ViewName

	q := snowflake.NewViewBuilder(view).WithDB(dbName).WithSchema(schema).Show()
	row := snowflake.QueryRow(db, q)
	v, err := snowflake.ScanView(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] view (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	if err = d.Set("name", v.Name.String); err != nil {
		return err
	}
	if err = d.Set("is_secure", v.IsSecure); err != nil {
		return err
	}
	if err = d.Set("copy_grants", v.HasCopyGrants()); err != nil {
		return err
	}
	if err = d.Set("comment", v.Comment.String); err != nil {
		return err
	}
	if err = d.Set("schema", v.SchemaName.String); err != nil {
		return err
	}
	if err = d.Set("created_on", v.CreatedOn.String()); err != nil {
		return err
	}

	// Want to only capture the Select part of the query because before that is the Create part of the view which we no longer care about

	extractor := snowflake.NewViewSelectStatementExtractor(v.Text.String)
	substringOfQuery, err := extractor.Extract()
	if err != nil {
		return err
	}
	if err = d.Set("statement", substringOfQuery); err != nil {
		return err
	}
	err = d.Set("database", v.DatabaseName.String)
	return err
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
		comment := d.Get("comment")

		if c := comment.(string); c == "" {
			q, err := builder.RemoveComment()
			if err != nil {
				return err
			}
			if err = snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting comment for view %v", d.Id())
			}
		} else {
			q, err := builder.ChangeComment(c)
			if err != nil {
				return err
			}
			if err = snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error updating comment for view %v", d.Id())
			}
		}
	}
	if d.HasChange("is_secure") {
		secure := d.Get("is_secure")

		if secure.(bool) {
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
	viewID, err := viewIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := viewID.DatabaseName
	schema := viewID.SchemaName
	view := viewID.ViewName

	q, err := snowflake.NewViewBuilder(view).WithDB(dbName).WithSchema(schema).Drop()
	if err != nil {
		return err
	}
	if err = snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting view %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
