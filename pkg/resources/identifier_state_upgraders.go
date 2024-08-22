package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
)

func migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	oldId := helpers.DecodeSnowflakeID(rawState["id"].(string))
	rawState["id"] = oldId.FullyQualifiedName()
	return rawState, nil
}
