package resources

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jmoiron/sqlx"
)

var warehouseProperties = []string{"comment", "warehouse_size"}
var warehouseSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"comment": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"warehouse_size": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			// TODO
			return
		},
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
			}
			return normalize(old) == normalize(new)
		},
	},
}

func Warehouse() *schema.Resource {
	return &schema.Resource{
		Create: CreateWarehouse,
		Read:   ReadWarehouse,
		Delete: DeleteWarehouse,
		Update: UpdateWarehouse,

		Schema: warehouseSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateWarehouse(data *schema.ResourceData, meta interface{}) error {
	return CreateResource("warehouse", warehouseProperties, warehouseSchema, snowflake.Warehouse, ReadWarehouse)(data, meta)

}

type warehouse struct {
	Name            sql.NullString `db:"name"`
	State           sql.NullString `db:"state"`
	Warehousetype   sql.NullString `db:"type"`
	Size            sql.NullString `db:"size"`
	MinClusterCount sql.NullString `db:"min_cluster_count"`
	MaxClusterCount sql.NullString `db:"max_cluster_count"`
	StartedClusters sql.NullString `db:"started_clusters"`
	Running         sql.NullString `db:"running"`
	Queued          sql.NullString `db:"queued"`
	IsDefault       sql.NullString `db:"is_default"`
	IsCurrent       sql.NullString `db:"is_current"`
	AutoSuspend     sql.NullString `db:"auto_suspend"`
	AutoResume      sql.NullString `db:"auto_resume"`
	Available       sql.NullString `db:"available"`
	Provisioning    sql.NullString `db:"provisioning"`
	Quiescing       sql.NullString `db:"quiescing"`
	Other           sql.NullString `db:"other"`
	CreatedOn       sql.NullString `db:"created_on"`
	ResumedOn       sql.NullString `db:"resumed_on"`
	UpdatedOn       sql.NullString `db:"updated_on"`
	Owner           sql.NullString `db:"owner"`
	Comment         sql.NullString `db:"comment"`
	ResourceMonitor sql.NullString `db:"resource_monitor"`
	Actives         sql.NullString `db:"actives"`
	Pendings        sql.NullString `db:"pendings"`
	Failed          sql.NullString `db:"failed"`
	Suspended       sql.NullString `db:"suspended"`
	Uuid            sql.NullString `db:"uuid"`
	ScalingPolidy   sql.NullString `db:"scaling_policy"`
}

func ReadWarehouse(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	sdb := sqlx.NewDb(db, "snowflake")

	name := data.Id()

	stmt := snowflake.Warehouse(name).Show()

	row := sdb.QueryRowx(stmt)
	warehouse := &warehouse{}
	err := row.StructScan(warehouse)
	if err != nil {
		return err
	}

	err = data.Set("name", warehouse.Name.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", warehouse.Comment.String)
	if err != nil {
		return err
	}

	err = data.Set("warehouse_size", warehouse.Size.String)

	return err
}

func UpdateWarehouse(data *schema.ResourceData, meta interface{}) error {
	return UpdateResource("warehouse", warehouseProperties, warehouseSchema, snowflake.Warehouse, ReadWarehouse)(data, meta)
}

func DeleteWarehouse(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("warehouse", snowflake.Warehouse)(data, meta)
}

func SetValidationFunc(set map[string]struct{}) func(val interface{}, key string) ([]string, []error) {
	keys := reflect.ValueOf(set).MapKeys()
	return func(val interface{}, key string) (warns []string, errors []error) {
		s := val.(string)
		_, ok := set[s]
		if !ok {
			errors = append(errors, fmt.Errorf("%s is not in {%v}", s, keys))
		}
		return
	}
}
