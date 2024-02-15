package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type v085FunctionId struct {
	DatabaseName string
	SchemaName   string
	FunctionName string
	ArgTypes     []string
}

func parseV085FunctionId(v string) (*v085FunctionId, error) {
	arr := strings.Split(v, "|")
	if len(arr) != 4 {
		return nil, sdk.NewError(fmt.Sprintf("ID %v is invalid", v))
	}

	// this is a bit different from V085 state, but it was buggy
	var args []string
	if arr[3] != "" {
		args = strings.Split(arr[3], "-")
	}

	return &v085FunctionId{
		DatabaseName: arr[0],
		SchemaName:   arr[1],
		FunctionName: arr[2],
		ArgTypes:     args,
	}, nil
}

func v085FunctionIdStateUpgrader(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldId := rawState["id"].(string)
	parsedV085FunctionId, err := parseV085FunctionId(oldId)
	if err != nil {
		return nil, err
	}

	argDataTypes := make([]sdk.DataType, len(parsedV085FunctionId.ArgTypes))
	for i, argType := range parsedV085FunctionId.ArgTypes {
		argDataType, err := sdk.ToDataType(argType)
		if err != nil {
			return nil, err
		}
		argDataTypes[i] = argDataType
	}

	schemaObjectIdentifierWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(parsedV085FunctionId.DatabaseName, parsedV085FunctionId.SchemaName, parsedV085FunctionId.FunctionName, argDataTypes)
	rawState["id"] = schemaObjectIdentifierWithArguments.FullyQualifiedName()

	return rawState, nil
}
