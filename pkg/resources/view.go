package resources

import (
	"database/sql"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var viewSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the view; must be unique for the schema in which the view is created.",
	},
	"is_secure": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies that the view is secure.",
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the view.",
	},
	"statement": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the query used to create the view. Arguments may be interpolated with a ? using the `statement_arguments` field",
	},
	"statement_arguments": &schema.Schema{
		Type:        schema.TypeList,
		Description: "Arguments for `statement` to be interpolated using the SQL engine.",
	},
}

// View returns a pointer to the resource representing a view
func View() *schema.Resource {
	return &schema.Resource{
		Create: CreateView,
		Read:   ReadView,
		Update: UpdateView,
		Delete: DeleteView,
		Exists: ViewExists,

		Schema: viewSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateView implements schema.CreateFunc
func CreateView(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	s := data.Get("statement").(string)
	args := data.Get("statement_arguments").(*schema.Set).List()

	builder := snowflake.View(name).WithStatement(s).WithStatementArgs(args)

	// Set optionals
	if v, ok := data.GetOk("is_secure"); ok && v.(bool) {
		builder.WithSecure()
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	q, args := builder.Create()

	err := DBExec(db, q, args)
	if err != nil {
		return errors.Wrapf(err, "error creating view %v", name)
	}

	data.SetId(name)

	return ReadView(data, meta)
}

// ReadView implements schema.ReadFunc
func ReadView(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt, args := snowflake.View(id).Show()
	row := db.QueryRow(stmt, args)
	var name, createdOn, loginName, displayName, firstName, lastName, email, minsToUnlock, daysToExpiry, comment, mustChangePassword, snowflakeLock, defaultWarehouse, defaultNamespace, defaultRole, extAuthnDuo, extAuthnUID, minsToBypassMfa, owner, lastSuccessLogin, expiresAtTime, lockedUntilTime, hasPassword sql.NullString
	var disabled, hasRsaPublicKey bool
	err := row.Scan(&name, &createdOn, &loginName, &displayName, &firstName, &lastName, &email, &minsToUnlock, &daysToExpiry, &comment, &disabled, &mustChangePassword, &snowflakeLock, &defaultWarehouse, &defaultNamespace, &defaultRole, &extAuthnDuo, &extAuthnUID, &minsToBypassMfa, &owner, &lastSuccessLogin, &expiresAtTime, &lockedUntilTime, &hasPassword, &hasRsaPublicKey)
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

	err = data.Set("login_name", loginName.String)
	if err != nil {
		return err
	}

	err = data.Set("disabled", disabled)
	if err != nil {
		return err
	}

	err = data.Set("default_role", defaultRole.String)
	if err != nil {
		return err
	}

	err = data.Set("default_namespace", defaultNamespace.String)
	if err != nil {
		return err
	}

	err = data.Set("default_warehouse", defaultWarehouse.String)
	if err != nil {
		return err
	}

	err = data.Set("has_rsa_public_key", hasRsaPublicKey)

	return err
}
