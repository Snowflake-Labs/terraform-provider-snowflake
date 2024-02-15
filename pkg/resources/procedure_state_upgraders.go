package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type v085ProcedureId struct {
	DatabaseName  string
	SchemaName    string
	ProcedureName string
	ArgTypes      []string
}

func parseV085ProcedureId(v string) (*v085ProcedureId, error) {
	arr := strings.Split(v, "|")
	if len(arr) != 4 {
		return nil, fmt.Errorf("ID %v is invalid", v)
	}

	// this is a bit different from V085 state, but it was buggy
	var args []string
	if arr[3] != "" {
		args = strings.Split(arr[3], "-")
	}

	return &v085ProcedureId{
		DatabaseName:  arr[0],
		SchemaName:    arr[1],
		ProcedureName: arr[2],
		ArgTypes:      args,
	}, nil
}

func v085ProcedureStateUpgrader(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldId := rawState["id"].(string)
	parsedV085ProcedureId, err := parseV085ProcedureId(oldId)
	if err != nil {
		return nil, err
	}

	argDataTypes := make([]sdk.DataType, len(parsedV085ProcedureId.ArgTypes))
	for i, argType := range parsedV085ProcedureId.ArgTypes {
		argDataType, err := sdk.ToDataType(argType)
		if err != nil {
			return nil, err
		}
		argDataTypes[i] = argDataType
	}

	schemaObjectIdentifierWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(parsedV085ProcedureId.DatabaseName, parsedV085ProcedureId.SchemaName, parsedV085ProcedureId.ProcedureName, argDataTypes)
	rawState["id"] = schemaObjectIdentifierWithArguments.FullyQualifiedName()

	return rawState, nil
}
