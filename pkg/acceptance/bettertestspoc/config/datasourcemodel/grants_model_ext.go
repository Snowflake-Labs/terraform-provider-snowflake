package datasourcemodel

import (
	"fmt"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func GrantsOnAccount(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"account": tfconfig.BoolVariable(true),
			}),
		)
}

func GrantsOnAccountObject(
	datasourceName string,
	id sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable(id.Name()),
				"object_type": tfconfig.StringVariable(fmt.Sprintf("%s", objectType)),
			}),
		)
}

func GrantsOnDatabaseObject(
	datasourceName string,
	id sdk.DatabaseObjectIdentifier,
	objectType sdk.ObjectType,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable(id.FullyQualifiedName()),
				"object_type": tfconfig.StringVariable(fmt.Sprintf("%s", objectType)),
			}),
		)
}

func GrantsOnSchemaObject(
	datasourceName string,
	id sdk.SchemaObjectIdentifier,
	objectType sdk.ObjectType,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable(id.FullyQualifiedName()),
				"object_type": tfconfig.StringVariable(fmt.Sprintf("%s", objectType)),
			}),
		)
}

func GrantsOnSchemaObjectWithArguments(
	datasourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	objectType sdk.ObjectType,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable(id.FullyQualifiedName()),
				"object_type": tfconfig.StringVariable(fmt.Sprintf("%s", objectType)),
			}),
		)
}

func GrantsOnMissingObjectType(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable("DATABASE"),
			}),
		)
}

func GrantsOnEmpty(
	datasourceName string,
) *GrantsModel {
	placeholder := fmt.Sprintf("%s", config.SnowflakeProviderConfigSingleAttributeWorkaround)
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"placeholder": tfconfig.StringVariable(placeholder),
			}),
		)
}
