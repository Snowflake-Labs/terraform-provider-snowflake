package provider_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	_ "github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	r := require.New(t)
	err := provider.Provider().InternalValidate()
	r.NoError(err)
}

// func TestConfigureProvider(t *testing.T) {
// 	// r := require.New(t)
// }

func TestDSN(t *testing.T) {
	type args struct {
		account,
		user,
		password string
		browserAuth bool
		region,
		role string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"simple", args{"acct", "user", "pass", false, "region", "role"},
			"user:pass@acct.region.snowflakecomputing.com:443?ocspFailOpen=true&region=region&role=role&validateDefaultParameters=true", false},
		{"us-west-2 special case", args{"acct2", "user2", "pass2", false, "us-west-2", "role2"},
			"user2:pass2@acct2.snowflakecomputing.com:443?ocspFailOpen=true&role=role2&validateDefaultParameters=true", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.DSN(tt.args.account, tt.args.user, tt.args.password, tt.args.browserAuth, "", "", tt.args.region, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("DSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func resourceData(t *testing.T, account, username, token, region, role string) *schema.ResourceData {
	r := require.New(t)

	in := map[string]interface{}{
		"account":  account,
		"username": username,
		"password": token,
		"region":   region,
		"role":     role,
	}

	d := schema.TestResourceDataRaw(t, provider.Provider().Schema, in)
	r.NotNil(d)
	return d
}

func TestOAuthDSN(t *testing.T) {
	type args struct {
		account,
		user,
		oauthAccessToken,
		region,
		role string
	}
	pseudorandom_access_token := "ETMsjLOLvQ-C/bzGmmdvbEM/RSQFFX-a+sefbQeQoJqwdFNXZ+ftBIdwlasApA+/MItZLNRRW-rYJiEZMvAAdzpGLxaghIoww+vDOuIeAFBDUxTAY-I+qGbQOXipkNcmzwuAaugjYtlTjPXGjqKw-OSsVacQXzsQyAMnbMyUrbdhRQEETIqTAdMuDqJBeaSj+LMsKDXzLd-guSlm-mmv+="
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"simple_oauth", args{"acct", "user", pseudorandom_access_token, "region", "role"},
			"user:@acct.region.snowflakecomputing.com:443?authenticator=oauth&ocspFailOpen=true&region=region&role=role&token=ETMsjLOLvQ-C%2FbzGmmdvbEM%2FRSQFFX-a%2BsefbQeQoJqwdFNXZ%2BftBIdwlasApA%2B%2FMItZLNRRW-rYJiEZMvAAdzpGLxaghIoww%2BvDOuIeAFBDUxTAY-I%2BqGbQOXipkNcmzwuAaugjYtlTjPXGjqKw-OSsVacQXzsQyAMnbMyUrbdhRQEETIqTAdMuDqJBeaSj%2BLMsKDXzLd-guSlm-mmv%2B%3D&validateDefaultParameters=true", false},
		{"oauth_over_password", args{"acct", "user", pseudorandom_access_token, "region", "role"},
			"user:@acct.region.snowflakecomputing.com:443?authenticator=oauth&ocspFailOpen=true&region=region&role=role&token=ETMsjLOLvQ-C%2FbzGmmdvbEM%2FRSQFFX-a%2BsefbQeQoJqwdFNXZ%2BftBIdwlasApA%2B%2FMItZLNRRW-rYJiEZMvAAdzpGLxaghIoww%2BvDOuIeAFBDUxTAY-I%2BqGbQOXipkNcmzwuAaugjYtlTjPXGjqKw-OSsVacQXzsQyAMnbMyUrbdhRQEETIqTAdMuDqJBeaSj%2BLMsKDXzLd-guSlm-mmv%2B%3D&validateDefaultParameters=true", false},
		{"empty_token_no_password_errors_out", args{"acct", "user", "", "region", "role"},
			"", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.DSN(tt.args.account, tt.args.user, "", false, "", tt.args.oauthAccessToken, tt.args.region, tt.args.role)

			if (err != nil) != tt.wantErr {
				t.Errorf("DSN() error = %v, dsn = %v, wantErr %v", err, got, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
