package resources

import (
	"database/sql"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var shareProperties = []string{
	"comment",
}

var shareSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the share; must be unique for the account in which the share is created.",
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the managed account.",
	},
}

// Share returns a pointer to the resource representing a share
func Share() *schema.Resource {
	return &schema.Resource{
		Create: CreateShare,
		Read:   ReadShare,
		Update: UpdateShare,
		Delete: DeleteShare,
		Exists: ShareExists,

		Schema: shareSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateShare implements schema.CreateFunc
func CreateShare(data *schema.ResourceData, meta interface{}) error {
	return CreateResource(
		"this does not seem to be used",
		shareProperties,
		shareSchema,
		snowflake.Share,
		ReadShare,
	)(data, meta)
}

// ReadShare implements schema.ReadFunc
func ReadShare(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.Share(id).Show()
	row := db.QueryRow(stmt)

	var createdOn, kind, name, databaseName, to, owner, comment sql.NullString
	err := row.Scan(&createdOn, &kind, &name, &databaseName, &to, &owner, &comment)
	if err != nil {
		return err
	}

	// TODO turn this into a loop after we switch to scaning in a struct
	err = data.Set("name", name.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", comment.String)

	return err
}

// UpdateShare implements schema.UpdateFunc
func UpdateShare(data *schema.ResourceData, meta interface{}) error {
	return UpdateResource("this does not seem to be used", shareProperties, shareSchema, snowflake.Share, ReadShare)(data, meta)
}

// DeleteShare implements schema.DeleteFunc
func DeleteShare(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("this does not seem to be used", snowflake.User)(data, meta)
}

// ShareExists implements schema.ExistsFunc
func ShareExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.Share(id).Show()
	rows, err := db.Query(stmt)
	if err != nil {
		return false, err
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
}
