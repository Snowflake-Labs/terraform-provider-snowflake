package resources

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/jmoiron/sqlx"
)

// warehouse is a go representation of a grant that can be used in conjunction
// with github.com/jmoiron/sqlx
type warehouse struct {
	Name            string    `db:"name"`
	State           string    `db:"state"`
	Type            string    `db:"type"`
	Size            string    `db:"size"`
	MinClusterCount int64     `db:"min_cluster_count"`
	MaxClusterCount int64     `db:"max_cluster_count"`
	StartedClusters int64     `db:"started_clusters"`
	Running         int64     `db:"running"`
	Queued          int64     `db:"queued"`
	IsDefault       string    `db:"is_default"`
	IsCurrent       string    `db:"is_current"`
	AutoSuspend     int64     `db:"auto_suspend"`
	AutoResume      bool      `db:"auto_resume"`
	Available       string    `db:"available"`
	Provisioning    string    `db:"provisioning"`
	Quiescing       string    `db:"quiescing"`
	Other           string    `db:"other"`
	CreatedOn       time.Time `db:"created_on"`
	ResumedOn       time.Time `db:"resumed_on"`
	UpdatedOn       time.Time `db:"updated_on"`
	Owner           string    `db:"owner"`
	Comment         string    `db:"comment"`
	ResourceMonitor string    `db:"resource_monitor"`
	Actives         int64     `db:"actives"`
	Pendings        int64     `db:"pendings"`
	Failed          int64     `db:"failed"`
	Suspended       int64     `db:"suspended"`
	UUID            string    `db:"uuid"`
	ScalingPolicy   string    `db:"scaling_policy"`
}

// warehouseCreateProperties are only available via the CREATE statement
var warehouseCreateProperties = []string{"initially_suspended", "wait_for_provisioning"}

var warehouseProperties = []string{
	"comment", "warehouse_size", "max_cluster_count", "min_cluster_count",
	"scaling_policy", "auto_suspend", "auto_resume",
	"resource_monitor",
}

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
		ValidateFunc: validation.StringInSlice([]string{
			"XSMALL", "X-SMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE",
			"X-LARGE", "XXLARGE", "X2LARGE", "2X-LARGE", "XXXLARGE", "X3LARGE",
			"3X-LARGE", "X4LARGE", "4X-LARGE",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
			}
			return normalize(old) == normalize(new)
		},
	},
	"max_cluster_count": &schema.Schema{
		Type:         schema.TypeInt,
		Description:  "Specifies the maximum number of server clusters for the warehouse.",
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"min_cluster_count": &schema.Schema{
		Type:         schema.TypeInt,
		Description:  "Specifies the minimum number of server clusters for the warehouse (only applies to multi-cluster warehouses).",
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"scaling_policy": &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode.",
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.StringInSlice([]string{"STANDARD", "ECONOMY"}, true),
	},
	"auto_suspend": &schema.Schema{
		Type:         schema.TypeInt,
		Description:  "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntAtLeast(60),
	},
	// @TODO add a disable_auto_suspend property that sets the value of auto_suspend to NULL
	"auto_resume": &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it.",
		Optional:    true,
		Computed:    true,
	},
	"initially_suspended": &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specifies whether the warehouse is created initially in the ‘Suspended’ state.",
		Optional:    true,
	},
	"resource_monitor": &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
		Optional:    true,
		Computed:    true,
	},
	"wait_for_provisioning": &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specifies whether the warehouse, after being resized, waits for all the servers to provision before executing any queued or new queries.",
		Optional:    true,
	},
}

// Warehouse returns a pointer to the resource representing a warehouse
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

// CreateWarehouse implements schema.CreateFunc
func CreateWarehouse(data *schema.ResourceData, meta interface{}) error {
	props := append(warehouseProperties, warehouseCreateProperties...)
	return CreateResource("warehouse", props, warehouseSchema, snowflake.Warehouse, ReadWarehouse)(data, meta)
}

// ReadWarehouse implements schema.ReadFunc
func ReadWarehouse(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	sdb := sqlx.NewDb(db, "snowflake")

	stmt := snowflake.Warehouse(data.Id()).Show()

	log.Printf("[DEBUG] stmt %v\n", stmt)
	row := sdb.QueryRowx(stmt)

	w := &warehouse{}

	err := row.StructScan(w)
	if err != nil {
		return err
	}

	err = data.Set("name", w.Name)
	if err != nil {
		return err
	}
	err = data.Set("comment", w.Comment)
	if err != nil {
		return err
	}
	err = data.Set("warehouse_size", w.Size)
	if err != nil {
		return err
	}
	err = data.Set("max_cluster_count", w.MaxClusterCount)
	if err != nil {
		return err
	}
	err = data.Set("min_cluster_count", w.MinClusterCount)
	if err != nil {
		return err
	}
	err = data.Set("scaling_policy", w.ScalingPolicy)
	if err != nil {
		return err
	}
	err = data.Set("auto_suspend", w.AutoSuspend)
	if err != nil {
		return err
	}
	err = data.Set("auto_resume", w.AutoResume)
	if err != nil {
		return err
	}
	err = data.Set("resource_monitor", w.ResourceMonitor)

	return err
}

// UpdateWarehouse implements schema.UpdateFunc
func UpdateWarehouse(data *schema.ResourceData, meta interface{}) error {
	return UpdateResource("warehouse", warehouseProperties, warehouseSchema, snowflake.Warehouse, ReadWarehouse)(data, meta)
}

// DeleteWarehouse implements schema.DeleteFunc
func DeleteWarehouse(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("warehouse", snowflake.Warehouse)(data, meta)
}
