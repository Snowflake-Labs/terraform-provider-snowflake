package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func StreamlitWithIds(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
	mainFile string,
	stageId sdk.SchemaObjectIdentifier,
) *StreamlitModel {
	return Streamlit(resourceName, id.DatabaseName(), mainFile, id.Name(), id.SchemaName(), stageId.FullyQualifiedName())
}

func (s *StreamlitModel) WithExternalAccessIntegrations(integrations ...sdk.AccountObjectIdentifier) *StreamlitModel {
	s.ExternalAccessIntegrations = tfconfig.SetVariable(
		collections.Map(integrations, func(role sdk.AccountObjectIdentifier) tfconfig.Variable {
			return tfconfig.StringVariable(role.Name())
		})...,
	)
	return s
}
