package resources

import (
	"database/sql"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func CreateResource(
	t string,
	properties []string,
	s map[string]*schema.Schema,
	builder func(string) *snowflake.Builder,
	read func(*schema.ResourceData, interface{}) error,
) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		db := meta.(*sql.DB)
		name := d.Get("name").(string)

		qb := builder(name).Create()

		for _, field := range properties {
			val, ok := d.GetOk(field)
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
		if v, ok := d.GetOk("tag"); ok {
			tags := getTags(v)
			qb.SetTags(tags.toSnowflakeTagValues())
		}
		err := snowflake.Exec(db, qb.Statement())

		if err != nil {
			return errors.Wrapf(err, "error creating %s", t)
		}

		d.SetId(name)

		return read(d, meta)
	}
}

func UpdateResource(
	t string,
	properties []string,
	s map[string]*schema.Schema,
	builder func(string) *snowflake.Builder,
	read func(*schema.ResourceData, interface{}) error,
) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		db := meta.(*sql.DB)
		if d.HasChange("name") {
			// I wish this could be done on one line.
			oldNameI, newNameI := d.GetChange("name")
			oldName := oldNameI.(string)
			newName := newNameI.(string)

			stmt := builder(oldName).Rename(newName)

			err := snowflake.Exec(db, stmt)
			if err != nil {
				return errors.Wrapf(err, "error renaming %s %s to %s", t, oldName, newName)
			}
			d.SetId(newName)
		}

		changes := []string{}
		for _, prop := range properties {
			if d.HasChange(prop) {
				changes = append(changes, prop)
			}
		}
		if len(changes) > 0 {
			name := d.Get("name").(string)
			qb := builder(name).Alter()

			for _, field := range changes {
				val := d.Get(field)
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
			if d.HasChange("tag") {
				log.Printf("[DEBUG] updating tags")
				v := d.Get("tag")
				tags := getTags(v)
				qb.SetTags(tags.toSnowflakeTagValues())
			}

			err := snowflake.Exec(db, qb.Statement())
			if err != nil {
				return errors.Wrapf(err, "error altering %s", t)
			}
		}
		log.Printf("[DEBUG] performing read")
		return read(d, meta)
	}
}

func DeleteResource(t string, builder func(string) *snowflake.Builder) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		db := meta.(*sql.DB)
		name := d.Get("name").(string)

		stmt := builder(name).Drop()

		err := snowflake.Exec(db, stmt)
		if err != nil {
			return errors.Wrapf(err, "error dropping %s %s", t, name)
		}

		d.SetId("")
		return nil
	}
}
