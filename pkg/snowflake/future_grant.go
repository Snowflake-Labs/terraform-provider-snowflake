package snowflake

import (
	"fmt"
)

type futureGrantType string
type futureGrantTarget string

const (
	futureSchemaType           futureGrantType = "SCHEMA"
	futureTableType            futureGrantType = "TABLE"
	futureViewType             futureGrantType = "VIEW"
	futureMaterializedViewType futureGrantType = "MATERIALIZED VIEW"
	futureStageType            futureGrantType = "STAGE"
	futureExternalTableType    futureGrantType = "EXTERNAL TABLE"
	futureFileFormatType       futureGrantType = "FILE FORMAT"
	futureFunctionType         futureGrantType = "FUNCTION"
	futureProcedureType        futureGrantType = "PROCEDURE"
	futureSequenceType         futureGrantType = "SEQUENCE"
	futureStreamType           futureGrantType = "STREAM"
	futurePipeType             futureGrantType = "PIPE"
	futureTaskType             futureGrantType = "TASK"
)

const (
	futureSchemaTarget   futureGrantTarget = "SCHEMA"
	futureDatabaseTarget futureGrantTarget = "DATABASE"
)

// FutureGrantBuilder abstracts the creation of FutureGrantExecutables
type FutureGrantBuilder struct {
	name              string
	qualifiedName     string
	futureGrantType   futureGrantType
	futureGrantTarget futureGrantTarget
}

func getNameAndQualifiedName(db, schema string) (string, string, futureGrantTarget) {
	name := schema
	futureTarget := futureSchemaTarget
	qualifiedName := fmt.Sprintf(`"%v"."%v"`, db, schema)

	if schema == "" {
		name = db
		futureTarget = futureDatabaseTarget
		qualifiedName = fmt.Sprintf(`"%v"`, db)
	}

	return name, qualifiedName, futureTarget
}

// Name returns the object name for this FutureGrantBuilder
func (fgb *FutureGrantBuilder) Name() string {
	return fgb.name
}

func (fgb *FutureGrantBuilder) GrantType() string {
	return string(fgb.futureGrantType)
}

// FutureSchemaGrant returns a pointer to a FutureGrantBuilder for a schema
func FutureSchemaGrant(db string) GrantBuilder {
	return &FutureGrantBuilder{
		name:              db,
		qualifiedName:     fmt.Sprintf(`"%v"`, db),
		futureGrantType:   futureSchemaType,
		futureGrantTarget: futureDatabaseTarget,
	}
}

// FutureTableGrant returns a pointer to a FutureGrantBuilder for a table
func FutureTableGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureTableType,
		futureGrantTarget: futureTarget,
	}
}

// FutureViewGrant returns a pointer to a FutureGrantBuilder for a view
func FutureViewGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureViewType,
		futureGrantTarget: futureTarget,
	}
}

// FutureMaterializedViewGrant returns a pointer to a FutureGrantBuilder for a view
func FutureMaterializedViewGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureMaterializedViewType,
		futureGrantTarget: futureTarget,
	}
}

// FutureStageGrant returns a pointer to a FutureGrantBuilder for a stage
func FutureStageGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureStageType,
		futureGrantTarget: futureTarget,
	}
}

// FutureExternalTableGrant returns a pointer to a FutureGrantBuilder for a external table
func FutureExternalTableGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureExternalTableType,
		futureGrantTarget: futureTarget,
	}
}

// FutureFileFormatGrant returns a pointer to a FutureGrantBuilder for a file format
func FutureFileFormatGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureFileFormatType,
		futureGrantTarget: futureTarget,
	}
}

// FutureFunctionGrant returns a pointer to a FutureGrantBuilder for a function
func FutureFunctionGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureFunctionType,
		futureGrantTarget: futureTarget,
	}
}

// FutureProcedureGrant returns a pointer to a FutureGrantBuilder for a procedure
func FutureProcedureGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureProcedureType,
		futureGrantTarget: futureTarget,
	}
}

// FutureSequenceGrant returns a pointer to a FutureGrantBuilder for a sequence
func FutureSequenceGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureSequenceType,
		futureGrantTarget: futureTarget,
	}
}

// FutureStreamGrant returns a pointer to a FutureGrantBuilder for a stream
func FutureStreamGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureStreamType,
		futureGrantTarget: futureTarget,
	}
}

// FuturePipeGrant returns a pointer to a FutureGrantBuilder for a pipe
func FuturePipeGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futurePipeType,
		futureGrantTarget: futureTarget,
	}
}

// FutureTaskGrant returns a pointer to a FutureGrantBuilder for a task
func FutureTaskGrant(db, schema string) GrantBuilder {
	name, qualifiedName, futureTarget := getNameAndQualifiedName(db, schema)
	return &FutureGrantBuilder{
		name:              name,
		qualifiedName:     qualifiedName,
		futureGrantType:   futureTaskType,
		futureGrantTarget: futureTarget,
	}
}

// Show returns the SQL that will show all privileges on the grant
func (fgb *FutureGrantBuilder) Show() string {
	return fmt.Sprintf(`SHOW FUTURE GRANTS IN %v %v`, fgb.futureGrantTarget, fgb.qualifiedName)
}

// FutureGrantExecutable abstracts the creation of SQL queries to build future grants for
// different future grant types.
type FutureGrantExecutable struct {
	grantName         string
	granteeName       string
	futureGrantType   futureGrantType
	futureGrantTarget futureGrantTarget
}

// Role returns a pointer to a FutureGrantExecutable for a role
func (fgb *FutureGrantBuilder) Role(n string) GrantExecutable {
	return &FutureGrantExecutable{
		granteeName:       n,
		grantName:         fgb.qualifiedName,
		futureGrantType:   fgb.futureGrantType,
		futureGrantTarget: fgb.futureGrantTarget,
	}
}

// Share is not implemented because future objects cannot be granted to shares.
func (gb *FutureGrantBuilder) Share(n string) GrantExecutable {
	return nil
}

// Grant returns the SQL that will grant future privileges on the grant to the grantee
func (fge *FutureGrantExecutable) Grant(p string, w bool) string {
	var template string
	if w {
		template = `GRANT %v ON FUTURE %vS IN %v %v TO ROLE "%v" WITH GRANT OPTION`
	} else {
		template = `GRANT %v ON FUTURE %vS IN %v %v TO ROLE "%v"`
	}
	return fmt.Sprintf(template,
		p, fge.futureGrantType, fge.futureGrantTarget, fge.grantName, fge.granteeName)
}

// Revoke returns the SQL that will revoke future privileges on the grant from the grantee
func (fge *FutureGrantExecutable) Revoke(p string) []string {
	return []string{
		fmt.Sprintf(`REVOKE %v ON FUTURE %vS IN %v %v FROM ROLE "%v"`,
			p, fge.futureGrantType, fge.futureGrantTarget, fge.grantName, fge.granteeName),
	}
}

// Show returns the SQL that will show all future grants on the schema
func (fge *FutureGrantExecutable) Show() string {
	return fmt.Sprintf(`SHOW FUTURE GRANTS IN %v %v`, fge.futureGrantTarget, fge.grantName)
}
