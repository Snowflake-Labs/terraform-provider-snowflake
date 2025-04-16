package sdk

import (
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/pelletier/go-toml/v2"
	"github.com/snowflakedb/gosnowflake"
	"github.com/youmark/pkcs8"
	"golang.org/x/crypto/ssh"
)

func DefaultConfig(verifyPermissions bool) *gosnowflake.Config {
	config, err := ProfileConfig("default", verifyPermissions)
	if err != nil || config == nil {
		log.Printf("[DEBUG] No Snowflake config file found, proceeding with empty config, err = %v", err)
		config = &gosnowflake.Config{}
	}
	return config
}

func ProfileConfig(profile string, verifyPermissions bool) (*gosnowflake.Config, error) {
	log.Printf("[DEBUG] Retrieving %s profile from a TOML file", profile)
	path, err := GetConfigFileName()
	if err != nil {
		return nil, err
	}

	configs, err := LoadConfigFile[LegacyConfigDTO](path, verifyPermissions)
	if err != nil {
		return nil, fmt.Errorf("could not load config file: %w", err)
	}

	if profile == "" {
		profile = "default"
	}
	var config *gosnowflake.Config
	if cfg, ok := configs[profile]; ok {
		log.Printf("[DEBUG] Loading config for profile: \"%s\"", profile)
		driverCfg, err := cfg.DriverConfig()
		if err != nil {
			return nil, fmt.Errorf("converting profile \"%s\" in file %s failed: %w", profile, path, err)
		}
		config = Pointer(driverCfg)
	}
	if config == nil {
		log.Printf("[DEBUG] No config found for profile: \"%s\"", profile)
		return nil, nil
	}

	// us-west-2 is Snowflake's default region, but if you actually specify that it won't trigger the default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if config.Region == "us-west-2" {
		config.Region = ""
	}

	return config, nil
}

func MergeConfig(baseConfig *gosnowflake.Config, mergeConfig *gosnowflake.Config) *gosnowflake.Config {
	if baseConfig == nil {
		return mergeConfig
	}
	if baseConfig.Account == "" {
		baseConfig.Account = mergeConfig.Account
	}
	if baseConfig.User == "" {
		baseConfig.User = mergeConfig.User
	}
	if baseConfig.Password == "" {
		baseConfig.Password = mergeConfig.Password
	}
	if baseConfig.Warehouse == "" {
		baseConfig.Warehouse = mergeConfig.Warehouse
	}
	if baseConfig.Role == "" {
		baseConfig.Role = mergeConfig.Role
	}
	if baseConfig.Region == "" {
		baseConfig.Region = mergeConfig.Region
	}
	if baseConfig.Host == "" {
		baseConfig.Host = mergeConfig.Host
	}
	if !configBoolSet(baseConfig.ValidateDefaultParameters) {
		baseConfig.ValidateDefaultParameters = mergeConfig.ValidateDefaultParameters
	}
	if mergedMap := collections.MergeMaps(mergeConfig.Params, baseConfig.Params); len(mergedMap) > 0 {
		baseConfig.Params = mergedMap
	}
	if baseConfig.ClientIP == nil {
		baseConfig.ClientIP = mergeConfig.ClientIP
	}
	if baseConfig.Protocol == "" {
		baseConfig.Protocol = mergeConfig.Protocol
	}
	if baseConfig.Host == "" {
		baseConfig.Host = mergeConfig.Host
	}
	if baseConfig.Port == 0 {
		baseConfig.Port = mergeConfig.Port
	}
	if baseConfig.Authenticator == gosnowflakeAuthTypeEmpty {
		baseConfig.Authenticator = mergeConfig.Authenticator
	}
	if baseConfig.Passcode == "" {
		baseConfig.Passcode = mergeConfig.Passcode
	}
	if !baseConfig.PasscodeInPassword {
		baseConfig.PasscodeInPassword = mergeConfig.PasscodeInPassword
	}
	if baseConfig.OktaURL == nil {
		baseConfig.OktaURL = mergeConfig.OktaURL
	}
	if baseConfig.LoginTimeout == 0 {
		baseConfig.LoginTimeout = mergeConfig.LoginTimeout
	}
	if baseConfig.RequestTimeout == 0 {
		baseConfig.RequestTimeout = mergeConfig.RequestTimeout
	}
	if baseConfig.JWTExpireTimeout == 0 {
		baseConfig.JWTExpireTimeout = mergeConfig.JWTExpireTimeout
	}
	if baseConfig.ClientTimeout == 0 {
		baseConfig.ClientTimeout = mergeConfig.ClientTimeout
	}
	if baseConfig.JWTClientTimeout == 0 {
		baseConfig.JWTClientTimeout = mergeConfig.JWTClientTimeout
	}
	if baseConfig.ExternalBrowserTimeout == 0 {
		baseConfig.ExternalBrowserTimeout = mergeConfig.ExternalBrowserTimeout
	}
	if baseConfig.MaxRetryCount == 0 {
		baseConfig.MaxRetryCount = mergeConfig.MaxRetryCount
	}
	if !baseConfig.InsecureMode { //nolint:staticcheck
		baseConfig.InsecureMode = mergeConfig.InsecureMode //nolint:staticcheck
	}
	if baseConfig.OCSPFailOpen == 0 {
		baseConfig.OCSPFailOpen = mergeConfig.OCSPFailOpen
	}
	if baseConfig.Token == "" {
		baseConfig.Token = mergeConfig.Token
	}
	if !baseConfig.KeepSessionAlive {
		baseConfig.KeepSessionAlive = mergeConfig.KeepSessionAlive
	}
	if baseConfig.PrivateKey == nil {
		baseConfig.PrivateKey = mergeConfig.PrivateKey
	}
	if !baseConfig.DisableTelemetry {
		baseConfig.DisableTelemetry = mergeConfig.DisableTelemetry
	}
	if baseConfig.Tracing == "" {
		baseConfig.Tracing = mergeConfig.Tracing
	}
	if baseConfig.TmpDirPath == "" {
		baseConfig.TmpDirPath = mergeConfig.TmpDirPath
	}
	if !configBoolSet(baseConfig.ClientRequestMfaToken) {
		baseConfig.ClientRequestMfaToken = mergeConfig.ClientRequestMfaToken
	}
	if !configBoolSet(baseConfig.ClientStoreTemporaryCredential) {
		baseConfig.ClientStoreTemporaryCredential = mergeConfig.ClientStoreTemporaryCredential
	}
	if !baseConfig.DisableQueryContextCache {
		baseConfig.DisableQueryContextCache = mergeConfig.DisableQueryContextCache
	}
	if !configBoolSet(baseConfig.IncludeRetryReason) {
		baseConfig.IncludeRetryReason = mergeConfig.IncludeRetryReason
	}
	if !configBoolSet(baseConfig.DisableConsoleLogin) {
		baseConfig.DisableConsoleLogin = mergeConfig.DisableConsoleLogin
	}
	return baseConfig
}

func configBoolSet(v gosnowflake.ConfigBool) bool {
	// configBoolNotSet is unexported in the driver, so we check if it's neither true nor false
	return slices.Contains([]gosnowflake.ConfigBool{gosnowflake.ConfigBoolFalse, gosnowflake.ConfigBoolTrue}, v)
}

func boolToConfigBool(v bool) gosnowflake.ConfigBool {
	if v {
		return gosnowflake.ConfigBoolTrue
	}
	return gosnowflake.ConfigBoolFalse
}

func GetConfigFileName() (string, error) {
	// has the user overridden the default config path?
	if configPath, ok := oswrapper.LookupEnv("SNOWFLAKE_CONFIG_PATH"); ok {
		if configPath != "" {
			return configPath, nil
		}
	}
	dir, err := oswrapper.UserHomeDir()
	if err != nil {
		return "", err
	}
	// default config path is ~/.snowflake/config.
	return filepath.Join(dir, ".snowflake", "config"), nil
}

func pointerAttributeSet[T any](src, dst *T) {
	if src != nil {
		*dst = *src
	}
}

func pointerTimeInSecondsAttributeSet(src *int, dst *time.Duration) {
	if src != nil {
		*dst = time.Second * time.Duration(*src)
	}
}

// TODO [SNOW-1827312]: fix this method
func pointerConfigBoolAttributeSet(src *bool, dst *gosnowflake.ConfigBool) {
	if src != nil {
		*dst = boolToConfigBool(*src)
	}
}

func pointerIpAttributeSet(src *string, dst *net.IP) {
	if src != nil {
		*dst = net.ParseIP(*src)
	}
}

func pointerUrlAttributeSet(src *string, dst **url.URL) error {
	if src != nil {
		url, err := url.Parse(*src)
		if err != nil {
			return err
		}
		*dst = url
	}
	return nil
}

func LoadConfigFile[T LegacyConfigDTO](path string, verifyPermissions bool) (map[string]T, error) {
	data, err := oswrapper.ReadFileSafe(path, verifyPermissions)
	if err != nil {
		return nil, err
	}
	var s map[string]T

	err = toml.Unmarshal(data, &s)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling config file %s: %w", path, err)
	}
	return s, nil
}

func ParsePrivateKey(privateKeyBytes []byte, passphrase []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("could not parse private key, key is not in PEM format")
	}

	if privateKeyBlock.Type == "ENCRYPTED PRIVATE KEY" {
		if len(passphrase) == 0 {
			return nil, fmt.Errorf("private key requires a passphrase, but private_key_passphrase was not supplied")
		}
		privateKey, err := pkcs8.ParsePKCS8PrivateKeyRSA(privateKeyBlock.Bytes, passphrase)
		if err != nil {
			return nil, fmt.Errorf("could not parse encrypted private key with passphrase, only ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc are supported err = %w", err)
		}
		return privateKey, nil
	}

	// TODO(SNOW-1754327): check if we can simply use ssh.ParseRawPrivateKeyWithPassphrase
	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key err = %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("privateKey not of type RSA")
	}
	return rsaPrivateKey, nil
}

type AuthenticationType string

const (
	AuthenticationTypeSnowflake           AuthenticationType = "SNOWFLAKE"
	AuthenticationTypeOauth               AuthenticationType = "OAUTH"
	AuthenticationTypeExternalBrowser     AuthenticationType = "EXTERNALBROWSER"
	AuthenticationTypeOkta                AuthenticationType = "OKTA"
	AuthenticationTypeJwt                 AuthenticationType = "SNOWFLAKE_JWT"
	AuthenticationTypeTokenAccessor       AuthenticationType = "TOKENACCESSOR"
	AuthenticationTypeUsernamePasswordMfa AuthenticationType = "USERNAMEPASSWORDMFA"

	AuthenticationTypeEmpty AuthenticationType = ""
)

var AllAuthenticationTypes = []AuthenticationType{
	AuthenticationTypeSnowflake,
	AuthenticationTypeOauth,
	AuthenticationTypeExternalBrowser,
	AuthenticationTypeOkta,
	AuthenticationTypeJwt,
	AuthenticationTypeTokenAccessor,
	AuthenticationTypeUsernamePasswordMfa,
}

func ToAuthenticatorType(s string) (gosnowflake.AuthType, error) {
	switch strings.ToUpper(s) {
	case string(AuthenticationTypeSnowflake):
		return gosnowflake.AuthTypeSnowflake, nil
	case string(AuthenticationTypeOauth):
		return gosnowflake.AuthTypeOAuth, nil
	case string(AuthenticationTypeExternalBrowser):
		return gosnowflake.AuthTypeExternalBrowser, nil
	case string(AuthenticationTypeOkta):
		return gosnowflake.AuthTypeOkta, nil
	case string(AuthenticationTypeJwt):
		return gosnowflake.AuthTypeJwt, nil
	case string(AuthenticationTypeTokenAccessor):
		return gosnowflake.AuthTypeTokenAccessor, nil
	case string(AuthenticationTypeUsernamePasswordMfa):
		return gosnowflake.AuthTypeUsernamePasswordMFA, nil
	default:
		return gosnowflake.AuthType(0), fmt.Errorf("invalid authenticator type: %s", s)
	}
}

const gosnowflakeAuthTypeEmpty = gosnowflake.AuthType(-1)

func ToExtendedAuthenticatorType(s string) (gosnowflake.AuthType, error) {
	switch strings.ToUpper(s) {
	case string(AuthenticationTypeEmpty):
		return gosnowflakeAuthTypeEmpty, nil
	default:
		return ToAuthenticatorType(s)
	}
}

type DriverLogLevel string

const (
	// these values are lower case on purpose to match gosnowflake case
	DriverLogLevelTrace   DriverLogLevel = "trace"
	DriverLogLevelDebug   DriverLogLevel = "debug"
	DriverLogLevelInfo    DriverLogLevel = "info"
	DriverLogLevelPrint   DriverLogLevel = "print"
	DriverLogLevelWarning DriverLogLevel = "warning"
	DriverLogLevelError   DriverLogLevel = "error"
	DriverLogLevelFatal   DriverLogLevel = "fatal"
	DriverLogLevelPanic   DriverLogLevel = "panic"
)

var AllDriverLogLevels = []DriverLogLevel{
	DriverLogLevelTrace,
	DriverLogLevelDebug,
	DriverLogLevelInfo,
	DriverLogLevelPrint,
	DriverLogLevelWarning,
	DriverLogLevelError,
	DriverLogLevelFatal,
	DriverLogLevelPanic,
}

func ToDriverLogLevel(s string) (DriverLogLevel, error) {
	lowerCase := strings.ToLower(s)
	switch lowerCase {
	case string(DriverLogLevelTrace),
		string(DriverLogLevelDebug),
		string(DriverLogLevelInfo),
		string(DriverLogLevelPrint),
		string(DriverLogLevelWarning),
		string(DriverLogLevelError),
		string(DriverLogLevelFatal),
		string(DriverLogLevelPanic):
		return DriverLogLevel(lowerCase), nil
	default:
		return "", fmt.Errorf("invalid driver log level: %s", s)
	}
}
