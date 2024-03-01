package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
	"allowed_values": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "List of allowed values for the tag.",
	},
}

var tagReferenceSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MinItems:    0,
	Description: "Definitions of a tag to associate with the resource.",
	Deprecated:  "Use the 'snowflake_tag_association' resource instead.",
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
				Optional:    true,
				Description: "Name of the database that the tag was created in.",
			},
			"schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the schema that the tag was created in.",
			},
		},
	},
}

type TagID struct {
	DatabaseName string
	SchemaName   string
	TagName      string
}

type TagBuilder interface {
	UnsetTag(snowflake.TagValue) string
	AddTag(snowflake.TagValue) string
	ChangeTag(snowflake.TagValue) string
}

// String() takes in a schemaID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|TagName.
func (ti *TagID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = schemaIDDelimiter
	dataIdentifiers := [][]string{{ti.DatabaseName, ti.SchemaName, ti.TagName}}

	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strTagID := strings.TrimSpace(buf.String())
	return strTagID, nil
}

// tagIDFromString() takes in a pipe-delimited string: DatabaseName|tagName
// and returns a tagID object.
func tagIDFromString(stringID string) (*TagID, error) {
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

	tagResult := &TagID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		TagName:      lines[0][2],
	}
	return tagResult, nil
}

// Schema returns a pointer to the resource representing a schema.
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

// CreateSchema implements schema.CreateFunc.
func CreateTag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)

	builder := snowflake.NewTagBuilder(name).WithDB(database).WithSchema(schema)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("allowed_values"); ok {
		builder.WithAllowedValues(expandStringList(v.([]interface{})))
	}

	q := builder.Create()

	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error creating tag %v", name)
	}

	tagID := &TagID{
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

// ReadSchema implements schema.ReadFunc.
func ReadTag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	tagID, err := tagIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tagID.DatabaseName
	schemaName := tagID.SchemaName
	tag := tagID.TagName

	q := snowflake.NewTagBuilder(tag).WithDB(dbName).WithSchema(schemaName).Show()
	row := snowflake.QueryRow(db, q)

	t, err := snowflake.ScanTag(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] tag (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("name", t.Name.String); err != nil {
		return err
	}

	if err := d.Set("database", t.DatabaseName.String); err != nil {
		return err
	}

	if err := d.Set("schema", t.SchemaName.String); err != nil {
		return err
	}

	if err := d.Set("comment", t.Comment.String); err != nil {
		return err
	}

	av := strings.ReplaceAll(t.AllowedValues.String, "\"", "")
	av = strings.TrimPrefix(av, "[")
	av = strings.TrimSuffix(av, "]")
	err = d.Set("allowed_values", helpers.StringListToList(av))
	return err
}

// UpdateTag implements schema.UpdateFunc.
func UpdateTag(d *schema.ResourceData, meta interface{}) error {
	tagID, err := tagIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tagID.DatabaseName
	schemaName := tagID.SchemaName
	tag := tagID.TagName

	builder := snowflake.NewTagBuilder(tag).WithDB(dbName).WithSchema(schemaName)

	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	if d.HasChange("comment") {
		comment, ok := d.GetOk("comment")
		var q string
		if ok {
			q = builder.ChangeComment(comment.(string))
		} else {
			q = builder.RemoveComment()
		}
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating tag comment on %v", d.Id())
		}
	}

	// If there is change in allowed_values field
	if d.HasChange("allowed_values") {
		if _, ok := d.GetOk("allowed_values"); ok {
			v := d.Get("allowed_values")

			ns := expandAllowedValues(v)

			q := builder.RemoveAllowedValues()

			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error removing ALLOWED_VALUES for tag %v", tag)
			}

			addQuery := builder.AddAllowedValues(ns)
			if err := snowflake.Exec(db, addQuery); err != nil {
				return fmt.Errorf("error adding ALLOWED_VALUES for tag %v", tag)
			}
		} else {
			q := builder.RemoveAllowedValues()
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error removing ALLOWED_VALUES for tag %v", tag)
			}
		}
	}

	return ReadTag(d, meta)
}

// Returns the slice of strings for inputed allowed values.
func expandAllowedValues(avChangeSet interface{}) []string {
	avList := avChangeSet.([]interface{})
	newAvs := make([]string, len(avList))
	for idx, value := range avList {
		newAvs[idx] = fmt.Sprintf("%v", value)
	}

	return newAvs
}

// DeleteTag implements schema.DeleteFunc.
func DeleteTag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	tagID, err := tagIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tagID.DatabaseName
	schemaName := tagID.SchemaName
	tag := tagID.TagName

	q := snowflake.NewTagBuilder(tag).WithDB(dbName).WithSchema(schemaName).Drop()

	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting tag %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}

type tags []tag

func (t tags) toSnowflakeTagValues() []snowflake.TagValue {
	sT := make([]snowflake.TagValue, len(t))
	for i, tag := range t {
		sT[i] = tag.toSnowflakeTagValue()
	}
	return sT
}

func (t tag) toSnowflakeTagValue() snowflake.TagValue {
	return snowflake.TagValue{
		Name:     t.name,
		Value:    t.value,
		Database: t.database,
		Schema:   t.schema,
	}
}

func (t tags) getNewIn(new tags) (added tags) {
	added = tags{}
	for _, t0 := range t {
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

func (t tags) getChangedTagProperties(new tags) (changed tags) {
	changed = tags{}
	for _, t0 := range t {
		for _, tN := range new {
			if t0.name == tN.name && t0.value != tN.value {
				changed = append(changed, tN)
			}
		}
	}
	return
}

func (t tags) diffs(new tags) (removed tags, added tags, changed tags) {
	return t.getNewIn(new), new.getNewIn(t), t.getChangedTagProperties(new)
}

func (t columns) getNewIn(new columns) (added columns) {
	added = columns{}
	for _, cO := range t {
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
