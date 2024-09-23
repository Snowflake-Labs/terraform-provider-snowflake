package sdk

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Secrets = (*secrets)(nil)

type secrets struct {
	client *Client
}

func (v *secrets) CreateWithOAuthClientCredentialsFlow(ctx context.Context, request *CreateWithOAuthClientCredentialsFlowSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) CreateWithOAuthAuthorizationCodeFlow(ctx context.Context, request *CreateWithOAuthAuthorizationCodeFlowSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) CreateWithBasicAuthentication(ctx context.Context, request *CreateWithBasicAuthenticationSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) CreateWithGenericString(ctx context.Context, request *CreateWithGenericStringSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) Alter(ctx context.Context, request *AlterSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) Drop(ctx context.Context, request *DropSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) Show(ctx context.Context, request *ShowSecretRequest) ([]Secret, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[secretDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[secretDBRow, Secret](dbRows)
	return resultList, nil
}

func (v *secrets) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Secret, error) {
	request := NewShowSecretRequest().WithIn(ExtendedIn{In: In{Schema: id.SchemaId()}}).WithLike(Like{String(id.Name())})
	secrets, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(secrets, func(r Secret) bool { return r.Name == id.Name() })
}

func (v *secrets) Describe(ctx context.Context, id SchemaObjectIdentifier) (*SecretDetails, error) {
	opts := &DescribeSecretOptions{
		name: id,
	}
	result, err := validateAndQueryOne[secretDetailsDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateWithOAuthClientCredentialsFlowSecretRequest) toOpts() *CreateWithOAuthClientCredentialsFlowSecretOptions {
	opts := &CreateWithOAuthClientCredentialsFlowSecretOptions{
		OrReplace:      r.OrReplace,
		IfNotExists:    r.IfNotExists,
		name:           r.name,
		ApiIntegration: r.ApiIntegration,
		OauthScopes:    r.OauthScopes,
		Comment:        r.Comment,
	}
	return opts
}

func (r *CreateWithOAuthAuthorizationCodeFlowSecretRequest) toOpts() *CreateWithOAuthAuthorizationCodeFlowSecretOptions {
	opts := &CreateWithOAuthAuthorizationCodeFlowSecretOptions{
		OrReplace:                   r.OrReplace,
		IfNotExists:                 r.IfNotExists,
		name:                        r.name,
		OauthRefreshToken:           r.OauthRefreshToken,
		OauthRefreshTokenExpiryTime: r.OauthRefreshTokenExpiryTime,
		ApiIntegration:              r.ApiIntegration,
		Comment:                     r.Comment,
	}
	return opts
}

func (r *CreateWithBasicAuthenticationSecretRequest) toOpts() *CreateWithBasicAuthenticationSecretOptions {
	opts := &CreateWithBasicAuthenticationSecretOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Username:    r.Username,
		Password:    r.Password,
		Comment:     r.Comment,
	}
	return opts
}

func (r *CreateWithGenericStringSecretRequest) toOpts() *CreateWithGenericStringSecretOptions {
	opts := &CreateWithGenericStringSecretOptions{
		OrReplace:    r.OrReplace,
		IfNotExists:  r.IfNotExists,
		name:         r.name,
		SecretString: r.SecretString,
		Comment:      r.Comment,
	}
	return opts
}

func (r *AlterSecretRequest) toOpts() *AlterSecretOptions {
	opts := &AlterSecretOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}

	if r.Set != nil {

		opts.Set = &SecretSet{
			Comment: r.Set.Comment,
		}

		if r.Set.SetForOAuthClientCredentialsFlow != nil {

			opts.Set.SetForOAuthClientCredentialsFlow = &SetForOAuthClientCredentialsFlow{
				OauthScopes: r.Set.SetForOAuthClientCredentialsFlow.OauthScopes,
			}

		}

		if r.Set.SetForOAuthAuthorizationFlow != nil {

			opts.Set.SetForOAuthAuthorizationFlow = &SetForOAuthAuthorizationFlow{
				OauthRefreshToken:           r.Set.SetForOAuthAuthorizationFlow.OauthRefreshToken,
				OauthRefreshTokenExpiryTime: r.Set.SetForOAuthAuthorizationFlow.OauthRefreshTokenExpiryTime,
			}

		}

		if r.Set.SetForBasicAuthentication != nil {

			opts.Set.SetForBasicAuthentication = &SetForBasicAuthentication{
				Username: r.Set.SetForBasicAuthentication.Username,
				Password: r.Set.SetForBasicAuthentication.Password,
			}

		}

		if r.Set.SetForGenericString != nil {

			opts.Set.SetForGenericString = &SetForGenericString{
				SecretString: r.Set.SetForGenericString.SecretString,
			}

		}

	}

	if r.Unset != nil {

		opts.Unset = &SecretUnset{
			Comment: r.Unset.Comment,
		}

	}

	return opts
}

func (r *DropSecretRequest) toOpts() *DropSecretOptions {
	opts := &DropSecretOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowSecretRequest) toOpts() *ShowSecretOptions {
	opts := &ShowSecretOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func getOauthScopes(scopesString string) []string {
	formatedScopes := make([]string, 0)
	scopesString = strings.TrimPrefix(scopesString, "[")
	scopesString = strings.TrimSuffix(scopesString, "]")
	for _, scope := range strings.Split(scopesString, ",") {
		formatedScopes = append(formatedScopes, strings.TrimSpace(scope))
	}
	return formatedScopes
}

func (r secretDBRow) convert() *Secret {
	s := &Secret{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		SchemaName:    r.SchemaName,
		DatabaseName:  r.DatabaseName,
		Owner:         r.Owner,
		SecretType:    r.SecretType,
		OwnerRoleType: r.OwnerRoleType,
	}
	if r.Comment.Valid {
		s.Comment = String(r.Comment.String)
	}
	if r.OauthScopes.Valid {
		s.OauthScopes = getOauthScopes(r.OauthScopes.String)
	}
	return s
}

func (r *DescribeSecretRequest) toOpts() *DescribeSecretOptions {
	opts := &DescribeSecretOptions{
		name: r.name,
	}
	return opts
}

func (r secretDetailsDBRow) convert() *SecretDetails {
	s := &SecretDetails{
		CreatedOn:                   r.CreatedOn,
		Name:                        r.Name,
		SchemaName:                  r.SchemaName,
		DatabaseName:                r.DatabaseName,
		Owner:                       r.Owner,
		SecretType:                  r.SecretType,
		OauthAccessTokenExpiryTime:  r.OauthAccessTokenExpiryTime,
		OauthRefreshTokenExpiryTime: r.OauthRefreshTokenExpiryTime,
	}
	if r.Username.Valid {
		s.Username = String(r.Username.String)
	}
	if r.Comment.Valid {
		s.Comment = String(r.Comment.String)
	}
	if r.OauthScopes.Valid {
		s.OauthScopes = getOauthScopes(r.OauthScopes.String)
	}
	if r.IntegrationName.Valid {
		s.IntegrationName = String(r.IntegrationName.String)
	}
	return s
}
