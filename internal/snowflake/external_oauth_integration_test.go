// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/stretchr/testify/require"
)

func TestCreateExternalOauthIntegration3(t *testing.T) {
	r := require.New(t)

	input := &snowflake.ExternalOauthIntegration3CreateInput{
		ExternalOauthIntegration3: snowflake.ExternalOauthIntegration3{
			TopLevelIdentifier: snowflake.TopLevelIdentifier{
				Name: "azure",
			},
			Type:                "EXTERNAL_OAUTH",
			TypeOk:              true,
			ExternalOauthType:   "AZURE",
			ExternalOauthTypeOk: true,
		},
	}

	mb, err := snowflake.NewExternalOauthIntegration3Manager()
	r.Nil(err)
	createStmt, err := mb.Create(input)
	r.Nil(err)
	r.Equal(`CREATE SECURITY INTEGRATION "azure" type = 'EXTERNAL_OAUTH' EXTERNAL_OAUTH_TYPE = 'AZURE';`, createStmt)
}

func TestAlterExternalOauthIntegration3(t *testing.T) {
	r := require.New(t)

	input := &snowflake.ExternalOauthIntegration3UpdateInput{
		ExternalOauthIntegration3: snowflake.ExternalOauthIntegration3{
			TopLevelIdentifier: snowflake.TopLevelIdentifier{
				Name: "azure",
			},
			ExternalOauthIssuer:             "someissuer",
			ExternalOauthIssuerOk:           true,
			ExternalOauthBlockedRolesList:   []string{"a", "b"},
			ExternalOauthBlockedRolesListOk: true,
		},

		IfExists:   true,
		IfExistsOk: true,
	}

	mb, err := snowflake.NewExternalOauthIntegration3Manager()
	r.Nil(err)
	alterStmt, err := mb.Update(input)
	r.Nil(err)
	r.Equal(
		`ALTER SECURITY INTEGRATION IF EXISTS "azure" SET EXTERNAL_OAUTH_ISSUER = 'someissuer' EXTERNAL_OAUTH_BLOCKED_ROLES_LIST = ('a', 'b');`,
		alterStmt,
	)
}

func TestUnsetExternalOauthIntegration3(t *testing.T) {
	r := require.New(t)

	input := &snowflake.ExternalOauthIntegration3UpdateInput{
		ExternalOauthIntegration3: snowflake.ExternalOauthIntegration3{
			TopLevelIdentifier: snowflake.TopLevelIdentifier{
				Name: "azure",
			},
			ExternalOauthTokenUserMappingClaimOk: true,
		},
	}

	mb, err := snowflake.NewExternalOauthIntegration3Manager()
	r.Nil(err)
	unsetStmt, err := mb.Unset(input)
	r.Nil(err)
	r.Equal(
		`ALTER SECURITY INTEGRATION "azure" UNSET EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM;`,
		unsetStmt,
	)
}

func TestDeleteExternalOauthIntegration3(t *testing.T) {
	r := require.New(t)

	input := &snowflake.ExternalOauthIntegration3DeleteInput{
		TopLevelIdentifier: snowflake.TopLevelIdentifier{
			Name: "azure",
		},
	}

	mb, err := snowflake.NewExternalOauthIntegration3Manager()
	r.Nil(err)
	dropStmt, err := mb.Delete(input)
	r.Nil(err)
	r.Equal(`DROP SECURITY INTEGRATION "azure";`, dropStmt)
}

func TestReadDescribeExternalOauthIntegration3(t *testing.T) {
	r := require.New(t)

	input := &snowflake.ExternalOauthIntegration3ReadInput{
		Name: "azure",
	}

	mb, err := snowflake.NewExternalOauthIntegration3Manager()
	r.Nil(err)
	describeStmt, err := mb.ReadDescribe(input)
	r.Nil(err)
	r.Equal(`DESCRIBE SECURITY INTEGRATION "azure";`, describeStmt)
}

func TestReadShowExternalOauthIntegration3(t *testing.T) {
	r := require.New(t)

	input := &snowflake.ExternalOauthIntegration3ReadInput{
		Name: "azure",
	}

	mb, err := snowflake.NewExternalOauthIntegration3Manager()
	r.Nil(err)
	describeStmt, err := mb.ReadShow(input)
	r.Nil(err)
	r.Equal(`SHOW SECURITY INTEGRATIONS LIKE 'azure';`, describeStmt)
}
