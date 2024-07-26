package helpers

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	IdDelimiter         = '.'
	ResourceIdDelimiter = '|'
)

func parseIdentifierStringWithOpts(identifier string, opts func(*csv.Reader)) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(identifier))
	if opts != nil {
		opts(reader)
	}
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read identifier: %s, err = %w", identifier, err)
	}
	if lines == nil {
		return make([]string, 0), nil
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("incompatible identifier: %s", identifier)
	}
	return lines[0], nil
}

func parseIdentifierString(identifier string) ([]string, error) {
	return parseIdentifierStringWithOpts(identifier, func(r *csv.Reader) {
		r.Comma = IdDelimiter
	})
}

func parseIdentifier[T sdk.ObjectIdentifier](identifier string, expectedParts int, expectedFormat string, constructFromParts func(parts []string) T) (T, error) {
	var emptyIdentifier T
	parts, err := parseIdentifierString(identifier)
	if err != nil {
		return emptyIdentifier, err
	}
	if len(parts) != expectedParts {
		return emptyIdentifier, fmt.Errorf(`unexpected number of parts %[1]d in identifier %[2]s, expected %[3]d in a form of "%[4]s"`, len(parts), identifier, expectedParts, expectedFormat)
	}
	return constructFromParts(parts), nil
}

func ParseResourceIdentifier(identifier string) []string {
	if identifier == "" {
		return make([]string, 0)
	}
	return strings.Split(identifier, string(ResourceIdDelimiter))
}

func EncodeResourceIdentifier(parts ...string) string {
	return strings.Join(parts, string(ResourceIdDelimiter))
}

func ParseAccountObjectIdentifier(identifier string) (sdk.AccountObjectIdentifier, error) {
	return parseIdentifier[sdk.AccountObjectIdentifier](
		identifier, 1, "<account_object_name>",
		func(parts []string) sdk.AccountObjectIdentifier {
			return sdk.NewAccountObjectIdentifier(parts[0])
		},
	)
}

func ParseDatabaseObjectIdentifier(identifier string) (sdk.DatabaseObjectIdentifier, error) {
	return parseIdentifier[sdk.DatabaseObjectIdentifier](
		identifier, 2, "<database_name>.<database_object_name>",
		func(parts []string) sdk.DatabaseObjectIdentifier {
			return sdk.NewDatabaseObjectIdentifier(parts[0], parts[1])
		},
	)
}

func ParseSchemaObjectIdentifier(identifier string) (sdk.SchemaObjectIdentifier, error) {
	return parseIdentifier[sdk.SchemaObjectIdentifier](
		identifier, 3, "<database_name>.<schema_name>.<schema_object_name>",
		func(parts []string) sdk.SchemaObjectIdentifier {
			return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], parts[2])
		},
	)
}

func ParseTableColumnIdentifier(identifier string) (sdk.TableColumnIdentifier, error) {
	return parseIdentifier[sdk.TableColumnIdentifier](
		identifier, 4, "<database_name>.<schema_name>.<table_name>.<table_column_name>",
		func(parts []string) sdk.TableColumnIdentifier {
			return sdk.NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3])
		},
	)
}

// ParseAccountIdentifier is implemented with an assumption that the recommended format is used that contains two parts,
// organization name and account name.
func ParseAccountIdentifier(identifier string) (sdk.AccountIdentifier, error) {
	return parseIdentifier[sdk.AccountIdentifier](
		identifier, 2, "<organization_name>.<account_name>",
		func(parts []string) sdk.AccountIdentifier {
			return sdk.NewAccountIdentifier(parts[0], parts[1])
		},
	)
}

// ParseExternalObjectIdentifier is implemented with an assumption that the identifier consists of three parts, because:
//   - After identifier rework, we expect account identifiers to always have two parts "<organization_name>.<account_name>".
//   - So far, the only external things that we referred to with external identifiers had only one part (not including the account identifier),
//     meaning it will always be represented as sdk.AccountObjectIdentifier. Documentation also doesn't describe any case where
//     account identifier would be used as part of the identifier that would refer to the "lower level" object.
//     Reference: https://docs.snowflake.com/en/user-guide/admin-account-identifier#where-are-account-identifiers-used.
func ParseExternalObjectIdentifier(identifier string) (sdk.ExternalObjectIdentifier, error) {
	return parseIdentifier[sdk.ExternalObjectIdentifier](
		identifier, 3, "<organization_name>.<account_name>.<external_object_name>",
		func(parts []string) sdk.ExternalObjectIdentifier {
			return sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(parts[0], parts[1]), sdk.NewAccountObjectIdentifier(parts[2]))
		},
	)
}
