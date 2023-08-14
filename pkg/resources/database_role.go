package resources

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

type databaseRoleID struct {
	DatabaseName string
	RoleName     string
}

// String() takes in a databaseRoleID object and returns a pipe-delimited string:
// DatabaseName|RoleName.
func (id *databaseRoleID) String() (string, error) {
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
// and returns a databaseRoleID object.
func databaseRoleIDFromString(stringID string) (*databaseRoleID, error) {
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

	dbRoleResult := &databaseRoleID{
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

	databaseName := dbRoleID.DatabaseName
	roleName := dbRoleID.RoleName

	roles, err := snowflake.ListDatabaseRoles(databaseName, db)
	if err != nil {
		return fmt.Errorf("error listing database roles err = %w", err)
	}

	var databaseRole snowflake.DatabaseRole
	// find the db role we are looking for by matching the name and db name
	for _, dbRole := range roles {
		if strings.EqualFold(dbRole.Name, roleName) {
			databaseRole = dbRole
			databaseRole.DatabaseName = databaseName // needs to be set as it's not included in the SHOW stmt results
			break
		}
	}

	if databaseRole.Name == "" {
		log.Printf("[DEBUG] database role (%v) not found when listing all database roles in database (%v)", roleName, databaseName)
		d.SetId("")
		return nil
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
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	databaseName := d.Get("database").(string)
	roleName := d.Get("name").(string)

	ctx := context.Background()
	objectIdentifier := sdk.NewDatabaseObjectIdentifier(databaseName, roleName)

	createRequest := sdk.NewCreateDatabaseRoleRequest(objectIdentifier)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	err := client.DatabaseRoles.Create(ctx, createRequest)
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))

	return ReadDatabaseRole(d, meta)
}

// UpdateDatabaseRole implements schema.UpdateFunc.
func UpdateDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	dbRoleID, err := databaseRoleIDFromString(d.Id())
	if err != nil {
		return err
	}

	db := meta.(*sql.DB)
	databaseName := dbRoleID.DatabaseName
	roleName := dbRoleID.RoleName
	builder := snowflake.NewDatabaseRoleBuilder(roleName, databaseName)

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

	roleName := dbRoleID.RoleName
	databaseName := dbRoleID.DatabaseName

	q := snowflake.NewDatabaseRoleBuilder(roleName, databaseName).Drop()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting database role %v err = %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}
