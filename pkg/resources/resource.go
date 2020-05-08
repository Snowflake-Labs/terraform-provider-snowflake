package resources

import (
	"database/sql"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

func CreateResource(
	t string,
	properties []string,
	s map[string]*schema.Schema,
	builder func(string) *snowflake.Builder,
	read func(*schema.ResourceData, interface{}) error,
) func(*schema.ResourceData, interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		db := meta.(*sql.DB)
		name := data.Get("name").(string)

		qb := builder(name).Create()

		for _, field := range properties {
			val, ok := data.GetOk(field)
			if ok {
				switch s[field].Type {
				case schema.TypeString:
					valStr := val.(string)
					qb.SetString(field, valStr)
				case schema.TypeBool:
					valBool := val.(bool)
					qb.SetBool(field, valBool)
				case schema.TypeInt:
					valInt := val.(int)
					qb.SetInt(field, valInt)
				}
			}
		}
		err := snowflake.Exec(db, qb.Statement())

		if err != nil {
			return errors.Wrapf(err, "error creating %s", t)
		}

		data.SetId(name)

		return read(data, meta)
	}
}

func UpdateResource(
	t string,
	properties []string,
	s map[string]*schema.Schema,
	builder func(string) *snowflake.Builder,
	read func(*schema.ResourceData, interface{}) error,
) func(*schema.ResourceData, interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
		data.Partial(true)

		db := meta.(*sql.DB)
		if data.HasChange("name") {
			// I wish this could be done on one line.
			oldNameI, newNameI := data.GetChange("name")
			oldName := oldNameI.(string)
			newName := newNameI.(string)

			stmt := builder(oldName).Rename(newName)

			err := snowflake.Exec(db, stmt)
			if err != nil {
				return errors.Wrapf(err, "error renaming %s %s to %s", t, oldName, newName)
			}

			data.SetId(newName)
			data.SetPartial("name")
		}
		data.Partial(false)

		changes := []string{}

		for _, prop := range properties {
			if data.HasChange(prop) {
				changes = append(changes, prop)
			}
		}
		if len(changes) > 0 {
			name := data.Get("name").(string)
			qb := builder(name).Alter()

			for _, field := range changes {
				val := data.Get(field)
				switch s[field].Type {
				case schema.TypeString:
					valStr := val.(string)
					qb.SetString(field, valStr)
				case schema.TypeBool:
					valBool := val.(bool)
					qb.SetBool(field, valBool)
				case schema.TypeInt:
					valInt := val.(int)
					qb.SetInt(field, valInt)
				}
			}

			err := snowflake.Exec(db, qb.Statement())
			if err != nil {
				return errors.Wrapf(err, "error altering %s", t)
			}
		}
		return read(data, meta)
	}
}

func DeleteResource(t string, builder func(string) *snowflake.Builder) func(*schema.ResourceData, interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		db := meta.(*sql.DB)
		name := data.Get("name").(string)

		stmt := builder(name).Drop()

		err := snowflake.Exec(db, stmt)
		if err != nil {
			return errors.Wrapf(err, "error dropping %s %s", t, name)
		}

		data.SetId("")
		return nil
	}
}
