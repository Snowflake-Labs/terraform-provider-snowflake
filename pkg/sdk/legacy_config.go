package sdk

import (
	"fmt"

	"github.com/snowflakedb/gosnowflake"
)

// TODO(SNOW-1787920): improve TOML parsing
type LegacyConfigDTO struct {
	AccountName            *string             `toml:"accountname"`
	OrganizationName       *string             `toml:"organizationname"`
	User                   *string             `toml:"user"`
	Username               *string             `toml:"username"`
	Password               *string             `toml:"password"`
	Host                   *string             `toml:"host"`
	Warehouse              *string             `toml:"warehouse"`
	Role                   *string             `toml:"role"`
	Params                 *map[string]*string `toml:"params"`
	ClientIp               *string             `toml:"clientip"`
	Protocol               *string             `toml:"protocol"`
	Passcode               *string             `toml:"passcode"`
	Port                   *int                `toml:"port"`
	PasscodeInPassword     *bool               `toml:"passcodeinpassword"`
	OktaUrl                *string             `toml:"oktaurl"`
	ClientTimeout          *int                `toml:"clienttimeout"`
	JwtClientTimeout       *int                `toml:"jwtclienttimeout"`
	LoginTimeout           *int                `toml:"logintimeout"`
	RequestTimeout         *int                `toml:"requesttimeout"`
	JwtExpireTimeout       *int                `toml:"jwtexpiretimeout"`
	ExternalBrowserTimeout *int                `toml:"externalbrowsertimeout"`
	MaxRetryCount          *int                `toml:"maxretrycount"`
	Authenticator          *string             `toml:"authenticator"`
	InsecureMode           *bool               `toml:"insecuremode"`
	OcspFailOpen           *bool               `toml:"ocspfailopen"`
	Token                  *string             `toml:"token"`
	KeepSessionAlive       *bool               `toml:"keepsessionalive"`
	PrivateKey             *string             `toml:"privatekey,multiline"`
	PrivateKeyPassphrase   *string             `toml:"privatekeypassphrase"`
	DisableTelemetry       *bool               `toml:"disabletelemetry"`
	// TODO [SNOW-1827312]: handle and test 3-value booleans properly from TOML
	ValidateDefaultParameters      *bool   `toml:"validatedefaultparameters"`
	ClientRequestMfaToken          *bool   `toml:"clientrequestmfatoken"`
	ClientStoreTemporaryCredential *bool   `toml:"clientstoretemporarycredential"`
	Tracing                        *string `toml:"tracing"`
	TmpDirPath                     *string `toml:"tmpdirpath"`
	DisableQueryContextCache       *bool   `toml:"disablequerycontextcache"`
	IncludeRetryReason             *bool   `toml:"includeretryreason"`
	DisableConsoleLogin            *bool   `toml:"disableconsolelogin"`
}

func (c *LegacyConfigDTO) DriverConfig() (gosnowflake.Config, error) {
	driverCfg := gosnowflake.Config{}
	if c.AccountName != nil && c.OrganizationName != nil {
		driverCfg.Account = fmt.Sprintf("%s-%s", *c.OrganizationName, *c.AccountName)
	}
	pointerAttributeSet(c.User, &driverCfg.User)
	pointerAttributeSet(c.Username, &driverCfg.User)
	pointerAttributeSet(c.Password, &driverCfg.Password)
	pointerAttributeSet(c.Host, &driverCfg.Host)
	pointerAttributeSet(c.Warehouse, &driverCfg.Warehouse)
	pointerAttributeSet(c.Role, &driverCfg.Role)
	pointerAttributeSet(c.Params, &driverCfg.Params)
	pointerIpAttributeSet(c.ClientIp, &driverCfg.ClientIP)
	pointerAttributeSet(c.Protocol, &driverCfg.Protocol)
	pointerAttributeSet(c.Passcode, &driverCfg.Passcode)
	pointerAttributeSet(c.Port, &driverCfg.Port)
	pointerAttributeSet(c.PasscodeInPassword, &driverCfg.PasscodeInPassword)
	err := pointerUrlAttributeSet(c.OktaUrl, &driverCfg.OktaURL)
	if err != nil {
		return gosnowflake.Config{}, err
	}
	pointerTimeInSecondsAttributeSet(c.ClientTimeout, &driverCfg.ClientTimeout)
	pointerTimeInSecondsAttributeSet(c.JwtClientTimeout, &driverCfg.JWTClientTimeout)
	pointerTimeInSecondsAttributeSet(c.LoginTimeout, &driverCfg.LoginTimeout)
	pointerTimeInSecondsAttributeSet(c.RequestTimeout, &driverCfg.RequestTimeout)
	pointerTimeInSecondsAttributeSet(c.JwtExpireTimeout, &driverCfg.JWTExpireTimeout)
	pointerTimeInSecondsAttributeSet(c.ExternalBrowserTimeout, &driverCfg.ExternalBrowserTimeout)
	pointerAttributeSet(c.MaxRetryCount, &driverCfg.MaxRetryCount)
	if c.Authenticator != nil {
		authenticator, err := ToAuthenticatorType(*c.Authenticator)
		if err != nil {
			return gosnowflake.Config{}, err
		}
		driverCfg.Authenticator = authenticator
	}
	pointerAttributeSet(c.InsecureMode, &driverCfg.InsecureMode) //nolint:staticcheck
	if c.OcspFailOpen != nil {
		if *c.OcspFailOpen {
			driverCfg.OCSPFailOpen = gosnowflake.OCSPFailOpenTrue
		} else {
			driverCfg.OCSPFailOpen = gosnowflake.OCSPFailOpenFalse
		}
	}
	pointerAttributeSet(c.Token, &driverCfg.Token)
	pointerAttributeSet(c.KeepSessionAlive, &driverCfg.KeepSessionAlive)
	if c.PrivateKey != nil {
		passphrase := make([]byte, 0)
		if c.PrivateKeyPassphrase != nil {
			passphrase = []byte(*c.PrivateKeyPassphrase)
		}
		privKey, err := ParsePrivateKey([]byte(*c.PrivateKey), passphrase)
		if err != nil {
			return gosnowflake.Config{}, err
		}
		driverCfg.PrivateKey = privKey
	}
	pointerAttributeSet(c.DisableTelemetry, &driverCfg.DisableTelemetry)
	pointerConfigBoolAttributeSet(c.ValidateDefaultParameters, &driverCfg.ValidateDefaultParameters)
	pointerConfigBoolAttributeSet(c.ClientRequestMfaToken, &driverCfg.ClientRequestMfaToken)
	pointerConfigBoolAttributeSet(c.ClientStoreTemporaryCredential, &driverCfg.ClientStoreTemporaryCredential)
	pointerAttributeSet(c.Tracing, &driverCfg.Tracing)
	pointerAttributeSet(c.TmpDirPath, &driverCfg.TmpDirPath)
	pointerAttributeSet(c.DisableQueryContextCache, &driverCfg.DisableQueryContextCache)
	pointerConfigBoolAttributeSet(c.IncludeRetryReason, &driverCfg.IncludeRetryReason)
	pointerConfigBoolAttributeSet(c.DisableConsoleLogin, &driverCfg.DisableConsoleLogin)

	return driverCfg, nil
}
