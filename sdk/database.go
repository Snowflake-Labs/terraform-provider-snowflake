package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	ResourceDatabase  = "DATABASE"
	ResourceDatabases = "DATABASES"
)

// Databases describes all the databases related methods that the
// Snowflake API supports.
type Databases interface {
	// List all the databases by pattern.
	List(ctx context.Context, options DatabaseListOptions) ([]*Database, error)
	// Create a new database with the given options.
	Create(ctx context.Context, options DatabaseCreateOptions) (*Database, error)
	// Read a database by its name.
	Read(ctx context.Context, database string) (*Database, error)
	// Update attributes of an existing database.
	Update(ctx context.Context, database string, options DatabaseUpdateOptions) (*Database, error)
	// Delete a database by its name.
	Delete(ctx context.Context, database string) error
	// Rename a database name.
	Rename(ctx context.Context, old string, new string) error
	// Create a copy of an existing database.
	Clone(ctx context.Context, source string, dest string) error
	// Use the active/current database for the session.
	Use(ctx context.Context, database string) error
}

// databases implements Databases
type databases struct {
	client *Client
}

// Database represents a Snowflake database.
type Database struct {
	Name          string
	IsDefault     string
	IsCurrent     string
	Origin        string
	Owner         string
	Comment       string
	Options       string
	RetentionTime string
	CreatedOn     time.Time
}

type databaseEntity struct {
	Name          sql.NullString `db:"name"`
	IsDefault     sql.NullString `db:"is_default"`
	IsCurrent     sql.NullString `db:"is_current"`
	Origin        sql.NullString `db:"origin"`
	Owner         sql.NullString `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime sql.NullString `db:"retention_time"`
	CreatedOn     sql.NullTime   `db:"created_on"`
}

func (d *databaseEntity) toDatabase() *Database {
	return &Database{
		Name:          d.Name.String,
		IsDefault:     d.IsDefault.String,
		IsCurrent:     d.IsCurrent.String,
		Origin:        d.Origin.String,
		Owner:         d.Owner.String,
		Comment:       d.Comment.String,
		Options:       d.Options.String,
		RetentionTime: d.RetentionTime.String,
		CreatedOn:     d.CreatedOn.Time,
	}
}

// DatabaseListOptions represents the options for listing databases.
type DatabaseListOptions struct {
	// Required: Filters the command output by object name
	Pattern string

	// Optional: Limits the maximum number of rows returned
	Limit *int
}

func (o DatabaseListOptions) validate() error {
	if o.Pattern == "" {
		return errors.New("pattern must not be empty")
	}
	return nil
}

type DatabaseProperties struct {
	// Optional: Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the database
	DataRetentionTimeInDays *int32

	// Optional: Specifies a comment for the database.
	Comment *string
}

type DatabaseCreateOptions struct {
	*DatabaseProperties

	// Required: Specifies the identifier for the database; must be unique for your account.
	Name string
}

func (o DatabaseCreateOptions) validate() error {
	if o.Name == "" {
		return errors.New("database name must not be empty")
	}
	return nil
}

// DatabaseUpdateOptions represents the options for updating a database.
type DatabaseUpdateOptions struct {
	*DatabaseProperties
}

// List all the databases by pattern.
func (d *databases) List(ctx context.Context, options DatabaseListOptions) ([]*Database, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate list options: %w", err)
	}

	sql := fmt.Sprintf("SHOW %s LIKE '%s'", ResourceDatabases, options.Pattern)
	rows, err := d.client.query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	entities := []*Database{}
	for rows.Next() {
		var entity databaseEntity
		if err := rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		entities = append(entities, entity.toDatabase())
	}
	return entities, nil
}

// Create a new database with the given options.
func (d *databases) Create(ctx context.Context, options DatabaseCreateOptions) (*Database, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate create options: %w", err)
	}
	sql := fmt.Sprintf("CREATE %s %s", ResourceDatabase, options.Name)
	if options.DatabaseProperties != nil {
		sql = sql + d.formatDatabaseProperties(options.DatabaseProperties)
	}
	if _, err := d.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	var entity databaseEntity
	if err := d.client.read(ctx, ResourceDatabases, options.Name, &entity); err != nil {
		return nil, err
	}
	return entity.toDatabase(), nil
}

// Read a database by its name.
func (d *databases) Read(ctx context.Context, database string) (*Database, error) {
	var entity databaseEntity
	if err := d.client.read(ctx, ResourceDatabases, database, &entity); err != nil {
		return nil, err
	}
	return entity.toDatabase(), nil
}

func (*databases) formatDatabaseProperties(properties *DatabaseProperties) string {
	var s string
	if properties.Comment != nil {
		s = s + " comment='" + *properties.Comment + "'"
	}
	if properties.DataRetentionTimeInDays != nil {
		s = s + fmt.Sprintf(" data_retention_time_in_days=%d", *properties.DataRetentionTimeInDays)
	}
	return s
}

// Update attributes of an existing database.
func (d *databases) Update(ctx context.Context, database string, options DatabaseUpdateOptions) (*Database, error) {
	if database == "" {
		return nil, errors.New("name must not be empty")
	}
	sql := fmt.Sprintf("ALTER %s %s SET", ResourceDatabase, database)
	if options.DatabaseProperties != nil {
		sql = sql + d.formatDatabaseProperties(options.DatabaseProperties)
	}
	if _, err := d.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	var entity databaseEntity
	if err := d.client.read(ctx, ResourceDatabases, database, &entity); err != nil {
		return nil, err
	}
	return entity.toDatabase(), nil
}

// Delete a database by its name.
func (d *databases) Delete(ctx context.Context, database string) error {
	return d.client.drop(ctx, ResourceDatabase, database)
}

// Rename a database name.
func (d *databases) Rename(ctx context.Context, old string, new string) error {
	return d.client.rename(ctx, ResourceDatabase, old, new)
}

// Create a copy of an existing database.
func (d *databases) Clone(ctx context.Context, source string, dest string) error {
	return d.client.clone(ctx, ResourceDatabase, source, dest)
}

// Use the active/current database for the session.
func (d *databases) Use(ctx context.Context, database string) error {
	return d.client.use(ctx, ResourceDatabase, database)
}
