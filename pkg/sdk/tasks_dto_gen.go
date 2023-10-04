package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateTaskOptions] = new(CreateTaskRequest)
	_ optionsProvider[AlterTaskOptions]  = new(AlterTaskRequest)
)

type CreateTaskRequest struct {
	OrReplace                   *bool
	IfNotExists                 *bool
	name                        SchemaObjectIdentifier // required
	Warehouse                   *CreateTaskWarehouseRequest
	Schedule                    *string
	Config                      *string
	AllowOverlappingExecution   *bool
	SessionParameters           *SessionParameters
	UserTaskTimeoutMs           *int
	SuspendTaskAfterNumFailures *int
	ErrorIntegration            *string
	CopyGrants                  *bool
	Comment                     *string
	After                       []SchemaObjectIdentifier
	Tag                         []TagAssociation
	When                        *string
	sql                         string // required
}

type CreateTaskWarehouseRequest struct {
	Warehouse                           *AccountObjectIdentifier
	UserTaskManagedInitialWarehouseSize *string
}

type AlterTaskRequest struct {
	IfExists    *bool
	name        SchemaObjectIdentifier // required
	Resume      *bool
	Suspend     *bool
	RemoveAfter []SchemaObjectIdentifier
	AddAfter    []SchemaObjectIdentifier
	Set         *TaskSetRequest
	Unset       *TaskUnsetRequest
	SetTags     []TagAssociation
	UnsetTags   []ObjectIdentifier
	ModifyAs    *string
	ModifyWhen  *string
}

type TaskSetRequest struct {
	Warehouse                   *AccountObjectIdentifier
	Schedule                    *string
	Config                      *string
	AllowOverlappingExecution   *bool
	UserTaskTimeoutMs           *int
	SuspendTaskAfterNumFailures *int
	Comment                     *string
	SessionParameters           *SessionParameters
}

type TaskUnsetRequest struct {
	Warehouse                   *bool
	Schedule                    *bool
	Config                      *bool
	AllowOverlappingExecution   *bool
	UserTaskTimeoutMs           *bool
	SuspendTaskAfterNumFailures *bool
	Comment                     *bool
	SessionParametersUnset      *SessionParametersUnset
}
