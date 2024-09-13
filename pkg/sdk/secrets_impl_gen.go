package sdk

import (
	"context"
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

func (r *CreateWithOAuthClientCredentialsFlowSecretRequest) toOpts() *CreateWithOAuthClientCredentialsFlowSecretOptions {
	opts := &CreateWithOAuthClientCredentialsFlowSecretOptions{
		OrReplace:           r.OrReplace,
		IfNotExists:         r.IfNotExists,
		name:                r.name,
		SecurityIntegration: r.SecurityIntegration,
		OauthScopes:         r.OauthScopes,
		Comment:             r.Comment,
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
		SecurityIntegration:         r.SecurityIntegration,
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

			OauthRefreshToken:           r.Set.OauthRefreshToken,
			OauthRefreshTokenExpiryTime: r.Set.OauthRefreshTokenExpiryTime,
			Username:                    r.Set.Username,
			Password:                    r.Set.Password,
			SecretString:                r.Set.SecretString,
		}

		if r.Set.OAuthScopes != nil {
			opts.Set.OAuthScopes = &OAuthScopes{
				OAuthScopes: r.Set.OAuthScopes.OAuthScopes,
			}

		}

	}

	if r.Unset != nil {
		opts.Unset = &SecretUnset{
			UnsetComment: r.Unset.UnsetComment,
		}

	}

	return opts
}
