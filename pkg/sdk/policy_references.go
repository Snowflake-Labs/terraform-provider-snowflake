package sdk

import (
	"context"
	"database/sql"
)

var _ convertibleRow[PolicyReference] = new(policyReferenceDBRow)

type PolicyReferences interface {
	GetForEntity(ctx context.Context, request *GetForEntityPolicyReferenceRequest) ([]PolicyReference, error)
}

type getForEntityPolicyReferenceOptions struct {
	selectEverythingFrom bool                       `ddl:"static" sql:"SELECT * FROM TABLE"`
	parameters           *policyReferenceParameters `ddl:"list,parentheses,no_comma"`
}

type policyReferenceParameters struct {
	functionFullyQualifiedName bool                              `ddl:"static" sql:"SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES"`
	arguments                  *policyReferenceFunctionArguments `ddl:"list,parentheses"`
}

type PolicyEntityDomain string

const (
	PolicyEntityDomainAccount     PolicyEntityDomain = "ACCOUNT"
	PolicyEntityDomainIntegration PolicyEntityDomain = "INTEGRATION"
	PolicyEntityDomainTable       PolicyEntityDomain = "TABLE"
	PolicyEntityDomainTag         PolicyEntityDomain = "TAG"
	PolicyEntityDomainUser        PolicyEntityDomain = "USER"
	PolicyEntityDomainView        PolicyEntityDomain = "VIEW"
)

type policyReferenceFunctionArguments struct {
	refEntityName   []ObjectIdentifier  `ddl:"parameter,single_quotes,arrow_equals" sql:"REF_ENTITY_NAME"`
	refEntityDomain *PolicyEntityDomain `ddl:"parameter,single_quotes,arrow_equals" sql:"REF_ENTITY_DOMAIN"`
}

type PolicyReference struct {
	PolicyDb          *string
	PolicySchema      *string
	PolicyName        string
	PolicyKind        string
	RefDatabaseName   *string
	RefSchemaName     *string
	RefEntityName     string
	RefEntityDomain   string
	RefColumnName     *string
	RefArgColumnNames *string
	TagDatabase       *string
	TagSchema         *string
	TagName           *string
	PolicyStatus      string
}

type policyReferenceDBRow struct {
	PolicyDb          sql.NullString `db:"POLICY_DB"`
	PolicySchema      sql.NullString `db:"POLICY_SCHEMA"`
	PolicyName        string         `db:"POLICY_NAME"`
	PolicyKind        string         `db:"POLICY_KIND"`
	RefDatabaseName   sql.NullString `db:"REF_DATABASE_NAME"`
	RefSchemaName     sql.NullString `db:"REF_SCHEMA_NAME"`
	RefEntityName     string         `db:"REF_ENTITY_NAME"`
	RefEntityDomain   string         `db:"REF_ENTITY_DOMAIN"`
	RefColumnName     sql.NullString `db:"REF_COLUMN_NAME"`
	RefArgColumnNames sql.NullString `db:"REF_ARG_COLUMN_NAMES"`
	TagDatabase       sql.NullString `db:"TAG_DATABASE"`
	TagSchema         sql.NullString `db:"TAG_SCHEMA"`
	TagName           sql.NullString `db:"TAG_NAME"`
	PolicyStatus      string         `db:"POLICY_STATUS"`
}

func (row policyReferenceDBRow) convert() *PolicyReference {
	policyReference := PolicyReference{
		PolicyName:      row.PolicyName,
		PolicyKind:      row.PolicyKind,
		RefEntityName:   row.RefEntityName,
		RefEntityDomain: row.RefEntityDomain,
		PolicyStatus:    row.PolicyStatus,
	}
	if row.PolicyDb.Valid {
		policyReference.PolicyDb = &row.PolicyDb.String
	}
	if row.PolicySchema.Valid {
		policyReference.PolicySchema = &row.PolicySchema.String
	}
	if row.RefDatabaseName.Valid {
		policyReference.RefDatabaseName = &row.RefDatabaseName.String
	}
	if row.RefSchemaName.Valid {
		policyReference.RefSchemaName = &row.RefSchemaName.String
	}
	if row.RefColumnName.Valid {
		policyReference.RefColumnName = &row.RefColumnName.String
	}
	if row.RefArgColumnNames.Valid {
		policyReference.RefArgColumnNames = &row.RefArgColumnNames.String
	}
	if row.TagDatabase.Valid {
		policyReference.TagDatabase = &row.TagDatabase.String
	}
	if row.TagSchema.Valid {
		policyReference.TagSchema = &row.TagSchema.String
	}
	if row.TagName.Valid {
		policyReference.TagName = &row.TagName.String
	}
	return &policyReference
}
