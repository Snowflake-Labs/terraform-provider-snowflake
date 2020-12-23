package resources

import (
	"database/sql"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
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

		Schema: tableSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateTable implements schema.CreateFunc
func CreateTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)

	// This type conversion is due to the test framework in the terraform-plugin-sdk having limited support
	// for data types in the HCL2ValueFromConfigValue method.
	columns := []snowflake.Column{}

	for _, column := range d.Get("column").([]interface{}) {
		typed := column.(map[string]interface{})
		columnDef := snowflake.Column{}
		columnDef.WithName(typed["name"].(string)).WithType(typed["type"].(string))
		columns = append(columns, columnDef)
	}
	builder := snowflake.TableWithColumnDefinitions(name, database, schema, columns)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	stmt := builder.Create()
	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating table %v", name)
	}

	tableID := &schemaScopedID{
		Database: database,
		Schema:   schema,
		Name:     name,
	}
	dataIDInput, err := tableID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadTable(d, meta)
}

// ReadTable implements schema.ReadFunc
func ReadTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tableID, err := idFromString(d.Id())
	if err != nil {
		return err
	}
	builder := snowflake.Table(tableID.Name, tableID.Database, tableID.Schema)

	row := snowflake.QueryRow(db, builder.Show())
	table, err := snowflake.ScanTable(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] table (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	// Describe the table to read the cols
	tableDescriptionRows, err := snowflake.Query(db, builder.ShowColumns())
	if err != nil {
		return err
	}

	tableDescription, err := snowflake.ScanTableDescription(tableDescriptionRows)
	if err != nil {
		return err
	}

	// Set the relevant data in the state
	toSet := map[string]interface{}{
		"name":     table.TableName.String,
		"owner":    table.Owner.String,
		"database": tableID.Database,
		"schema":   tableID.Schema,
		"comment":  table.Comment.String,
		"column":   snowflake.NewColumns(tableDescription).Flatten(),
	}

	for key, val := range toSet {
		err = d.Set(key, val) //lintignore:R001
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateTable implements schema.UpdateFunc
func UpdateTable(d *schema.ResourceData, meta interface{}) error {
	tableID, err := idFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tableID.Database
	schema := tableID.Schema
	tableName := tableID.Name

	builder := snowflake.Table(tableName, dbName, schema)

	db := meta.(*sql.DB)
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating table comment on %v", d.Id())
		}
	}

	return ReadTable(d, meta)
}

// DeleteTable implements schema.DeleteFunc
func DeleteTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tableID, err := idFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tableID.Database
	schema := tableID.Schema
	tableName := tableID.Name

	q := snowflake.Table(tableName, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting pipe %v", d.Id())
	}

	d.SetId("")

	return nil
}

// TableExists implements schema.ExistsFunc
func TableExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	tableID, err := idFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := tableID.Database
	schema := tableID.Schema
	tableName := tableID.Name

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
