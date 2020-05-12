package resources

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

var space = regexp.MustCompile(`\s+`)

var viewSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the view; must be unique for the schema in which the view is created. Don't use the | character.",
	},
	"database": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the view. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "PUBLIC",
		Description: "The schema in which to create the view. Don't use the | character.",
		ForceNew:    true,
	},
	"is_secure": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies that the view is secure.",
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the view.",
	},
	"statement": &schema.Schema{
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
// queries and either comparing ASTs or emiting a canonical serialization for comparison. I couldnt'
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
		Exists: ViewExists,

		Schema: viewSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateView implements schema.CreateFunc
func CreateView(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	schema := data.Get("schema").(string)
	database := data.Get("database").(string)
	s := data.Get("statement").(string)

	builder := snowflake.View(name).WithDB(database).WithSchema(schema).WithStatement(s)

	// Set optionals
	if v, ok := data.GetOk("is_secure"); ok && v.(bool) {
		builder.WithSecure()
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := data.GetOk("schema"); ok {
		builder.WithSchema(v.(string))
	}

	q := builder.Create()
	log.Print("[DEBUG] xxx ", q)
	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating view %v", name)
	}

	data.SetId(fmt.Sprintf("%v|%v|%v", database, schema, name))

	return ReadView(data, meta)
}

// ReadView implements schema.ReadFunc
func ReadView(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName, schema, view, err := splitViewID(data.Id())
	if err != nil {
		return err
	}

	q := snowflake.View(view).WithDB(dbName).WithSchema(schema).Show()
	row := snowflake.QueryRow(db, q)
	v, err := snowflake.ScanView(row)
	if err != nil {
		return err
	}

	err = data.Set("name", v.Name.String)
	if err != nil {
		return err
	}

	err = data.Set("is_secure", v.IsSecure)
	if err != nil {
		return err
	}

	err = data.Set("comment", v.Comment.String)
	if err != nil {
		return err
	}

	err = data.Set("schema", v.SchemaName.String)
	if err != nil {
		return err
	}

	// Want to only capture the Select part of the query because before that is the Create part of the view which we no longer care about

	extractor := snowflake.NewViewSelectStatementExtractor(v.Text.String)
	substringOfQuery, err := extractor.Extract()
	if err != nil {
		return err
	}

	err = data.Set("statement", substringOfQuery)
	if err != nil {
		return err
	}

	return data.Set("database", v.DatabaseName.String)
}

// UpdateView implements schema.UpdateFunc
func UpdateView(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	dbName, schema, view, err := splitViewID(data.Id())
	if err != nil {
		return err
	}

	builder := snowflake.View(view).WithDB(dbName).WithSchema(schema)

	db := meta.(*sql.DB)
	if data.HasChange("name") {
		_, name := data.GetChange("name")

		q := builder.Rename(name.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error renaming view %v", data.Id())
		}

		data.SetId(fmt.Sprintf("%v|%v|%v", dbName, schema, name.(string)))
		data.SetPartial("name")
	}

	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")

		if c := comment.(string); c == "" {
			q := builder.RemoveComment()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting comment for view %v", data.Id())
			}
		} else {
			q := builder.ChangeComment(c)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error updating comment for view %v", data.Id())
			}
		}

		data.SetPartial("comment")
	}

	data.Partial(false)
	if data.HasChange("is_secure") {
		_, secure := data.GetChange("is_secure")

		if secure.(bool) {
			q := builder.Secure()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error setting secure for view %v", data.Id())
			}
		} else {
			q := builder.Unsecure()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting secure for view %v", data.Id())
			}
		}
	}

	return ReadView(data, meta)
}

// DeleteView implements schema.DeleteFunc
func DeleteView(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName, schema, view, err := splitViewID(data.Id())
	if err != nil {
		return err
	}

	q := snowflake.View(view).WithDB(dbName).WithSchema(schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting view %v", data.Id())
	}

	data.SetId("")

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
