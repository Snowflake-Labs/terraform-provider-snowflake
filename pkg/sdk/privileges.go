package sdk

type Privilege string

const (
	PrivilegeUsage          Privilege = "USAGE"
	PrivilegeSelect         Privilege = "SELECT"
	PrivilegeReferenceUsage Privilege = "REFERENCE_USAGE"
)

func (p Privilege) String() string {
	return string(p)
}
