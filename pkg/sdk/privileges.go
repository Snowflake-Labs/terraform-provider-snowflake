package sdk

type GlobalPrivilege string

const (
	// CREATE {
	//	ACCOUNT | DATA EXCHANGE LISTING | DATABASE | FAILOVER GROUP | INTEGRATION
	//	| NETWORK POLICY | REPLICATION GROUP | ROLE | SHARE | USER | WAREHOUSE
	// }
	GlobalPrivilegeCreateAccount             GlobalPrivilege = "CREATE ACCOUNT"
	GlobalPrivilegeCreateDataExchangeListing GlobalPrivilege = "CREATE DATA EXCHANGE LISTING"
	GlobalPrivilegeCreateDatabase            GlobalPrivilege = "CREATE DATABASE"
	GlobalPrivilegeCreateFailoverGroup       GlobalPrivilege = "CREATE FAILOVER GROUP"
	GlobalPrivilegeCreateIntegration         GlobalPrivilege = "CREATE INTEGRATION"
	GlobalPrivilegeCreateNetworkPolicy       GlobalPrivilege = "CREATE NETWORK POLICY"
	GlobalPrivilegeCreateExternalVolume      GlobalPrivilege = "CREATE EXTERNAL VOLUME"
	GlobalPrivilegeCreateReplicationGroup    GlobalPrivilege = "CREATE REPLICATION GROUP"
	GlobalPrivilegeCreateRole                GlobalPrivilege = "CREATE ROLE"
	GlobalPrivilegeCreateShare               GlobalPrivilege = "CREATE SHARE"
	GlobalPrivilegeCreateUser                GlobalPrivilege = "CREATE USER"
	GlobalPrivilegeCreateWarehouse           GlobalPrivilege = "CREATE WAREHOUSE"

	// | APPLY { { MASKING | PASSWORD | ROW ACCESS | SESSION } POLICY | TAG }
	GlobalPrivilegeApplyMaskingPolicy   GlobalPrivilege = "APPLY MASKING POLICY"
	GlobalPrivilegeApplyPasswordPolicy  GlobalPrivilege = "APPLY PASSWORD POLICY"
	GlobalPrivilegeApplyRowAccessPolicy GlobalPrivilege = "APPLY ROW ACCESS POLICY"
	GlobalPrivilegeApplySessionPolicy   GlobalPrivilege = "APPLY SESSION POLICY"
	GlobalPrivilegeApplyTag             GlobalPrivilege = "APPLY TAG"

	// | ATTACH POLICY | AUDIT |
	GlobalPrivilegeAttachPolicy GlobalPrivilege = "ATTACH POLICY"
	GlobalPrivilegeAudit        GlobalPrivilege = "AUDIT"

	// | EXECUTE { ALERT | TASK }
	GlobalPrivilegeExecuteAlert GlobalPrivilege = "EXECUTE ALERT"
	GlobalPrivilegeExecuteTask  GlobalPrivilege = "EXECUTE TASK"
	// | IMPORT SHARE
	GlobalPrivilegeImportShare GlobalPrivilege = "IMPORT SHARE"
	// | MANAGE GRANTS
	GlobalPrivilegeManageGrants GlobalPrivilege = "MANAGE GRANTS"
	// | MANAGE WAREHOUSES
	GlobalPrivilegeManageWarehouses GlobalPrivilege = "MANAGE WAREHOUSES"

	// | MODIFY { LOG LEVEL | TRACE LEVEL | SESSION LOG LEVEL | SESSION TRACE LEVEL }
	GlobalPrivilegeModifyLogLevel          GlobalPrivilege = "MODIFY LOG LEVEL"
	GlobalPrivilegeModifyTraceLevel        GlobalPrivilege = "MODIFY TRACE LEVEL"
	GlobalPrivilegeModifySessionLogLevel   GlobalPrivilege = "MODIFY SESSION LOG LEVEL"
	GlobalPrivilegeModifySessionTraceLevel GlobalPrivilege = "MODIFY SESSION TRACE LEVEL"

	// | MONITOR { EXECUTION | USAGE }
	GlobalPrivilegeMonitorExecution GlobalPrivilege = "MONITOR EXECUTION"
	GlobalPrivilegeMonitorUsage     GlobalPrivilege = "MONITOR USAGE"

	// | OVERRIDE SHARE RESTRICTIONS | RESOLVE ALL
	GlobalPrivilegeOverrideShareRestrictions GlobalPrivilege = "OVERRIDE SHARE RESTRICTIONS"
	GlobalPrivilegeResolveAll                GlobalPrivilege = "RESOLVE ALL"
)

func (p GlobalPrivilege) String() string {
	return string(p)
}

type AccountObjectPrivilege string

const (
	// -- For DATABASE
	// { CREATE { DATABASE ROLE | SCHEMA } | IMPORTED PRIVILEGES | MODIFY | MONITOR | USAGE } [ , ... ]
	AccountObjectPrivilegeCreateDatabaseRole AccountObjectPrivilege = "CREATE DATABASE ROLE"
	AccountObjectPrivilegeCreateSchema       AccountObjectPrivilege = "CREATE SCHEMA"
	AccountObjectPrivilegeImportedPrivileges AccountObjectPrivilege = "IMPORTED PRIVILEGES"
	AccountObjectPrivilegeModify             AccountObjectPrivilege = "MODIFY"
	AccountObjectPrivilegeMonitor            AccountObjectPrivilege = "MONITOR"
	AccountObjectPrivilegeUsage              AccountObjectPrivilege = "USAGE"

	// -- For EXTERNAL VOLUME
	// AccountObjectPrivilegeUsage              AccountObjectPrivilege = "USAGE" (duplicate)

	// -- For FAILOVER GROUP
	// { FAILOVER | MODIFY | MONITOR | REPLICATE } [ , ... ]
	AccountObjectPrivilegeFailover AccountObjectPrivilege = "FAILOVER"
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor AccountObjectPrivilege = "MONITOR" (duplicate)
	// AccountObjectPrivilegeReplicate AccountObjectPrivilege = "REPLICATE" (duplicate)

	// -- For INTEGRATION
	// { USAGE | USE_ANY_ROLE } [ , ... ]
	// AccountObjectPrivilegeUsage AccountObjectPrivilege = "USAGE" (duplicate)
	AccountObjectPrivilegeUseAnyRole AccountObjectPrivilege = "USE_ANY_ROLE"

	// -- For REPLICATION GROUP
	// { MODIFY | MONITOR | REPLICATE } [ , ... ]
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor AccountObjectPrivilege = "MONITOR" (duplicate)
	AccountObjectPrivilegeReplicate AccountObjectPrivilege = "REPLICATE"

	//-- For RESOURCE MONITOR
	// { MODIFY | MONITOR } [ , ... ]
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor AccountObjectPrivilege = "MONITOR" (duplicate)

	// -- For USER
	// { MONITOR } [ , ... ]
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)

	// -- For WAREHOUSE
	// { MODIFY | MONITOR | USAGE | OPERATE } [ , ... ]
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor AccountObjectPrivilege = "MONITOR" (duplicate)
	// AccountObjectPrivilegeUsage AccountObjectPrivilege = "USAGE" (duplicate)
	AccountObjectPrivilegeOperate AccountObjectPrivilege = "OPERATE"
)

func (p AccountObjectPrivilege) String() string {
	return string(p)
}

type SchemaPrivilege string

const (
	/*
		ADD SEARCH OPTIMIZATION
		| CREATE {
			ALERT | EXTERNAL TABLE | FILE FORMAT | FUNCTION
			| MATERIALIZED VIEW | PIPE | PROCEDURE
			| { MASKING | PASSWORD | ROW ACCESS | SESSION } POLICY
			| SECRET | SEQUENCE | STAGE | STREAM
			| TAG | TABLE | TASK | VIEW
		  }
		| MODIFY | MONITOR | USAGE
		[ , ... ]
	*/
	SchemaPrivilegeAddSearchOptimization  SchemaPrivilege = "ADD SEARCH OPTIMIZATION"
	SchemaPrivilegeApplyBudget            SchemaPrivilege = "APPLYBUDGET"
	SchemaPrivilegeCreateAlert            SchemaPrivilege = "CREATE ALERT"
	SchemaPrivilegeCreateDynamicTable     SchemaPrivilege = "CREATE DYNAMIC TABLE"
	SchemaPrivilegeCreateExternalTable    SchemaPrivilege = "CREATE EXTERNAL TABLE"
	SchemaPrivilegeCreateFileFormat       SchemaPrivilege = "CREATE FILE FORMAT"
	SchemaPrivilegeCreateFunction         SchemaPrivilege = "CREATE FUNCTION"
	SchemaPrivilegeCreateIcebergTable     SchemaPrivilege = "CREATE ICEBERG TABLE"
	SchemaPrivilegeCreateMaterializedView SchemaPrivilege = "CREATE MATERIALIZED VIEW"
	SchemaPrivilegeCreatePipe             SchemaPrivilege = "CREATE PIPE"
	SchemaPrivilegeCreateProcedure        SchemaPrivilege = "CREATE PROCEDURE"
	SchemaPrivilegeCreateMaskingPolicy    SchemaPrivilege = "CREATE MASKING POLICY"
	SchemaPrivilegeCreatePasswordPolicy   SchemaPrivilege = "CREATE PASSWORD POLICY"
	SchemaPrivilegeCreateRowAccessPolicy  SchemaPrivilege = "CREATE ROW ACCESS POLICY"
	SchemaPrivilegeCreateSessionPolicy    SchemaPrivilege = "CREATE SESSION POLICY"
	SchemaPrivilegeCreateSecret           SchemaPrivilege = "CREATE SECRET"
	SchemaPrivilegeCreateSequence         SchemaPrivilege = "CREATE SEQUENCE"
	SchemaPrivilegeCreateStage            SchemaPrivilege = "CREATE STAGE"
	SchemaPrivilegeCreateStream           SchemaPrivilege = "CREATE STREAM"
	SchemaPrivilegeCreateTag              SchemaPrivilege = "CREATE TAG"
	SchemaPrivilegeCreateTable            SchemaPrivilege = "CREATE TABLE"
	SchemaPrivilegeCreateTask             SchemaPrivilege = "CREATE TASK"
	SchemaPrivilegeCreateView             SchemaPrivilege = "CREATE VIEW"
	SchemaPrivilegeModify                 SchemaPrivilege = "MODIFY"
	SchemaPrivilegeMonitor                SchemaPrivilege = "MONITOR"
	SchemaPrivilegeUsage                  SchemaPrivilege = "USAGE"
)

func (p SchemaPrivilege) String() string {
	return string(p)
}

type SchemaObjectPrivilege string

const (
	SchemaObjectOwnership SchemaObjectPrivilege = "OWNERSHIP"

	// -- For ALERT
	// OPERATE [ , ... ]
	SchemaObjectPrivilegeOperate SchemaObjectPrivilege = "OPERATE"

	// -- FOR DYNAMIC TABLE
	//  OPERATE, SELECT [ , ...]
	// SchemaObjectPrivilegeOperate SchemaObjectPrivilege = "OPERATE" (duplicate)
	// SchemaObjectPrivilegeSelect SchemaObjectPrivilege = "SELECT" (duplicate)

	// -- For EVENT TABLE
	// { SELECT | INSERT } [ , ... ]
	SchemaObjectPrivilegeSelect SchemaObjectPrivilege = "SELECT"
	SchemaObjectPrivilegeInsert SchemaObjectPrivilege = "INSERT"

	// -- For FILE FORMAT, FUNCTION (UDF or external function), PROCEDURE, SECRET, or SEQUENCE
	// USAGE [ , ... ]
	SchemaObjectPrivilegeUsage SchemaObjectPrivilege = "USAGE"

	// -- For ICEBERG TABLE
	SchemaObjectPrivilegeApplyBudget SchemaObjectPrivilege = "APPLYBUDGET"
	//SchemaObjectPrivilegeDelete      SchemaObjectPrivilege = "DELETE" (duplicate)
	//SchemaObjectPrivilegeInsert      SchemaObjectPrivilege = "INSERT" (duplicate)
	//SchemaObjectPrivilegeReferences  SchemaObjectPrivilege = "REFERENCES" (duplicate)
	//SchemaObjectPrivilegeSelect      SchemaObjectPrivilege = "SELECT" (duplicate)
	//SchemaObjectPrivilegeTruncate      SchemaObjectPrivilege = "Truncate" (duplicate)
	//SchemaObjectPrivilegeUpdate      SchemaObjectPrivilege = "Update" (duplicate)

	// -- For PIPE
	// { MONITOR | OPERATE } [ , ... ]
	SchemaObjectPrivilegeMonitor SchemaObjectPrivilege = "MONITOR"
	// SchemaObjectPrivilegeOperate SchemaObjectPrivilege = "OPERATE" (duplicate)

	// -- For { MASKING | PASSWORD | ROW ACCESS | SESSION } POLICY or TAG
	// APPLY [ , ... ]
	SchemaObjectPrivilegeApply SchemaObjectPrivilege = "APPLY"

	// -- For external STAGE
	// USAGE [ , ... ]
	// SchemaObjectPrivilegeUsage SchemaObjectPrivilege = "USAGE" (duplicate)

	// -- For internal STAGE
	// READ [ , WRITE ] [ , ... ]
	SchemaObjectPrivilegeRead  SchemaObjectPrivilege = "READ"
	SchemaObjectPrivilegeWrite SchemaObjectPrivilege = "WRITE"

	// -- For STREAM
	// SELECT [ , ... ]
	// SchemaObjectPrivilegeSelect SchemaObjectPrivilege = "SELECT" (duplicate)

	// -- For TABLE
	// { SELECT | INSERT | UPDATE | DELETE | TRUNCATE | REFERENCES } [ , ... ]
	// SchemaObjectPrivilegeSelect SchemaObjectPrivilege = "SELECT" (duplicate)
	// SchemaObjectPrivilegeInsert SchemaObjectPrivilege = "INSERT" (duplicate)
	SchemaObjectPrivilegeUpdate     SchemaObjectPrivilege = "UPDATE"
	SchemaObjectPrivilegeDelete     SchemaObjectPrivilege = "DELETE"
	SchemaObjectPrivilegeTruncate   SchemaObjectPrivilege = "TRUNCATE"
	SchemaObjectPrivilegeReferences SchemaObjectPrivilege = "REFERENCES"

	// -- For TASK
	// { MONITOR | OPERATE } [ , ... ]
	// SchemaObjectPrivilegeMonitor SchemaObjectPrivilege = "MONITOR" (duplicate)
	// SchemaObjectPrivilegeOperate SchemaObjectPrivilege = "OPERATE" (duplicate)

	// -- For VIEW or MATERIALIZED VIEW
	// { SELECT | REFERENCES } [ , ... ]
	// SchemaObjectPrivilegeSelect SchemaObjectPrivilege = "SELECT" (duplicate)
	// SchemaObjectPrivilegeReferences SchemaObjectPrivilege = "REFERENCES" (duplicate)
)

func (p SchemaObjectPrivilege) String() string {
	return string(p)
}

type ObjectPrivilege string

const (
	ObjectPrivilegeUsage          ObjectPrivilege = "USAGE"
	ObjectPrivilegeSelect         ObjectPrivilege = "SELECT"
	ObjectPrivilegeReferenceUsage ObjectPrivilege = "REFERENCE_USAGE"
)

func (p ObjectPrivilege) String() string {
	return string(p)
}
