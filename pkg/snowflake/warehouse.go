package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type WarehouseBuilder struct {
	*Builder
}

func (wb *WarehouseBuilder) Show() string {
	return wb.Builder.Show()
}

func (wb *WarehouseBuilder) Describe() string {
	return wb.Builder.Describe()
}

func (wb *WarehouseBuilder) Drop() string {
	return wb.Builder.Drop()
}

func (wb *WarehouseBuilder) Rename(newName string) string {
	return wb.Builder.Rename(newName)
}

func (wb *WarehouseBuilder) Alter() *AlterPropertiesBuilder {
	return wb.Builder.Alter()
}

func (wb *WarehouseBuilder) Create() *CreateBuilder {
	return wb.Builder.Create()
}

// ShowParameters returns the query to show the parameters for the warehouse
func (wb *WarehouseBuilder) ShowParameters() string {
	return fmt.Sprintf("SHOW PARAMETERS IN WAREHOUSE %v", wb.Builder.name)
}

func Warehouse(name string) *WarehouseBuilder {
	return &WarehouseBuilder{
		&Builder{
			name:       name,
			entityType: WarehouseType,
		},
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

// warehouseParams struct to represent a row of parameters
type warehouseParams struct {
	Key          string `db:"key"`
	Value        string `db:"value"`
	DefaultValue string `db:"default"`
	Level        string `db:"level"`
	Description  string `db:"description"`
	Type         string `db:"type"`
}

func ScanWarehouse(row *sqlx.Row) (*warehouse, error) {
	w := &warehouse{}
	err := row.StructScan(w)
	return w, err
}

// ScanWarehouseParameters takes a database row and converts it to a warehouse parameter pointer
func ScanWarehouseParameters(rows *sqlx.Rows) ([]*warehouseParams, error) {
	params := []*warehouseParams{}

	for rows.Next() {
		w := &warehouseParams{}
		err := rows.StructScan(w)
		if err != nil {
			return nil, err
		}
		params = append(params, w)
	}
	return params, nil
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
