package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (s *Saml2SecurityIntegrationModel) WithAllowedEmailPatterns(values ...string) *Saml2SecurityIntegrationModel {
	s.AllowedEmailPatterns = tfconfig.SetVariable(
		collections.Map(values, func(value string) tfconfig.Variable {
			return tfconfig.StringVariable(value)
		})...,
	)
	return s
}

func (s *Saml2SecurityIntegrationModel) WithAllowedUserDomains(values ...string) *Saml2SecurityIntegrationModel {
	s.AllowedUserDomains = tfconfig.SetVariable(
		collections.Map(values, func(value string) tfconfig.Variable {
			return tfconfig.StringVariable(value)
		})...,
	)
	return s
}
