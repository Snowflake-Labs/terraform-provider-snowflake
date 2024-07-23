package helpers

import (
	"encoding/csv"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strings"
)

const (
	IdDelimiter         = '.'
	ResourceIdDelimiter = '|'
)

type identifierParsingFunc func(string) ([]string, error)

func parseIdentifierStringWithOpts(identifier string, opts func(*csv.Reader)) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(identifier))
	if opts != nil {
		opts(reader)
	}
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read identifier: %s, err = %w", identifier, err)
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("incompatible identifier: %s", identifier)
	}
	return lines[0], nil
}

func ParseIdentifierString(identifier string) ([]string, error) {
	return parseIdentifierStringWithOpts(identifier, func(r *csv.Reader) {
		r.Comma = IdDelimiter
	})
}

func EncodeResourceIdentifier(parts ...string) string {
	return strings.Join(parts, string(ResourceIdDelimiter))
}

func ParseResourceIdentifier(identifier string) ([]string, error) {
	return parseIdentifierStringWithOpts(identifier, func(r *csv.Reader) {
		r.Comma = ResourceIdDelimiter
	})
}

func parseAccountObjectIdentifier(identifier string, parser identifierParsingFunc) (sdk.AccountObjectIdentifier, error) {
	parts, err := parser(identifier)
	if err != nil {
		return sdk.AccountObjectIdentifier{}, err
	}
	if len(parts) != 1 {
		return sdk.AccountObjectIdentifier{}, fmt.Errorf(`unexpected number of parts %d in identifier %s, expected 1 in a form of "<account_object_name>"`, len(parts), identifier)
	}
	return sdk.NewAccountObjectIdentifier(parts[0]), nil
}

func ParseAccountObjectIdentifier(identifier string) (sdk.AccountObjectIdentifier, error) {
	return parseAccountObjectIdentifier(identifier, ParseIdentifierString)
}

func ParseAccountObjectResourceIdentifier(identifier string) (sdk.AccountObjectIdentifier, error) {
	return parseAccountObjectIdentifier(identifier, ParseResourceIdentifier)
}

func parseDatabaseObjectIdentifier(identifier string, parser identifierParsingFunc, delimiter rune) (sdk.DatabaseObjectIdentifier, error) {
	parts, err := parser(identifier)
	if err != nil {
		return sdk.DatabaseObjectIdentifier{}, err
	}
	if len(parts) != 2 {
		return sdk.DatabaseObjectIdentifier{}, fmt.Errorf(`unexpected number of parts %d in identifier %s, expected 2 in a form of "<database_name>%c<database_object_name>"`, len(parts), identifier, delimiter)
	}
	return sdk.NewDatabaseObjectIdentifier(parts[0], parts[1]), nil
}

func ParseDatabaseObjectIdentifier(identifier string) (sdk.DatabaseObjectIdentifier, error) {
	return parseDatabaseObjectIdentifier(identifier, ParseIdentifierString, IdDelimiter)
}

func ParseDatabaseObjectResourceIdentifier(identifier string) (sdk.DatabaseObjectIdentifier, error) {
	return parseDatabaseObjectIdentifier(identifier, ParseResourceIdentifier, ResourceIdDelimiter)
}

func parseSchemaObjectIdentifier(identifier string, parser identifierParsingFunc, delimiter rune) (sdk.SchemaObjectIdentifier, error) {
	parts, err := parser(identifier)
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, err
	}
	if len(parts) != 3 {
		return sdk.SchemaObjectIdentifier{}, fmt.Errorf(`unexpected number of parts %[1]d in identifier %[2]s, expected 3 in a form of "<database_name>%[3]c<schema_name>%[3]c<schema_object_name>"`, len(parts), identifier, delimiter)
	}
	return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], parts[2]), nil
}

func ParseSchemaObjectIdentifier(identifier string) (sdk.SchemaObjectIdentifier, error) {
	return parseSchemaObjectIdentifier(identifier, ParseIdentifierString, IdDelimiter)
}

func ParseSchemaObjectResourceIdentifier(identifier string) (sdk.SchemaObjectIdentifier, error) {
	return parseSchemaObjectIdentifier(identifier, ParseResourceIdentifier, ResourceIdDelimiter)
}

func parseTableColumnObjectIdentifier(identifier string, parser identifierParsingFunc, delimiter rune) (sdk.TableColumnIdentifier, error) {
	parts, err := parser(identifier)
	if err != nil {
		return sdk.TableColumnIdentifier{}, err
	}
	if len(parts) != 4 {
		return sdk.TableColumnIdentifier{}, fmt.Errorf(`unexpected number of parts %[1]d in identifier %[2]s, expected 4 in a form of "<database_name>%[3]c<schema_name>%[3]c<table_name>%[3]c<table_column_name>"`, len(parts), identifier, delimiter)
	}
	return sdk.NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3]), nil
}

func ParseTableColumnIdentifier(identifier string) (sdk.TableColumnIdentifier, error) {
	return parseTableColumnObjectIdentifier(identifier, ParseIdentifierString, IdDelimiter)
}

func ParseTableColumnResourceIdentifier(identifier string) (sdk.TableColumnIdentifier, error) {
	return parseTableColumnObjectIdentifier(identifier, ParseResourceIdentifier, ResourceIdDelimiter)
}

func parseAccountIdentifier(identifier string, parser identifierParsingFunc, delimiter rune) (sdk.AccountIdentifier, error) {
	parts, err := parser(identifier)
	if err != nil {
		return sdk.AccountIdentifier{}, err
	}
	if len(parts) != 2 {
		return sdk.AccountIdentifier{}, fmt.Errorf(`unexpected number of parts %d in identifier %s, expected 2 in a form of "<organization_name>%c<account_name>"`, len(parts), identifier, delimiter)
	}
	return sdk.NewAccountIdentifier(parts[0], parts[1]), nil
}

func ParseAccountIdentifier(identifier string) (sdk.AccountIdentifier, error) {
	return parseAccountIdentifier(identifier, ParseIdentifierString, IdDelimiter)
}

func ParseAccountResourceIdentifier(identifier string) (sdk.AccountIdentifier, error) {
	return parseAccountIdentifier(identifier, ParseResourceIdentifier, ResourceIdDelimiter)
}

// parseExternalObjectIdentifier is implemented with an assumption that the identifier consists of three parts, because:
//   - After identifier rework, we expect account identifiers to always have two parts "<organization_name>.<account_name>".
//   - So far, the only external things that we referred to with external identifiers had only one part (not including the account identifier),
//     meaning it will always be represented as sdk.AccountObjectIdentifier. Documentation also doesn't describe any case where
//     account identifier would be used as part of the identifier that would refer to the "lower level" object.
//     Reference: https://docs.snowflake.com/en/user-guide/admin-account-identifier#where-are-account-identifiers-used.
func parseExternalObjectIdentifier(identifier string, parser identifierParsingFunc, delimiter rune) (sdk.ExternalObjectIdentifier, error) {
	parts, err := parser(identifier)
	if err != nil {
		return sdk.ExternalObjectIdentifier{}, err
	}
	if len(parts) != 3 {
		return sdk.ExternalObjectIdentifier{}, fmt.Errorf(`unexpected number of parts %[1]d in identifier %[2]s, expected 3 in a form of "<organization_name>%[3]c<account_name>%[3]c<external_object_name>"`, len(parts), identifier, delimiter)
	}
	return sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(parts[0], parts[1]), sdk.NewAccountObjectIdentifier(parts[2])), nil
}

func ParseExternalObjectIdentifier(identifier string) (sdk.ExternalObjectIdentifier, error) {
	return parseExternalObjectIdentifier(identifier, ParseIdentifierString, IdDelimiter)
}

func ParseExternalObjectResourceIdentifier(identifier string) (sdk.ExternalObjectIdentifier, error) {
	return parseExternalObjectIdentifier(identifier, ParseResourceIdentifier, ResourceIdDelimiter)
}
