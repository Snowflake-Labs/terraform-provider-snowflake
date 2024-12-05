package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (t *TagAssociationModel) WithObjectIdentifiers(objectIdentifiers ...sdk.ObjectIdentifier) *TagAssociationModel {
	objectIdentifiersStringVariables := make([]tfconfig.Variable, len(objectIdentifiers))
	for i, v := range objectIdentifiers {
		objectIdentifiersStringVariables[i] = tfconfig.StringVariable(v.FullyQualifiedName())
	}

	t.ObjectIdentifiers = tfconfig.SetVariable(objectIdentifiersStringVariables...)
	return t
}
