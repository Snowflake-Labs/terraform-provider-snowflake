package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateTaskOptions] = new(CreateTaskRequest)
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
