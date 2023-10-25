// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/helpers"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk"
)

const (
	schemaIDDelimiter = '|'
)

var schemaSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the schema; must be unique for the database in which the schema is created.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the schema.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the schema.",
	},
	"is_transient": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies a schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
		ForceNew:    true,
	},
	"is_managed": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner.",
	},
	"data_retention_days": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the schema, as well as specifying the default Time Travel retention time for all tables created in the schema.",
		ValidateFunc: validation.IntBetween(0, 90),
	},
	"tag": tagReferenceSchema,
}

// Schema returns a pointer to the resource representing a schema.
func Schema() *schema.Resource {
	return &schema.Resource{
		Create: CreateSchema,
		Read:   ReadSchema,
		Update: UpdateSchema,
		Delete: DeleteSchema,

		Schema: schemaSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSchema implements schema.CreateFunc.
func CreateSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	database := d.Get("database").(string)

	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	err := client.Schemas.Create(ctx, sdk.NewDatabaseObjectIdentifier(database, name), &sdk.CreateSchemaOptions{
		Transient:               GetPropertyAsPointer[bool](d, "is_transient"),
		WithManagedAccess:       GetPropertyAsPointer[bool](d, "is_managed"),
		DataRetentionTimeInDays: GetPropertyAsPointer[int](d, "data_retention_days"),
		Tag:                     getPropertyTags(d, "tag"),
		Comment:                 GetPropertyAsPointer[string](d, "comment"),
	})
	if err != nil {
		return fmt.Errorf("error creating schema %v err = %w", name, err)
	}

	d.SetId(helpers.EncodeSnowflakeID(database, name))

	return ReadSchema(d, meta)
}

// ReadSchema implements schema.ReadFunc.
func ReadSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)

	_, err := client.Databases.ShowByID(ctx, sdk.NewAccountObjectIdentifier(id.DatabaseName()))
	if err != nil {
		d.SetId("")
	}

	s, err := client.Schemas.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	var retentionTime int64
	// "retention_time" may sometimes be empty string instead of an integer
	{
		rt := s.RetentionTime
		if rt == "" {
			rt = "0"
		}

		retentionTime, err = strconv.ParseInt(rt, 10, 64)
		if err != nil {
			return err
		}
	}

	values := map[string]any{
		"name":                s.Name,
		"database":            s.DatabaseName,
		"data_retention_days": retentionTime,
		// reset the options before reading back from the DB
		"is_transient": false,
		"is_managed":   false,
	}
	if s.Comment != nil {
		values["comment"] = *s.Comment
	}

	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	if opts := s.Options; opts != nil && *opts != "" {
		for _, opt := range strings.Split(*opts, ", ") {
			switch opt {
			case "TRANSIENT":
				if err := d.Set("is_transient", true); err != nil {
					return err
				}
			case "MANAGED ACCESS":
				if err := d.Set("is_managed", true); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// UpdateSchema implements schema.UpdateFunc.
func UpdateSchema(d *schema.ResourceData, meta interface{}) error {
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	if d.HasChange("name") {
		newName := d.Get("name")
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			NewName: sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), newName.(string)),
		})
		if err != nil {
			return fmt.Errorf("error updating schema name on %v err = %w", d.Id(), err)
		}
		d.SetId(helpers.EncodeSnowflakeID(id.DatabaseName(), newName))
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Set: &sdk.SchemaSet{
				Comment: sdk.String(comment.(string)),
			},
		})
		if err != nil {
			return fmt.Errorf("error updating schema comment on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("is_managed") {
		managed := d.Get("is_managed")
		var err error
		if managed.(bool) {
			err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
				EnableManagedAccess: sdk.Bool(true),
			})
		} else {
			err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
				DisableManagedAccess: sdk.Bool(true),
			})
		}
		if err != nil {
			return fmt.Errorf("error changing management state on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("data_retention_days") {
		days := d.Get("data_retention_days")
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Set: &sdk.SchemaSet{
				DataRetentionTimeInDays: sdk.Int(days.(int)),
			},
		})
		if err != nil {
			return fmt.Errorf("error updating data retention days on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("tag") {
		o, n := d.GetChange("tag")
		removed, added, changed := getTags(o).diffs(getTags(n))

		unsetTags := make([]sdk.ObjectIdentifier, len(removed))
		for i, t := range removed {
			unsetTags[i] = sdk.NewDatabaseObjectIdentifier(t.database, t.name)
		}
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Unset: &sdk.SchemaUnset{
				Tag: unsetTags,
			},
		})
		if err != nil {
			return fmt.Errorf("error dropping tags on %v", d.Id())
		}

		setTags := make([]sdk.TagAssociation, len(added)+len(changed))
		for i, t := range added {
			setTags[i] = sdk.TagAssociation{
				Name:  sdk.NewSchemaObjectIdentifier(t.database, t.schema, t.name),
				Value: t.value,
			}
		}
		for i, t := range changed {
			setTags[i] = sdk.TagAssociation{
				Name:  sdk.NewSchemaObjectIdentifier(t.database, t.schema, t.name),
				Value: t.value,
			}
		}
		err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Set: &sdk.SchemaSet{
				Tag: setTags,
			},
		})
		if err != nil {
			return fmt.Errorf("error setting tags on %v", d.Id())
		}
	}

	return ReadSchema(d, meta)
}

// DeleteSchema implements schema.DeleteFunc.
func DeleteSchema(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)

	err := client.Schemas.Drop(ctx, id, new(sdk.DropSchemaOptions))
	if err != nil {
		return fmt.Errorf("error deleting schema %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
