package sdk

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/youmark/pkcs8"
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

	m, err := loadConfigFile(configPath)
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

func TestProfileConfig(t *testing.T) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	unencryptedDer, err := x509.MarshalPKCS8PrivateKey(rsaPrivateKey)
	require.NoError(t, err)
	privBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: unencryptedDer,
	}
	unencryptedKey := string(pem.EncodeToMemory(&privBlock))

	encryptedDer, err := pkcs8.MarshalPrivateKey(rsaPrivateKey, []byte("password"), &pkcs8.Opts{
		Cipher: pkcs8.AES256CBC,
		KDFOpts: pkcs8.PBKDF2Opts{
			SaltSize: 16, IterationCount: 2000, HMACHash: crypto.SHA256,
		},
	})
	require.NoError(t, err)
	privEncryptedBlock := pem.Block{
		Type:  "ENCRYPTED PRIVATE KEY",
		Bytes: encryptedDer,
	}
	encryptedKey := string(pem.EncodeToMemory(&privEncryptedBlock))

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

	[securityadmin.params]
	foo = 'bar'
	`, encryptedKey, "password")
	configPath := testhelpers.TestFile(t, "config", []byte(c))

	t.Run("with found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("securityadmin")
		require.NoError(t, err)
		require.NotNil(t, config.PrivateKey)

		gotKey, err := x509.MarshalPKCS8PrivateKey(config.PrivateKey)
		require.NoError(t, err)
		gotUnencryptedKey := pem.EncodeToMemory(
			&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: gotKey,
			},
		)

		assert.Equal(t, "organizationname-accountname", config.Account)
		assert.Equal(t, "user", config.User)
		assert.Equal(t, "password", config.Password)
		assert.Equal(t, "warehouse", config.Warehouse)
		assert.Equal(t, "role", config.Role)
		assert.Equal(t, map[string]*string{"foo": Pointer("bar")}, config.Params)
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
		assert.Equal(t, "token", config.Token)
		assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
		assert.Equal(t, true, config.KeepSessionAlive)
		assert.Equal(t, unencryptedKey, string(gotUnencryptedKey))
		assert.Equal(t, true, config.DisableTelemetry)
		assert.Equal(t, "tracing", config.Tracing)
		assert.Equal(t, ".", config.TmpDirPath)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
		assert.Equal(t, true, config.DisableQueryContextCache)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
	})

	t.Run("with not found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		_, err := ProfileConfig("orgadmin")
		require.ErrorContains(t, err, "profile \"orgadmin\" not found in file")
	})

	t.Run("with not found config", func(t *testing.T) {
		name := random.AlphaN(8)
		t.Setenv(snowflakeenvs.ConfigPath, name)

		_, err = ProfileConfig("orgadmin")
		require.ErrorContains(t, err, fmt.Sprintf("open %s: no such file or directory", name))
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
		IncludeRetryReason:             gosnowflake.ConfigBoolTrue,
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
