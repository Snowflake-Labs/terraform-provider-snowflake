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
	"aws_us_gov_east_1":      "us-east-1-gov.aws",
	"aws_ca_central_1":       "ca-central-1.aws",
	"aws_sa_east_1":          "sa-east-1.aws",
	"aws_eu_west_1":          "eu-west-1",
	"aws_eu_west_2":          "eu-west-2.aws",
	"aws_eu_west_3":          "eu-west-3.aws",
	"aws_eu_central_1":       "eu-central-1",
	"aws_eu_north_1":         "eu-north-1.aws",
	"aws_ap_northeast_1":     "ap-northeast-1.aws",
	"aws_ap_northeast_2":     "ap-northeast-2.aws",
	"aws_ap_northeast_3":     "ap-northeast-3.aws",
	"aws_ap_south_1":         "ap-south-1.aws",
	"aws_ap_southeast_1":     "ap-southeast-1",
	"aws_ap_southeast_2":     "ap-southeast-2",
	"gcp_us_central1":        "us-central1.gcp",
	"gcp_us_east4":           "us-east4.gcp",
	"gcp_europe_west2":       "europe-west2.gcp",
	"gcp_europe_west4":       "europe-west4.gcp",
	"azure_westus2":          "west-us-2.azure",
	"azure_centralus":        "central-us.azure",
	"azure_southcentralus":   "south-central-us.azure",
	"azure_eastus2":          "east-us-2.azure",
	"azure_usgovvirginia":    "us-gov-virginia.azure",
	"azure_canadacentral":    "canada-central.azure",
	"azure_uksouth":          "uk-south.azure",
	"azure_northeurope":      "north-europe.azure",
	"azure_westeurope":       "west-europe.azure",
	"azure_southeastasia":    "southeast-asia.azure",
	"azure_switzerlandnorth": "switzerland-north.azure",
	"azure_uaenorth":         "uae-north.azure",
	"azure_centralindia":     "central-india.azure",
	"azure_japaneast":        "japan-east.azure",
	"azure_australiaeast":    "australia-east.azure",
}

// SelectCurrentAccount returns the query that will return the current account, region_group, and region
// the CURRENT_REGION() function returns the format region_group.region (e.g. PUBLIC.AWS_US_WEST_2) only when organizations have accounts in multiple region groups. Otherwise, this function returns the snowflake region without the region_group.
func SelectCurrentAccount() string {
	return `SELECT CURRENT_ACCOUNT() as "account",  CASE WHEN CONTAINS(CURRENT_REGION(), '.') THEN LEFT(CURRENT_REGION(), POSITION('.' IN CURRENT_REGION()) - 1) ELSE 'PUBLIC' END AS "region_group", CASE WHEN CONTAINS(CURRENT_REGION(), '.') THEN RIGHT(CURRENT_REGION(), LENGTH(CURRENT_REGION()) - POSITION('.' IN CURRENT_REGION())) ELSE CURRENT_REGION() END AS "region";`
}

type CurrentAccount struct {
	Account     string `db:"account"`
	RegionGroup string `db:"region_group"`
	Region      string `db:"region"`
}

func ScanCurrentAccount(row *sqlx.Row) (*CurrentAccount, error) {
	acc := &CurrentAccount{}
	err := row.StructScan(acc)
	return acc, err
}

func ReadCurrentAccount(db *sql.DB) (*CurrentAccount, error) {
	row := QueryRow(db, SelectCurrentAccount())
	return ScanCurrentAccount(row)
}

func (acc *CurrentAccount) AccountURL() (string, error) {
	if regionID, ok := regionMapping[strings.ToLower(acc.Region)]; ok {
		accountID := acc.Account
		if len(regionID) > 0 {
			accountID = fmt.Sprintf("%s.%s", accountID, regionID)
		}
		return fmt.Sprintf("https://%s.snowflakecomputing.com", accountID), nil
	}
	return "", fmt.Errorf("failed to map Snowflake account region %s to a region_id", acc.Region)
}
