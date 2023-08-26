package resources

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
)

const (
	schemaIDDelimiter = '|'
)

var schemaSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the schema; must be unique for the database in which the schema is created.",
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
// DatabaseName|schemaName.
func (si *schemaID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = schemaIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strSchemaID := strings.TrimSpace(buf.String())
	return strSchemaID, nil
}

// schemaIDFromString() takes in a pipe-delimited string: DatabaseName|schemaName
// and returns a schemaID object.
func schemaIDFromString(stringID string) (*schemaID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = schemaIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
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

// Schema returns a pointer to the resource representing a schema.
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

// CreateSchema implements schema.CreateFunc.
func CreateSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	database := d.Get("database").(string)

	builder := snowflake.NewSchemaBuilder(name).WithDB(database)

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
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error creating schema %v err = %w", name, err)
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

// ReadSchema implements schema.ReadFunc.
func ReadSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	schemaID, err := schemaIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	// Checks if the corresponding database still exists; if not, than the schema also cannot exist
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	_, err = client.Databases.ShowByID(ctx, sdk.NewAccountObjectIdentifier(dbName))
	if err != nil {
		d.SetId("")
	}

	q := snowflake.NewSchemaBuilder(schema).WithDB(dbName).Show()
	row := snowflake.QueryRow(db, q)

	s, err := snowflake.ScanSchema(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("name", s.Name.String); err != nil {
		return err
	}

	if err := d.Set("database", s.DatabaseName.String); err != nil {
		return err
	}

	if err := d.Set("comment", s.Comment.String); err != nil {
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

		if err := d.Set("data_retention_days", i); err != nil {
			return err
		}
	}

	// reset the options before reading back from the DB
	if err := d.Set("is_transient", false); err != nil {
		return err
	}

	if err := d.Set("is_managed", false); err != nil {
		return err
	}

	if opts := s.Options.String; opts != "" {
		for _, opt := range strings.Split(opts, ", ") {
			switch opt {
			case "TRANSIENT":
				if err := d.Set("is_transient", true); err != nil {
					return err
				}
			case "MANAGED ACCESS":
				if err := d.Set("is_managed", true); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// UpdateSchema implements schema.UpdateFunc.
func UpdateSchema(d *schema.ResourceData, meta interface{}) error {
	sid, err := schemaIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := sid.DatabaseName
	schema := sid.SchemaName

	builder := snowflake.NewSchemaBuilder(schema).WithDB(dbName)

	db := meta.(*sql.DB)
	if d.HasChange("name") {
		name := d.Get("name")
		q := builder.Rename(name.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating schema name on %v err = %w", d.Id(), err)
		}

		schemaID := &schemaID{
			DatabaseName: dbName,
			SchemaName:   name.(string),
		}
		dataIDInput, err := schemaID.String()
		if err != nil {
			return err
		}
		d.SetId(dataIDInput)
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating schema comment on %v err = %w", d.Id(), err)
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

		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error changing management state on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("data_retention_days") {
		days := d.Get("data_retention_days")

		q := builder.ChangeDataRetentionDays(days.(int))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating data retention days on %v err = %w", d.Id(), err)
		}
	}

	tagChangeErr := handleTagChanges(db, d, builder)
	if tagChangeErr != nil {
		return tagChangeErr
	}

	return ReadSchema(d, meta)
}

// DeleteSchema implements schema.DeleteFunc.
func DeleteSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	schemaID, err := schemaIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	q := snowflake.NewSchemaBuilder(schema).WithDB(dbName).Drop()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting schema %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
