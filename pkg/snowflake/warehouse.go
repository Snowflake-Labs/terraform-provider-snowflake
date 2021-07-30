package snowflake

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	Name                            string    `db:"name"`
	State                           string    `db:"state"`
	Type                            string    `db:"type"`
	Size                            string    `db:"size"`
	MinClusterCount                 int64     `db:"min_cluster_count"`
	MaxClusterCount                 int64     `db:"max_cluster_count"`
	StartedClusters                 int64     `db:"started_clusters"`
	Running                         int64     `db:"running"`
	Queued                          int64     `db:"queued"`
	IsDefault                       string    `db:"is_default"`
	IsCurrent                       string    `db:"is_current"`
	AutoSuspend                     int64     `db:"auto_suspend"`
	AutoResume                      bool      `db:"auto_resume"`
	Available                       string    `db:"available"`
	Provisioning                    string    `db:"provisioning"`
	Quiescing                       string    `db:"quiescing"`
	Other                           string    `db:"other"`
	CreatedOn                       time.Time `db:"created_on"`
	ResumedOn                       time.Time `db:"resumed_on"`
	UpdatedOn                       time.Time `db:"updated_on"`
	Owner                           string    `db:"owner"`
	Comment                         string    `db:"comment"`
	ResourceMonitor                 string    `db:"resource_monitor"`
	StatementTimeoutInSeconds       int64     `db:"statement_timeout_in_seconds"`
	StatementQueuedTimeoutInSeconds int64     `db:"statement_queued_timeout_in_seconds"`
	MaxConcurrencyLevel             int64     `db:"max_concurrency_level"`
	Actives                         int64     `db:"actives"`
	Pendings                        int64     `db:"pendings"`
	Failed                          int64     `db:"failed"`
	Suspended                       int64     `db:"suspended"`
	UUID                            string    `db:"uuid"`
	ScalingPolicy                   string    `db:"scaling_policy"`
}

func ScanWarehouse(row *sqlx.Row) (*warehouse, error) {
	w := &warehouse{}
	err := row.StructScan(w)
	return w, err
}

func ListWarehouses(db *sql.DB) ([]warehouse, error) {
	stmt := "SHOW WAREHOUSES"
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []warehouse{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no warehouses found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
