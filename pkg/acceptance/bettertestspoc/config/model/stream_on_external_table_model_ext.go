package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func StreamOnExternalTableBase(resourceName string, id, externalTableId sdk.SchemaObjectIdentifier) *StreamOnExternalTableModel {
	return StreamOnExternalTable(resourceName, id.DatabaseName(), externalTableId.FullyQualifiedName(), id.Name(), id.SchemaName()).WithInsertOnly("true")
}
