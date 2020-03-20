package resources

type privilege string

func (p privilege) string() string {
	return string(p)
}

const (
	privilegeAll                    privilege = "ALL"
	privilegeSelect                 privilege = "SELECT"
	privilegeInsert                 privilege = "INSERT"
	privilegeUpdate                 privilege = "UPDATE"
	privilegeDelete                 privilege = "DELETE"
	privilegeTruncate               privilege = "TRUNCATE"
	privilegeReferences             privilege = "REFERENCES"
	privilegeCreateSchema           privilege = "CREATE SCHEMA"
	privilegeImportedPrivileges     privilege = "IMPORTED PRIVILEGES"
	privilegeModify                 privilege = "MODIFY"
	privilegeOperate                privilege = "OPERATE"
	privilegeMonitor                privilege = "MONITOR"
	privilegeOwnership              privilege = "OWNERSHIP"
	privilegeRead                   privilege = "READ"
	privilegeReferenceUsage         privilege = "REFERENCE_USAGE"
	privilegeUsage                  privilege = "USAGE"
	privilegeWrite                  privilege = "WRITE"
	privilegeCreateTable            privilege = "CREATE TABLE"
	privilegeCreateView             privilege = "CREATE VIEW"
	privilegeCreateFileFormat       privilege = "CREATE FILE FORMAT"
	privilegeCreateStage            privilege = "CREATE STAGE"
	privilegeCreatePipe             privilege = "CREATE PIPE"
	privilegeCreateStream           privilege = "CREATE STREAM"
	privilegeCreateTask             privilege = "CREATE TASK"
	privilegeCreateSequence         privilege = "CREATE SEQUENCE"
	privilegeCreateFunction         privilege = "CREATE FUNCTION"
	privilegeCreateProcedure        privilege = "CREATE PROCEDURE"
	privilegeCreateExternalTable    privilege = "CREATE EXTERNAL TABLE"
	privilegeCreateMaterializedView privilege = "CREATE MATERIALIZED VIEW"
	privilegeCreateTemporaryTable   privilege = "CREATE TEMPORARY TABLE"

	privilegeCreateRole        privilege = "CREATE ROLE"
	privilegeCreateUser        privilege = "CREATE USER"
	privilegeCreateWarehouse   privilege = "CREATE WAREHOUSE"
	privilegeCreateDatabase    privilege = "CREATE DATABASE"
	privilegeCreateIntegration privilege = "CREATE INTEGRATION"
	privilegeManageGrants      privilege = "MANAGE GRANTS"
	privilegeMonitorUsage      privilege = "MONITOR USAGE"
)

type privilegeSet map[privilege]struct{}

func newPrivilegeSet(privileges ...privilege) privilegeSet {
	ps := privilegeSet{}
	for _, priv := range privileges {
		ps[priv] = struct{}{}
	}
	return ps
}

func (ps privilegeSet) toList() []string {
	privs := []string{}
	for p := range ps {
		privs = append(privs, string(p))
	}
	return privs
}

func (ps privilegeSet) addString(s string) {
	ps[privilege(s)] = struct{}{}
}

func (ps privilegeSet) hasString(s string) bool {
	_, ok := ps[privilege(s)]
	return ok
}

func (ps privilegeSet) ALLPrivsPresent(validPrivs privilegeSet) bool {
	for p := range validPrivs {
		if p == privilegeAll || p == privilegeOwnership || p == privilegeCreateStream {
			continue
		}
		if _, ok := ps[p]; !ok {
			return false
		}
	}
	return true
}
