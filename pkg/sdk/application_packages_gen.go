package sdk

import (
	"context"
	"database/sql"
)

type ApplicationPackages interface {
	Create(ctx context.Context, request *CreateApplicationPackageRequest) error
	Alter(ctx context.Context, request *AlterApplicationPackageRequest) error
	Drop(ctx context.Context, request *DropApplicationPackageRequest) error
	Show(ctx context.Context, request *ShowApplicationPackageRequest) ([]ApplicationPackage, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ApplicationPackage, error)
}

// CreateApplicationPackageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-application-package.
type CreateApplicationPackageOptions struct {
	create                     bool                    `ddl:"static" sql:"CREATE"`
	applicationPackage         bool                    `ddl:"static" sql:"APPLICATION PACKAGE"`
	IfNotExists                *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       AccountObjectIdentifier `ddl:"identifier"`
	DataRetentionTimeInDays    *int                    `ddl:"parameter,no_quotes" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int                    `ddl:"parameter,no_quotes" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDdlCollation        *string                 `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Distribution               *Distribution           `ddl:"parameter" sql:"DISTRIBUTION"`
	Tag                        []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
}

// AlterApplicationPackageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-application-package.
type AlterApplicationPackageOptions struct {
	alter                      bool                        `ddl:"static" sql:"ALTER"`
	applicationPackage         bool                        `ddl:"static" sql:"APPLICATION PACKAGE"`
	IfExists                   *bool                       `ddl:"keyword" sql:"IF EXISTS"`
	name                       AccountObjectIdentifier     `ddl:"identifier"`
	Set                        *ApplicationPackageSet      `ddl:"keyword" sql:"SET"`
	Unset                      *ApplicationPackageUnset    `ddl:"list,no_parentheses" sql:"UNSET"`
	ModifyReleaseDirective     *ModifyReleaseDirective     `ddl:"keyword" sql:"MODIFY RELEASE DIRECTIVE"`
	SetDefaultReleaseDirective *SetDefaultReleaseDirective `ddl:"keyword" sql:"SET DEFAULT RELEASE DIRECTIVE"`
	SetReleaseDirective        *SetReleaseDirective        `ddl:"keyword" sql:"SET RELEASE DIRECTIVE"`
	UnsetReleaseDirective      *UnsetReleaseDirective      `ddl:"keyword" sql:"UNSET RELEASE DIRECTIVE"`
	AddVersion                 *AddVersion                 `ddl:"keyword" sql:"ADD VERSION"`
	DropVersion                *DropVersion                `ddl:"keyword" sql:"DROP VERSION"`
	AddPatchForVersion         *AddPatchForVersion         `ddl:"keyword" sql:"ADD PATCH FOR VERSION"`
	SetTags                    []TagAssociation            `ddl:"keyword" sql:"SET TAG"`
	UnsetTags                  []ObjectIdentifier          `ddl:"keyword" sql:"UNSET TAG"`
}

type ApplicationPackageSet struct {
	DataRetentionTimeInDays    *int          `ddl:"parameter,no_quotes" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int          `ddl:"parameter,no_quotes" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDdlCollation        *string       `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *string       `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Distribution               *Distribution `ddl:"parameter" sql:"DISTRIBUTION"`
}

type ApplicationPackageUnset struct {
	DataRetentionTimeInDays    *bool `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *bool `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDdlCollation        *bool `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
	Distribution               *bool `ddl:"keyword" sql:"DISTRIBUTION"`
}

type ModifyReleaseDirective struct {
	ReleaseDirective string `ddl:"keyword,no_quotes"`
	Version          string `ddl:"parameter,no_quotes" sql:"VERSION"`
	Patch            int    `ddl:"parameter,no_quotes" sql:"PATCH"`
}

type SetDefaultReleaseDirective struct {
	Version string `ddl:"parameter,no_quotes" sql:"VERSION"`
	Patch   int    `ddl:"parameter,no_quotes" sql:"PATCH"`
}

type SetReleaseDirective struct {
	ReleaseDirective string   `ddl:"keyword,no_quotes"`
	Accounts         []string `ddl:"parameter,no_quotes,must_parentheses" sql:"ACCOUNTS"`
	Version          string   `ddl:"parameter,no_quotes" sql:"VERSION"`
	Patch            int      `ddl:"parameter,no_quotes" sql:"PATCH"`
}

type UnsetReleaseDirective struct {
	ReleaseDirective string `ddl:"keyword,no_quotes"`
}

type AddVersion struct {
	VersionIdentifier *string `ddl:"keyword,no_quotes"`
	Using             string  `ddl:"parameter,single_quotes,no_equals" sql:"USING"`
	Label             *string `ddl:"parameter,single_quotes" sql:"LABEL"`
}

type DropVersion struct {
	VersionIdentifier string `ddl:"keyword,no_quotes"`
}

type AddPatchForVersion struct {
	VersionIdentifier *string `ddl:"keyword,no_quotes"`
	Using             string  `ddl:"parameter,single_quotes,no_equals" sql:"USING"`
	Label             *string `ddl:"parameter,single_quotes" sql:"LABEL"`
}

// DropApplicationPackageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-application-package.
type DropApplicationPackageOptions struct {
	drop               bool                    `ddl:"static" sql:"DROP"`
	applicationPackage bool                    `ddl:"static" sql:"APPLICATION PACKAGE"`
	name               AccountObjectIdentifier `ddl:"identifier"`
}

// ShowApplicationPackageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-application-packages.
type ShowApplicationPackageOptions struct {
	show                bool       `ddl:"static" sql:"SHOW"`
	applicationPackages bool       `ddl:"static" sql:"APPLICATION PACKAGES"`
	Like                *Like      `ddl:"keyword" sql:"LIKE"`
	StartsWith          *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit               *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type applicationPackageRow struct {
	CreatedOn        string         `db:"created_on"`
	Name             string         `db:"name"`
	IsDefault        string         `db:"is_default"`
	IsCurrent        string         `db:"is_current"`
	Distribution     string         `db:"distribution"`
	Owner            string         `db:"owner"`
	Comment          string         `db:"comment"`
	RetentionTime    int            `db:"retention_time"`
	Options          string         `db:"options"`
	DroppedOn        sql.NullString `db:"dropped_on"`
	ApplicationClass sql.NullString `db:"application_class"`
}

type ApplicationPackage struct {
	CreatedOn        string
	Name             string
	IsDefault        bool
	IsCurrent        bool
	Distribution     string
	Owner            string
	Comment          string
	RetentionTime    int
	Options          string
	DroppedOn        string
	ApplicationClass string
}
