package snowflake

import (
	"fmt"
)

type (
	allGrantType   string
	allGrantTarget string
)

const (
	allSchemaType           allGrantType = "SCHEMA"
	allTableType            allGrantType = "TABLE"
	allViewType             allGrantType = "VIEW"
	allMaterializedViewType allGrantType = "MATERIALIZED VIEW"
	allStageType            allGrantType = "STAGE"
	allExternalTableType    allGrantType = "EXTERNAL TABLE"
	allFileFormatType       allGrantType = "FILE FORMAT"
	allFunctionType         allGrantType = "FUNCTION"
	allProcedureType        allGrantType = "PROCEDURE"
	allSequenceType         allGrantType = "SEQUENCE"
	allStreamType           allGrantType = "STREAM"
	allPipeType             allGrantType = "PIPE"
	allTaskType             allGrantType = "TASK"
)

const (
	allSchemaTarget   allGrantTarget = "SCHEMA"
	allDatabaseTarget allGrantTarget = "DATABASE"
)

// AllGrantBuilder abstracts the creation of ExistingGrantExecutables.
type AllGrantBuilder struct {
	name           string
	qualifiedName  string
	allGrantType   allGrantType
	allGrantTarget allGrantTarget
}

func getNameAndQualifiedNameForAllGrants(db, schema string) (string, string, allGrantTarget) {
	name := schema
	allTarget := allSchemaTarget
	qualifiedName := fmt.Sprintf(`"%v"."%v"`, db, schema)

	if schema == "" {
		name = db
		allTarget = allDatabaseTarget
		qualifiedName = fmt.Sprintf(`"%v"`, db)
	}

	return name, qualifiedName, allTarget
}

// Name returns the object name for this FutureGrantBuilder.
func (fgb *AllGrantBuilder) Name() string {
	return fgb.name
}

func (fgb *AllGrantBuilder) GrantType() string {
	return string(fgb.allGrantType)
}

// ExistingSchemaGrant returns a pointer to a AllGrantBuilder for a schema.
func ExistingSchemaGrant(db string) GrantBuilder {
	return &AllGrantBuilder{
		name:           db,
		qualifiedName:  fmt.Sprintf(`"%v"`, db),
		allGrantType:   allSchemaType,
		allGrantTarget: allDatabaseTarget,
	}
}

// AllTableGrant returns a pointer to a AllGrantBuilder for a table.
func AllTableGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allTableType,
		allGrantTarget: allTarget,
	}
}

// ExistingViewGrant returns a pointer to a AllGrantBuilder for a view.
func ExistingViewGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allViewType,
		allGrantTarget: allTarget,
	}
}

// ExistingMaterializedViewGrant returns a pointer to a AllGrantBuilder for a view.
func ExistingMaterializedViewGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allMaterializedViewType,
		allGrantTarget: allTarget,
	}
}

// ExistingStageGrant returns a pointer to a AllGrantBuilder for a stage.
func ExistingStageGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allStageType,
		allGrantTarget: allTarget,
	}
}

// ExistingExternalTableGrant returns a pointer to a AllGrantBuilder for a external table.
func ExistingExternalTableGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allExternalTableType,
		allGrantTarget: allTarget,
	}
}

// ExistingFileFormatGrant returns a pointer to a AllGrantBuilder for a file format.
func ExistingFileFormatGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allFileFormatType,
		allGrantTarget: allTarget,
	}
}

// ExistingFunctionGrant returns a pointer to a AllGrantBuilder for a function.
func ExistingFunctionGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allFunctionType,
		allGrantTarget: allTarget,
	}
}

// ExistingProcedureGrant returns a pointer to a AllGrantBuilder for a procedure.
func ExistingProcedureGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allProcedureType,
		allGrantTarget: allTarget,
	}
}

// ExistingSequenceGrant returns a pointer to a AllGrantBuilder for a sequence.
func ExistingSequenceGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allSequenceType,
		allGrantTarget: allTarget,
	}
}

// ExistingStreamGrant returns a pointer to a AllGrantBuilder for a stream.
func ExistingStreamGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allStreamType,
		allGrantTarget: allTarget,
	}
}

// ExistingPipeGrant returns a pointer to a AllGrantBuilder for a pipe.
func ExistingPipeGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allPipeType,
		allGrantTarget: allTarget,
	}
}

// ExistingTaskGrant returns a pointer to a AllGrantBuilder for a task.
func ExistingTaskGrant(db, schema string) GrantBuilder {
	name, qualifiedName, allTarget := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   allTaskType,
		allGrantTarget: allTarget,
	}
}

// Show returns the SQL that will show all privileges on the grant.
func (fgb *AllGrantBuilder) Show() string {
	return fmt.Sprintf(`SHOW ALL GRANTS IN %v %v`, fgb.allGrantTarget, fgb.qualifiedName)
}

// ExistingGrantExecutable abstracts the creation of SQL queries to build all grants for
// different all grant types.
type ExistingGrantExecutable struct {
	grantName      string
	granteeName    string
	allGrantType   allGrantType
	allGrantTarget allGrantTarget
}

// Role returns a pointer to a FutureGrantExecutable for a role.
func (fgb *AllGrantBuilder) Role(n string) GrantExecutable {
	return &ExistingGrantExecutable{
		granteeName:    n,
		grantName:      fgb.qualifiedName,
		allGrantType:   fgb.allGrantType,
		allGrantTarget: fgb.allGrantTarget,
	}
}

// Share is not implemented because all objects cannot be granted to shares.
func (fgb *AllGrantBuilder) Share(n string) GrantExecutable {
	return nil
}

// Grant returns the SQL that will grant all privileges on the grant to the grantee.
func (fge *ExistingGrantExecutable) Grant(p string, w bool) string {
	var template string
	if w {
		template = `GRANT %v ON ALL %vS IN %v %v TO ROLE "%v" WITH GRANT OPTION`
	} else {
		template = `GRANT %v ON ALL %vS IN %v %v TO ROLE "%v"`
	}
	return fmt.Sprintf(template,
		p, fge.allGrantType, fge.allGrantTarget, fge.grantName, fge.granteeName)
}

// Revoke returns the SQL that will revoke all privileges on the grant from the grantee.
func (fge *ExistingGrantExecutable) Revoke(p string) []string {
	// TODO: has no effect for ALL GRANTS
	return []string{
		fmt.Sprintf(`REVOKE %v ON ALL %vS IN %v %v FROM ROLE "%v"`,
			p, fge.allGrantType, fge.allGrantTarget, fge.grantName, fge.granteeName),
	}
}

// Show returns the SQL that will show all all grants on the schema.
func (fge *ExistingGrantExecutable) Show() string {
	return fmt.Sprintf(`SHOW ALL GRANTS IN %v %v`, fge.allGrantTarget, fge.grantName)
}
