package sdk

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type CortexSearchServices interface {
	Create(ctx context.Context, request *CreateCortexSearchServiceRequest) error
	Alter(ctx context.Context, request *AlterCortexSearchServiceRequest) error
	Describe(ctx context.Context, request *DescribeCortexSearchServiceRequest) (*CortexSearchServiceDetails, error)
	Drop(ctx context.Context, request *DropCortexSearchServiceRequest) error
	Show(ctx context.Context, request *ShowCortexSearchServiceRequest) ([]CortexSearchService, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*CortexSearchService, error)
}

// createCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/create-cortex-search
type createCortexSearchServiceOptions struct {
	create              bool                    `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	cortexSearchService bool                    `ddl:"static" sql:"CORTEX SEARCH SERVICE"`
	IfNotExists         *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                SchemaObjectIdentifier  `ddl:"identifier"`
	on                  string                  `ddl:"parameter,no_equals,no_quotes" sql:"ON"`
	attributes          *Attributes             `ddl:"keyword"`
	warehouse           AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	targetLag           string                  `ddl:"parameter,single_quotes" sql:"TARGET_LAG"`
	Comment             *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	query               string                  `ddl:"parameter,no_equals,no_quotes" sql:"AS"`
}

type Attributes struct {
	attributes bool     `ddl:"static" sql:"ATTRIBUTES"`
	columns    []string `ddl:"list,no_parentheses,no_equals"`
}

type CortexSearchServiceSet struct {
	TargetLag *string                  `ddl:"parameter,single_quotes" sql:"TARGET_LAG"`
	Warehouse *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
}

// alterCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/alter-cortex-search
type alterCortexSearchServiceOptions struct {
	alter               bool                    `ddl:"static" sql:"ALTER"`
	cortexSearchService bool                    `ddl:"static" sql:"CORTEX SEARCH SERVICE"`
	IfExists            *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                SchemaObjectIdentifier  `ddl:"identifier"`
	Set                 *CortexSearchServiceSet `ddl:"keyword" sql:"SET"`
}

// dropCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/drop-cortex-search
type dropCortexSearchServiceOptions struct {
	drop                bool                   `ddl:"static" sql:"DROP"`
	cortexSearchService bool                   `ddl:"static" sql:"CORTEX SEARCH SERVICE"`
	IfExists            *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
}

// showCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/show-cortex-search
type showCortexSearchServiceOptions struct {
	show                bool       `ddl:"static" sql:"SHOW"`
	cortexSearchService bool       `ddl:"static" sql:"CORTEX SEARCH SERVICES"`
	Like                *Like      `ddl:"keyword" sql:"LIKE"`
	In                  *In        `ddl:"keyword" sql:"IN"`
	StartsWith          *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit               *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type CortexSearchService struct {
	CreatedOn    time.Time
	Name         string
	DatabaseName string
	SchemaName   string
	Comment      string
}

func (dt *CortexSearchService) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(dt.DatabaseName, dt.SchemaName, dt.Name)
}

type cortexSearchServiceRow struct {
	CreatedOn    time.Time `db:"created_on"`
	Name         string    `db:"name"`
	DatabaseName string    `db:"database_name"`
	SchemaName   string    `db:"schema_name"`
	Comment      string    `db:"comment"`
}

func (cssr cortexSearchServiceRow) convert() *CortexSearchService {
	css := &CortexSearchService{
		CreatedOn:    cssr.CreatedOn,
		Name:         cssr.Name,
		DatabaseName: cssr.DatabaseName,
		SchemaName:   cssr.SchemaName,
		Comment:      cssr.Comment,
	}
	return css
}

// describeCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/desc-cortex-search
type describeCortexSearchServiceOptions struct {
	describe            bool                   `ddl:"static" sql:"DESCRIBE"`
	cortexSearchService bool                   `ddl:"static" sql:"CORTEX SEARCH SERVICE"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
}

type CortexSearchServiceDetails struct {
	Name           string
	Schema         string
	Database       string
	Warehouse      string
	TargetLag      string
	On             string
	Attributes     []string
	ServiceUrl     string
	NumRowsIndexed *int
	Comment        string
}

type cortexSearchServiceDetailsRow struct {
	Name           string         `db:"name"`
	Schema         string         `db:"schema"`
	Database       string         `db:"database"`
	Warehouse      string         `db:"warehouse"`
	TargetLag      string         `db:"target_lag"`
	On             string         `db:"search_column"`
	Attributes     sql.NullString `db:"included_columns"`
	ServiceUrl     sql.NullString `db:"service_url"`
	NumRowsIndexed sql.NullInt64  `db:"num_rows_indexed"`
	Comment        sql.NullString `db:"comment"`
}

func (row cortexSearchServiceDetailsRow) convert() *CortexSearchServiceDetails {
	cssd := CortexSearchServiceDetails{
		Name:      row.Name,
		Schema:    row.Schema,
		Database:  row.Database,
		Warehouse: row.Warehouse,
		TargetLag: row.TargetLag,
		On:        row.On,
	}
	if row.Attributes.Valid {
		for _, elem := range strings.Split(row.Attributes.String, ",") {
			if strings.TrimSpace(elem) != "" {
				cssd.Attributes = append(cssd.Attributes, strings.TrimSpace(elem))
			}
		}
	}
	if row.ServiceUrl.Valid {
		cssd.ServiceUrl = row.ServiceUrl.String
	}
	if row.NumRowsIndexed.Valid {
		cssd.NumRowsIndexed = Int(int(row.NumRowsIndexed.Int64))
	}
	if row.Comment.Valid {
		cssd.Comment = row.Comment.String
	}
	return &cssd
}
