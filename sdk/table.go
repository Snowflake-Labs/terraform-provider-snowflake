package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	ResourceTable  = "TABLE"
	ResourceTables = "TABLES"
)

// Roles describes all the tables related methods that the
// Snowflake API supports.
type Tables interface {
	// List all the tables by pattern.
	List(ctx context.Context, options TableListOptions) ([]*Table, error)
	// Create a new table with the given options.
	Create(ctx context.Context, options TableCreateOptions) (*Table, error)
	// Read an table by its name.
	Read(ctx context.Context, table string) (*Table, error)
	// Update attributes of an existing table.
	// Update(ctx context.Context, table string, options TableUpdateOptions) (*Table, error)
	// Delete a table by its name.
	Delete(ctx context.Context, table string) error
}

// tables implements Tables
type tables struct {
	client *Client
}

// Table represents a Snowflake table.
type Table struct {
	CreatedOn                  string
	Name                       string
	DatabaseName               string
	SchemaName                 string
	Kind                       string
	Comment                    string
	ClusterBy                  string
	Rows                       int32
	Bytes                      int32
	Owner                      string
	RetentionTime              int32
	AutomaticClustering        string
	ChangeTracking             string
	SearchOptimization         string
	SearchOptimizationProgress string
	SearchOptimizationBytes    string
	IsExternal                 string
}

type tableEntity struct {
	CreatedOn                  sql.NullString `db:"created_on"`
	Name                       sql.NullString `db:"name"`
	DatabaseName               sql.NullString `db:"database_name"`
	SchemaName                 sql.NullString `db:"schema_name"`
	Kind                       sql.NullString `db:"kind"`
	Comment                    sql.NullString `db:"comment"`
	ClusterBy                  sql.NullString `db:"cluster_by"`
	Rows                       sql.NullInt32  `db:"rows"`
	Bytes                      sql.NullInt32  `db:"bytes"`
	Owner                      sql.NullString `db:"owner"`
	RetentionTime              sql.NullInt32  `db:"retention_time"`
	AutomaticClustering        sql.NullString `db:"automatic_clustering"`
	ChangeTracking             sql.NullString `db:"change_tracking"`
	SearchOptimization         sql.NullString `db:"search_optimization"`
	SearchOptimizationProgress sql.NullString `db:"search_optimization_progress"`
	SearchOptimizationBytes    sql.NullString `db:"search_optimization_bytes"`
	IsExternal                 sql.NullString `db:"is_external"`
}

func (t *tableEntity) toTable() *Table {
	return &Table{
		CreatedOn:                  t.CreatedOn.String,
		Name:                       t.Name.String,
		DatabaseName:               t.DatabaseName.String,
		SchemaName:                 t.SchemaName.String,
		Kind:                       t.Kind.String,
		Comment:                    t.Comment.String,
		ClusterBy:                  t.ClusterBy.String,
		Rows:                       t.Rows.Int32,
		Bytes:                      t.Bytes.Int32,
		Owner:                      t.Owner.String,
		RetentionTime:              t.RetentionTime.Int32,
		AutomaticClustering:        t.AutomaticClustering.String,
		ChangeTracking:             t.ChangeTracking.String,
		SearchOptimization:         t.SearchOptimization.String,
		SearchOptimizationProgress: t.SearchOptimizationProgress.String,
		SearchOptimizationBytes:    t.SearchOptimizationBytes.String,
		IsExternal:                 t.IsExternal.String,
	}
}

// TableListOptions represents the options for listing tables.
type TableListOptions struct {
	// Required: Filters the command output by object name
	Pattern string

	// Optional: Limits the maximum number of rows returned
	Limit *int
}

func (o TableListOptions) validate() error {
	if o.Pattern == "" {
		return errors.New("pattern must not be empty")
	}
	return nil
}

type TableProperties struct {
}

// TableCreateOptions represents the options for creating a table.
type TableCreateOptions struct {
	*TableProperties

	// Required: Identifier for the table; must be unique for your account.
	Name string

	// Required: Specifies the columns of the table
	Columns []string

	// Required: Specifies the database name
	DatabaseName string
}

func (o TableCreateOptions) validate() error {
	if o.Name == "" {
		return errors.New("name must not be empty")
	}
	if o.DatabaseName == "" {
		return errors.New("database name must not be empty")
	}
	if len(o.Columns) == 0 {
		return errors.New("columns must not be empty")
	}
	return nil
}

// TableUpdateOptions represents the options for updating a table.
type TableUpdateOptions struct {
	*TableProperties
}

// List all the tables by pattern.
func (t *tables) List(ctx context.Context, options TableListOptions) ([]*Table, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate list options: %w", err)
	}

	sql := fmt.Sprintf("SHOW %s LIKE '%s'", ResourceTables, options.Pattern)
	rows, err := t.client.query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	entities := []*Table{}
	for rows.Next() {
		var entity tableEntity
		if err := rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		entities = append(entities, entity.toTable())
	}
	return entities, nil
}

// Read a table by its name.
func (t *tables) Read(ctx context.Context, table string) (*Table, error) {
	var entity tableEntity
	if err := t.client.read(ctx, ResourceTables, table, &entity); err != nil {
		return nil, err
	}
	return entity.toTable(), nil
}

func (t *tables) formatTableProperties(properties *TableProperties) string {
	var s string
	return s
}

// Create a new table with the given options.
func (t *tables) Create(ctx context.Context, options TableCreateOptions) (*Table, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate create options: %w", err)
	}
	if err := t.client.use(ctx, ResourceDatabase, options.DatabaseName); err != nil {
		return nil, fmt.Errorf("use database: %w", err)
	}
	sql := fmt.Sprintf("CREATE %s %s", ResourceTable, options.Name)
	sql = sql + "(" + strings.Join(options.Columns, ",") + ")"
	if options.TableProperties != nil {
		sql = sql + " " + t.formatTableProperties(options.TableProperties)
	}
	if _, err := t.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	var entity tableEntity
	if err := t.client.read(ctx, ResourceTables, options.Name, &entity); err != nil {
		return nil, err
	}
	return entity.toTable(), nil
}

// Delete a table by its name.
func (t *tables) Delete(ctx context.Context, table string) error {
	return t.client.drop(ctx, ResourceTable, table)
}
