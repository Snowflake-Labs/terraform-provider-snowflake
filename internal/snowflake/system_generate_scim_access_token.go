// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// SystemGenerateSCIMAccessTokenBuilder abstracts calling the SYSTEM$GENERATE_SCIM_ACCESS_TOKEN system function.
type SystemGenerateSCIMAccessTokenBuilder struct {
	integrationName string
}

// SystemGenerateSCIMAccessToken returns a pointer to a builder that abstracts calling the the SYSTEM$GENERATE_SCIM_ACCESS_TOKEN system function.
func NewSystemGenerateSCIMAccessTokenBuilder(integrationName string) *SystemGenerateSCIMAccessTokenBuilder {
	return &SystemGenerateSCIMAccessTokenBuilder{
		integrationName: integrationName,
	}
}

// Select generates the select statement for obtaining the scim access token.
func (pb *SystemGenerateSCIMAccessTokenBuilder) Select() string {
	return fmt.Sprintf(`SELECT SYSTEM$GENERATE_SCIM_ACCESS_TOKEN('%v') AS "TOKEN"`, pb.integrationName)
}

type SCIMAccessToken struct {
	Token string `db:"TOKEN"`
}

// ScanSCIMAccessToken convert a result into a.
func ScanSCIMAccessToken(row *sqlx.Row) (*SCIMAccessToken, error) {
	p := &SCIMAccessToken{}
	e := row.StructScan(p)
	return p, e
}
