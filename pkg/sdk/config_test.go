package sdk

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFile(t *testing.T) {
	c := `
	[default]
	account='TEST_ACCOUNT'
	user='TEST_USER'
	password='abcd1234'
	role='ACCOUNTADMIN'

	[securityadmin]
	account='TEST_ACCOUNT'
	user='TEST_USER'
	password='abcd1234'
	role='SECURITYADMIN'
	`
	configPath := testhelpers.TestFile(t, "config", []byte(c))
	t.Setenv(snowflakeenvs.ConfigPath, configPath)

	m, err := loadConfigFile()
	require.NoError(t, err)
	assert.Equal(t, "TEST_ACCOUNT", *m["default"].Account)
	assert.Equal(t, "TEST_USER", *m["default"].User)
	assert.Equal(t, "abcd1234", *m["default"].Password)
	assert.Equal(t, "ACCOUNTADMIN", *m["default"].Role)
	assert.Equal(t, "TEST_ACCOUNT", *m["securityadmin"].Account)
	assert.Equal(t, "TEST_USER", *m["securityadmin"].User)
	assert.Equal(t, "abcd1234", *m["securityadmin"].Password)
	assert.Equal(t, "SECURITYADMIN", *m["securityadmin"].Role)
}

// These keys were generated with the following commands:
// openssl genrsa -aes128 -passout pass:password -out private.key 2048
// openssl rsa -in private.key -out decrypted_private.key
// <enter password>
// TODO(): generate keys dynamically using Go libraries
const privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDyXQvU2Bffntwm
2RMntwBagl8Q/yco8L+Gotd7jq1JldamvxJbpO2tj7Iq/nhAA9N4gDQ3AHbEmWgV
f/+xYKuWs9Itp5OAU0fLRxPkUIsapX8PO67NlTIQLXsp7Jg60R36r6aSD68pGie1
NWSfqw5EB6eezn/7iBjJQfnvZ5TdTyHm+L0+qsAj5PhWSvgLyzGXgXw5Cz0QwbmD
bkt5gkPIvjzgPc/1V8kclknyc7mU5Eu2F06xvK+MQHvR9usLZBG2Q7BD2jDPBHGe
M/aVsUA1v9F8bw9CdPQEa+7TdHmsdvDi8m3OhYC2RNBoqpZ+yy+A6JiWMYBCK82e
Itg00SdjAgMBAAECggEAAeuHUzKJD2u7kejo4IQDvm7k+OjEkLDP4Is3VH7VNoia
FXZTDKNaMtXDBGMtxzz5IwZdma9VS2K8Mu8Y7dF2qU3JAXAf7CMkJeMUdSWzlkLM
/yg+Yp6iFBeT1xhNKiCVSi3bNgDIzlerICum2LEhL9dSm0f356tXJ9nv6pvjWUrb
Nch5bm084m8As1La99HPLOhr4zN6YONAnb2SkuTy/TJJ3BCWPci3obrtGePSbnf5
3DSqfRF/BBcq1yPMyHKvmdEOFBzSNtH8PKT//9Yg5IpDj3WWOGlWTv2eQIriRmQx
Ha36RbzXMn3BSPPkVSGPbDGrp3WBLTVTx6ST5fnIyQKBgQD82qCdNwCEK8Aya9wL
303BK4KFnnAfNT91pS7sgnrFvTKAnCtYbofNUsQwtzicKSLJ6rOLPsxfvr1hwVZE
avUgfTXSjWuJ3kQq9MoXKtCb2XKDcUk+HrN5l9XObZLfB5mpQgBZmXGbMaHujcA9
1gP/bJ0k4whOyQDHloxjSH83mwKBgQD1YQFLmhmqydTUPlxr1N7haxNpttC6RBIb
oRbXQNihzAeyYccab/FdSoWyzPjYvYm0EwFWYhZBetnSb9D9vZP0opvdtLLzzySI
CzW053eePyvb63FNJgWpEQsrL3GAOUIIrGkuC2BRL5XtF+mUcU2jPQDTgDhfrefz
0dFvPtHf2QKBgQDuSBmUDoEuDQzSd1Km3YkowRf/U4/V2Rg0hbXyrAOG1QUCrikq
7P6NP7IjNobiouFl5wfL8SIoGFfgB5KEZ0cZluVhxmPRSOR0lrrbmj18oS6JL/kV
0VjQ/YU/Q4NlKoRkPQ6XYULuPZecd3jyzPx3eKOeX1U06bcSX41tAqTggQKBgQCB
6AFPjR3ZlVDfrMQxMlls7csxRF/svOz5Q6db/jCyN9o7ThiinnEh+rodlvaHiJDG
jOlAWl19/RQknJ4AN8WE1jG+hlPXT+r/OzALvh9N4BPQMi2hsmd8wlEvY8arI6Ua
Am0Mu2kakh7FjstSk0mPClTNpCw0O1V5d7NxOcjSwQKBgF6HPWRTkE/no2NNwLle
epeaH9iCeFMKZWsSF1t4T09i4H+izBmTQqgJKm202FXSOEngHyiGdXZNzcujzQLb
V+mI/CeDv1PZTa2ju5D/fniqgeRRgfixqR0wq7sDqJIlpBzHnQnANBoJW2t4pirE
HPOlkuMp4rxcPXwQgV5LMrao
-----END PRIVATE KEY-----
`

const encryptedPrivateKey = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIFNTBfBgkqhkiG9w0BBQ0wUjAxBgkqhkiG9w0BBQwwJAQQ7R0F4yA4ir5buYp7
UzK+vAICCAAwDAYIKoZIhvcNAgkFADAdBglghkgBZQMEAQIEEIaGrNdggYwWopje
YurbXlcEggTQpb/QNPnjHiZ5WNeM/bnYjK0/W0oKBiCmDz7ZgyHypIyAqfOTesh5
rEbI7+EnkEOqCE4OwpGWZgkJsbrxbmbv1eGm1fTZXxZBONpmeuv72FRL1NZ2kTUJ
JV84tc7uyYbJcKnfPjcqjDOxxK5gVkArm2uaVLSZkXpCjL2LWgurF4ajZfotOiDo
ziYxDnLXzkhhdNu/itTJ0Qoo5/Eo91QO6+zZWx6/T1mNtDW4gUXN0T8ev5FYpFdR
aJVetwMLJYprbsbatHMjazsdUoMQVt74e4pDfkZROidLgfSP2ud3ZgTp50uIbyfG
h9mTAWAv1Cv8KCvc2BbCQc3B6IJdiA9oO6P1H9bmSlQE7UYwwqJ4TVN+RCc2Uosi
QshRpeAoSOHGqQw0LeXbr6wFYsYbbjTTlMEK9dor3vLDXfsuSBzf9rVyMpLeL+s2
GCag90Bd+MjHJCz+hQHEXlDtbSLNEp8oIOj+Y0FodDKzfMBrKghUsIiE64rVVLOQ
SwUrQTYWsbu9O8bmcKthzedz5ZCJc6JujR1jgnTeLsLNTjOazHiZEoNJAUyqds3p
StpTlRFBAy+UGYqSuX3aNkVL2hXQPzZ4Xe3QMGrCyQfzLxH0QaSYdCdH37k5EzEb
MbaHMt9ktt0aGJfAAKzqagcvOnwgvq8lQHPMJZcyKCpUCrf1yKxpIutVD8l0l+70
rtMr3McALKOfhKqokiZBR6VCL1l78Ifu//qhzsCXW3HLcdwfBSwHwoAwCkip9tFq
7ZjfMZ4/6466kvt1ZwrrjClkfS5Kaz0ZuzFdUSpUoE4sXzgvL6i/HvnBCwEA476E
yTJKu9QmGFE2Jz2PVuhhRyGmDy4Tpg68wTBhZZLbf+190G0eSvXdhus/IZfdh5bH
pkbcs7ApsR9I+2HwPqUpNzclktlblsXPpwzBsFvE48qYs9ybSej1wamzx3DmDmuJ
surY5eF5lZJtKEdaFASORZbtgwmr5OnycYK0Qrzm/P5lTwlZPaJSN8sJCsAV/++w
r76vxmVUaIvhaIs2DR3fp8FTXz1+obulivP6j7qNGcW72Yk1Ssx9wuroS/4PZ/T3
x6hgWqdeaLCDK9pBibf3R/7wZLr3UygkUrA95bHzLaEsKHA7h1qwPaWO6a/255Nx
mJBeDruhULaCMjrcMK5SHf/Iwo6wRbDSbb8uQkpdH9Si5lDaBtf0FtNVCefzu+B8
MFVQ9VrQPLv5h0AZnKNynLe0JwdB+mdadKCebF+2FhB5X7h8IfyhAhsQ+vST61P5
meRXgyuIukr3BjAWenuZkYHc7BSxwTSfWcxptDW9BzM3aEA+wP2HR+dzbv+mZrw/
ABUOA1WweCTKrKMnvOaJPuXoZMe8BfN4YvAHYRn14f4M/cU2D41iZQV0VAq0jnaZ
Oye8TLd5QgqVXIoXKsSTzcVhs90ga3c7UJWixUK8K5d5tspbY3JCfHY24P8t6vlt
MQOU9wX+FwliEXEGuntPwKksyzOYD0P+olzfqc3U+xI+60jCK6EoPHlJE7BeDBkI
L63OUALjyIijmvsKJ6NnQ8RzPqR6qrLocTQA/dEE1RUy7K8RQduWjHmASSe8tiCT
EGmrLRxfSjAOuO14x/WbIh+p88t6u4ewNid1LrQNWbE+xFaEQqqDpmo=
-----END ENCRYPTED PRIVATE KEY-----
`

func TestProfileConfig(t *testing.T) {
	c := fmt.Sprintf(`
	[securityadmin]
	account='account'
	accountname='accountname'
	organizationname='organizationname'
	user='user'
	password='password'
	host='host'
	warehouse='warehouse'
	role='role'
	clientip='1.1.1.1'
	protocol='http'
	passcode='passcode'
	port=1
	passcodeinpassword=true
	oktaurl='https://example.com'
	clienttimeout=10
	jwtclienttimeout=20
	logintimeout=30
	requesttimeout=40
	jwtexpiretimeout=50
	externalbrowsertimeout=60
	maxretrycount=1
	authenticator='jwt'
	insecuremode=true
	ocspfailopen=true
	token='token'
	keepsessionalive=true
	privatekey="""%s"""
	privatekeypassphrase='%s'
	disabletelemetry=true
	validatedefaultparameters=true
	clientrequestmfatoken=true
	clientstoretemporarycredential=true
	tracing='tracing'
	tmpdirpath='.'
	disablequerycontextcache=true
	includeretryreason=true
	disableconsolelogin=true
	`, encryptedPrivateKey, "password")
	configPath := testhelpers.TestFile(t, "config", []byte(c))

	t.Run("with found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("securityadmin")
		require.NoError(t, err)

		key, err := x509.MarshalPKCS8PrivateKey(config.PrivateKey)
		require.NoError(t, err)
		pemdata := pem.EncodeToMemory(
			&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: key,
			},
		)

		assert.Equal(t, "organizationname-accountname", config.Account)
		assert.Equal(t, "user", config.User)
		assert.Equal(t, "password", config.Password)
		assert.Equal(t, "warehouse", config.Warehouse)
		assert.Equal(t, "role", config.Role)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
		assert.Equal(t, "1.1.1.1", config.ClientIP.String())
		assert.Equal(t, "http", config.Protocol)
		assert.Equal(t, "host", config.Host)
		assert.Equal(t, 1, config.Port)
		assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
		assert.Equal(t, "passcode", config.Passcode)
		assert.Equal(t, true, config.PasscodeInPassword)
		assert.Equal(t, "https://example.com", config.OktaURL.String())
		assert.Equal(t, 10*time.Second, config.ClientTimeout)
		assert.Equal(t, 20*time.Second, config.JWTClientTimeout)
		assert.Equal(t, 30*time.Second, config.LoginTimeout)
		assert.Equal(t, 40*time.Second, config.RequestTimeout)
		assert.Equal(t, 50*time.Second, config.JWTExpireTimeout)
		assert.Equal(t, 60*time.Second, config.ExternalBrowserTimeout)
		assert.Equal(t, 1, config.MaxRetryCount)
		assert.Equal(t, true, config.InsecureMode)
		assert.Equal(t, true, config.InsecureMode)
		assert.Equal(t, "token", config.Token)
		assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
		assert.Equal(t, true, config.KeepSessionAlive)
		assert.Equal(t, privateKey, string(pemdata))
		assert.Equal(t, true, config.DisableTelemetry)
		assert.Equal(t, "tracing", config.Tracing)
		assert.Equal(t, ".", config.TmpDirPath)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
		assert.Equal(t, true, config.DisableQueryContextCache)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
	})

	t.Run("with not found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("orgadmin")
		require.NoError(t, err)
		require.Nil(t, config)
	})

	t.Run("with not found config", func(t *testing.T) {
		dir, err := os.UserHomeDir()
		require.NoError(t, err)
		t.Setenv(snowflakeenvs.ConfigPath, dir)

		config, err := ProfileConfig("orgadmin")
		require.Error(t, err)
		require.Nil(t, config)
	})
}

func Test_MergeConfig(t *testing.T) {
	oktaUrl1, err := url.Parse("https://example1.com")
	require.NoError(t, err)
	oktaUrl2, err := url.Parse("https://example2.com")
	require.NoError(t, err)

	config1 := &gosnowflake.Config{
		Account:                   "account1",
		User:                      "user1",
		Password:                  "password1",
		Warehouse:                 "warehouse1",
		Role:                      "role1",
		ValidateDefaultParameters: 1,
		Params: map[string]*string{
			"foo": Pointer("1"),
		},
		ClientIP:                       net.ParseIP("1.1.1.1"),
		Protocol:                       "protocol1",
		Host:                           "host1",
		Port:                           1,
		Authenticator:                  1,
		Passcode:                       "passcode1",
		PasscodeInPassword:             false,
		OktaURL:                        oktaUrl1,
		LoginTimeout:                   1,
		RequestTimeout:                 1,
		JWTExpireTimeout:               1,
		ClientTimeout:                  1,
		JWTClientTimeout:               1,
		ExternalBrowserTimeout:         1,
		MaxRetryCount:                  1,
		InsecureMode:                   false,
		OCSPFailOpen:                   1,
		Token:                          "token1",
		KeepSessionAlive:               false,
		PrivateKey:                     random.GenerateRSAPrivateKey(t),
		DisableTelemetry:               false,
		Tracing:                        "tracing1",
		TmpDirPath:                     "tmpdirpath1",
		ClientRequestMfaToken:          gosnowflake.ConfigBoolFalse,
		ClientStoreTemporaryCredential: gosnowflake.ConfigBoolFalse,
		DisableQueryContextCache:       false,
		IncludeRetryReason:             1,
		DisableConsoleLogin:            gosnowflake.ConfigBoolFalse,
	}

	config2 := &gosnowflake.Config{
		Account:                   "account2",
		User:                      "user2",
		Password:                  "password2",
		Warehouse:                 "warehouse2",
		Role:                      "role2",
		ValidateDefaultParameters: 1,
		Params: map[string]*string{
			"foo": Pointer("2"),
		},
		ClientIP:                       net.ParseIP("2.2.2.2"),
		Protocol:                       "protocol2",
		Host:                           "host2",
		Port:                           2,
		Authenticator:                  2,
		Passcode:                       "passcode2",
		PasscodeInPassword:             true,
		OktaURL:                        oktaUrl2,
		LoginTimeout:                   2,
		RequestTimeout:                 2,
		JWTExpireTimeout:               2,
		ClientTimeout:                  2,
		JWTClientTimeout:               2,
		ExternalBrowserTimeout:         2,
		MaxRetryCount:                  2,
		InsecureMode:                   true,
		OCSPFailOpen:                   2,
		Token:                          "token2",
		KeepSessionAlive:               true,
		PrivateKey:                     random.GenerateRSAPrivateKey(t),
		DisableTelemetry:               true,
		Tracing:                        "tracing2",
		TmpDirPath:                     "tmpdirpath2",
		ClientRequestMfaToken:          gosnowflake.ConfigBoolTrue,
		ClientStoreTemporaryCredential: gosnowflake.ConfigBoolTrue,
		DisableQueryContextCache:       true,
		IncludeRetryReason:             2,
		DisableConsoleLogin:            gosnowflake.ConfigBoolTrue,
	}

	t.Run("base config empty", func(t *testing.T) {
		config := MergeConfig(&gosnowflake.Config{}, config1)

		require.Equal(t, config1, config)
	})

	t.Run("merge config empty", func(t *testing.T) {
		config := MergeConfig(config1, &gosnowflake.Config{})

		require.Equal(t, config1, config)
	})

	t.Run("both configs filled - base config takes precedence", func(t *testing.T) {
		config := MergeConfig(config1, config2)
		require.Equal(t, config1, config)
	})
}

func Test_toAuthenticationType(t *testing.T) {
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
			got, err := ToAuthenticatorType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToAuthenticatorType(tc.input)
			require.Error(t, err)
		})
	}
}
