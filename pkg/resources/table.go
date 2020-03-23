package resources

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"reflect"

	"github.com/FindHotel/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

var space = regexp.MustCompile(`\s+`)

var tableSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the table; must be unique for the schema in which the table is created. Don't use the | character.",
	},
	"database": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the table. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "PUBLIC",
		Description: "The schema in which to create the table. Don't use the | character.",
		ForceNew:    true,
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the table.",
	},
	"columns": &schema.Schema{
		Type:        schema.TypeMap,
		Elem:        &schema.Schema{Type: schema.TypeString}
		Required:    true,
		Description: "Specifies the column names and column types used to create the table.",
		ForceNew:    true,
		DiffSuppressFunc: func(k, old, new map[string]string, d *schema.ResourceData) bool {
			return reflect.DeepEqual(old, new)
		},
	},
}

// Table returns a pointer to the resource representing a table
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

// CreateTable implements schema.CreateFunc
func CreateTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	schema := data.Get("schema").(string)
	database := data.Get("database").(string)
	columns := data.Get("columns").(map[string]string)

	builder := snowflake.Table(name).WithDB(database).WithSchema(schema).WithColumns(columns)

	// Set optionals
	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := data.GetOk("schema"); ok {
		builder.WithSchema(v.(string))
	}

	q := builder.Create()

	err := DBExec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating table %v", name)
	}

	data.SetId(fmt.Sprintf("%v|%v|%v", database, schema, name))

	return ReadTable(data, meta)
}

// ReadTable implements schema.ReadFunc
func ReadTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName, schema, table, err := splitTableID(data.Id())
	if err != nil {
		return err
	}

	// Check table
	q := snowflake.Table(table).WithDB(dbName).WithSchema(schema).Show()
	row := db.QueryRow(q)
	var createdOn, name, databaseName, schemaName, kind, comment, clusterBy, rows, bytes, owner, retentionTime, changeTracking sql.NullString
	err = row.Scan(&createdOn, &name, &databaseName, &schemaName, &kind, &comment, &clusterBy, &rows, &bytes, &owner, &retentionTime, &changeTracking)
	if err != nil {
		return err
	}

	// TODO turn this into a loop after we switch to scaning in a struct
	err = data.Set("name", name.String)
	if err != nil {
		return err
	}

	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	err = data.Set("schema", schemaName.String)
	if err != nil {
		return err
	}

	// Check table columns
	q := snowflake.Table(table).WithDB(dbName).WithSchema(schema).ShowColumns()
	rows, err := db.Query(q)
	if err != nil {
		return err
	}
	var tableName, schemaName, columnName, dataType, null, default, kind, expression, comment, databaseName, autoincrement sql.NullString
	var columns map[string]string
	for rows.Next() {
		err = row.Scan(&tableName, &schemaName, &columnName, &dataType, &null, &default, &kind, &expression, &comment, &databaseName, &autoincrement)
		if err != nil {
			return err
		}
		// TODO convert dataType object to a string
		columns[columnName.String] = dataType.String
	}

	err = data.Set("columns", columns)
	if err != nil {
		return err
	}
	
	return data.Set("database", databaseName.String)
}

// UpdateTable implements schema.UpdateFunc
func UpdateTable(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	dbName, schema, table, err := splitTableID(data.Id())
	if err != nil {
		return err
	}

	builder := snowflake.Table(table).WithDB(dbName).WithSchema(schema)

	db := meta.(*sql.DB)
	if data.HasChange("name") {
		_, name := data.GetChange("name")

		q := builder.Rename(name.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error renaming table %v", data.Id())
		}

		data.SetId(fmt.Sprintf("%v|%v|%v", dbName, schema, name.(string)))
		data.SetPartial("name")
	}

	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")

		if c := comment.(string); c == "" {
			q := builder.RemoveComment()
			err := DBExec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting comment for table %v", data.Id())
			}
		} else {
			q := builder.ChangeComment(c)
			err := DBExec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error updating comment for table %v", data.Id())
			}
		}

		data.SetPartial("comment")
	}

	data.Partial(false)
	if data.HasChange("is_secure") {
		_, secure := data.GetChange("is_secure")

		if secure.(bool) {
			q := builder.Secure()
			err := DBExec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error setting secure for table %v", data.Id())
			}
		} else {
			q := builder.Unsecure()
			err := DBExec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting secure for table %v", data.Id())
			}
		}
	}

	return ReadTable(data, meta)
}

// DeleteTable implements schema.DeleteFunc
func DeleteTable(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName, schema, table, err := splitTableID(data.Id())
	if err != nil {
		return err
	}

	q := snowflake.Table(table).WithDB(dbName).WithSchema(schema).Drop()

	err = DBExec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting table %v", data.Id())
	}

	data.SetId("")

	return nil
}

// TableExists implements schema.ExistsFunc
func TableExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	dbName, schema, table, err := splitTableID(data.Id())
	if err != nil {
		return false, err
	}

	q := snowflake.Table(table).WithDB(dbName).WithSchema(schema).Show()
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

// splitTableID takes the <database_name>|<schema_name>|<table_name> ID and returns the database
// name, schema name and table name.
func splitTableID(v string) (string, string, string, error) {
	arr := strings.Split(v, "|")
	if len(arr) != 3 {
		return "", "", "", fmt.Errorf("ID %v is invalid", v)
	}

	return arr[0], arr[1], arr[2], nil
}
