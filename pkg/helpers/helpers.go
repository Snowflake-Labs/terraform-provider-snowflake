package helpers

import (
	"encoding/csv"
	"fmt"
	"log"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	IDDelimiter = "|"
)

// ListContentToString strips list elements of double quotes or brackets.
func ListContentToString(listString string) string {
	re := regexp.MustCompile(`[\"\[\]]`)
	return re.ReplaceAllString(listString, "")
}

// StringToBool converts a string to a bool.
func StringToBool(s string) bool {
	return strings.ToLower(s) == "true"
}

// EncodeSnowflakeID generates a unique ID for a resource.
func EncodeSnowflakeID(attributes ...interface{}) string {
	// is attribute already an object identifier?
	if len(attributes) == 1 {
		if id, ok := attributes[0].(sdk.ObjectIdentifier); ok {
			if val := reflect.ValueOf(id); val.Kind() == reflect.Ptr && val.IsNil() {
				log.Panicf("Nil object identifier received")
			}
			parts := make([]string, 0)
			switch v := id.(type) {
			case sdk.AccountObjectIdentifier:
				parts = append(parts, v.Name())
			case *sdk.AccountObjectIdentifier:
				parts = append(parts, v.Name())
			case sdk.DatabaseObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.Name())
			case *sdk.DatabaseObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.Name())
			case sdk.SchemaObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.Name())
			case *sdk.SchemaObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.Name())
			case sdk.TableColumnIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.TableName(), v.Name())
			case *sdk.TableColumnIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.TableName(), v.Name())
			default:
				log.Panicf("Unsupported object identifier: %v", id)
			}
			return strings.Join(parts, IDDelimiter)
		}
	}
	var parts []string
	for i, attr := range attributes {
		if attr == nil {
			attributes[i] = ""
		}
		switch reflect.TypeOf(attr).Kind() {
		case reflect.String:
			parts = append(parts, attr.(string))
		case reflect.Bool:
			parts = append(parts, strconv.FormatBool(attr.(bool)))
		case reflect.Slice:
			parts = append(parts, strings.Join(attr.([]string), ","))
		}
	}
	return strings.Join(parts, "|")
}

func DecodeSnowflakeID(id string) sdk.ObjectIdentifier {
	parts := strings.Split(id, IDDelimiter)
	switch len(parts) {
	case 1:
		return sdk.NewAccountObjectIdentifier(parts[0])
	case 2:
		return sdk.NewDatabaseObjectIdentifier(parts[0], parts[1])
	case 3:
		return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], parts[2])
	case 4:
		return sdk.NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3])
	default:
		return nil
	}
}

// DecodeSnowflakeParameterID decodes identifier (usually passed as one of the parameter in tf configuration) into sdk.ObjectIdentifier.
// identifier can be specified in two ways: quoted and unquoted, e.g.
//
// quoted { "some_identifier": "\"database.name\".\"schema.name\".\"test.name\" }
// (note that here dots as part of the name are allowed)
//
// unquoted { "some_identifier": "database_name.schema_name.test_name" }
// (note that here dots as part of the name are NOT allowed, because they're treated in this case as dividers)
//
// The following configuration { "some_identifier": "db.name" } will be parsed as an object called "name" that lives
// inside database called "db", not a database called "db.name". In this case quotes should be used.
func DecodeSnowflakeParameterID(identifier string) (sdk.ObjectIdentifier, error) {
	parts, err := ParseIdentifierString(identifier)
	if err != nil {
		return nil, err
	}
	switch len(parts) {
	case 1:
		return sdk.NewAccountObjectIdentifier(parts[0]), nil
	case 2:
		return sdk.NewDatabaseObjectIdentifier(parts[0], parts[1]), nil
	case 3:
		return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], parts[2]), nil
	case 4:
		return sdk.NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3]), nil
	default:
		return nil, fmt.Errorf("unable to classify identifier: %s", identifier)
	}
}

// DecodeSnowflakeAccountIdentifier decodes account identifier (usually passed as one of the parameter in tf configuration) into sdk.AccountIdentifier.
// Check more in https://docs.snowflake.com/en/sql-reference/sql/create-account#required-parameters.
func DecodeSnowflakeAccountIdentifier(identifier string) (sdk.AccountIdentifier, error) {
	parts, err := ParseIdentifierString(identifier)
	if err != nil {
		return sdk.AccountIdentifier{}, err
	}
	switch len(parts) {
	case 1:
		return sdk.AccountIdentifier{}, fmt.Errorf("identifier: %s seems to be account locator and these are not allowed - please use <organization_name>.<account_name>", identifier)
	case 2:
		return sdk.NewAccountIdentifier(parts[0], parts[1]), nil
	default:
		return sdk.AccountIdentifier{}, fmt.Errorf("unable to classify account identifier: %s, expected format: <organization_name>.<account_name>", identifier)
	}
}

// TODO(SNOW-1479870): Test
// MergeMaps takes any number of maps (of the same type) and concatenates them.
// In case of key collision, the value will be selected from the map that is provided
// later in the src function parameter.
func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

// TODO: use slices.Concat in Go 1.22
func ConcatSlices[T any](slices ...[]T) []T {
	var tmp []T
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}

// TODO(SNOW-999049): address during identifiers rework
func ParseRootLocation(location string) (sdk.SchemaObjectIdentifier, string, error) {
	location = strings.TrimPrefix(location, "@")
	parts, err := parseIdentifierStringWithOpts(location, func(r *csv.Reader) {
		r.Comma = '.'
		r.LazyQuotes = true
	})
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, "", err
	}
	if len(parts) < 3 {
		return sdk.SchemaObjectIdentifier{}, "", fmt.Errorf("expected 3 parts for location %s, got %d", location, len(parts))
	}
	parts[2] = strings.Join(parts[2:], ".")
	lastParts := strings.Split(parts[2], "/")
	return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], lastParts[0]), path.Join(lastParts[1:]...), nil
}
