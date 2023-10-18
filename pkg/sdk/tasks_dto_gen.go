package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateTaskOptions]   = new(CreateTaskRequest)
	_ optionsProvider[CloneTaskOptions]    = new(CloneTaskRequest)
	_ optionsProvider[AlterTaskOptions]    = new(AlterTaskRequest)
	_ optionsProvider[DropTaskOptions]     = new(DropTaskRequest)
	_ optionsProvider[ShowTaskOptions]     = new(ShowTaskRequest)
	_ optionsProvider[DescribeTaskOptions] = new(DescribeTaskRequest)
	_ optionsProvider[ExecuteTaskOptions]  = new(ExecuteTaskRequest)
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

func (r *CreateTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type CreateTaskWarehouseRequest struct {
	Warehouse                           *AccountObjectIdentifier
	UserTaskManagedInitialWarehouseSize *WarehouseSize
}

type CloneTaskRequest struct {
	OrReplace  *bool
	name       SchemaObjectIdentifier // required
	sourceTask SchemaObjectIdentifier // required
	CopyGrants *bool
}

func (r *CloneTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
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
	ErrorIntegration            *string
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
	ErrorIntegration            *bool
	Comment                     *bool
	SessionParametersUnset      *SessionParametersUnset
}

type DropTaskRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowTaskRequest struct {
	Terse      *bool
	Like       *Like
	In         *In
	StartsWith *string
	RootOnly   *bool
	Limit      *LimitFrom
}

type DescribeTaskRequest struct {
	name SchemaObjectIdentifier // required
}

type ExecuteTaskRequest struct {
	name      SchemaObjectIdentifier // required
	RetryLast *bool
}
