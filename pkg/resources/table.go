package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

const (
	tableIDDelimiter = '|'
)

var tableSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the table; must be unique for the database and schema in which the table is created.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the table.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the table.",
	},
	"column": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		ForceNew:    true,
		Description: "Definitions of a column to create in the table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column name",
				},
				"type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column type, e.g. VARIANT",
				},
			},
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the table.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the table.",
	},
}

func Table() *schema.Resource {
	return &schema.Resource{
		Create: CreateTable,
		Read:   ReadTable,
		Update: UpdateTable,
		Delete: DeleteTable,
		Exists: TableExists,

		Schema: tableSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

type tableID struct {
	DatabaseName string
	SchemaName   string
	TableName    string
}

//String() takes in a tableID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|TableName
func (si *tableID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = tableIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.TableName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strTableID := strings.TrimSpace(buf.String())
	return strTableID, nil
}

// tableIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TableName
// and returns a tableID object
func tableIDFromString(stringID string) (*tableID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = tableIDDelimiter
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

	tableResult := &tableID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		TableName:    lines[0][2],
	}
	return tableResult, nil
}

// CreateTable implements schema.CreateFunc
func CreateTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	schema := data.Get("schema").(string)
	name := data.Get("name").(string)

	// This type conversion is due to the test framework in the terraform-plugin-sdk having limited support
	// for data types in the HCL2ValueFromConfigValue method.
	columns := []map[string]string{}
	for _, column := range data.Get("column").([]interface{}) {
		columnDef := map[string]string{}
		for key, val := range column.(map[string]interface{}) {
			columnDef[key] = val.(string)
		}
		columns = append(columns, columnDef)
	}
	builder := snowflake.TableWithColumnDefinitions(name, database, schema, columns)

	// Set optionals
	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	stmt := builder.Create()
	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating table %v", name)
	}

	tableID := &tableID{
		DatabaseName: database,
		SchemaName:   schema,
		TableName:    name,
	}
	dataIDInput, err := tableID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadTable(data, meta)
}

// ReadTable implements schema.ReadFunc
func ReadTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tableID, err := tableIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := tableID.DatabaseName
	schema := tableID.SchemaName
	name := tableID.TableName

	stmt := snowflake.Table(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, stmt)
	table, err := snowflake.ScanTable(row)
	if err != nil {
		return err
	}

	err = data.Set("name", table.TableName.String)
	if err != nil {
		return err
	}

	err = data.Set("owner", table.Owner.String)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTable implements schema.UpdateFunc
func UpdateTable(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	tableID, err := tableIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := tableID.DatabaseName
	schema := tableID.SchemaName
	tableName := tableID.TableName

	builder := snowflake.Table(tableName, dbName, schema)

	db := meta.(*sql.DB)
	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating table comment on %v", data.Id())
		}

		data.SetPartial("comment")
	}

	return ReadTable(data, meta)
}

// DeleteTable implements schema.DeleteFunc
func DeleteTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tableID, err := tableIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := tableID.DatabaseName
	schema := tableID.SchemaName
	tableName := tableID.TableName

	q := snowflake.Table(tableName, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting pipe %v", data.Id())
	}

	data.SetId("")

	return nil
}

// TableExists implements schema.ExistsFunc
func TableExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	tableID, err := tableIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := tableID.DatabaseName
	schema := tableID.SchemaName
	tableName := tableID.TableName

	q := snowflake.Table(tableName, dbName, schema).Show()
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
