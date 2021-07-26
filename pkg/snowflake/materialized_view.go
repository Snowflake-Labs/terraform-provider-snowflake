package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	pe "github.com/pkg/errors"
)

// MaterializedViewBuilder abstracts the creation of SQL queries for a Snowflake Materialized View
type MaterializedViewBuilder struct {
	name      string
	db        string
	schema    string
	warehouse string
	secure    bool
	replace   bool
	comment   string
	statement string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (vb *MaterializedViewBuilder) QualifiedName() string {
	var n strings.Builder

	if vb.db != "" && vb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, vb.db, vb.schema))
	}

	if vb.db != "" && vb.schema == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, vb.db))
	}

	if vb.db == "" && vb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, vb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, vb.name))

	return n.String()
}

// WithComment adds a comment to the MaterializedViewBuilder
func (vb *MaterializedViewBuilder) WithComment(c string) *MaterializedViewBuilder {
	vb.comment = c
	return vb
}

// WithDB adds the name of the database to the MaterializedViewBuilder
func (vb *MaterializedViewBuilder) WithDB(db string) *MaterializedViewBuilder {
	vb.db = db
	return vb
}

// WithReplace adds the "OR REPLACE" option to the MaterializedViewBuilder
func (vb *MaterializedViewBuilder) WithReplace() *MaterializedViewBuilder {
	vb.replace = true
	return vb
}

// WithSchema adds the name of the schema to the MaterializedViewBuilder
func (vb *MaterializedViewBuilder) WithSchema(s string) *MaterializedViewBuilder {
	vb.schema = s
	return vb
}

// WithWarehouse adds the name of the warehouse to the MaterializedViewBuilder
func (vb *MaterializedViewBuilder) WithWarehouse(s string) *MaterializedViewBuilder {
	vb.warehouse = s
	return vb
}

// WithSecure sets the secure boolean to true
// [Snowflake Reference](https://docs.snowflake.net/manuals/user-guide/views-secure.html)
func (vb *MaterializedViewBuilder) WithSecure() *MaterializedViewBuilder {
	vb.secure = true
	return vb
}

// WithStatement adds the SQL statement to be used for the view
func (vb *MaterializedViewBuilder) WithStatement(s string) *MaterializedViewBuilder {
	vb.statement = s
	return vb
}

// View returns a pointer to a Builder that abstracts the DDL operations for a view.
//
// Supported DDL operations are:
//   - CREATE MATERIALIZED VIEW
//   - ALTER MATERIALIZED VIEW
//   - DROP MATERIALIZED VIEW
//   - SHOW MATERIALIZED VIEWS
//   - DESCRIBE MATERIALIZED VIEW
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-table.html#materialized-view-management)
func MaterializedView(name string) *MaterializedViewBuilder {
	return &MaterializedViewBuilder{
		name: name,
	}
}

// Create returns the SQL query that will create a new view.
func (vb *MaterializedViewBuilder) Create() []string {
	var q0 strings.Builder

	q0.WriteString(fmt.Sprintf(`USE WAREHOUSE %v;`, vb.warehouse))

	var q1 strings.Builder

	q1.WriteString("CREATE")

	if vb.replace {
		q1.WriteString(" OR REPLACE")
	}

	if vb.secure {
		q1.WriteString(" SECURE")
	}

	q1.WriteString(fmt.Sprintf(` MATERIALIZED VIEW %v`, vb.QualifiedName()))

	if vb.comment != "" {
		q1.WriteString(fmt.Sprintf(" COMMENT = '%v'", EscapeString(vb.comment)))
	}

	q1.WriteString(fmt.Sprintf(" AS %v", vb.statement))

	s := make([]string, 2)
	s[0] = q0.String()
	s[1] = q1.String()
	return s
}

// Rename returns the SQL query that will rename the view.
func (vb *MaterializedViewBuilder) Rename(newName string) string {
	oldName := vb.QualifiedName()
	vb.name = newName
	return fmt.Sprintf(`ALTER MATERIALIZED VIEW %v RENAME TO %v`, oldName, vb.QualifiedName())
}

// Secure returns the SQL query that will change the view to a secure view.
func (vb *MaterializedViewBuilder) Secure() string {
	return fmt.Sprintf(`ALTER MATERIALIZED VIEW %v SET SECURE`, vb.QualifiedName())
}

// Unsecure returns the SQL query that will change the view to a normal (unsecured) view.
func (vb *MaterializedViewBuilder) Unsecure() string {
	return fmt.Sprintf(`ALTER MATERIALIZED VIEW %v UNSET SECURE`, vb.QualifiedName())
}

// ChangeComment returns the SQL query that will update the comment on the view.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (vb *MaterializedViewBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER MATERIALIZED VIEW %v SET COMMENT = '%v'`, vb.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the view.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (vb *MaterializedViewBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER MATERIALIZED VIEW %v UNSET COMMENT`, vb.QualifiedName())
}

// Show returns the SQL query that will show the row representing this view.
func (vb *MaterializedViewBuilder) Show() string {
	if vb.db == "" {
		return fmt.Sprintf(`SHOW MATERIALIZED VIEWS LIKE '%v'`, vb.name)
	}
	return fmt.Sprintf(`SHOW MATERIALIZED VIEWS LIKE '%v' IN DATABASE "%v"`, vb.name, vb.db)
}

// Drop returns the SQL query that will drop the row representing this view.
func (vb *MaterializedViewBuilder) Drop() string {
	return fmt.Sprintf(`DROP MATERIALIZED VIEW %v`, vb.QualifiedName())
}

type materializedView struct {
	Comment       sql.NullString `db:"comment"`
	IsSecure      bool           `db:"is_secure"`
	Name          sql.NullString `db:"name"`
	SchemaName    sql.NullString `db:"schema_name"`
	Text          sql.NullString `db:"text"`
	DatabaseName  sql.NullString `db:"database_name"`
	WarehouseName sql.NullString `db:"warehouse_name"`
}

func ScanMaterializedView(row *sqlx.Row) (*materializedView, error) {
	r := &materializedView{}
	err := row.StructScan(r)
	return r, err
}

func ListMaterializedViews(databaseName string, schemaName string, db *sql.DB) ([]materializedView, error) {
	stmt := fmt.Sprintf(`SHOW MATERIALIZED VIEWS IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []materializedView{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no materialized views found")
		return nil, nil
	}
	return dbs, pe.Wrapf(err, "unable to scan row for %s", stmt)
}
