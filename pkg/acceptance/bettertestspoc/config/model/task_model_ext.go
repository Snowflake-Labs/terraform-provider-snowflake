package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TaskWithId(resourceName string, id sdk.SchemaObjectIdentifier, sqlStatement string) *TaskModel {
	t := &TaskModel{ResourceModelMeta: config.Meta(resourceName, resources.Task)}
	t.WithDatabase(id.DatabaseName())
	t.WithSchema(id.SchemaName())
	t.WithName(id.Name())
	t.WithSqlStatement(sqlStatement)
	return t
}
