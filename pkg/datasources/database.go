package datasources

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
)

var databaseSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return its metadata.",
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_default": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_current": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"origin": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"retention_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"options": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

// Database the Snowflake Database resource.
func Database() *schema.Resource {
	return &schema.Resource{
		Read:   ReadDatabase,
		Schema: databaseSchema,
	}
}

// ReadDatabase read the database meta-data information.
func ReadDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbx := sqlx.NewDb(db, "snowflake")
	log.Printf("[DEBUG] database: %v", d.Get("name"))
	dbData, err := snowflake.ListDatabase(dbx, d.Get("name").(string))
	if err != nil {
		log.Println("[DEBUG] list database failed to decode")
		d.SetId("")
		return nil
	}
	if dbData == nil || !dbData.DBName.Valid {
		log.Println("[DEBUG] database not found")
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] list database: %v", dbData)
	d.SetId(dbData.DBName.String)
	commentErr := d.Set("comment", dbData.Comment.String)
	if commentErr != nil {
		return commentErr
	}
	ownerErr := d.Set("owner", dbData.Owner.String)
	if ownerErr != nil {
		return ownerErr
	}
	isDefaultErr := d.Set("is_default", dbData.IsDefault.String == "Y")
	if isDefaultErr != nil {
		return isDefaultErr
	}
	isCurrentErr := d.Set("is_current", dbData.IsCurrent.String == "Y")
	if isCurrentErr != nil {
		return isCurrentErr
	}
	originErr := d.Set("origin", dbData.Origin.String)
	if originErr != nil {
		return originErr
	}
	createdOnErr := d.Set("created_on", dbData.CreatedOn.String)
	if createdOnErr != nil {
		return createdOnErr
	}
	optionsErr := d.Set("options", dbData.Options.String)
	if optionsErr != nil {
		return optionsErr
	}
	retentionTimeErr := d.Set("retention_time", -1)
	if retentionTimeErr != nil {
		return retentionTimeErr
	}
	if dbData.RetentionTime.Valid {
		v, err := strconv.Atoi(dbData.RetentionTime.String)
		if err == nil {
			retentionTimeErr := d.Set("retention_time", v)
			if retentionTimeErr != nil {
				return retentionTimeErr
			}
		}
	}
	return nil
}
