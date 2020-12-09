package resources

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
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
		Optional:    true,
		Default:     "PUBLIC",
		Description: "The schema in which to create the view. Don't use the | character.",
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
}

func normalizeQuery(str string) string {
	return strings.TrimSpace(space.ReplaceAllString(str, " "))
}

// DiffSuppressStatement will suppress diffs between statemens if they differ in only case or in
// runs of whitespace (\s+ = \s). This is needed because the snowflake api does not faithfully
// round-trip queries so we cannot do a simple character-wise comparison to detect changes.
//
// Warnings: We will have false positives in cases where a change in case or run of whitespace is
// semantically significant.
//
// If we can find a sql parser that can handle the snowflake dialect then we should switch to parsing
// queries and either comparing ASTs or emiting a canonical serialization for comparison. I couldn't
// find such a library.
func DiffSuppressStatement(_, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(normalizeQuery(old), normalizeQuery(new))
}

// View returns a pointer to the resource representing a view
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

// CreateView implements schema.CreateFunc
func CreateView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	s := d.Get("statement").(string)

	builder := snowflake.View(name).WithDB(database).WithSchema(schema).WithStatement(s)

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

	q := builder.Create()
	log.Print("[DEBUG] xxx ", q)
	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating view %v", name)
	}

	d.SetId(fmt.Sprintf("%v|%v|%v", database, schema, name))

	return ReadView(d, meta)
}

// ReadView implements schema.ReadFunc
func ReadView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName, schema, view, err := splitViewID(d.Id())
	if err != nil {
		return err
	}

	q := snowflake.View(view).WithDB(dbName).WithSchema(schema).Show()
	row := snowflake.QueryRow(db, q)
	v, err := snowflake.ScanView(row)
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
	substringOfQuery, err := extractor.Extract()
	if err != nil {
		return err
	}

	err = d.Set("statement", substringOfQuery)
	if err != nil {
		return err
	}

	return d.Set("database", v.DatabaseName.String)
}

// UpdateView implements schema.UpdateFunc
func UpdateView(d *schema.ResourceData, meta interface{}) error {
	dbName, schema, view, err := splitViewID(d.Id())
	if err != nil {
		return err
	}

	builder := snowflake.View(view).WithDB(dbName).WithSchema(schema)

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
				return errors.Wrapf(err, "error unsetting secure for view %v", d.Id())
			}
		}
	}

	return ReadView(d, meta)
}

// DeleteView implements schema.DeleteFunc
func DeleteView(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName, schema, view, err := splitViewID(d.Id())
	if err != nil {
		return err
	}

	q := snowflake.View(view).WithDB(dbName).WithSchema(schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting view %v", d.Id())
	}

	d.SetId("")

	return nil
}

// ViewExists implements schema.ExistsFunc
func ViewExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	dbName, schema, view, err := splitViewID(data.Id())
	if err != nil {
		return false, err
	}

	q := snowflake.View(view).WithDB(dbName).WithSchema(schema).Show()
	rows, err := db.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

// splitViewID takes the <database_name>|<schema_name>|<view_name> ID and returns the database
// name, schema name and view name.
func splitViewID(v string) (string, string, string, error) {
	arr := strings.Split(v, "|")
	if len(arr) != 3 {
		return "", "", "", fmt.Errorf("ID %v is invalid", v)
	}

	return arr[0], arr[1], arr[2], nil
}
