package provider

import (
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

func Test_Provider_toAuthenticationType(t *testing.T) {
	type test struct {
		input string
		want  gosnowflake.AuthType
	}

	valid := []test{
		// Case insensitive.
		{input: "snowflake", want: gosnowflake.AuthTypeSnowflake},

		// Supported Values.
		{input: "SNOWFLAKE", want: gosnowflake.AuthTypeSnowflake},
		{input: "OAUTH", want: gosnowflake.AuthTypeOAuth},
		{input: "EXTERNALBROWSER", want: gosnowflake.AuthTypeExternalBrowser},
		{input: "OKTA", want: gosnowflake.AuthTypeOkta},
		{input: "JWT", want: gosnowflake.AuthTypeJwt},
		{input: "SNOWFLAKE_JWT", want: gosnowflake.AuthTypeJwt},
		{input: "TOKENACCESSOR", want: gosnowflake.AuthTypeTokenAccessor},
		{input: "USERNAMEPASSWORDMFA", want: gosnowflake.AuthTypeUsernamePasswordMFA},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := toAuthenticatorType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := toAuthenticatorType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Provider_toProtocol(t *testing.T) {
	type test struct {
		input string
		want  protocol
	}

	valid := []test{
		// Case insensitive.
		{input: "http", want: protocolHttp},

		// Supported Values.
		{input: "HTTP", want: protocolHttp},
		{input: "HTTPS", want: protocolHttps},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := toProtocol(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := toProtocol(tc.input)
			require.Error(t, err)
		})
	}
}
