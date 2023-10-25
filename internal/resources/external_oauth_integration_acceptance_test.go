// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/acceptance"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ExternalOauthIntegration(t *testing.T) {
	oauthIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	integrationType := "AZURE"

	issuer := fmt.Sprintf("https://sts.windows.net/%s", uuid.NewString())

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalOauthIntegrationConfig(oauthIntName, integrationType, issuer, "test resource"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", oauthIntName),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "type", integrationType),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "snowflake_user_mapping_attribute", "LOGIN_NAME"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.#", "2"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.0", "test"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.1", "upn"),
				),
			},
			{
				ResourceName:      "snowflake_external_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ExternalOauthIntegrationEmptyComment(t *testing.T) {
	oauthIntName := strings.ToLower(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	integrationType := "AZURE"

	issuer := fmt.Sprintf("https://sts.windows.net/%s", uuid.NewString())

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalOauthIntegrationConfig(oauthIntName, integrationType, issuer, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", oauthIntName),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "type", integrationType),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "snowflake_user_mapping_attribute", "LOGIN_NAME"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.#", "2"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.0", "test"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.1", "upn"),
				),
			},
			{
				ResourceName:      "snowflake_external_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ExternalOauthIntegrationLowercaseName(t *testing.T) {
	oauthIntName := strings.ToLower(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	integrationType := "AZURE"

	issuer := fmt.Sprintf("https://sts.windows.net/%s", uuid.NewString())

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalOauthIntegrationConfig(oauthIntName, integrationType, issuer, "test resource"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", oauthIntName),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "type", integrationType),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "snowflake_user_mapping_attribute", "LOGIN_NAME"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.#", "2"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.0", "test"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.1", "upn"),
				),
			},
			{
				ResourceName:      "snowflake_external_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ExternalOauthIntegrationCustom(t *testing.T) {
	oauthIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	integrationType := "CUSTOM"

	issuer := fmt.Sprintf("https://sts.windows.net/%s", uuid.NewString())

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "snowflake_external_oauth_integration" "test" {
					name = "%s"
					type = "%s"
					enabled = true
					issuer = "%s"
					snowflake_user_mapping_attribute = "LOGIN_NAME"
					jws_keys_urls = ["https://login.windows.net/common/discovery/keys"]
					audience_urls = ["https://analysis.windows.net/powerbi/connector/Snowflake"]
					token_user_mapping_claims = ["upn", "test"]
					scope_mapping_attribute = "scp"
					comment = "hey"
				}
				`, oauthIntName, integrationType, issuer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", oauthIntName),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "type", integrationType),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "snowflake_user_mapping_attribute", "LOGIN_NAME"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "scope_mapping_attribute", "scp"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.#", "2"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.0", "test"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.1", "upn"),
				),
			},
			{
				ResourceName:      "snowflake_external_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func externalOauthIntegrationConfig(name, integrationType, issuer, comment string) string {
	return fmt.Sprintf(`
	resource "snowflake_external_oauth_integration" "test" {
		name = "%s"
		type = "%s"
		enabled = true
  		issuer = "%s"
  		snowflake_user_mapping_attribute = "LOGIN_NAME"
		jws_keys_urls = ["https://login.windows.net/common/discovery/keys"]
		audience_urls = ["https://analysis.windows.net/powerbi/connector/Snowflake"]
  		token_user_mapping_claims = ["upn", "test"]
		comment = "%s"
	}
	`, name, integrationType, issuer, comment)
}
