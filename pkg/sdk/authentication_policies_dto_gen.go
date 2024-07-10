package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateAuthenticationPolicyOptions]   = new(CreateAuthenticationPolicyRequest)
	_ optionsProvider[AlterAuthenticationPolicyOptions]    = new(AlterAuthenticationPolicyRequest)
	_ optionsProvider[DropAuthenticationPolicyOptions]     = new(DropAuthenticationPolicyRequest)
	_ optionsProvider[ShowAuthenticationPolicyOptions]     = new(ShowAuthenticationPolicyRequest)
	_ optionsProvider[DescribeAuthenticationPolicyOptions] = new(DescribeAuthenticationPolicyRequest)
)

type CreateAuthenticationPolicyRequest struct {
	OrReplace                *bool
	name                     AccountObjectIdentifier // required
	AuthenticationMethods    []AuthenticationMethods
	MfaAuthenticationMethods []MfaAuthenticationMethods
	MfaEnrollment            *string
	ClientTypes              []ClientTypes
	SecurityIntegrations     []SchemaObjectIdentifier
	Comment                  *string
}

type AlterAuthenticationPolicyRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
	Set      *AuthenticationPolicySetRequest
	Unset    *AuthenticationPolicyUnsetRequest
	RenameTo *AccountObjectIdentifier
}

type AuthenticationPolicySetRequest struct {
	AuthenticationMethods    []AuthenticationMethods
	MfaAuthenticationMethods []MfaAuthenticationMethods
	MfaEnrollment            *string
	ClientTypes              []ClientTypes
	SecurityIntegrations     []SchemaObjectIdentifier
	Comment                  *string
}

type AuthenticationPolicyUnsetRequest struct {
	ClientTypes              *bool
	AuthenticationMethods    *bool
	SecurityIntegrations     *bool
	MfaAuthenticationMethods *bool
	MfaEnrollment            *bool
	Comment                  *bool
}

type DropAuthenticationPolicyRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowAuthenticationPolicyRequest struct {
}

type DescribeAuthenticationPolicyRequest struct {
	name AccountObjectIdentifier // required
}
