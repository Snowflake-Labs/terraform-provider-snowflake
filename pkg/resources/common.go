package resources

import (
	"context"
	"regexp"
	"strings"

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

func ImportName[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier](ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	switch any(new(T)).(type) {
	case *sdk.AccountObjectIdentifier:
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}
	case *sdk.DatabaseObjectIdentifier:
		id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}

		if err := d.Set("database", id.DatabaseName()); err != nil {
			return nil, err
		}
	case *sdk.SchemaObjectIdentifier:
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}

		if err := d.Set("database", id.DatabaseName()); err != nil {
			return nil, err
		}

		if err := d.Set("schema", id.SchemaName()); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}
