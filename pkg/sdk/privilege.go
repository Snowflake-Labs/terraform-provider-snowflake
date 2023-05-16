package sdk

type Privilege string

const (
	PrivilegeUsage Privilege = "USAGE"
	PrivilegeSelect Privilege = "SELECT"
	PrivilegeReferenceUsage Privilege = "REFERENCE USAGE"
)

func (p Privilege) String() string {
	return string(p)
}
