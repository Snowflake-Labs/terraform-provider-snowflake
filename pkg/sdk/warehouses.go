package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Warehouses interface {
	// Create creates a warehouse.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateWarehouseOptions) error
	// Alter modifies an existing warehouse
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterWarehouseOptions) error
	// Drop removes a warehouse.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropWarehouseOptions) error
	// Show returns a list of warehouses.
	Show(ctx context.Context, opts *ShowWarehouseOptions) ([]*Warehouse, error)
	// ShowByID returns a warehouse by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error)
	// Describe returns the details of a warehouse.
	Describe(ctx context.Context, id AccountObjectIdentifier) (*WarehouseDetails, error)
}

var _ Warehouses = (*warehouses)(nil)

type warehouses struct {
	client *Client
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

func ToWarehouseSize(s string) (WarehouseSize, error) {
	s = strings.ToUpper(s)
	switch s {
	case "XSMALL", "X-SMALL":
		return WarehouseSizeXSmall, nil
	case "SMALL":
		return WarehouseSizeSmall, nil
	case "MEDIUM":
		return WarehouseSizeMedium, nil
	case "LARGE":
		return WarehouseSizeLarge, nil
	case "XLARGE", "X-LARGE":
		return WarehouseSizeXLarge, nil
	case "XXLARGE", "X2LARGE", "2X-LARGE":
		return WarehouseSizeXXLarge, nil
	case "XXXLARGE", "X3LARGE", "3X-LARGE":
		return WarehouseSizeXXXLarge, nil
	case "X4LARGE", "4X-LARGE":
		return WarehouseSizeX4Large, nil
	case "X5LARGE", "5X-LARGE":
		return WarehouseSizeX5Large, nil
	case "X6LARGE", "6X-LARGE":
		return WarehouseSizeX6Large, nil
	default:
		return "", fmt.Errorf("invalid warehouse size: %s", s)
	}
}

type ScalingPolicy string

var (
	ScalingPolicyStandard ScalingPolicy = "STANDARD"
	ScalingPolicyEconomy  ScalingPolicy = "ECONOMY"
)

type CreateWarehouseOptions struct {
	create      bool                    `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	warehouse   bool                    `ddl:"static" sql:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"`

	// Object properties
	WarehouseType                   *WarehouseType `ddl:"parameter,single_quotes" sql:"WAREHOUSE_TYPE"`
	WarehouseSize                   *WarehouseSize `ddl:"parameter,single_quotes" sql:"WAREHOUSE_SIZE"`
	MaxClusterCount                 *int           `ddl:"parameter" sql:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *int           `ddl:"parameter" sql:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *ScalingPolicy `ddl:"parameter,single_quotes" sql:"SCALING_POLICY"`
	AutoSuspend                     *int           `ddl:"parameter" sql:"AUTO_SUSPEND"`
	AutoResume                      *bool          `ddl:"parameter" sql:"AUTO_RESUME"`
	InitiallySuspended              *bool          `ddl:"parameter" sql:"INITIALLY_SUSPENDED"`
	ResourceMonitor                 *string        `ddl:"parameter,double_quotes" sql:"RESOURCE_MONITOR"`
	Comment                         *string        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	EnableQueryAcceleration         *bool          `ddl:"parameter" sql:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *int           `ddl:"parameter" sql:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *int             `ddl:"parameter" sql:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *int             `ddl:"parameter" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *int             `ddl:"parameter" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	Tag                             []TagAssociation `ddl:"keyword,parentheses" sql:"TAG"`
}

func (opts *CreateWarehouseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if valueSet(opts.MaxClusterCount) && !validateIntInRange(*opts.MaxClusterCount, 1, 10) {
		return fmt.Errorf("MaxClusterCount must be between 1 and 10")
	}
	if valueSet(opts.MinClusterCount) && !validateIntInRange(*opts.MinClusterCount, 1, 10) {
		return fmt.Errorf("MinClusterCount must be between 1 and 10")
	}
	if valueSet(opts.MinClusterCount) && valueSet(opts.MaxClusterCount) && !validateIntGreaterThanOrEqual(*opts.MaxClusterCount, *opts.MinClusterCount) {
		return fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount")
	}
	if valueSet(opts.QueryAccelerationMaxScaleFactor) && !validateIntInRange(*opts.QueryAccelerationMaxScaleFactor, 0, 100) {
		return fmt.Errorf("QueryAccelerationMaxScaleFactor must be between 0 and 100")
	}
	return nil
}

func (c *warehouses) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateWarehouseOptions) error {
	if opts == nil {
		opts = &CreateWarehouseOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	stmt, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, stmt)
	return err
}

type AlterWarehouseOptions struct {
	alter     bool                    `ddl:"static" sql:"ALTER"`     //lint:ignore U1000 This is used in the ddl tag
	warehouse bool                    `ddl:"static" sql:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists  *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name      AccountObjectIdentifier `ddl:"identifier"`

	Suspend         *bool                   `ddl:"keyword" sql:"SUSPEND"`
	Resume          *bool                   `ddl:"keyword" sql:"RESUME"`
	IfSuspended     *bool                   `ddl:"keyword" sql:"IF SUSPENDED"`
	AbortAllQueries *bool                   `ddl:"keyword" sql:"ABORT ALL QUERIES"`
	NewName         AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`

	Set   *WarehouseSet   `ddl:"keyword" sql:"SET"`
	Unset *WarehouseUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

func (opts *AlterWarehouseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if ok := exactlyOneValueSet(
		opts.Suspend,
		opts.Resume,
		opts.AbortAllQueries,
		opts.NewName,
		opts.Set,
		opts.Unset); !ok {
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
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			return err
		}
	}
	return nil
}

type WarehouseSet struct {
	// Object properties
	WarehouseType                   *WarehouseType          `ddl:"parameter,single_quotes" sql:"WAREHOUSE_TYPE"`
	WarehouseSize                   *WarehouseSize          `ddl:"parameter,single_quotes" sql:"WAREHOUSE_SIZE"`
	WaitForCompletion               *bool                   `ddl:"parameter" sql:"WAIT_FOR_COMPLETION"`
	MaxClusterCount                 *int                    `ddl:"parameter" sql:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *int                    `ddl:"parameter" sql:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *ScalingPolicy          `ddl:"parameter,single_quotes" sql:"SCALING_POLICY"`
	AutoSuspend                     *int                    `ddl:"parameter" sql:"AUTO_SUSPEND"`
	AutoResume                      *bool                   `ddl:"parameter" sql:"AUTO_RESUME"`
	ResourceMonitor                 AccountObjectIdentifier `ddl:"identifier,equals" sql:"RESOURCE_MONITOR"`
	Comment                         *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	EnableQueryAcceleration         *bool                   `ddl:"parameter" sql:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *int                    `ddl:"parameter" sql:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *int `ddl:"parameter" sql:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *int `ddl:"parameter" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *int `ddl:"parameter" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`

	Tag []TagAssociation `ddl:"keyword" sql:"TAG"`
}

func (v *WarehouseSet) validate() error {
	if v.MaxClusterCount != nil {
		if ok := validateIntInRange(*v.MaxClusterCount, 1, 10); !ok {
			return fmt.Errorf("MaxClusterCount must be between 1 and 10")
		}
	}
	if v.MinClusterCount != nil {
		if ok := validateIntInRange(*v.MinClusterCount, 1, 10); !ok {
			return fmt.Errorf("MinClusterCount must be between 1 and 10")
		}
	}
	if v.AutoSuspend != nil {
		if ok := validateIntGreaterThanOrEqual(*v.AutoSuspend, 0); !ok {
			return fmt.Errorf("AutoSuspend must be greater than or equal to 0")
		}
	}
	if v.QueryAccelerationMaxScaleFactor != nil {
		if ok := validateIntInRange(*v.QueryAccelerationMaxScaleFactor, 0, 100); !ok {
			return fmt.Errorf("QueryAccelerationMaxScaleFactor must be between 0 and 100")
		}
	}
	if valueSet(v.Tag) && !everyValueNil(v.AutoResume, v.EnableQueryAcceleration, v.MaxClusterCount, v.MinClusterCount, v.AutoSuspend, v.QueryAccelerationMaxScaleFactor) {
		return fmt.Errorf("Tag cannot be set with any other Set parameter")
	}
	return nil
}

type WarehouseUnset struct {
	// Object properties
	WarehouseType                   *bool `ddl:"keyword" sql:"WAREHOUSE_TYPE"`
	WarehouseSize                   *bool `ddl:"keyword" sql:"WAREHOUSE_SIZE"`
	WaitForCompletion               *bool `ddl:"keyword" sql:"WAIT_FOR_COMPLETION"`
	MaxClusterCount                 *bool `ddl:"keyword" sql:"MAX_CLUSTER_COUNT"`
	MinClusterCount                 *bool `ddl:"keyword" sql:"MIN_CLUSTER_COUNT"`
	ScalingPolicy                   *bool `ddl:"keyword" sql:"SCALING_POLICY"`
	AutoSuspend                     *bool `ddl:"keyword" sql:"AUTO_SUSPEND"`
	AutoResume                      *bool `ddl:"keyword" sql:"AUTO_RESUME"`
	ResourceMonitor                 *bool `ddl:"keyword" sql:"RESOURCE_MONITOR"`
	Comment                         *bool `ddl:"keyword" sql:"COMMENT"`
	EnableQueryAcceleration         *bool `ddl:"keyword" sql:"ENABLE_QUERY_ACCELERATION"`
	QueryAccelerationMaxScaleFactor *bool `ddl:"keyword" sql:"QUERY_ACCELERATION_MAX_SCALE_FACTOR"`

	// Object params
	MaxConcurrencyLevel             *bool              `ddl:"keyword" sql:"MAX_CONCURRENCY_LEVEL"`
	StatementQueuedTimeoutInSeconds *bool              `ddl:"keyword" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *bool              `ddl:"keyword" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	Tag                             []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

func (v *WarehouseUnset) validate() error {
	if valueSet(v.Tag) && !everyValueNil(v.AutoResume, v.EnableQueryAcceleration, v.MaxClusterCount, v.MinClusterCount, v.AutoSuspend, v.QueryAccelerationMaxScaleFactor) {
		return fmt.Errorf("Tag cannot be unset with any other Unset parameter")
	}
	return nil
}

func (c *warehouses) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterWarehouseOptions) error {
	if opts == nil {
		opts = &AlterWarehouseOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, sql)
	return err
}

type DropWarehouseOptions struct {
	drop      bool                    `ddl:"static" sql:"DROP"`      //lint:ignore U1000 This is used in the ddl tag
	warehouse bool                    `ddl:"static" sql:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists  *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name      AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *DropWarehouseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (c *warehouses) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropWarehouseOptions) error {
	if opts == nil {
		opts = &DropWarehouseOptions{
			name: id,
		}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return err
}

type ShowWarehouseOptions struct {
	show       bool  `ddl:"static" sql:"SHOW"`       //lint:ignore U1000 This is used in the ddl tag
	warehouses bool  `ddl:"static" sql:"WAREHOUSES"` //lint:ignore U1000 This is used in the ddl tag
	Like       *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowWarehouseOptions) validate() error {
	return nil
}

type WarehouseState string

const (
	WarehouseStateSuspended  WarehouseState = "SUSPENDED"
	WarehouseStateSuspending WarehouseState = "SUSPENDING"
	WarehouseStateStarted    WarehouseState = "STARTED"
	WarehouseStateResizing   WarehouseState = "RESIZING"
	WarehouseStateResuming   WarehouseState = "RESUMING"
)

type Warehouse struct {
	Name                            string
	State                           WarehouseState
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
		State:                           WarehouseState(row.State),
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

func (c *warehouses) Show(ctx context.Context, opts *ShowWarehouseOptions) ([]*Warehouse, error) {
	if opts == nil {
		opts = &ShowWarehouseOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []warehouseDBRow{}
	err = c.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]*Warehouse, len(dest))
	for i, row := range dest {
		resultList[i] = row.toWarehouse()
	}

	return resultList, nil
}

func (c *warehouses) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error) {
	warehouses, err := c.Show(ctx, &ShowWarehouseOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, warehouse := range warehouses {
		if warehouse.ID().name == id.Name() {
			return warehouse, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type warehouseDescribeOptions struct {
	describe  bool                    `ddl:"static" sql:"DESCRIBE"`  //lint:ignore U1000 This is used in the ddl tag
	warehouse bool                    `ddl:"static" sql:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	name      AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *warehouseDescribeOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
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

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := warehouseDetailsRow{}
	err = c.client.queryOne(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return dest.toWarehouseDetails(), nil
}

func (v *Warehouse) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Warehouse) ObjectType() ObjectType {
	return ObjectTypeWarehouse
}
