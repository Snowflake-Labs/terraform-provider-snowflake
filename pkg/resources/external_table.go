package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

const (
	externalTableIDDelimiter = '|'
)

var externalTableSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the external table; must be unique for the database and schema in which the externalTable is created.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the external table.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the external table.",
	},
	"column": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		ForceNew:    true,
		Description: "Definitions of a column to create in the external table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column name",
					ForceNew:    true,
				},
				"type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column type, e.g. VARIANT",
					ForceNew:    true,
				},
				"as": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "String that specifies the expression for the column. When queried, the column returns results derived from this expression.",
					ForceNew:    true,
				},
			},
		},
	},
	"location": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies a location for the external table.",
	},
	"file_format": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the file format for the external table.",
	},
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the file names and/or paths on the external stage to match.",
	},
	"aws_sns_topic": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the aws sns topic for the external table.",
	},
	"partition_by": {
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		ForceNew:    true,
		Description: "Specifies any partition columns to evaluate for the external table.",
	},
	"refresh_on_create": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies weather to refresh when an external table is created.",
		Default:     true,
		ForceNew:    true,
	},
	"auto_refresh": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to automatically refresh the external table metadata once, immediately after the external table is created.",
		Default:     true,
		ForceNew:    true,
	},
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies to retain the access permissions from the original table when an external table is recreated using the CREATE OR REPLACE TABLE variant",
		Default:     false,
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a comment for the external table.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		ForceNew:    true,
		Description: "Name of the role that owns the external table.",
	},
	"tag": tagReferenceSchema,
}

func ExternalTable() *schema.Resource {
	return &schema.Resource{
		Create: CreateExternalTable,
		Read:   ReadExternalTable,
		Update: UpdateExternalTable,
		Delete: DeleteExternalTable,

		Schema: externalTableSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

type externalTableID struct {
	DatabaseName      string
	SchemaName        string
	ExternalTableName string
}

//String() takes in a externalTableID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|ExternalTableName
func (si *externalTableID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = externalTableIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.ExternalTableName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strExternalTableID := strings.TrimSpace(buf.String())
	return strExternalTableID, nil
}

// externalTableIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|ExternalTableName
// and returns a externalTableID object
func externalTableIDFromString(stringID string) (*externalTableID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = externalTableIDDelimiter
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

	externalTableResult := &externalTableID{
		DatabaseName:      lines[0][0],
		SchemaName:        lines[0][1],
		ExternalTableName: lines[0][2],
	}
	return externalTableResult, nil
}

// CreateExternalTable implements schema.CreateFunc
func CreateExternalTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	dbSchema := data.Get("schema").(string)
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
	builder := snowflake.ExternalTable(name, database, dbSchema)
	builder.WithColumns(columns)
	builder.WithFileFormat(data.Get("file_format").(string))
	builder.WithLocation(data.Get("location").(string))

	builder.WithAutoRefresh(data.Get("auto_refresh").(bool))
	builder.WithRefreshOnCreate(data.Get("refresh_on_create").(bool))
	builder.WithCopyGrants(data.Get("copy_grants").(bool))

	// Set optionals
	if v, ok := data.GetOk("partition_by"); ok {
		partitionBys := expandStringList(v.([]interface{}))
		builder.WithPartitionBys(partitionBys)
	}

	if v, ok := data.GetOk("pattern"); ok {
		builder.WithPattern(v.(string))
	}

	if v, ok := data.GetOk("aws_sns_topic"); ok {
		builder.WithAwsSNSTopic(v.(string))
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := data.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	stmt := builder.Create()
	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating externalTable %v", name)
	}

	externalTableID := &externalTableID{
		DatabaseName:      database,
		SchemaName:        dbSchema,
		ExternalTableName: name,
	}
	dataIDInput, err := externalTableID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadExternalTable(data, meta)
}

// ReadExternalTable implements schema.ReadFunc
func ReadExternalTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	externalTableID, err := externalTableIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := externalTableID.DatabaseName
	schema := externalTableID.SchemaName
	name := externalTableID.ExternalTableName

	stmt := snowflake.ExternalTable(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, stmt)
	externalTable, err := snowflake.ScanExternalTable(row)
	if err != nil {
		return err
	}

	err = data.Set("name", externalTable.ExternalTableName.String)
	if err != nil {
		return err
	}

	err = data.Set("owner", externalTable.Owner.String)
	if err != nil {
		return err
	}

	return nil
}

// UpdateExternalTable implements schema.UpdateFunc
func UpdateExternalTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	dbSchema := data.Get("schema").(string)
	name := data.Get("name").(string)

	builder := snowflake.ExternalTable(name, database, dbSchema)

	if data.HasChange("tag") {
		v := data.Get("tag")
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	stmt := builder.Update()
	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error updating externalTable %v", name)
	}

	externalTableID := &externalTableID{
		DatabaseName:      database,
		SchemaName:        dbSchema,
		ExternalTableName: name,
	}
	dataIDInput, err := externalTableID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadExternalTable(data, meta)
}

// DeleteExternalTable implements schema.DeleteFunc
func DeleteExternalTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	externalTableID, err := externalTableIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := externalTableID.DatabaseName
	schema := externalTableID.SchemaName
	externalTableName := externalTableID.ExternalTableName

	q := snowflake.ExternalTable(externalTableName, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting pipe %v", data.Id())
	}

	data.SetId("")

	return nil
}

// ExternalTableExists implements schema.ExistsFunc
func ExternalTableExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	externalTableID, err := externalTableIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := externalTableID.DatabaseName
	schema := externalTableID.SchemaName
	externalTableName := externalTableID.ExternalTableName

	q := snowflake.ExternalTable(externalTableName, dbName, schema).Show()
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
