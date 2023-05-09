package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

type externalTableID struct {
	DatabaseName      string
	SchemaName        string
	ExternalTableName string
}

// String() takes in a externalTableID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|ExternalTableName.
func (si *externalTableID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = externalTableIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.ExternalTableName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strExternalTableID := strings.TrimSpace(buf.String())
	return strExternalTableID, nil
}

// externalTableIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|ExternalTableName
// and returns a externalTableID object.
func externalTableIDFromString(stringID string) (*externalTableID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = externalTableIDDelimiter
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

	externalTableResult := &externalTableID{
		DatabaseName:      lines[0][0],
		SchemaName:        lines[0][1],
		ExternalTableName: lines[0][2],
	}
	return externalTableResult, nil
}

// CreateExternalTable implements schema.CreateFunc.
func CreateExternalTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	dbSchema := d.Get("schema").(string)
	name := d.Get("name").(string)

	// This type conversion is due to the test framework in the terraform-plugin-sdk having limited support
	// for data types in the HCL2ValueFromConfigValue method.
	columns := []map[string]string{}
	for _, column := range d.Get("column").([]interface{}) {
		columnDef := map[string]string{}
		for key, val := range column.(map[string]interface{}) {
			columnDef[key] = val.(string)
		}
		columns = append(columns, columnDef)
	}
	builder := snowflake.NewExternalTableBuilder(name, database, dbSchema)
	builder.WithColumns(columns)
	builder.WithFileFormat(d.Get("file_format").(string))
	builder.WithLocation(d.Get("location").(string))

	builder.WithAutoRefresh(d.Get("auto_refresh").(bool))
	builder.WithRefreshOnCreate(d.Get("refresh_on_create").(bool))
	builder.WithCopyGrants(d.Get("copy_grants").(bool))

	// Set optionals
	if v, ok := d.GetOk("partition_by"); ok {
		partitionBys := expandStringList(v.([]interface{}))
		builder.WithPartitionBys(partitionBys)
	}

	if v, ok := d.GetOk("pattern"); ok {
		builder.WithPattern(v.(string))
	}

	if v, ok := d.GetOk("aws_sns_topic"); ok {
		builder.WithAwsSNSTopic(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	stmt := builder.Create()
	if err := snowflake.Exec(db, stmt); err != nil {
		return fmt.Errorf("error creating externalTable %v err = %w", name, err)
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
	d.SetId(dataIDInput)

	return ReadExternalTable(d, meta)
}

// ReadExternalTable implements schema.ReadFunc.
func ReadExternalTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	externalTableID, err := externalTableIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := externalTableID.DatabaseName
	schema := externalTableID.SchemaName
	name := externalTableID.ExternalTableName

	stmt := snowflake.NewExternalTableBuilder(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, stmt)
	externalTable, err := snowflake.ScanExternalTable(row)
	if err != nil {
		if err.Error() == snowflake.ErrNoRowInRS {
			log.Printf("[DEBUG] external table (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("name", externalTable.ExternalTableName.String); err != nil {
		return err
	}

	if err := d.Set("owner", externalTable.Owner.String); err != nil {
		return err
	}
	return nil
}

// UpdateExternalTable implements schema.UpdateFunc.
func UpdateExternalTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	dbSchema := d.Get("schema").(string)
	name := d.Get("name").(string)

	builder := snowflake.NewExternalTableBuilder(name, database, dbSchema)

	if d.HasChange("tag") {
		v := d.Get("tag")
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	stmt := builder.Update()
	if err := snowflake.Exec(db, stmt); err != nil {
		return fmt.Errorf("error updating externalTable %v err = %w", name, err)
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
	d.SetId(dataIDInput)

	return ReadExternalTable(d, meta)
}

// DeleteExternalTable implements schema.DeleteFunc.
func DeleteExternalTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	externalTableID, err := externalTableIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := externalTableID.DatabaseName
	schema := externalTableID.SchemaName
	externalTableName := externalTableID.ExternalTableName

	q := snowflake.NewExternalTableBuilder(externalTableName, dbName, schema).Drop()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting pipe %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
