package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// taken from https://docs.snowflake.com/en/user-guide/admin-account-identifier.html#snowflake-region-ids
var regionMapping = map[string]string{
	"aws_us_west_2":          "", // left black as this is the default
	"aws_us_east_2":          "us-east-2.aws",
	"aws_us_east_1":          "us-east-1",
	"aws_us_east_1_gov":      "us-east-1-gov.aws",
	"aws_ca_central_1":       "ca-central-1.aws",
	"aws_eu_west_1":          "eu-west-1",
	"aws_eu_west_2":          "eu-west-2.aws",
	"aws_eu_central_1":       "eu-central-1",
	"aws_ap_northeast_1":     "ap-northeast-1.aws",
	"aws_ap_south_1":         "ap-south-1.aws",
	"aws_ap_southeast_1":     "ap-southeast-1",
	"aws_ap_southeast_2":     "ap-southeast-2",
	"gcp_us_central1":        "us-central1.gcp",
	"gcp_europe_west2":       "europe-west2.gcp",
	"gcp_europe_west4":       "europe-west4.gcp",
	"azure_westus2":          "west-us-2.azure",
	"azure_eastus2":          "east-us-2.azure",
	"azure_usgovvirginia":    "us-gov-virginia.azure",
	"azure_canadacentral":    "canada-central.azure",
	"azure_westeurope":       "west-europe.azure",
	"azure_southeastasia":    "southeast-asia.azure",
	"azure_switzerlandnorth": "switzerland-north.azure",
	"azure_australiaeast":    "australia-east.azure",
}

func SelectCurrentAccount() string {
	return `SELECT CURRENT_ACCOUNT() AS "account", CURRENT_REGION() AS "region";`
}

type account struct {
	Account string `db:"account"`
	Region  string `db:"region"`
}

func ScanCurrentAccount(row *sqlx.Row) (*account, error) {
	acc := &account{}
	err := row.StructScan(acc)
	return acc, err
}

func ReadCurrentAccount(db *sql.DB) (*account, error) {
	row := QueryRow(db, SelectCurrentAccount())
	return ScanCurrentAccount(row)
}

func (acc *account) AccountURL() (string, error) {
	if region_id, ok := regionMapping[strings.ToLower(acc.Region)]; ok {
		account_id := acc.Account
		if len(region_id) > 0 {
			account_id = fmt.Sprintf("%s.%s", account_id, region_id)
		}
		return fmt.Sprintf("https://%s.snowflakecomputing.com", account_id), nil
	}

	return "", fmt.Errorf("Failed to map Snowflake account region %s to a region_id", acc.Region)
}
