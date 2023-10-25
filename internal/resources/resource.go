// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				case schema.TypeSet:
					valList := expandStringList(val.(*schema.Set).List())
					qb.SetStringList(field, valList)
				}
			}
		}
		if v, ok := d.GetOk("tag"); ok {
			tags := getTags(v)
			qb.SetTags(tags.toSnowflakeTagValues())
		}
		if err := snowflake.Exec(db, qb.Statement()); err != nil {
			return fmt.Errorf("error creating %s err = %w", t, err)
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
				return fmt.Errorf("error renaming %s %s to %s err = %w", t, oldName, newName, err)
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
				case schema.TypeSet:
					valList := expandStringList(val.(*schema.Set).List())
					qb.SetStringList(field, valList)
				}
			}
			if d.HasChange("tag") {
				log.Println("[DEBUG] updating tags")
				v := d.Get("tag")
				tags := getTags(v)
				qb.SetTags(tags.toSnowflakeTagValues())
			}

			if err := snowflake.Exec(db, qb.Statement()); err != nil {
				return fmt.Errorf("error altering %s err = %w", t, err)
			}
		}
		log.Println("[DEBUG] performing read")
		return read(d, meta)
	}
}

func DeleteResource(t string, builder func(string) *snowflake.Builder) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		db := meta.(*sql.DB)
		name := d.Get("name").(string)

		stmt := builder(name).Drop()
		if err := snowflake.Exec(db, stmt); err != nil {
			return fmt.Errorf("error dropping %s %s err = %w", t, name, err)
		}

		d.SetId("")
		return nil
	}
}
