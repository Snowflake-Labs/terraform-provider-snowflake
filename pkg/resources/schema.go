package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

const (
	schemaIDDelimiter = '|'
)

var schemaSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the schema; must be unique for the database in which the schema is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the schema.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the schema.",
	},
	"is_transient": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies a schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
		ForceNew:    true,
	},
	"is_managed": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner.",
	},
	"data_retention_days": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the schema, as well as specifying the default Time Travel retention time for all tables created in the schema.",
		ValidateFunc: validation.IntBetween(0, 90),
	},
	"tag": tagReferenceSchema,
}

type schemaID struct {
	DatabaseName string
	SchemaName   string
}

// String() takes in a schemaID object and returns a pipe-delimited string:
// DatabaseName|schemaName
func (si *schemaID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = schemaIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strSchemaID := strings.TrimSpace(buf.String())
	return strSchemaID, nil
}

// schemaIDFromString() takes in a pipe-delimited string: DatabaseName|schemaName
// and returns a schemaID object
func schemaIDFromString(stringID string) (*schemaID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = schemaIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per schema")
	}
	if len(lines[0]) != 2 {
		return nil, fmt.Errorf("2 fields allowed")
	}

	schemaResult := &schemaID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
	}
	return schemaResult, nil
}

// Schema returns a pointer to the resource representing a schema
func Schema() *schema.Resource {
	return &schema.Resource{
		Create: CreateSchema,
		Read:   ReadSchema,
		Update: UpdateSchema,
		Delete: DeleteSchema,

		Schema: schemaSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSchema implements schema.CreateFunc
func CreateSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	database := d.Get("database").(string)

	builder := snowflake.Schema(name).WithDB(database)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("is_transient"); ok && v.(bool) {
		builder.Transient()
	}

	if v, ok := d.GetOk("is_managed"); ok && v.(bool) {
		builder.Managed()
	}

	if v, ok := d.GetOk("data_retention_days"); ok {
		builder.WithDataRetentionDays(v.(int))
	}

	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	q := builder.Create()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating schema %v", name)
	}

	schemaID := &schemaID{
		DatabaseName: database,
		SchemaName:   name,
	}
	dataIDInput, err := schemaID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadSchema(d, meta)
}

// ReadSchema implements schema.ReadFunc
func ReadSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	schemaID, err := schemaIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	q := snowflake.Schema(schema).WithDB(dbName).Show()
	row := snowflake.QueryRow(db, q)

	s, err := snowflake.ScanSchema(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("name", s.Name.String)
	if err != nil {
		return err
	}

	err = d.Set("database", s.DatabaseName.String)
	if err != nil {
		return err
	}

	err = d.Set("comment", s.Comment.String)
	if err != nil {
		return err
	}

	// "retention_time" may sometimes be empty string instead of an integer
	{
		retentionTime := s.RetentionTime.String
		if retentionTime == "" {
			retentionTime = "0"
		}

		i, err := strconv.ParseInt(retentionTime, 10, 64)
		if err != nil {
			return err
		}

		err = d.Set("data_retention_days", i)
		if err != nil {
			return err
		}
	}

	// reset the options before reading back from the DB
	err = d.Set("is_transient", false)
	if err != nil {
		return err
	}

	err = d.Set("is_managed", false)
	if err != nil {
		return err
	}

	if opts := s.Options.String; opts != "" {
		for _, opt := range strings.Split(opts, ", ") {
			switch opt {
			case "TRANSIENT":
				err = d.Set("is_transient", true)
				if err != nil {
					return err
				}
			case "MANAGED ACCESS":
				err = d.Set("is_managed", true)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// UpdateSchema implements schema.UpdateFunc
func UpdateSchema(d *schema.ResourceData, meta interface{}) error {
	schemaID, err := schemaIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	builder := snowflake.Schema(schema).WithDB(dbName)

	db := meta.(*sql.DB)
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating schema comment on %v", d.Id())
		}
	}

	if d.HasChange("is_managed") {
		managed := d.Get("is_managed")
		var q string
		if managed.(bool) {
			q = builder.Manage()
		} else {
			q = builder.Unmanage()
		}

		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error changing management state on %v", d.Id())
		}
	}

	if d.HasChange("data_retention_days") {
		days := d.Get("data_retention_days")

		q := builder.ChangeDataRetentionDays(days.(int))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating data retention days on %v", d.Id())
		}
	}

	handleTagChanges(db, d, builder)

	return ReadSchema(d, meta)
}

// DeleteSchema implements schema.DeleteFunc
func DeleteSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	schemaID, err := schemaIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	q := snowflake.Schema(schema).WithDB(dbName).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting schema %v", d.Id())
	}

	d.SetId("")

	return nil
}

// SchemaExists implements schema.ExistsFunc
func SchemaExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	schemaID, err := schemaIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	q := snowflake.Schema(schema).WithDB(dbName).Show()
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
