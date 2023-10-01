package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSessionPolicyOptions] = new(CreateSessionPolicyRequest)
)

type CreateSessionPolicyRequest struct {
	OrReplace                *bool
	IfNotExists              *bool
	name                     AccountObjectIdentifier // required
	SessionIdleTimeoutMins   *int
	SessionUiIdleTimeoutMins *int
	Comment                  *string
}
