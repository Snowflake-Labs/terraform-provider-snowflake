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
	name                     SchemaObjectIdentifier // required
	AuthenticationMethods    []AuthenticationMethods
	MfaAuthenticationMethods []MfaAuthenticationMethods
	MfaEnrollment            *string
	ClientTypes              []ClientTypes
	SecurityIntegrations     []SecurityIntegrationsOption
	Comment                  *string
}

type AlterAuthenticationPolicyRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Set      *AuthenticationPolicySetRequest
	Unset    *AuthenticationPolicyUnsetRequest
	RenameTo *SchemaObjectIdentifier
}

type AuthenticationPolicySetRequest struct {
	AuthenticationMethods    []AuthenticationMethods
	MfaAuthenticationMethods []MfaAuthenticationMethods
	MfaEnrollment            *string
	ClientTypes              []ClientTypes
	SecurityIntegrations     []SecurityIntegrationsOption
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
	name     SchemaObjectIdentifier // required
}

type ShowAuthenticationPolicyRequest struct {
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeAuthenticationPolicyRequest struct {
	name SchemaObjectIdentifier // required
}
