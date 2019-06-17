package resources

import (
	"database/sql"
	"fmt"

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
		ForceNew:    true,
	},
	"statement_arguments": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Arguments for `statement` to be interpolated using the SQL engine.",
		ForceNew:    true,
		Optional:    true,
	},
	"view_text": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The interpolated text of the SQL statement for this view.",
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

	fmt.Printf("[DEBUG] args: %+v\n", args)

	builder := snowflake.View(name).WithStatement(s).WithStatementArgs(args)

	// Set optionals
	if v, ok := data.GetOk("is_secure"); ok && v.(bool) {
		builder.WithSecure()
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	q, qArgs := builder.Create()

	err := DBExec(db, q, qArgs...)
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
	row := db.QueryRow(stmt, args...)
	var createdOn, name, reserved, databaseName, schemaName, owner, comment, text sql.NullString
	var isSecure bool
	err := row.Scan(&createdOn, &name, &reserved, &databaseName, &schemaName, &owner, &comment, &text, &isSecure)
	if err != nil {
		return err
	}

	// TODO turn this into a loop after we switch to scaning in a struct
	err = data.Set("name", name.String)
	if err != nil {
		return err
	}

	err = data.Set("is_secure", isSecure)
	if err != nil {
		return err
	}

	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	err = data.Set("view_text", text.String)
	if err != nil {
		return err
	}

	return err
}

// UpdateView implements schema.UpdateFunc
func UpdateView(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	db := meta.(*sql.DB)
	if data.HasChange("name") {
		_, name := data.GetChange("name")

		q, args := snowflake.View(data.Id()).Rename(name.(string))
		err := DBExec(db, q, args)
		if err != nil {
			return errors.Wrapf(err, "error renaming view %v", data.Id())
		}

		data.SetId(name.(string))
		data.SetPartial("name")
	}

	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")

		if c := comment.(string); c == "" {
			q, args := snowflake.View(data.Id()).RemoveComment()
			err := DBExec(db, q, args)
			if err != nil {
				return errors.Wrapf(err, "error unsetting comment for view %v", data.Id())
			}
		} else {
			q, args := snowflake.View(data.Id()).ChangeComment(c)
			err := DBExec(db, q, args)
			if err != nil {
				return errors.Wrapf(err, "error updating comment for view %v", data.Id())
			}
		}

		data.SetPartial("comment")
	}

	data.Partial(false)
	if data.HasChange("is_secure") {
		_, secure := data.GetChange("is_secure")

		if secure.(bool) {
			q, args := snowflake.View(data.Id()).Secure()
			err := DBExec(db, q, args)
			if err != nil {
				return errors.Wrapf(err, "error setting secure for view %v", data.Id())
			}
		} else {
			q, args := snowflake.View(data.Id()).Unsecure()
			err := DBExec(db, q, args)
			if err != nil {
				return errors.Wrapf(err, "error unsetting secure for view %v", data.Id())
			}
		}
	}

	return ReadView(data, meta)
}

// DeleteView implements schema.DeleteFunc
func DeleteView(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	q, args := snowflake.View(data.Id()).Drop()
	err := DBExec(db, q, args)
	if err != nil {
		return errors.Wrapf(err, "error deleting view %v", data.Id())
	}

	data.SetId("")

	return nil
}

// ViewExists implements schema.ExistsFunc
func ViewExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)

	q, args := snowflake.View(data.Id()).Show()
	rows, err := db.Query(q, args)
	if err != nil {
		return false, err
	}

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
