package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSessionPolicyOptions] = new(CreateSessionPolicyRequest)
	_ optionsProvider[AlterSessionPolicyOptions]  = new(AlterSessionPolicyRequest)
)

type CreateSessionPolicyRequest struct {
	OrReplace                *bool
	IfNotExists              *bool
	name                     AccountObjectIdentifier // required
	SessionIdleTimeoutMins   *int
	SessionUiIdleTimeoutMins *int
	Comment                  *string
}

type AlterSessionPolicyRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	RenameTo  *AccountObjectIdentifier
	Set       *SessionPolicySetRequest
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
	Unset     *SessionPolicyUnsetRequest
}

type SessionPolicySetRequest struct {
	SessionIdleTimeoutMins   *int
	SessionUiIdleTimeoutMins *int
	Comment                  *string
}

type SessionPolicyUnsetRequest struct {
	SessionIdleTimeoutMins   *bool
	SessionUiIdleTimeoutMins *bool
	Comment                  *bool
}
