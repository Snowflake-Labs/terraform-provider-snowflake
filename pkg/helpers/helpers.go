package helpers

import (
	"encoding/csv"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	IDDelimiter          = "|"
	ParameterIDDelimiter = '.'
)

// ToDo: We can merge these two functions together and also add more functions here with similar functionality

// This function converts list of string into snowflake formated string like 'ele1', 'ele2'.
func ListToSnowflakeString(list []string) string {
	for index, element := range list {
		list[index] = fmt.Sprintf(`'%v'`, strings.ReplaceAll(element, "'", "\\'"))
	}

	return fmt.Sprintf("%v", strings.Join(list, ", "))
}

// ListContentToString strips list elements of double quotes or brackets.
func ListContentToString(listString string) string {
	re := regexp.MustCompile(`[\"\[\]]`)
	return re.ReplaceAllString(listString, "")
}

// StringListToList splits a string into a slice of strings, separated by a separator. It also removes empty strings and trims whitespace.
func StringListToList(s string) []string {
	var v []string
	for _, elem := range strings.Split(s, ",") {
		if strings.TrimSpace(elem) != "" {
			v = append(v, strings.TrimSpace(elem))
		}
	}
	return v
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

func SafelyDecodeSnowflakeID[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier](stringIdentifier string) (id T, err error) {
	// TODO(SNOW-1163071): Right now we have to skip validation for AccountObjectIdentifier to handle a case where identifier contains dots
	if _, ok := any(sdk.AccountObjectIdentifier{}).(T); ok {
		var accountObjectIdentifier any = sdk.NewAccountObjectIdentifier(stringIdentifier)
		return accountObjectIdentifier.(T), nil
	}

	objectIdentifier, err := DecodeSnowflakeParameterID(stringIdentifier)
	if err != nil {
		return id, fmt.Errorf(
			"Unable to parse the identifier: %s. Make sure you are using the correct form of the fully qualified name for this field: %s.\nOriginal Error: %w",
			stringIdentifier,
			GetExpectedIdentifierRepresentationFromGeneric[T](),
			err,
		)
	}

	if _, ok := objectIdentifier.(T); !ok {
		return id, fmt.Errorf(
			"expected %s identifier type, but got: %T. The correct form of the fully qualified name for this field is: %s, but was %s",
			reflect.TypeOf(new(T)).Elem().Name(),
			objectIdentifier,
			GetExpectedIdentifierRepresentationFromGeneric[T](),
			GetExpectedIdentifierRepresentationFromParam(objectIdentifier),
		)
	}

	return objectIdentifier.(T), nil
}

func GetExpectedIdentifierRepresentationFromGeneric[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier]() string {
	return getExpectedIdentifierForm(new(T))
}

func GetExpectedIdentifierRepresentationFromParam(id sdk.ObjectIdentifier) string {
	return getExpectedIdentifierForm(id)
}

func getExpectedIdentifierForm(id any) string {
	switch id.(type) {
	case sdk.AccountObjectIdentifier, *sdk.AccountObjectIdentifier:
		return "<name>"
	case sdk.DatabaseObjectIdentifier, *sdk.DatabaseObjectIdentifier:
		return "<database_name>.<name>"
	case sdk.SchemaObjectIdentifier, *sdk.SchemaObjectIdentifier:
		return "<database_name>.<schema_name>.<name>"
	case sdk.TableColumnIdentifier, *sdk.TableColumnIdentifier:
		return "<database_name>.<schema_name>.<table_name>.<column_name>"
	}
	return ""
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
	reader := csv.NewReader(strings.NewReader(identifier))
	reader.Comma = ParameterIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read identifier: %s, err = %w", identifier, err)
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("incompatible identifier: %s", identifier)
	}
	parts := lines[0]
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

func Retry(attempts int, sleepDuration time.Duration, f func() (error, bool)) error {
	for i := 0; i < attempts; i++ {
		err, done := f()
		if err != nil {
			return err
		}
		if done {
			return nil
		} else {
			log.Printf("[INFO] operation not finished yet, retrying in %v seconds\n", sleepDuration.Seconds())
			time.Sleep(sleepDuration)
		}
	}
	return fmt.Errorf("giving up after %v attempts", attempts)
}
