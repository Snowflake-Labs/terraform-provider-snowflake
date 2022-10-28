package datasources

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
)

var databasesSchema = map[string]*schema.Schema{
	"databases": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Snowflake databases",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
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
				"replication_configuration": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"accounts": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"ignore_edition_check": {
								Type:     schema.TypeBool,
								Computed: true,
							},
						},
					},
				},
			},
		},
	},
}

// Databases the Snowflake current account resource.
func Databases() *schema.Resource {
	return &schema.Resource{
		Read:   ReadDatabases,
		Schema: databasesSchema,
	}
}

// ReadDatabases read the current snowflake account information.
func ReadDatabases(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbx := sqlx.NewDb(db, "snowflake")
	dbs, err := snowflake.ListDatabases(dbx)
	if err != nil {
		log.Println("[DEBUG] list databases failed to decode")
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] list databases: %v", dbs)
	d.SetId("databases_read")
	databases := []map[string]interface{}{}
	for _, db := range dbs {
		dbR := map[string]interface{}{}
		if !db.DBName.Valid {
			continue
		}
		dbR["name"] = db.DBName.String
		dbR["comment"] = db.Comment.String
		dbR["owner"] = db.Owner.String
		dbR["is_default"] = db.IsDefault.String == "Y"
		dbR["is_current"] = db.IsCurrent.String == "Y"
		dbR["origin"] = db.Origin.String
		dbR["created_on"] = db.CreatedOn.String
		dbR["options"] = db.Options.String
		dbR["retention_time"] = -1
		if db.RetentionTime.Valid {
			v, err := strconv.Atoi(db.RetentionTime.String)
			if err == nil {
				dbR["retention_time"] = v
			}
		}
		databases = append(databases, dbR)

	}
	databasesErr := d.Set("databases", databases)
	if databasesErr != nil {
		return databasesErr
	}
	return nil
}
