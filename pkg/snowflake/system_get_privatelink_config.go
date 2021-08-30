package snowflake

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func SystemGetPrivateLinkConfigQuery() string {
	return `SELECT SYSTEM$GET_PRIVATELINK_CONFIG() AS "config"`
}

type RawPrivateLinkConfig struct {
	Config string `db:"config"`
}

type privateLinkConfigInternal struct {
	AccountName               string `json:"privatelink-account-name"`
	AwsVpceID                 string `json:"privatelink-vpce-id,omitempty"`
	AzurePrivateLinkServiceID string `json:"privatelink-pls-id,omitempty"`
	AccountURL                string `json:"privatelink-account-url"`
	OCSPURL                   string `json:"privatelink-ocsp-url,omitempty"`
	TypodOCSPURL              string `json:"privatelink_ocsp-url,omitempty"` // because snowflake returns this for AWS, but don't have an Azure account to verify against
}

type PrivateLinkConfig struct {
	AccountName               string
	AwsVpceID                 string
	AzurePrivateLinkServiceID string
	AccountURL                string
	OCSPURL                   string
}

func ScanPrivateLinkConfig(row *sqlx.Row) (*RawPrivateLinkConfig, error) {
	config := &RawPrivateLinkConfig{}
	err := row.StructScan(config)
	return config, err
}

func (r *RawPrivateLinkConfig) GetStructuredConfig() (*PrivateLinkConfig, error) {
	config := &privateLinkConfigInternal{}
	err := json.Unmarshal([]byte(r.Config), config)
	if err != nil {
		return nil, err
	}

	return config.getPrivateLinkConfig()
}

func (i *privateLinkConfigInternal) getPrivateLinkConfig() (*PrivateLinkConfig, error) {
	config := &PrivateLinkConfig{
		i.AccountName,
		i.AwsVpceID,
		i.AzurePrivateLinkServiceID,
		i.AccountURL,
		i.OCSPURL,
	}

	if i.TypodOCSPURL != "" {
		config.OCSPURL = i.TypodOCSPURL
	}

	return config, nil
}
