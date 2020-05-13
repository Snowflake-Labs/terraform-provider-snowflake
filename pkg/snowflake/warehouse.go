package snowflake

import (
	"time"

	"github.com/jmoiron/sqlx"
)

func Warehouse(name string) *Builder {
	return &Builder{
		name:       name,
		entityType: WarehouseType,
	}
}

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

func ScanWarehouse(row *sqlx.Row) (*warehouse, error) {
	w := &warehouse{}
	err := row.StructScan(w)
	return w, err
}
