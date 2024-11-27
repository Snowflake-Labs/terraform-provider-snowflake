package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// FullTomlConfigForServiceUser is a temporary function used to test provider configuration
// TODO [next PR]: use toml marshalling from "github.com/pelletier/go-toml/v2"
// TODO [next PR]: add builders for our toml config struct
func FullTomlConfigForServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string) string {
	t.Helper()

	return fmt.Sprintf(`
[%[1]s]
user = '%[2]s'
privatekey = '''%[7]s'''
role = '%[3]s'
organizationname = '%[5]s'
accountname = '%[6]s'
warehouse = '%[4]s'
clientip = '1.2.3.4'
protocol = 'https'
port = 443
oktaurl = 'https://example.com'
clienttimeout = 10
jwtclienttimeout = 20
logintimeout = 30
requesttimeout = 40
jwtexpiretimeout = 50
externalbrowsertimeout = 60
maxretrycount = 1
authenticator = 'SNOWFLAKE_JWT'
insecuremode = true
ocspfailopen = true
token = 'token'
keepsessionalive = true
disabletelemetry = true
validatedefaultparameters = true
clientrequestmfatoken = true
clientstoretemporarycredential = true
tracing = 'warning'
tmpdirpath = '.'
disablequerycontextcache = true
includeretryreason = true
disableconsolelogin = true

[%[1]s.params]
foo = 'bar'
`, profile, userId.Name(), roleId.Name(), warehouseId.Name(), accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey)
}
