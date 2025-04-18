package sdk

import (
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
	DriverTracing                  *string `toml:"tracing"`
	TmpDirPath                     *string `toml:"tmpdirpath"`
	DisableQueryContextCache       *bool   `toml:"disablequerycontextcache"`
	IncludeRetryReason             *bool   `toml:"includeretryreason"`
	DisableConsoleLogin            *bool   `toml:"disableconsolelogin"`
}

func (c *LegacyConfigDTO) DriverConfig() (gosnowflake.Config, error) {
	// Simply fallback to ConfigDTO behavior, as LegacyConfigDTO has compliant fields.
	configDTO := ConfigDTO(*c)
	return configDTO.DriverConfig()
}
