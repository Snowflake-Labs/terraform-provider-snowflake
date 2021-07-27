package snowflake

import (
	"fmt"
	"strings"
)

type grantType string

const (
	accountType grantType = "ACCOUNT"

	resourceMonitorType grantType = "RESOURCE MONITOR"
	integrationType     grantType = "INTEGRATION"

	databaseType         grantType = "DATABASE"
	schemaType           grantType = "SCHEMA"
	stageType            grantType = "STAGE"
	viewType             grantType = "VIEW"
	materializedViewType grantType = "MATERIALIZED VIEW"
	tableType            grantType = "TABLE"
	warehouseType        grantType = "WAREHOUSE"
	externalTableType    grantType = "EXTERNAL TABLE"
	fileFormatType       grantType = "FILE FORMAT"
	functionType         grantType = "FUNCTION"
	procedureType        grantType = "PROCEDURE"
	sequenceType         grantType = "SEQUENCE"
	streamType           grantType = "STREAM"
	maskingPolicyType    grantType = "MASKING POLICY"
	pipeType             grantType = "PIPE"
	taskType             grantType = "TASK"
)

type GrantExecutable interface {
	Grant(p string, w bool) string
	Revoke(p string) []string
	Show() string
}

type GrantBuilder interface {
	Name() string
	GrantType() string
	Role(string) GrantExecutable
	Share(string) GrantExecutable
	Show() string
}

// CurrentGrantBuilder abstracts the creation of GrantExecutables
type CurrentGrantBuilder struct {
	name          string
	qualifiedName string
	grantType     grantType
}

// Name returns the object name for this CurrentGrantBuilder
func (gb *CurrentGrantBuilder) Name() string {
	return gb.name
}

func (gb *CurrentGrantBuilder) GrantType() string {
	return string(gb.grantType)
}

// Show returns the SQL that will show all privileges on the grant
func (gb *CurrentGrantBuilder) Show() string {
	return fmt.Sprintf(`SHOW GRANTS ON %v %v`, gb.grantType, gb.qualifiedName)
}

///////////////////////////////////////////////
// START CurrentMaterializedViewGrantBuilder //
///////////////////////////////////////////////
type CurrentMaterializedViewGrantBuilder struct {
	name          string
	qualifiedName string
	grantType     grantType
}

// Name returns the object name for this CurrentGrantBuilder
func (gb *CurrentMaterializedViewGrantBuilder) Name() string {
	return gb.name
}

func (gb *CurrentMaterializedViewGrantBuilder) GrantType() string {
	return string(gb.grantType)
}

// Show returns the SQL that will show all privileges on the grant
func (gb *CurrentMaterializedViewGrantBuilder) Show() string {
	return fmt.Sprintf(`SHOW GRANTS ON %v %v`, gb.grantType, gb.qualifiedName)
}

// Role returns a pointer to a CurrentGrantExecutable for a role
func (gb *CurrentMaterializedViewGrantBuilder) Role(n string) GrantExecutable {
	return &CurrentGrantExecutable{
		grantName:   gb.qualifiedName,
		grantType:   viewType,
		granteeName: n,
		granteeType: roleType,
	}
}

// Share returns a pointer to a CurrentGrantExecutable for a share
func (gb *CurrentMaterializedViewGrantBuilder) Share(n string) GrantExecutable {
	return &CurrentGrantExecutable{
		grantName:   gb.qualifiedName,
		grantType:   viewType,
		granteeName: n,
		granteeType: shareType,
	}
}

///////////////////////////////////////////////
/// END CurrentMaterializedViewGrantBuilder ///
///////////////////////////////////////////////

// AccountGrant returns a pointer to a CurrentGrantBuilder for an account
func AccountGrant() GrantBuilder {
	return &CurrentGrantBuilder{
		grantType: accountType,
	}
}

// DatabaseGrant returns a pointer to a CurrentGrantBuilder for a database
func DatabaseGrant(name string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          name,
		qualifiedName: fmt.Sprintf(`"%v"`, name),
		grantType:     databaseType,
	}
}

// SchemaGrant returns a pointer to a CurrentGrantBuilder for a schema
func SchemaGrant(db, schema string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          schema,
		qualifiedName: fmt.Sprintf(`"%v"."%v"`, db, schema),
		grantType:     schemaType,
	}
}

// StageGrant returns a pointer to a CurrentGrantBuilder for a stage
func StageGrant(db, schema, stage string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          stage,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, stage),
		grantType:     stageType,
	}
}

// ViewGrant returns a pointer to a CurrentGrantBuilder for a view
func ViewGrant(db, schema, view string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          view,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, view),
		grantType:     viewType,
	}
}

// MaterializedViewGrant returns a pointer to a CurrentGrantBuilder for a materialized view
func MaterializedViewGrant(db, schema, view string) GrantBuilder {
	return &CurrentMaterializedViewGrantBuilder{
		name:          view,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, view),
		grantType:     materializedViewType,
	}
}

// TableGrant returns a pointer to a CurrentGrantBuilder for a table
func TableGrant(db, schema, table string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          table,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, table),
		grantType:     tableType,
	}
}

// ResourceMonitorGrant returns a pointer to a CurrentGrantBuilder for a resource monitor
func ResourceMonitorGrant(w string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          w,
		qualifiedName: fmt.Sprintf(`"%v"`, w),
		grantType:     resourceMonitorType,
	}
}

// IntegrationGrant returns a pointer to a CurrentGrantBuilder for an integration
func IntegrationGrant(w string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          w,
		qualifiedName: fmt.Sprintf(`"%v"`, w),
		grantType:     integrationType,
	}
}

// WarehouseGrant returns a pointer to a CurrentGrantBuilder for a warehouse
func WarehouseGrant(w string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          w,
		qualifiedName: fmt.Sprintf(`"%v"`, w),
		grantType:     warehouseType,
	}
}

// ExternalTableGrant returns a pointer to a CurrentGrantBuilder for an external table
func ExternalTableGrant(db, schema, externalTable string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          externalTable,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, externalTable),
		grantType:     externalTableType,
	}
}

// FileFormatGrant returns a pointer to a CurrentGrantBuilder for a file format
func FileFormatGrant(db, schema, fileFormat string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          fileFormat,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, fileFormat),
		grantType:     fileFormatType,
	}
}

// FunctionGrant returns a pointer to a CurrentGrantBuilder for a view
func FunctionGrant(db, schema, function string, argumentTypes []string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          function,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"(%v)`, db, schema, function, strings.Join(argumentTypes, ", ")),
		grantType:     functionType,
	}
}

// ProcedureGrant returns a pointer to a CurrentGrantBuilder for a procedure
func ProcedureGrant(db, schema, procedure string, argumentTypes []string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          procedure,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"(%v)`, db, schema, procedure, strings.Join(argumentTypes, ", ")),
		grantType:     procedureType,
	}
}

// SequenceGrant returns a pointer to a CurrentGrantBuilder for a sequence
func SequenceGrant(db, schema, sequence string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          sequence,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, sequence),
		grantType:     sequenceType,
	}
}

// StreamGrant returns a pointer to a CurrentGrantBuilder for a stream
func StreamGrant(db, schema, stream string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          stream,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, stream),
		grantType:     streamType,
	}
}

// MaskingPolicyGrant returns a pointer to a CurrentGrantBuilder for a masking policy
func MaskingPolicyGrant(db, schema, maskingPolicy string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          maskingPolicy,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, maskingPolicy),
		grantType:     maskingPolicyType,
	}
}

// PipeGrant returns a pointer to a CurrentGrantBuilder for a pipe
func PipeGrant(db, schema, pipe string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          pipe,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, pipe),
		grantType:     pipeType,
	}
}

// TaskGrant returns a pointer to a CurrentGrantBuilder for a task
func TaskGrant(db, schema, task string) GrantBuilder {
	return &CurrentGrantBuilder{
		name:          task,
		qualifiedName: fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, task),
		grantType:     taskType,
	}
}

type granteeType string

const (
	roleType  granteeType = "ROLE"
	shareType granteeType = "SHARE"
	userType  granteeType = "USER" // user is only supported for RoleGrants.
)

// CurrentGrantExecutable abstracts the creation of SQL queries to build grants for
// different resources
type CurrentGrantExecutable struct {
	grantName   string
	grantType   grantType
	granteeName string
	granteeType granteeType
}

// Role returns a pointer to a CurrentGrantExecutable for a role
func (gb *CurrentGrantBuilder) Role(n string) GrantExecutable {
	return &CurrentGrantExecutable{
		grantName:   gb.qualifiedName,
		grantType:   gb.grantType,
		granteeName: n,
		granteeType: roleType,
	}
}

// Share returns a pointer to a CurrentGrantExecutable for a share
func (gb *CurrentGrantBuilder) Share(n string) GrantExecutable {
	return &CurrentGrantExecutable{
		grantName:   gb.qualifiedName,
		grantType:   gb.grantType,
		granteeName: n,
		granteeType: shareType,
	}
}

// Grant returns the SQL that will grant privileges on the grant to the grantee
func (ge *CurrentGrantExecutable) Grant(p string, w bool) string {
	var template string
	if p == `OWNERSHIP` {
		template = `GRANT %v ON %v %v TO %v "%v" COPY CURRENT GRANTS`
	} else if w {
		template = `GRANT %v ON %v %v TO %v "%v" WITH GRANT OPTION`
	} else {
		template = `GRANT %v ON %v %v TO %v "%v"`
	}
	return fmt.Sprintf(template,
		p, ge.grantType, ge.grantName, ge.granteeType, ge.granteeName)
}

// Revoke returns the SQL that will revoke privileges on the grant from the grantee
func (ge *CurrentGrantExecutable) Revoke(p string) []string {
	// Since 10/2020 Snowflake dropped support for REVOKE OWNERSHIP.
	// It's only possible to transfer it to another role now, so we grant it to Terraform's role.
	if p == `OWNERSHIP` {
		return []string{
			"SET currentRole=CURRENT_ROLE()",
			fmt.Sprintf(`GRANT %v ON %v %v TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, p, ge.grantType, ge.grantName),
		}
	}
	return []string{
		fmt.Sprintf(`REVOKE %v ON %v %v FROM %v "%v"`,
			p, ge.grantType, ge.grantName, ge.granteeType, ge.granteeName),
	}
}

// Show returns the SQL that will show all grants of the grantee
func (ge *CurrentGrantExecutable) Show() string {
	return fmt.Sprintf(`SHOW GRANTS OF %v "%v"`, ge.granteeType, ge.granteeName)
}
