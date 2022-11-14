package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	ResourceWarehouse  = "WAREHOUSE"
	ResourceWarehouses = "WAREHOUSES"
)

// Compile-time proof of interface implementation.
var _ Warehouses = (*warehouses)(nil)

// Warehouses describes all the warehouses related methods that the
// Snowflake API supports.
type Warehouses interface {
	// List all the warehouses by pattern.
	List(ctx context.Context, options WarehouseListOptions) ([]*Warehouse, error)
	// Create a new warehouse with the given options.
	Create(ctx context.Context, options WarehouseCreateOptions) (*Warehouse, error)
	// Read an warehouse by its name.
	Read(ctx context.Context, warehouse string) (*Warehouse, error)
	// Update attributes of an existing warehouse.
	Update(ctx context.Context, warehouse string, options WarehouseUpdateOptions) (*Warehouse, error)
	// Delete a warehouse by its name.
	Delete(ctx context.Context, warehouse string) error
	// Rename a warehouse name.
	Rename(ctx context.Context, old string, new string) error
	// Use the active/current warehouse for the session.
	Use(ctx context.Context, warehouse string) error
}

// warehouses implements Warehouses
type warehouses struct {
	client *Client
}

type Warehouse struct {
	Name                            string
	State                           string
	Type                            string
	WarehouseSize                   string
	MinClusterCount                 int32
	MaxClusterCount                 int32
	StartedClusters                 int32
	Running                         int32
	Queued                          int32
	IsDefault                       string
	IsCurrent                       string
	AutoSuspend                     int32
	AutoResume                      bool
	Available                       string
	Provisioning                    string
	Quiescing                       string
	Other                           string
	CreatedOn                       time.Time
	ResumedOn                       time.Time
	UpdatedOn                       time.Time
	Owner                           string
	Comment                         string
	EnableQueryAcceleration         bool
	QueryAccelerationMaxScaleFactor int32
	ResourceMonitor                 string
	Actives                         int32
	Pendings                        int32
	Failed                          int32
	Suspended                       int32
	Uuid                            string
	ScalingPolicy                   string
}

type warehouseEntity struct {
	Name                            sql.NullString `db:"name"`
	State                           sql.NullString `db:"state"`
	Type                            sql.NullString `db:"type"`
	WarehouseSize                   sql.NullString `db:"size"`
	MinClusterCount                 sql.NullInt32  `db:"min_cluster_count"`
	MaxClusterCount                 sql.NullInt32  `db:"max_cluster_count"`
	StartedClusters                 sql.NullInt32  `db:"started_clusters"`
	Running                         sql.NullInt32  `db:"running"`
	Queued                          sql.NullInt32  `db:"queued"`
	IsDefault                       sql.NullString `db:"is_default"`
	IsCurrent                       sql.NullString `db:"is_current"`
	AutoSuspend                     sql.NullInt32  `db:"auto_suspend"`
	AutoResume                      sql.NullBool   `db:"auto_resume"`
	Available                       sql.NullString `db:"available"`
	Provisioning                    sql.NullString `db:"provisioning"`
	Quiescing                       sql.NullString `db:"quiescing"`
	Other                           sql.NullString `db:"other"`
	CreatedOn                       sql.NullTime   `db:"created_on"`
	ResumedOn                       sql.NullTime   `db:"resumed_on"`
	UpdatedOn                       sql.NullTime   `db:"updated_on"`
	Owner                           sql.NullString `db:"owner"`
	Comment                         sql.NullString `db:"comment"`
	EnableQueryAcceleration         sql.NullBool   `db:"enable_query_acceleration"`
	QueryAccelerationMaxScaleFactor sql.NullInt32  `db:"query_acceleration_max_scale_factor"`
	ResourceMonitor                 sql.NullString `db:"resource_monitor"`
	Actives                         sql.NullInt32  `db:"actives"`
	Pendings                        sql.NullInt32  `db:"pendings"`
	Failed                          sql.NullInt32  `db:"failed"`
	Suspended                       sql.NullInt32  `db:"suspended"`
	Uuid                            sql.NullString `db:"uuid"`
	ScalingPolicy                   sql.NullString `db:"scaling_policy"`
}

func (w *warehouseEntity) toWarehouse() *Warehouse {
	return &Warehouse{
		Name:                            w.Name.String,
		State:                           w.State.String,
		Type:                            w.Type.String,
		WarehouseSize:                   w.WarehouseSize.String,
		MinClusterCount:                 w.MinClusterCount.Int32,
		MaxClusterCount:                 w.MaxClusterCount.Int32,
		StartedClusters:                 w.StartedClusters.Int32,
		Running:                         w.Running.Int32,
		Queued:                          w.Queued.Int32,
		IsDefault:                       w.IsDefault.String,
		IsCurrent:                       w.IsCurrent.String,
		AutoSuspend:                     w.AutoSuspend.Int32,
		AutoResume:                      w.AutoResume.Bool,
		Available:                       w.Available.String,
		Provisioning:                    w.Provisioning.String,
		Quiescing:                       w.Quiescing.String,
		Other:                           w.Other.String,
		CreatedOn:                       w.CreatedOn.Time,
		ResumedOn:                       w.ResumedOn.Time,
		UpdatedOn:                       w.UpdatedOn.Time,
		Owner:                           w.Owner.String,
		Comment:                         w.Comment.String,
		EnableQueryAcceleration:         w.EnableQueryAcceleration.Bool,
		QueryAccelerationMaxScaleFactor: w.QueryAccelerationMaxScaleFactor.Int32,
		ResourceMonitor:                 w.ResourceMonitor.String,
		Actives:                         w.Actives.Int32,
		Pendings:                        w.Pendings.Int32,
		Failed:                          w.Failed.Int32,
		Suspended:                       w.Suspended.Int32,
		Uuid:                            w.Uuid.String,
		ScalingPolicy:                   w.ScalingPolicy.String,
	}
}

type WarehouseProperties struct {
	// Optional: Specifies the warehouse type.
	WarehouseType *string

	// Optional: Specifies the size of the virtual warehouse.
	WarehouseSize *string

	// Optional: Specifies the maximum number of clusters for a multi-cluster warehouse. For a single-cluster warehouse, this value is always 1.
	MaxClusterCount *int32

	// Optional: Specifies the minimum number of clusters for a multi-cluster warehouse (only applies to multi-cluster warehouses).
	MinClusterCount *int32

	// Optional: Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode.
	ScalingPolicy *string

	// Optional: Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.
	AutoSuspend *int32

	// Optional: Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it.
	AutoResume *bool

	// Optional: Specifies a comment for the warehouse.
	Comment *string
}

// WarehouseListOptions represents the options for listing warehouses.
type WarehouseListOptions struct {
	// Required: Filters the command output by object name
	Pattern string

	// Optional: Limits the maximum number of rows returned
	Limit *int
}

func (o WarehouseListOptions) validate() error {
	if o.Pattern == "" {
		return errors.New("pattern must not be empty")
	}
	return nil
}

// WarehouseCreateOptions represents the options for creating a warehouse.
type WarehouseCreateOptions struct {
	*WarehouseProperties

	// Required: Name of the warehouse
	Name string
}

func (o WarehouseCreateOptions) validate() error {
	if o.Name == "" {
		return errors.New("warehouse name must not be empty")
	}
	return nil
}

// WarehouseUpdateOptions represents the options for updating a warehouse.
type WarehouseUpdateOptions struct {
	*WarehouseProperties
}

// List all the warehouses by pattern.
func (w *warehouses) List(ctx context.Context, options WarehouseListOptions) ([]*Warehouse, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate list options: %w", err)
	}

	sql := fmt.Sprintf("SHOW %s LIKE '%s'", ResourceWarehouses, options.Pattern)
	if options.Limit != nil {
		sql = sql + fmt.Sprintf(" LIMIT %d", *options.Limit)
	}
	rows, err := w.client.query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	entities := []*Warehouse{}
	for rows.Next() {
		var entity warehouseEntity
		if err := rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		entities = append(entities, entity.toWarehouse())
	}
	return entities, nil
}

func (*warehouses) formatWarehouseProperties(properties *WarehouseProperties) string {
	var s string
	if properties.WarehouseType != nil {
		s = s + " warehouse_type='" + *properties.WarehouseType + "'"
	}
	if properties.WarehouseSize != nil {
		s = s + " warehouse_size='" + *properties.WarehouseSize + "'"
	}
	if properties.MaxClusterCount != nil {
		s = s + fmt.Sprintf(" max_cluster_count=%d", *properties.MaxClusterCount)
	}
	if properties.MinClusterCount != nil {
		s = s + fmt.Sprintf(" min_cluster_count=%d", *properties.MinClusterCount)
	}
	if properties.ScalingPolicy != nil {
		s = s + " scaling_policy='" + *properties.ScalingPolicy + "'"
	}
	if properties.AutoSuspend != nil {
		s = s + fmt.Sprintf(" auto_suspend=%d", *properties.AutoSuspend)
	}
	if properties.AutoResume != nil {
		s = s + fmt.Sprintf(" auto_resume=%t", *properties.AutoResume)
	}
	if properties.Comment != nil {
		s = s + " comment='" + *properties.Comment + "'"
	}
	return s
}

// Create a new warehouse with the given options.
func (w *warehouses) Create(ctx context.Context, options WarehouseCreateOptions) (*Warehouse, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate create options: %w", err)
	}
	sql := fmt.Sprintf("CREATE %s %s", ResourceWarehouse, options.Name)
	if options.WarehouseProperties != nil {
		sql = sql + w.formatWarehouseProperties(options.WarehouseProperties)
	}
	if _, err := w.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	var entity warehouseEntity
	if err := w.client.read(ctx, ResourceWarehouses, options.Name, &entity); err != nil {
		return nil, err
	}
	return entity.toWarehouse(), nil
}

// Read an warehouse by its name.
func (w *warehouses) Read(ctx context.Context, warehouse string) (*Warehouse, error) {
	var entity warehouseEntity
	if err := w.client.read(ctx, ResourceWarehouses, warehouse, &entity); err != nil {
		return nil, err
	}
	return entity.toWarehouse(), nil
}

// Update attributes of an existing warehouse.
func (w *warehouses) Update(ctx context.Context, warehouse string, options WarehouseUpdateOptions) (*Warehouse, error) {
	if warehouse == "" {
		return nil, errors.New("name must not be empty")
	}
	sql := fmt.Sprintf("ALTER %s %s SET", ResourceWarehouse, warehouse)
	if options.WarehouseProperties != nil {
		sql = sql + w.formatWarehouseProperties(options.WarehouseProperties)
	}
	if _, err := w.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	var entity warehouseEntity
	if err := w.client.read(ctx, ResourceWarehouses, warehouse, &entity); err != nil {
		return nil, err
	}
	return entity.toWarehouse(), nil
}

// Delete a warehouse by its name.
func (w *warehouses) Delete(ctx context.Context, warehouse string) error {
	return w.client.drop(ctx, ResourceWarehouse, warehouse)
}

// Rename a warehouse name.
func (w *warehouses) Rename(ctx context.Context, old string, new string) error {
	return w.client.rename(ctx, ResourceWarehouse, old, new)
}

// Use the active/current warehouse for the session.
func (w *warehouses) Use(ctx context.Context, warehouse string) error {
	return w.client.use(ctx, ResourceWarehouse, warehouse)
}
