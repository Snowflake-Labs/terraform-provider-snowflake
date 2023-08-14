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
	client := sdk.NewClientFromDB(db)

	// TODO: what to do with decode snowflake id method in case of SchemaIdentifier and DatabaseObjectIdentifier?
	schemaIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaIdentifier)
	objectIdentifier := sdk.NewDatabaseObjectIdentifier(schemaIdentifier.DatabaseName(), schemaIdentifier.Name())

	ctx := context.Background()
	databaseRole, err := client.DatabaseRoles.ShowByID(ctx, objectIdentifier)
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] database role (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("name", databaseRole.Name); err != nil {
		return err
	}

	if err := d.Set("database", objectIdentifier.DatabaseName()); err != nil {
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

	objectIdentifier := sdk.NewDatabaseObjectIdentifier(databaseName, roleName)
	createRequest := sdk.NewCreateDatabaseRoleRequest(objectIdentifier)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	ctx := context.Background()
	err := client.DatabaseRoles.Create(ctx, createRequest)
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))

	return ReadDatabaseRole(d, meta)
}

// UpdateDatabaseRole implements schema.UpdateFunc.
func UpdateDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	schemaIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaIdentifier)
	objectIdentifier := sdk.NewDatabaseObjectIdentifier(schemaIdentifier.DatabaseName(), schemaIdentifier.Name())

	if d.HasChange("comment") {
		_, newVal := d.GetChange("comment")

		ctx := context.Background()
		alterRequest := sdk.NewAlterDatabaseRoleRequest(objectIdentifier).WithSet(sdk.NewDatabaseRoleSetRequest(newVal.(string)))
		err := client.DatabaseRoles.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating database role %v: %w", objectIdentifier.Name(), err)
		}
	}

	return ReadDatabaseRole(d, meta)
}

// DeleteDatabaseRole implements schema.DeleteFunc.
func DeleteDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	schemaIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaIdentifier)
	objectIdentifier := sdk.NewDatabaseObjectIdentifier(schemaIdentifier.DatabaseName(), schemaIdentifier.Name())

	ctx := context.Background()
	dropRequest := sdk.NewDropDatabaseRoleRequest(objectIdentifier)
	err := client.DatabaseRoles.Drop(ctx, dropRequest)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
