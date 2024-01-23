package sdk

import (
	"context"
	"database/sql"
	"time"
)

type MaterializedViews interface {
	Create(ctx context.Context, request *CreateMaterializedViewRequest) error
	Alter(ctx context.Context, request *AlterMaterializedViewRequest) error
	Drop(ctx context.Context, request *DropMaterializedViewRequest) error
	Show(ctx context.Context, request *ShowMaterializedViewRequest) ([]MaterializedView, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*MaterializedView, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]MaterializedViewDetails, error)
}

// CreateMaterializedViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-materialized-view.
type CreateMaterializedViewOptions struct {
	create                 bool                                  `ddl:"static" sql:"CREATE"`
	OrReplace              *bool                                 `ddl:"keyword" sql:"OR REPLACE"`
	Secure                 *bool                                 `ddl:"keyword" sql:"SECURE"`
	materializedView       bool                                  `ddl:"static" sql:"MATERIALIZED VIEW"`
	IfNotExists            *bool                                 `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                   SchemaObjectIdentifier                `ddl:"identifier"`
	CopyGrants             *bool                                 `ddl:"keyword" sql:"COPY GRANTS"`
	Columns                []MaterializedViewColumn              `ddl:"list,parentheses"`
	ColumnsMaskingPolicies []MaterializedViewColumnMaskingPolicy `ddl:"list,no_parentheses,no_equals"`
	Comment                *string                               `ddl:"parameter,single_quotes" sql:"COMMENT"`
	RowAccessPolicy        *MaterializedViewRowAccessPolicy      `ddl:"keyword"`
	Tag                    []TagAssociation                      `ddl:"keyword,parentheses" sql:"TAG"`
	ClusterBy              *MaterializedViewClusterBy            `ddl:"keyword"`
	as                     bool                                  `ddl:"static" sql:"AS"`
	sql                    string                                `ddl:"keyword,no_quotes"`
}

type MaterializedViewColumn struct {
	Name    string  `ddl:"keyword,double_quotes"`
	Comment *string `ddl:"parameter,single_quotes,no_equals" sql:"COMMENT"`
}

type MaterializedViewColumnMaskingPolicy struct {
	Name          string                 `ddl:"keyword"`
	MaskingPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"MASKING POLICY"`
	Using         []string               `ddl:"keyword,parentheses" sql:"USING"`
	Tag           []TagAssociation       `ddl:"keyword,parentheses" sql:"TAG"`
}

type MaterializedViewRowAccessPolicy struct {
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
	On              []string               `ddl:"keyword,parentheses" sql:"ON"`
}

type MaterializedViewClusterBy struct {
	clusterBy   bool                                  `ddl:"static" sql:"CLUSTER BY"`
	Expressions []MaterializedViewClusterByExpression `ddl:"list,parentheses"`
}

type MaterializedViewClusterByExpression struct {
	Name string `ddl:"keyword,double_quotes"`
}

// AlterMaterializedViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-materialized-view.
type AlterMaterializedViewOptions struct {
	alter             bool                       `ddl:"static" sql:"ALTER"`
	materializedView  bool                       `ddl:"static" sql:"MATERIALIZED VIEW"`
	name              SchemaObjectIdentifier     `ddl:"identifier"`
	RenameTo          *SchemaObjectIdentifier    `ddl:"identifier" sql:"RENAME TO"`
	ClusterBy         *MaterializedViewClusterBy `ddl:"keyword"`
	DropClusteringKey *bool                      `ddl:"keyword" sql:"DROP CLUSTERING KEY"`
	SuspendRecluster  *bool                      `ddl:"keyword" sql:"SUSPEND RECLUSTER"`
	ResumeRecluster   *bool                      `ddl:"keyword" sql:"RESUME RECLUSTER"`
	Suspend           *bool                      `ddl:"keyword" sql:"SUSPEND"`
	Resume            *bool                      `ddl:"keyword" sql:"RESUME"`
	Set               *MaterializedViewSet       `ddl:"keyword" sql:"SET"`
	Unset             *MaterializedViewUnset     `ddl:"keyword" sql:"UNSET"`
}

type MaterializedViewSet struct {
	Secure  *bool   `ddl:"keyword" sql:"SECURE"`
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type MaterializedViewUnset struct {
	Secure  *bool `ddl:"keyword" sql:"SECURE"`
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropMaterializedViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-materialized-view.
type DropMaterializedViewOptions struct {
	drop             bool                   `ddl:"static" sql:"DROP"`
	materializedView bool                   `ddl:"static" sql:"MATERIALIZED VIEW"`
	IfExists         *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name             SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowMaterializedViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-materialized-views.
type ShowMaterializedViewOptions struct {
	show              bool  `ddl:"static" sql:"SHOW"`
	materializedViews bool  `ddl:"static" sql:"MATERIALIZED VIEWS"`
	Like              *Like `ddl:"keyword" sql:"LIKE"`
	In                *In   `ddl:"keyword" sql:"IN"`
}

type materializedViewDBRow struct {
	CreatedOn           string         `db:"created_on"`
	Name                string         `db:"name"`
	Reserved            sql.NullString `db:"reserved"`
	DatabaseName        string         `db:"database_name"`
	SchemaName          string         `db:"schema_name"`
	ClusterBy           sql.NullString `db:"cluster_by"`
	Rows                int            `db:"rows"`
	Bytes               int            `db:"bytes"`
	SourceDatabaseName  string         `db:"source_database_name"`
	SourceSchemaName    string         `db:"source_schema_name"`
	SourceTableName     string         `db:"source_table_name"`
	RefreshedOn         time.Time      `db:"refreshed_on"`
	CompactedOn         time.Time      `db:"compacted_on"`
	Owner               string         `db:"owner"`
	Invalid             bool           `db:"invalid"`
	InvalidReason       sql.NullString `db:"invalid_reason"`
	BehindBy            string         `db:"behind_by"`
	Comment             sql.NullString `db:"comment"`
	Text                string         `db:"text"`
	IsSecure            bool           `db:"is_secure"`
	AutomaticClustering string         `db:"automatic_clustering"`
	OwnerRoleType       sql.NullString `db:"owner_role_type"`
	Budget              sql.NullString `db:"budget"`
}

type MaterializedView struct {
	CreatedOn           string
	Name                string
	Reserved            *string
	DatabaseName        string
	SchemaName          string
	ClusterBy           string
	Rows                int
	Bytes               int
	SourceDatabaseName  string
	SourceSchemaName    string
	SourceTableName     string
	RefreshedOn         time.Time
	CompactedOn         time.Time
	Owner               string
	Invalid             bool
	InvalidReason       string
	BehindBy            string
	Comment             string
	Text                string
	IsSecure            bool
	AutomaticClustering bool
	OwnerRoleType       string
	Budget              string
}

func (v *MaterializedView) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

// DescribeMaterializedViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-materialized-view.
type DescribeMaterializedViewOptions struct {
	describe         bool                   `ddl:"static" sql:"DESCRIBE"`
	materializedView bool                   `ddl:"static" sql:"MATERIALIZED VIEW"`
	name             SchemaObjectIdentifier `ddl:"identifier"`
}

type materializedViewDetailsRow struct {
	Name       string         `db:"name"`
	Type       DataType       `db:"type"`
	Kind       string         `db:"kind"`
	Null       string         `db:"null?"`
	Default    sql.NullString `db:"default"`
	PrimaryKey string         `db:"primary key"`
	UniqueKey  string         `db:"unique key"`
	Check      sql.NullString `db:"check"`
	Expression sql.NullString `db:"expression"`
	Comment    sql.NullString `db:"comment"`
}

type MaterializedViewDetails struct {
	Name       string
	Type       DataType
	Kind       string
	IsNullable bool
	Default    *string
	IsPrimary  bool
	IsUnique   bool
	Check      *bool
	Expression *string
	Comment    *string
}
