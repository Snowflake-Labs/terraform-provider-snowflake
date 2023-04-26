package snowflake

import (
	"fmt"
)

type (
	AllGrantType   string
	AllGrantTarget string
)

const (
	AllGrantTypeSchema           AllGrantType = "SCHEMA"
	AllGrantTypeTable            AllGrantType = "TABLE"
	AllGrantTypeView             AllGrantType = "VIEW"
	AllGrantTypeMaterializedView AllGrantType = "MATERIALIZED VIEW"
	AllGrantTypeStage            AllGrantType = "STAGE"
	AllGrantTypeExternalTable    AllGrantType = "EXTERNAL TABLE"
	AllGrantTypeFileFormat       AllGrantType = "FILE FORMAT"
	AllGrantTypeFunction         AllGrantType = "FUNCTION"
	AllGrantTypeProcedure        AllGrantType = "PROCEDURE"
	AllGrantTypeSequence         AllGrantType = "SEQUENCE"
	AllGrantTypeStream           AllGrantType = "STREAM"
	AllGrantTypePipe             AllGrantType = "PIPE"
	AllGrantTypeTask             AllGrantType = "TASK"
)

const (
	AllGrantTargetSchema   AllGrantTarget = "SCHEMA"
	AllGrantTargetDatabase AllGrantTarget = "DATABASE"
)

// AllGrantBuilder abstracts the creation of ExistingGrantExecutables.
type AllGrantBuilder struct {
	name           string
	qualifiedName  string
	allGrantType   AllGrantType
	allGrantTarget AllGrantTarget
}

func getNameAndQualifiedNameForAllGrants(db, schema string) (string, string, AllGrantTarget) {
	name := schema
	AllGrantTarget := AllGrantTargetSchema
	qualifiedName := fmt.Sprintf(`"%v"."%v"`, db, schema)

	if schema == "" {
		name = db
		AllGrantTarget = AllGrantTargetDatabase
		qualifiedName = fmt.Sprintf(`"%v"`, db)
	}

	return name, qualifiedName, AllGrantTarget
}

// Name returns the object name for this FutureGrantBuilder.
func (agb *AllGrantBuilder) Name() string {
	return agb.name
}

func (agb *AllGrantBuilder) GrantType() string {
	return string(agb.allGrantType)
}

// AllSchemaGrant returns a pointer to a AllGrantBuilder for a schema.
func AllSchemaGrant(db string) GrantBuilder {
	return &AllGrantBuilder{
		name:           db,
		qualifiedName:  fmt.Sprintf(`"%v"`, db),
		allGrantType:   AllGrantTypeSchema,
		allGrantTarget: AllGrantTargetDatabase,
	}
}

// AllTableGrant returns a pointer to a AllGrantBuilder for a table.
func AllTableGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeTable,
		allGrantTarget: target,
	}
}

// AllViewGrant returns a pointer to a AllGrantBuilder for a view.
func AllViewGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeView,
		allGrantTarget: target,
	}
}

// AllMaterializedViewGrant returns a pointer to a AllGrantBuilder for a view.
func AllMaterializedViewGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeMaterializedView,
		allGrantTarget: target,
	}
}

// AllStageGrant returns a pointer to a AllGrantBuilder for a stage.
func AllStageGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeStage,
		allGrantTarget: target,
	}
}

// ExistingExternalTableGrant returns a pointer to a AllGrantBuilder for a external table.
func ExistingExternalTableGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeExternalTable,
		allGrantTarget: target,
	}
}

// ExistingFileFormatGrant returns a pointer to a AllGrantBuilder for a file format.
func ExistingFileFormatGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeFileFormat,
		allGrantTarget: target,
	}
}

// ExistingFunctionGrant returns a pointer to a AllGrantBuilder for a function.
func ExistingFunctionGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeFunction,
		allGrantTarget: target,
	}
}

// ExistingProcedureGrant returns a pointer to a AllGrantBuilder for a procedure.
func ExistingProcedureGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeProcedure,
		allGrantTarget: target,
	}
}

// ExistingSequenceGrant returns a pointer to a AllGrantBuilder for a sequence.
func ExistingSequenceGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeSequence,
		allGrantTarget: target,
	}
}

// ExistingStreamGrant returns a pointer to a AllGrantBuilder for a stream.
func ExistingStreamGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeStream,
		allGrantTarget: target,
	}
}

// ExistingPipeGrant returns a pointer to a AllGrantBuilder for a pipe.
func ExistingPipeGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypePipe,
		allGrantTarget: target,
	}
}

// ExistingTaskGrant returns a pointer to a AllGrantBuilder for a task.
func ExistingTaskGrant(db, schema string) GrantBuilder {
	name, qualifiedName, target := getNameAndQualifiedNameForAllGrants(db, schema)
	return &AllGrantBuilder{
		name:           name,
		qualifiedName:  qualifiedName,
		allGrantType:   AllGrantTypeTask,
		allGrantTarget: target,
	}
}

// Show returns the SQL that will show all privileges on the grant.
func (agb *AllGrantBuilder) Show() string {
	return fmt.Sprintf(`SHOW ALL GRANTS IN %v %v`, agb.allGrantTarget, agb.qualifiedName)
}

// Role returns a pointer to a ExistingGrantExecutable for a role.
func (agb *AllGrantBuilder) Role(n string) GrantExecutable {
	return &ExistingGrantExecutable{
		granteeName:    n,
		grantName:      agb.qualifiedName,
		allGrantType:   agb.allGrantType,
		allGrantTarget: agb.allGrantTarget,
	}
}

// Share is not implemented because all objects cannot be granted to shares.
func (agb *AllGrantBuilder) Share(n string) GrantExecutable {
	return nil
}

// ExistingGrantExecutable abstracts the creation of SQL queries to build all grants for
// different all grant types.
type ExistingGrantExecutable struct {
	grantName      string
	granteeName    string
	allGrantType   AllGrantType
	allGrantTarget AllGrantTarget
}

// Grant returns the SQL that will grant all privileges on the grant to the grantee.
func (ege *ExistingGrantExecutable) Grant(p string, w bool) string {
	var template string
	if w {
		template = `GRANT %v ON ALL %vS IN %v %v TO ROLE "%v" WITH GRANT OPTION`
	} else {
		template = `GRANT %v ON ALL %vS IN %v %v TO ROLE "%v"`
	}
	return fmt.Sprintf(template,
		p, ege.allGrantType, ege.allGrantTarget, ege.grantName, ege.granteeName)
}

// Revoke returns the SQL that will revoke all privileges on the grant from the grantee.
func (ege *ExistingGrantExecutable) Revoke(p string) []string {
	// Note: has no effect for ALL GRANTS
	return []string{
		fmt.Sprintf(`REVOKE %v ON ALL %vS IN %v %v FROM ROLE "%v"`,
			p, ege.allGrantType, ege.allGrantTarget, ege.grantName, ege.granteeName),
	}
}

// Revoke returns the SQL that will revoke ownership privileges on the grant from the grantee.
// Note: returns the same SQL as Revoke.
func (ege *ExistingGrantExecutable) RevokeOwnership(r string) []string {
	// Note: has no effect for ALL GRANTS
	return []string{
		fmt.Sprintf(`REVOKE OWNERSHIP ON ALL %vS IN %v %v FROM ROLE "%v"`,
			ege.allGrantType, ege.allGrantTarget, ege.grantName, ege.granteeName),
	}
}

// Show returns the SQL that will show all grants on the schema.
func (ege *ExistingGrantExecutable) Show() string {
	// Note: There is no `SHOW ALL GRANTS IN \"test_db\"`, therefore changed the query to `SHOW ALL GRANTS IN \"test_db\"` to have a command, which runs in snowflake.
	return fmt.Sprintf(`SHOW GRANTS ON %v %v`, ege.allGrantTarget, ege.grantName)
}
