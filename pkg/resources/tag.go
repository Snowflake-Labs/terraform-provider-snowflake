package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

const (
	tagIDDelimiter = '|'
)

var tagSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the tag; must be unique for the database in which the tag is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the tag.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the tag.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the tag.",
	},
}

var tagReferenceSchema = &schema.Schema{
	Type:        schema.TypeList,
	Required:    false,
	Optional:    true,
	MinItems:    0,
	Description: "Definitions of a tag to associate with the resource.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag name, e.g. department.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag value, e.g. marketing_info.",
			},
			"database": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Name of the database that the tag was created in.",
			},
			"schema": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Name of the schema that the tag was created in.",
			},
		},
	},
}

type tagID struct {
	DatabaseName string
	SchemaName   string
	TagName      string
}

type TagBuilder interface {
	UnsetTag(snowflake.TagValue) string
	AddTag(snowflake.TagValue) string
	ChangeTag(snowflake.TagValue) string
}

func handleTagChanges(db *sql.DB, d *schema.ResourceData, builder TagBuilder) error {
	if d.HasChange("tag") {
		old, new := d.GetChange("tag")
		removed, added, changed := getTags(old).diffs(getTags(new))
		for _, tA := range removed {
			q := builder.UnsetTag(tA.toSnowflakeTagValue())
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error dropping tag on %v", d.Id())
			}
		}
		for _, tA := range added {
			q := builder.AddTag(tA.toSnowflakeTagValue())

			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error adding column on %v", d.Id())
			}
		}
		for _, tA := range changed {
			q := builder.ChangeTag(tA.toSnowflakeTagValue())
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error changing property on %v", d.Id())
			}
		}
	}
	return nil
}

// String() takes in a schemaID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|TagName
func (ti *tagID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = schemaIDDelimiter
	dataIdentifiers := [][]string{{ti.DatabaseName, ti.SchemaName, ti.TagName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strTagID := strings.TrimSpace(buf.String())
	return strTagID, nil
}

// tagIDFromString() takes in a pipe-delimited string: DatabaseName|tagName
// and returns a tagID object
func tagIDFromString(stringID string) (*tagID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = tagIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per schema")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	tagResult := &tagID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		TagName:      lines[0][2],
	}
	return tagResult, nil
}

// Schema returns a pointer to the resource representing a schema
func Tag() *schema.Resource {
	return &schema.Resource{
		Create: CreateTag,
		Read:   ReadTag,
		Update: UpdateTag,
		Delete: DeleteTag,

		Schema: tagSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSchema implements schema.CreateFunc
func CreateTag(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)

	builder := snowflake.Tag(name).WithDB(database).WithSchema(schema)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	q := builder.Create()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating tag %v", name)
	}

	tagID := &tagID{
		DatabaseName: database,
		SchemaName:   schema,
		TagName:      name,
	}
	dataIDInput, err := tagID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadTag(d, meta)
}

// ReadSchema implements schema.ReadFunc
func ReadTag(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tagID, err := tagIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tagID.DatabaseName
	schemaName := tagID.SchemaName
	tag := tagID.TagName

	q := snowflake.Tag(tag).WithDB(dbName).WithSchema(schemaName).Show()
	row := snowflake.QueryRow(db, q)

	t, err := snowflake.ScanTag(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] tag (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("name", t.Name.String)
	if err != nil {
		return err
	}

	err = d.Set("database", t.DatabaseName.String)
	if err != nil {
		return err
	}

	err = d.Set("schema", t.SchemaName.String)
	if err != nil {
		return err
	}

	err = d.Set("comment", t.Comment.String)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTag implements schema.UpdateFunc
func UpdateTag(d *schema.ResourceData, meta interface{}) error {
	tagID, err := tagIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tagID.DatabaseName
	schemaName := tagID.SchemaName
	tag := tagID.TagName

	builder := snowflake.Tag(tag).WithDB(dbName).WithSchema(schemaName)

	db := meta.(*sql.DB)
	if d.HasChange("comment") {
		comment, ok := d.GetOk("comment")
		var q string
		if ok {
			q = builder.ChangeComment(comment.(string))
		} else {
			q = builder.RemoveComment()
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating tag comment on %v", d.Id())
		}
	}

	return ReadTag(d, meta)
}

// DeleteTag implements schema.DeleteFunc
func DeleteTag(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tagID, err := tagIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tagID.DatabaseName
	schemaName := tagID.SchemaName
	tag := tagID.TagName

	q := snowflake.Tag(tag).WithDB(dbName).WithSchema(schemaName).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting tag %v", d.Id())
	}

	d.SetId("")

	return nil
}

// SchemaExists implements schema.ExistsFunc
func TagExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	tagID, err := tagIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := tagID.DatabaseName
	schemaName := tagID.SchemaName
	tag := tagID.TagName

	q := snowflake.Tag(tag).WithDB(dbName).WithSchema(schemaName).Show()
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

type tags []tag

func (t tags) toSnowflakeTagValues() []snowflake.TagValue {
	sT := make([]snowflake.TagValue, len(t))
	for i, tag := range t {
		sT[i] = tag.toSnowflakeTagValue()
	}
	return sT
}

func (tag tag) toSnowflakeTagValue() snowflake.TagValue {
	return snowflake.TagValue{
		Name:     tag.name,
		Value:    tag.value,
		Database: tag.database,
		Schema:   tag.schema,
	}
}

func (old tags) getNewIn(new tags) (added tags) {
	added = tags{}
	for _, t0 := range old {
		found := false
		for _, cN := range new {
			if t0.name == cN.name {
				found = true
				break
			}
		}
		if !found {
			added = append(added, t0)
		}
	}
	return
}

func (old tags) getChangedTagProperties(new tags) (changed tags) {
	changed = tags{}
	for _, t0 := range old {
		for _, tN := range new {
			if t0.name == tN.name && t0.value != tN.value {
				changed = append(changed, tN)
			}
		}
	}
	return
}

func (old tags) diffs(new tags) (removed tags, added tags, changed tags) {
	return old.getNewIn(new), new.getNewIn(old), old.getChangedTagProperties(new)
}

func (old columns) getNewIn(new columns) (added columns) {
	added = columns{}
	for _, cO := range old {
		found := false
		for _, cN := range new {
			if cO.name == cN.name {
				found = true
				break
			}
		}
		if !found {
			added = append(added, cO)
		}
	}
	return
}

type tag struct {
	name     string
	value    string
	database string
	schema   string
}

func getTags(from interface{}) (to tags) {
	tags := from.([]interface{})
	to = make([]tag, len(tags))
	for i, t := range tags {
		v := t.(map[string]interface{})
		to[i] = tag{
			name:     v["name"].(string),
			value:    v["value"].(string),
			database: v["database"].(string),
			schema:   v["schema"].(string),
		}
	}
	return to
}
