package resources

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func listValueToSlice(value string, trimQuotes bool) []string {
	value = strings.TrimLeft(value, "[")
	value = strings.TrimRight(value, "]")
	if value == "" {
		return nil
	}
	elems := strings.Split(value, ",")
	for i := range elems {
		if trimQuotes {
			elems[i] = strings.Trim(elems[i], " '")
		} else {
			elems[i] = strings.Trim(elems[i], " ")
		}
	}
	return elems
}

func ctyValToSliceString(valueElems []cty.Value) []string {
	elems := make([]string, len(valueElems))
	for i, v := range valueElems {
		elems[i] = v.AsString()
	}
	return elems
}
