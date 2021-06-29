package snowflake

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// SystemGenerateSCIMAccessTokenBuilder abstracts calling the SYSTEM$GENERATE_SCIM_ACCESS_TOKEN system function
type SystemGenerateSCIMAccessTokenBuilder struct {
	integrationName string
}

// SystemGenerateSCIMAccessToken returns a pointer to a builder that abstracts calling the the SYSTEM$GENERATE_SCIM_ACCESS_TOKEN system function
func SystemGenerateSCIMAccessToken(integrationName string) *SystemGenerateSCIMAccessTokenBuilder {
	return &SystemGenerateSCIMAccessTokenBuilder{
		integrationName: integrationName,
	}
}

// Select generates the select statement for obtaining the scim access token
func (pb *SystemGenerateSCIMAccessTokenBuilder) Select() string {
	return fmt.Sprintf(`SELECT SYSTEM$GENERATE_SCIM_ACCESS_TOKEN('%v') AS "token"`, pb.integrationName)
}

type scimAccessToken struct {
	Token string `db:"token"`
}

// ScanSCIMAccessToken convert a result into a
func ScanSCIMAccessToken(row *sqlx.Row) (*scimAccessToken, error) {
	p := &scimAccessToken{}
	e := row.StructScan(p)
	return p, e
}
