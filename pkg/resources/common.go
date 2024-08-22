package resources

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const FullyQualifiedNameAttributeName = "fully_qualified_name"

// DiffSuppressStatement will suppress diffs between statements if they differ in only case or in
// runs of whitespace (\s+ = \s). This is needed because the snowflake api does not faithfully
// round-trip queries, so we cannot do a simple character-wise comparison to detect changes.
//
// Warnings: We will have false positives in cases where a change in case or run of whitespace is
// semantically significant.
//
// If we can find a sql parser that can handle the snowflake dialect then we should switch to parsing
// queries and either comparing ASTs or emitting a canonical serialization for comparison. I couldn't
// find such a library.
func DiffSuppressStatement(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(normalizeQuery(old), normalizeQuery(new))
}

var space = regexp.MustCompile(`\s+`)

func normalizeQuery(str string) string {
	return strings.TrimSpace(space.ReplaceAllString(str, " "))
}

// TODO [SNOW-999049]: address during identifiers rework
func suppressIdentifierQuoting(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return false
	} else {
		oldId, err := helpers.DecodeSnowflakeParameterID(oldValue)
		if err != nil {
			return false
		}
		newId, err := helpers.DecodeSnowflakeParameterID(newValue)
		if err != nil {
			return false
		}
		return oldId.FullyQualifiedName() == newId.FullyQualifiedName()
	}
}

// TODO [SNOW-1325214]: address during stage resource rework
func suppressQuoting(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return false
	} else {
		oldWithoutQuotes := strings.ReplaceAll(oldValue, "'", "")
		newWithoutQuotes := strings.ReplaceAll(newValue, "'", "")
		return oldWithoutQuotes == newWithoutQuotes
	}
}

func ctyValToSliceString(valueElems []cty.Value) []string {
	elems := make([]string, len(valueElems))
	for i, v := range valueElems {
		elems[i] = v.AsString()
	}
	return elems
}

func ensureWarehouse(ctx context.Context, client *sdk.Client) (func(), error) {
	warehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	if err != nil {
		return nil, err
	}
	if len(warehouse) > 0 {
		// everything is fine, return a no-op function to avoid checking by callers
		return func() {}, nil
	}
	randomWarehouseName := fmt.Sprintf("terraform-provider-snowflake-%v", helpers.RandomString())
	log.Printf("[DEBUG] no current warehouse set, creating a temporary warehouse %s", randomWarehouseName)
	wid := sdk.NewAccountObjectIdentifier(randomWarehouseName)
	if err := client.Warehouses.Create(ctx, wid, nil); err != nil {
		return nil, err
	}
	cleanup := func() {
		if err := client.Warehouses.Drop(ctx, wid, nil); err != nil {
			log.Printf("[WARN] error cleaning up temp warehouse %v", err)
		}
	}
	if err := client.Sessions.UseWarehouse(ctx, wid); err != nil {
		return cleanup, err
	}
	return cleanup, nil
}
