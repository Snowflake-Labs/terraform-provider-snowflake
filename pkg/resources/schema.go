package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

const (
	schemaIDDelimiter = '|'
)

var schemaSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the schema; must be unique for the database in which the schema is created.",
		ForceNew:    true,
	},
	"database": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the schema.",
		ForceNew:    true,
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the schema.",
	},
	"is_transient": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies a schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
		ForceNew:    true,
	},
	"is_managed": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner.",
	},
	"data_retention_days": &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the schema, as well as specifying the default Time Travel retention time for all tables created in the schema.",
		ValidateFunc: validation.IntBetween(0, 90),
	},
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
		Exists: SchemaExists,

		Schema: schemaSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateSchema implements schema.CreateFunc
func CreateSchema(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	database := data.Get("database").(string)

	builder := snowflake.Schema(name).WithDB(database)

	// Set optionals
	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := data.GetOk("is_transient"); ok && v.(bool) {
		builder.Transient()
	}

	if v, ok := data.GetOk("is_managed"); ok && v.(bool) {
		builder.Managed()
	}

	if v, ok := data.GetOk("data_retention_days"); ok {
		builder.WithDataRetentionDays(v.(int))
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
	data.SetId(dataIDInput)

	return ReadSchema(data, meta)
}

// ReadSchema implements schema.ReadFunc
func ReadSchema(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	schemaID, err := schemaIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	q := snowflake.Schema(schema).WithDB(dbName).Show()
	row := snowflake.QueryRow(db, q)

	s, err := snowflake.ScanSchema(row)
	if err != nil {
		return err
	}

	err = data.Set("name", s.Name.String)
	if err != nil {
		return err
	}

	err = data.Set("database", s.DatabaseName.String)
	if err != nil {
		return err
	}

	err = data.Set("comment", s.Comment.String)
	if err != nil {
		return err
	}

	err = data.Set("data_retention_days", s.RetentionTime.Int64)
	if err != nil {
		return err
	}

	// reset the options before reading back from the DB
	err = data.Set("is_transient", false)
	if err != nil {
		return err
	}

	err = data.Set("is_managed", false)
	if err != nil {
		return err
	}

	if opts := s.Options.String; opts != "" {
		for _, opt := range strings.Split(opts, ", ") {
			switch opt {
			case "TRANSIENT":
				err = data.Set("is_transient", true)
				if err != nil {
					return err
				}
			case "MANAGED ACCESS":
				err = data.Set("is_managed", true)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// UpdateSchema implements schema.UpdateFunc
func UpdateSchema(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	schemaID, err := schemaIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	builder := snowflake.Schema(schema).WithDB(dbName)

	db := meta.(*sql.DB)
	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating schema comment on %v", data.Id())
		}

		data.SetPartial("comment")
	}

	if data.HasChange("is_managed") {
		_, managed := data.GetChange("is_managed")
		var q string
		if managed.(bool) {
			q = builder.Manage()
		} else {
			q = builder.Unmanage()
		}

		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error changing management state on %v", data.Id())
		}

		data.SetPartial("is_managed")
	}

	data.Partial(false)
	if data.HasChange("data_retention_days") {
		_, days := data.GetChange("data_retention_days")

		q := builder.ChangeDataRetentionDays(days.(int))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating data retention days on %v", data.Id())
		}
	}

	return ReadSchema(data, meta)
}

// DeleteSchema implements schema.DeleteFunc
func DeleteSchema(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	schemaID, err := schemaIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := schemaID.DatabaseName
	schema := schemaID.SchemaName

	q := snowflake.Schema(schema).WithDB(dbName).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting schema %v", data.Id())
	}

	data.SetId("")

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
