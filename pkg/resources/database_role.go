package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	databaseRoleIDDelimiter = '|'
)

var databaseRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the database role.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the database role.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database role.",
	},
}

type dbRoleID struct {
	DatabaseName string
	RoleName     string
}

// String() takes in a dbRoleID object and returns a pipe-delimited string:
// DatabaseName|RoleName.
func (id *dbRoleID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = databaseRoleIDDelimiter
	dataIdentifiers := [][]string{{id.DatabaseName, id.RoleName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strDbRoleID := strings.TrimSpace(buf.String())
	return strDbRoleID, nil
}

// databaseRoleIDFromString() takes in a pipe-delimited string: DatabaseName|RoleName
// and returns a dbRoleID object.
func databaseRoleIDFromString(stringID string) (*dbRoleID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = pipeIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per database role")
	}
	if len(lines[0]) != 2 {
		return nil, fmt.Errorf("2 fields allowed")
	}

	dbRoleResult := &dbRoleID{
		DatabaseName: lines[0][0],
		RoleName:     lines[0][1],
	}
	return dbRoleResult, nil
}

// DatabaseRole returns a pointer to the resource representing a database role.
func DatabaseRole() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabaseRole,
		Read:   ReadDatabaseRole,
		Update: UpdateDatabaseRole,
		Delete: DeleteDatabaseRole,

		Schema: databaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadDatabaseRole implements schema.ReadFunc.
func ReadDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbRoleID, err := databaseRoleIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := dbRoleID.DatabaseName
	name := dbRoleID.RoleName

	builder := snowflake.NewDatabaseRoleBuilder(name, database)
	qry := builder.Show()
	row := snowflake.QueryRow(db, qry)
	// FIXME scan for name as there is LIKE pattern syntax
	databaseRole, err := snowflake.ScanDatabaseRole(row)
	databaseRole.DatabaseName = database // Nasty.
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] database role (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("name", databaseRole.Name); err != nil {
		return err
	}

	if err := d.Set("database", databaseRole.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("comment", databaseRole.Comment); err != nil {
		return err
	}

	return nil
}

// CreateDatabaseRole implements schema.CreateFunc.
func CreateDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	var err error
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	name := d.Get("name").(string)

	builder := snowflake.NewDatabaseRoleBuilder(name, database)
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	q := builder.Create()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error creating database role %v err = %w", name, err)
	}

	dbRoleID := &dbRoleID{
		DatabaseName: database,
		RoleName:     name,
	}
	dataIDInput, err := dbRoleID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadDatabaseRole(d, meta)
}

// UpdateDatabaseRole implements schema.UpdateFunc.
func UpdateDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	dbRoleID, err := databaseRoleIDFromString(d.Id())
	if err != nil {
		return err
	}

	db := meta.(*sql.DB)
	database := dbRoleID.DatabaseName
	name := dbRoleID.RoleName
	builder := snowflake.NewDatabaseRoleBuilder(name, database)

	if d.HasChange("comment") {
		var q string
		_, newVal := d.GetChange("comment")
		q = builder.ChangeComment(newVal.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating comment on database role %v", d.Id())
		}
	}

	return ReadDatabaseRole(d, meta)
}

// DeleteDatabaseRole implements schema.DeleteFunc.
func DeleteDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbRoleID, err := databaseRoleIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := dbRoleID.DatabaseName
	name := dbRoleID.RoleName

	q := snowflake.NewDatabaseRoleBuilder(name, database).Drop()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting database role %v err = %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}
