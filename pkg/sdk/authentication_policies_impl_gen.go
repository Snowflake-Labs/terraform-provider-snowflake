package sdk

import (
	"context"
)

var _ AuthenticationPolicies = (*authenticationPolicies)(nil)

type authenticationPolicies struct {
	client *Client
}

func (v *authenticationPolicies) Create(ctx context.Context, request *CreateAuthenticationPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *authenticationPolicies) Alter(ctx context.Context, request *AlterAuthenticationPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *authenticationPolicies) Drop(ctx context.Context, request *DropAuthenticationPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *authenticationPolicies) Show(ctx context.Context, request *ShowAuthenticationPolicyRequest) ([]AuthenticationPolicy, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showAuthenticationPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showAuthenticationPolicyDBRow, AuthenticationPolicy](dbRows)
	return resultList, nil
}

func (v *authenticationPolicies) Describe(ctx context.Context, id AccountObjectIdentifier) ([]AuthenticationPolicyDescription, error) {
	opts := &DescribeAuthenticationPolicyOptions{
		name: id,
	}
	rows, err := validateAndQuery[describeAuthenticationPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[describeAuthenticationPolicyDBRow, AuthenticationPolicyDescription](rows), nil
}

func (r *CreateAuthenticationPolicyRequest) toOpts() *CreateAuthenticationPolicyOptions {
	opts := &CreateAuthenticationPolicyOptions{
		OrReplace:                r.OrReplace,
		name:                     r.name,
		AuthenticationMethods:    r.AuthenticationMethods,
		MfaAuthenticationMethods: r.MfaAuthenticationMethods,
		MfaEnrollment:            r.MfaEnrollment,
		ClientTypes:              r.ClientTypes,
		SecurityIntegrations:     r.SecurityIntegrations,
		Comment:                  r.Comment,
	}
	return opts
}

func (r *AlterAuthenticationPolicyRequest) toOpts() *AlterAuthenticationPolicyOptions {
	opts := &AlterAuthenticationPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,

		RenameTo: r.RenameTo,
	}

	if r.Set != nil {

		opts.Set = &AuthenticationPolicySet{
			AuthenticationMethods:    r.Set.AuthenticationMethods,
			MfaAuthenticationMethods: r.Set.MfaAuthenticationMethods,
			MfaEnrollment:            r.Set.MfaEnrollment,
			ClientTypes:              r.Set.ClientTypes,
			SecurityIntegrations:     r.Set.SecurityIntegrations,
			Comment:                  r.Set.Comment,
		}

	}

	if r.Unset != nil {

		opts.Unset = &AuthenticationPolicyUnset{
			ClientTypes:              r.Unset.ClientTypes,
			AuthenticationMethods:    r.Unset.AuthenticationMethods,
			SecurityIntegrations:     r.Unset.SecurityIntegrations,
			MfaAuthenticationMethods: r.Unset.MfaAuthenticationMethods,
			MfaEnrollment:            r.Unset.MfaEnrollment,
			Comment:                  r.Unset.Comment,
		}

	}

	return opts
}

func (r *DropAuthenticationPolicyRequest) toOpts() *DropAuthenticationPolicyOptions {
	opts := &DropAuthenticationPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowAuthenticationPolicyRequest) toOpts() *ShowAuthenticationPolicyOptions {
	opts := &ShowAuthenticationPolicyOptions{}
	return opts
}

func (r showAuthenticationPolicyDBRow) convert() *AuthenticationPolicy {
	// TODO: Mapping
	return &AuthenticationPolicy{}
}

func (r *DescribeAuthenticationPolicyRequest) toOpts() *DescribeAuthenticationPolicyOptions {
	opts := &DescribeAuthenticationPolicyOptions{
		name: r.name,
	}
	return opts
}

func (r describeAuthenticationPolicyDBRow) convert() *AuthenticationPolicyDescription {
	// TODO: Mapping
	return &AuthenticationPolicyDescription{}
}
