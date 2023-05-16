package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Warehouses interface {
	// Create creates a warehouse.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseCreateOptions) error
	// Alter modifies an existing warehouse
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseAlterOptions) error
	// Drop removes a warehouse.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseDropOptions) error
	// Show returns a list of warehouses.
	Show(ctx context.Context, opts *WarehouseShowOptions) ([]*Warehouse, error)
	// ShowByID returns a warehouse by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error)
	// Describe returns the details of a warehouse.
	Describe(ctx context.Context, id AccountObjectIdentifier) (*WarehouseDetails, error)
}

var _ Warehouses = (*warehouses)(nil)

type warehouses struct {
	client  *Client
	builder *sqlBuilder
}

type WarehouseType string

var (
	WarehouseTypeStandard          WarehouseType = "STANDARD"
	WarehouseTypeSnowparkOptimized WarehouseType = "SNOWPARK-OPTIMIZED"
)

type WarehouseSize string

var (
	WarehouseSizeXSmall   WarehouseSize = "XSMALL"
	WarehouseSizeSmall    WarehouseSize = "SMALL"
	WarehouseSizeMedium   WarehouseSize = "MEDIUM"
	WarehouseSizeLarge    WarehouseSize = "LARGE"
	WarehouseSizeXLarge   WarehouseSize = "XLARGE"
	WarehouseSizeXXLarge  WarehouseSize = "XXLARGE"
	WarehouseSizeXXXLarge WarehouseSize = "XXXLARGE"
	WarehouseSizeX4Large  WarehouseSize = "X4LARGE"
	WarehouseSizeX5Large  WarehouseSize = "X5LARGE"
	WarehouseSizeX6Large  WarehouseSize = "X6LARGE"
)

type ScalingPolicy string

var (
	ScalingPolicyStandard ScalingPolicy = "STANDARD"
	ScalingPolicyEconomy  ScalingPolicy = "ECONOMY"
)

type WarehouseCreateOptions struct {
	create      bool                    `ddl:"static" db:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                   `ddl:"keyword" db:"OR REPLACE"`
	warehouse   bool                    `ddl:"static" db:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                   `ddl:"keyword" db:"IF NOT EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"`

	// Object properties
	WarehouseType                   *WarehouseType `ddl:"parameter,single_quotes" db:"WAREHOUSE_TYPE"`
	WarehouseSize                   *WarehouseSize `ddl:"parameter,single_quotes" db:"WAREHOUSE_SIZE"`
	MaxClusterCount                 *int           `ddl:"parameter" db:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *int           `ddl:"parameter" db:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *ScalingPolicy `ddl:"parameter,single_quotes" db:"SCALING_POLICY"`
	AutoSuspend                     *int           `ddl:"parameter" db:"AUTO_SUSPEND"`
	AutoResume                      *bool          `ddl:"parameter" db:"AUTO_RESUME"`
	InitiallySuspended              *bool          `ddl:"parameter" db:"INITIALLY_SUSPENDED"`
	ResourceMonitor                 *string        `ddl:"parameter,double_quotes" db:"RESOURCE_MONITOR"`
	Comment                         *string        `ddl:"parameter,single_quotes" db:"COMMENT"`
	EnableQueryAcceleration         *bool          `ddl:"parameter" db:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *int           `ddl:"parameter" db:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *int             `ddl:"parameter" db:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *int             `ddl:"parameter" db:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *int             `ddl:"parameter" db:"STATEMENT_TIMEOUT_IN_SECONDS"`
	Tag                             []TagAssociation `ddl:"keyword,parentheses" db:"TAG"`
}

func (opts *WarehouseCreateOptions) validate() error {
	if opts.MaxClusterCount != nil && ((*opts.MaxClusterCount < 1) || (10 < *opts.MaxClusterCount)) {
		return fmt.Errorf("MaxClusterCount must be between 1 and 10")
	}
	if opts.MinClusterCount != nil && ((*opts.MinClusterCount < 1) || (10 < *opts.MinClusterCount)) {
		return fmt.Errorf("MinClusterCount must be between 1 and 10")
	}
	if opts.MinClusterCount != nil && opts.MaxClusterCount != nil && *opts.MaxClusterCount < *opts.MinClusterCount {
		return fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount")
	}
	if opts.QueryAccelerationMaxScaleFactor != nil && 100 < *opts.QueryAccelerationMaxScaleFactor {
		return fmt.Errorf("QueryAccelerationMaxScaleFactor must be less than or equal to 100")
	}
	return nil
}

func (c *warehouses) Create(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseCreateOptions) error {
	if opts == nil {
		opts = &WarehouseCreateOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}

	clauses, err := c.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := c.builder.sql(clauses...)
	_, err = c.client.exec(ctx, stmt)
	return err
}

type WarehouseAlterOptions struct {
	alter     bool                    `ddl:"static" db:"ALTER"`     //lint:ignore U1000 This is used in the ddl tag
	warehouse bool                    `ddl:"static" db:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists  *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name      AccountObjectIdentifier `ddl:"identifier"`

	Suspend         *bool                    `ddl:"keyword" db:"SUSPEND"`
	Resume          *bool                    `ddl:"keyword" db:"RESUME"`
	IfSuspended     *bool                    `ddl:"keyword" db:"IF SUSPENDED"`
	AbortAllQueries *bool                    `ddl:"keyword" db:"ABORT ALL QUERIES"`
	NewName         *AccountObjectIdentifier `ddl:"identifier" db:"RENAME TO"`

	Set   *WarehouseSet   `ddl:"keyword" db:"SET"`
	Unset *WarehouseUnset `ddl:"list,no_parentheses" db:"UNSET"`
}

type WarehouseSet struct {
	// Object properties
	WarehouseType                   *WarehouseType          `ddl:"parameter,single_quotes" db:"WAREHOUSE_TYPE"`
	WarehouseSize                   *WarehouseSize          `ddl:"parameter,single_quotes" db:"WAREHOUSE_SIZE"`
	WaitForCompletion               *bool                   `ddl:"parameter" db:"WAIT_FOR_COMPLETION"`
	MaxClusterCount                 *int                    `ddl:"parameter" db:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *int                    `ddl:"parameter" db:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *ScalingPolicy          `ddl:"parameter,single_quotes" db:"SCALING_POLICY"`
	AutoSuspend                     *int                    `ddl:"parameter" db:"AUTO_SUSPEND"`
	AutoResume                      *bool                   `ddl:"parameter" db:"AUTO_RESUME"`
	ResourceMonitor                 AccountObjectIdentifier `ddl:"identifier,equals" db:"RESOURCE_MONITOR"`
	Comment                         *string                 `ddl:"parameter,single_quotes" db:"COMMENT"`
	EnableQueryAcceleration         *bool                   `ddl:"parameter" db:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *int                    `ddl:"parameter" db:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *int `ddl:"parameter" db:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *int `ddl:"parameter" db:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *int `ddl:"parameter" db:"STATEMENT_TIMEOUT_IN_SECONDS"`

	Tag []TagAssociation `ddl:"keyword" db:"TAG"`
}

type WarehouseUnset struct {
	// Object properties
	WarehouseType                   *bool `ddl:"keyword" db:"WAREHOUSE_TYPE"`
	WarehouseSize                   *bool `ddl:"keyword" db:"WAREHOUSE_SIZE"`
	WaitForCompletion               *bool `ddl:"keyword" db:"WAIT_FOR_COMPLETION"`
	MaxClusterCount                 *bool `ddl:"keyword" db:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *bool `ddl:"keyword" db:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *bool `ddl:"keyword" db:"SCALING_POLICY"`
	AutoSuspend                     *bool `ddl:"keyword" db:"AUTO_SUSPEND"`
	AutoResume                      *bool `ddl:"keyword" db:"AUTO_RESUME"`
	ResourceMonitor                 *bool `ddl:"keyword" db:"RESOURCE_MONITOR"`
	Comment                         *bool `ddl:"keyword" db:"COMMENT"`
	EnableQueryAcceleration         *bool `ddl:"keyword" db:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *bool `ddl:"keyword" db:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *bool              `ddl:"keyword" db:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *bool              `ddl:"keyword" db:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *bool              `ddl:"keyword" db:"STATEMENT_TIMEOUT_IN_SECONDS"`
	Tag                             []ObjectIdentifier `ddl:"keyword" db:"TAG"`
}

func (opts *WarehouseAlterOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name must not be empty")
	}

	if err := exactlyOneValueSet(
		opts.Suspend,
		opts.Resume,
		opts.AbortAllQueries,
		opts.NewName,
		opts.Set,
		opts.Unset); err != nil {
		return fmt.Errorf("exactly one of Suspend, Resume, AbortAllQueries, NewName, Set, Unset must be set")
	}
	if everyValueSet(opts.Suspend, opts.Resume) && (*opts.Suspend && *opts.Resume) {
		return fmt.Errorf("Suspend and Resume cannot both be true")
	}
	if (valueSet(opts.IfSuspended) && *opts.IfSuspended) && (!valueSet(opts.Resume) || !*opts.Resume) {
		return fmt.Errorf(`"Resume" has to be set when using "IfSuspended"`)
	}
	if everyValueSet(opts.Set, opts.Unset) {
		return fmt.Errorf("Set and Unset cannot both be set")
	}
	if opts.Set != nil {
		if opts.Set.MaxClusterCount != nil {
			if ok := validateIntInRange(*opts.Set.MaxClusterCount, 1, 10); !ok {
				return fmt.Errorf("MaxClusterCount must be between 1 and 10")
			}
		}
		if opts.Set.MinClusterCount != nil {
			if ok := validateIntInRange(*opts.Set.MinClusterCount, 1, 10); !ok {
				return fmt.Errorf("MinClusterCount must be between 1 and 10")
			}
		}
		if opts.Set.AutoSuspend != nil {
			if ok := validateIntGreaterThanOrEqual(*opts.Set.AutoSuspend, 0); !ok {
				return fmt.Errorf("AutoSuspend must be greater than or equal to 0")
			}
		}
		if opts.Set.QueryAccelerationMaxScaleFactor != nil {
			if ok := validateIntInRange(*opts.Set.QueryAccelerationMaxScaleFactor, 0, 100); !ok {
				return fmt.Errorf("QueryAccelerationMaxScaleFactor must be between 0 and 100")
			}
		}
		if valueSet(opts.Set.Tag) && !everyValueNil(opts.Set.AutoResume, opts.Set.EnableQueryAcceleration, opts.Set.MaxClusterCount, opts.Set.MinClusterCount, opts.Set.AutoSuspend, opts.Set.QueryAccelerationMaxScaleFactor) {
			return fmt.Errorf("Tag cannot be set with any other Set parameter")
		}
	}
	return nil
}

func (c *warehouses) Alter(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseAlterOptions) error {
	if opts == nil {
		opts = &WarehouseAlterOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := c.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := c.builder.sql(clauses...)
	_, err = c.client.exec(ctx, stmt)
	return err
}

type WarehouseDropOptions struct {
	drop      bool                    `ddl:"static" db:"DROP"`      //lint:ignore U1000 This is used in the ddl tag
	warehouse bool                    `ddl:"static" db:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists  *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name      AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *WarehouseDropOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (c *warehouses) Drop(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseDropOptions) error {
	if opts == nil {
		opts = &WarehouseDropOptions{
			name: id,
		}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := c.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := c.builder.sql(clauses...)
	_, err = c.client.exec(ctx, stmt)
	if err != nil {
		return decodeDriverError(err)
	}
	return err
}

type WarehouseShowOptions struct {
	show       bool  `ddl:"static" db:"SHOW"`       //lint:ignore U1000 This is used in the ddl tag
	warehouses bool  `ddl:"static" db:"WAREHOUSES"` //lint:ignore U1000 This is used in the ddl tag
	Like       *Like `ddl:"keyword" db:"LIKE"`
}

func (opts *WarehouseShowOptions) validate() error {
	return nil
}

type Warehouse struct {
	Name                            string
	State                           string
	Type                            WarehouseType
	Size                            WarehouseSize
	MinClusterCount                 int
	MaxClusterCount                 int
	StartedClusters                 int
	Running                         int
	Queued                          int
	IsDefault                       bool
	IsCurrent                       bool
	AutoSuspend                     int
	AutoResume                      bool
	Available                       float64
	Provisioning                    float64
	Quiescing                       float64
	Other                           float64
	CreatedOn                       time.Time
	ResumedOn                       time.Time
	UpdatedOn                       time.Time
	Owner                           string
	Comment                         string
	EnableQueryAcceleration         bool
	QueryAccelerationMaxScaleFactor int
	ResourceMonitor                 string
	Actives                         string
	Pendings                        string
	Failed                          string
	Suspended                       string
	UUID                            string
	ScalingPolicy                   ScalingPolicy
}

type warehouseDBRow struct {
	Name                            string        `db:"name"`
	State                           string        `db:"state"`
	Type                            string        `db:"type"`
	Size                            string        `db:"size"`
	MinClusterCount                 int           `db:"min_cluster_count"`
	MaxClusterCount                 int           `db:"max_cluster_count"`
	StartedClusters                 int           `db:"started_clusters"`
	Running                         int           `db:"running"`
	Queued                          int           `db:"queued"`
	IsDefault                       string        `db:"is_default"`
	IsCurrent                       string        `db:"is_current"`
	AutoSuspend                     sql.NullInt64 `db:"auto_suspend"`
	AutoResume                      bool          `db:"auto_resume"`
	Available                       string        `db:"available"`
	Provisioning                    string        `db:"provisioning"`
	Quiescing                       string        `db:"quiescing"`
	Other                           string        `db:"other"`
	CreatedOn                       time.Time     `db:"created_on"`
	ResumedOn                       time.Time     `db:"resumed_on"`
	UpdatedOn                       time.Time     `db:"updated_on"`
	Owner                           string        `db:"owner"`
	Comment                         string        `db:"comment"`
	EnableQueryAcceleration         bool          `db:"enable_query_acceleration"`
	QueryAccelerationMaxScaleFactor int           `db:"query_acceleration_max_scale_factor"`
	ResourceMonitor                 string        `db:"resource_monitor"`
	Actives                         string        `db:"actives"`
	Pendings                        string        `db:"pendings"`
	Failed                          string        `db:"failed"`
	Suspended                       string        `db:"suspended"`
	UUID                            string        `db:"uuid"`
	ScalingPolicy                   string        `db:"scaling_policy"`
}

func (row warehouseDBRow) toWarehouse() *Warehouse {
	wh := &Warehouse{
		Name:                            row.Name,
		State:                           row.State,
		Type:                            WarehouseType(row.Type),
		Size:                            WarehouseSize(strings.ReplaceAll(strings.ToUpper(row.Size), "-", "")),
		MinClusterCount:                 row.MinClusterCount,
		MaxClusterCount:                 row.MaxClusterCount,
		StartedClusters:                 row.StartedClusters,
		Running:                         row.Running,
		Queued:                          row.Queued,
		IsDefault:                       row.IsDefault == "Y",
		IsCurrent:                       row.IsCurrent == "Y",
		AutoResume:                      row.AutoResume,
		CreatedOn:                       row.CreatedOn,
		ResumedOn:                       row.ResumedOn,
		UpdatedOn:                       row.UpdatedOn,
		Owner:                           row.Owner,
		Comment:                         row.Comment,
		EnableQueryAcceleration:         row.EnableQueryAcceleration,
		QueryAccelerationMaxScaleFactor: row.QueryAccelerationMaxScaleFactor,
		ResourceMonitor:                 row.ResourceMonitor,
		Actives:                         row.Actives,
		Pendings:                        row.Pendings,
		Failed:                          row.Failed,
		Suspended:                       row.Suspended,
		UUID:                            row.UUID,
		ScalingPolicy:                   ScalingPolicy(row.ScalingPolicy),
	}
	if val, err := strconv.ParseFloat(row.Available, 64); err != nil {
		wh.Available = val
	}
	if val, err := strconv.ParseFloat(row.Provisioning, 64); err != nil {
		wh.Provisioning = val
	}
	if val, err := strconv.ParseFloat(row.Quiescing, 64); err != nil {
		wh.Quiescing = val
	}
	if val, err := strconv.ParseFloat(row.Other, 64); err != nil {
		wh.Other = val
	}
	if row.AutoSuspend.Valid {
		wh.AutoSuspend = int(row.AutoSuspend.Int64)
	}
	return wh
}

func (c *warehouses) Show(ctx context.Context, opts *WarehouseShowOptions) ([]*Warehouse, error) {
	if opts == nil {
		opts = &WarehouseShowOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	clauses, err := c.builder.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	stmt := c.builder.sql(clauses...)
	dest := []warehouseDBRow{}

	err = c.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}
	resultList := make([]*Warehouse, len(dest))
	for i, row := range dest {
		resultList[i] = row.toWarehouse()
	}

	return resultList, nil
}

func (c *warehouses) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error) {
	results, err := c.Show(ctx, &WarehouseShowOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, res := range results {
		if res.ID().name == id.Name() {
			return res, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type warehouseDescribeOptions struct {
	describe  bool                    `ddl:"static" db:"DESCRIBE"`  //lint:ignore U1000 This is used in the ddl tag
	warehouse bool                    `ddl:"static" db:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	name      AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *warehouseDescribeOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type warehouseDetailsRow struct {
	CreatedOn time.Time `db:"created_on"`
	Name      string    `db:"name"`
	Kind      string    `db:"kind"`
}

func (row *warehouseDetailsRow) toWarehouseDetails() *WarehouseDetails {
	return &WarehouseDetails{
		CreatedOn: row.CreatedOn,
		Name:      row.Name,
		Kind:      row.Kind,
	}
}

type WarehouseDetails struct {
	CreatedOn time.Time
	Name      string
	Kind      string
}

func (c *warehouses) Describe(ctx context.Context, id AccountObjectIdentifier) (*WarehouseDetails, error) {
	opts := &warehouseDescribeOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	clauses, err := c.builder.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	stmt := c.builder.sql(clauses...)
	dest := warehouseDetailsRow{}
	err = c.client.queryOne(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}

	return dest.toWarehouseDetails(), nil
}

func (v *Warehouse) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}
