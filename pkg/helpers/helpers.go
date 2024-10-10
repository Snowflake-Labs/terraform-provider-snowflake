package helpers

import (
	"encoding/csv"
	"encoding/json"
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
	parts, err := sdk.ParseIdentifierString(identifier)
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
	parts, err := sdk.ParseIdentifierString(identifier)
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

// Structs for parsing external volume desribe output
type S3StorageLocation struct {
	Name                    string   `json:"NAME"`
	StorageProvider         string   `json:"STORAGE_PROVIDER"`
	StorageBaseUrl          string   `json:"STORAGE_BASE_URL"`
	StorageAllowedLocations []string `json:"-"`
	StorageAwsRoleArn       string   `json:"STORAGE_AWS_ROLE_ARN"`
	StroageAwsIamUserArn    string   `json:"-"`
	StorageAwsExternalId    string   `json:"STORAGE_AWS_EXTERNAL_ID"`
	EncryptionType          string   `json:"ENCRYPTION_TYPE,omitempty"`
	EncryptionKmsId         string   `json:"ENCRYPTION_KMS_KEY_ID,omitempty"`
}

type GCSStorageLocation struct {
	Name                     string   `json:"NAME"`
	StorageProvider          string   `json:"STORAGE_PROVIDER"`
	StorageBaseUrl           string   `json:"STORAGE_BASE_URL"`
	StorageAllowedLocations  []string `json:"-"`
	StorageGcpServiceAccount string   `json:"-"`
	EncryptionType           string   `json:"ENCRYPTION_TYPE,omitempty"`
	EncryptionKmsId          string   `json:"ENCRYPTION_KMS_KEY_ID,omitempty"`
}

type AzureStorageLocation struct {
	Name                    string   `json:"NAME"`
	StorageProvider         string   `json:"STORAGE_PROVIDER"`
	StorageBaseUrl          string   `json:"STORAGE_BASE_URL"`
	StorageAllowedLocations []string `json:"-"`
	AzureTenantId           string   `json:"AZURE_TENANT_ID"`
	AzureMultiTenantAppName string   `json:"-"`
	AzureConsentUrl         string   `json:"-"`
	EncryptionType          string   `json:"-"`
	EncryptionKmsId         string   `json:"-"`
}

type StorageLocation struct {
	Name                 string `json:"NAME"`
	StorageProvider      string `json:"STORAGE_PROVIDER"`
	StorageBaseUrl       string `json:"STORAGE_BASE_URL"`
	StorageAwsRoleArn    string `json:"STORAGE_AWS_ROLE_ARN,omitempty"`
	StorageAwsExternalId string `json:"STORAGE_AWS_EXTERNAL_ID,omitempty"`
	EncryptionType       string `json:"ENCRYPTION_TYPE,omitempty"`
	EncryptionKmsKeyId   string `json:"ENCRYPTION_KMS_KEY_ID,omitempty"`
	AzureTenantId        string `json:"AZURE_TENANT_ID,omitempty"`
}

func storageLocationsEqual(s1 StorageLocation, s2 StorageLocation) bool {
	return s1.Name == s2.Name &&
		s1.StorageProvider == s2.StorageProvider &&
		s1.StorageBaseUrl == s2.StorageBaseUrl &&
		s1.StorageAwsRoleArn == s2.StorageAwsRoleArn &&
		s1.StorageAwsExternalId == s2.StorageAwsExternalId &&
		s1.EncryptionType == s2.EncryptionType &&
		s1.EncryptionKmsKeyId == s2.EncryptionKmsKeyId &&
		s1.AzureTenantId == s2.AzureTenantId
}

func validateParsedExternalVolumeDescribed(p ParsedExternalVolumeDescribed) error {
	if len(p.StorageLocations) == 0 {
		return fmt.Errorf("No storage locations could be parsed from the external volume.")
	}
	if len(p.AllowWrites) == 0 {
		return fmt.Errorf("The external volume AllowWrites property could not be parsed.")
	}

	for _, s := range p.StorageLocations {
		if len(s.Name) == 0 {
			return fmt.Errorf("A storage location's Name in this volume could not be parsed.")
		}
		if len(s.StorageProvider) == 0 {
			return fmt.Errorf("A storage location's StorageProvider in this volume could not be parsed.")
		}
		if len(s.StorageBaseUrl) == 0 {
			return fmt.Errorf("A storage location's StorageBaseUrl in this volume could not be parsed.")
		}

		storageProvider, err := sdk.ToStorageProvider(s.StorageProvider)
		if err != nil {
			return err
		}

		switch storageProvider {
		case sdk.StorageProviderS3, sdk.StorageProviderS3GOV:
			if len(s.StorageAwsRoleArn) == 0 {
				return fmt.Errorf("An S3 storage location's StorageAwsRoleArn in this volume could not be parsed.")
			}
		case sdk.StorageProviderAzure:
			if len(s.AzureTenantId) == 0 {
				return fmt.Errorf("An Azure storage location's AzureTenantId in this volume could not be parsed.")
			}
		}
	}

	return nil
}

type ParsedExternalVolumeDescribed struct {
	StorageLocations []StorageLocation
	Active           string
	Comment          string
	AllowWrites      string
}

func ParsedExternalVolumesDescribedEqual(p1 ParsedExternalVolumeDescribed, p2 ParsedExternalVolumeDescribed) bool {
	attributesEqual := p1.Active == p2.Active && p1.Comment == p2.Comment && p1.AllowWrites == p2.AllowWrites
	if attributesEqual && (len(p1.StorageLocations) == len(p2.StorageLocations)) {
		for i := range p1.StorageLocations {
			if !storageLocationsEqual(p1.StorageLocations[i], p2.StorageLocations[i]) {
				return false
			}
		}

		return true
	}
	return false
}

func ParseExternalVolumeDescribed(props []sdk.ExternalVolumeProperty) (ParsedExternalVolumeDescribed, error) {
	parsedExternalVolumeDescribed := ParsedExternalVolumeDescribed{}
	var storageLocations []StorageLocation
	for _, p := range props {
		switch {
		case p.Name == "COMMENT":
			parsedExternalVolumeDescribed.Comment = p.Value
		case p.Name == "ACTIVE":
			parsedExternalVolumeDescribed.Active = p.Value
		case p.Name == "ALLOW_WRITES":
			parsedExternalVolumeDescribed.AllowWrites = p.Value
		case strings.Contains(p.Name, "STORAGE_LOCATION_"):
			switch {
			case strings.Contains(p.Value, `"STORAGE_PROVIDER":"S3"`):
				s3StorageLocation := S3StorageLocation{}
				err := json.Unmarshal([]byte(p.Value), &s3StorageLocation)
				if err != nil {
					return ParsedExternalVolumeDescribed{}, err
				}
				storageLocation := StorageLocation{
					Name:                 s3StorageLocation.Name,
					StorageProvider:      s3StorageLocation.StorageProvider,
					StorageBaseUrl:       s3StorageLocation.StorageBaseUrl,
					StorageAwsRoleArn:    s3StorageLocation.StorageAwsRoleArn,
					StorageAwsExternalId: s3StorageLocation.StorageAwsExternalId,
					EncryptionType:       s3StorageLocation.EncryptionType,
					EncryptionKmsKeyId:   s3StorageLocation.EncryptionKmsId,
				}
				storageLocations = append(
					storageLocations,
					storageLocation,
				)
			case strings.Contains(p.Value, `"STORAGE_PROVIDER":"GCS"`):
				gcsStorageLocation := GCSStorageLocation{}
				err := json.Unmarshal([]byte(p.Value), &gcsStorageLocation)
				if err != nil {
					return ParsedExternalVolumeDescribed{}, err
				}

				storageLocation := StorageLocation{
					Name:               gcsStorageLocation.Name,
					StorageProvider:    gcsStorageLocation.StorageProvider,
					StorageBaseUrl:     gcsStorageLocation.StorageBaseUrl,
					EncryptionType:     gcsStorageLocation.EncryptionType,
					EncryptionKmsKeyId: gcsStorageLocation.EncryptionKmsId,
				}
				storageLocations = append(
					storageLocations,
					storageLocation,
				)
			case strings.Contains(p.Value, `"STORAGE_PROVIDER":"AZURE"`):
				azureStorageLocation := AzureStorageLocation{}
				err := json.Unmarshal([]byte(p.Value), &azureStorageLocation)
				if err != nil {
					return ParsedExternalVolumeDescribed{}, err
				}

				storageLocation := StorageLocation{
					Name:            azureStorageLocation.Name,
					StorageProvider: azureStorageLocation.StorageProvider,
					StorageBaseUrl:  azureStorageLocation.StorageBaseUrl,
					AzureTenantId:   azureStorageLocation.AzureTenantId,
				}
				storageLocations = append(
					storageLocations,
					storageLocation,
				)
			default:
				return ParsedExternalVolumeDescribed{}, fmt.Errorf("Unrecognized storage provider in storage location property: %s", p.Value)
			}
		default:
			return ParsedExternalVolumeDescribed{}, fmt.Errorf("Unrecognized external volume property: %s", p.Name)
		}
	}

	parsedExternalVolumeDescribed.StorageLocations = storageLocations
	validated := validateParsedExternalVolumeDescribed(parsedExternalVolumeDescribed)
	if validated != nil {
		return ParsedExternalVolumeDescribed{}, validated
	}

	return parsedExternalVolumeDescribed, nil
}

// Generate input to the ParseExternalVolumeDescribedInput, useful for testing purposes
func GenerateParseExternalVolumeDescribedInput(comment string, allowWrites string, storageLocations []string, active string) []sdk.ExternalVolumeProperty {
	storageLocationProperties := make([]sdk.ExternalVolumeProperty, len(storageLocations))
	allowWritesProperty := sdk.ExternalVolumeProperty{
		Parent:  "",
		Name:    "ALLOW_WRITES",
		Type:    "Boolean",
		Value:   allowWrites,
		Default: "true",
	}

	commentProperty := sdk.ExternalVolumeProperty{
		Parent:  "",
		Name:    "COMMENT",
		Type:    "String",
		Value:   comment,
		Default: "",
	}

	activeProperty := sdk.ExternalVolumeProperty{
		Parent:  "STORAGE_LOCATIONS",
		Name:    "ACTIVE",
		Type:    "String",
		Value:   active,
		Default: "",
	}

	for i, property := range storageLocations {
		storageLocationProperties[i] = sdk.ExternalVolumeProperty{
			Parent:  "STORAGE_LOCATIONS",
			Name:    fmt.Sprintf("STORAGE_LOCATION_%s", strconv.Itoa(i+1)),
			Type:    "String",
			Value:   property,
			Default: "",
		}
	}

	return append(append([]sdk.ExternalVolumeProperty{allowWritesProperty, commentProperty}, storageLocationProperties...), activeProperty)
}

// TODO(SNOW-1569530): address during identifiers rework follow-up
func ParseRootLocation(location string) (sdk.SchemaObjectIdentifier, string, error) {
	location = strings.TrimPrefix(location, "@")
	parts, err := sdk.ParseIdentifierStringWithOpts(location, func(r *csv.Reader) {
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

// ContainsIdentifierIgnoringQuotes takes ids (a slice of Snowflake identifiers represented as strings), and
// id (a string representing Snowflake id). It checks if id is contained within ids ignoring quotes around identifier parts.
//
// The original quoting should be retrieved to avoid situations like "object" == "\"object\"" (true)
// where that should not be a truthful comparison (different ids). Right now, we assume this case won't happen because the quoting difference would only appear
// in cases where the identifier parts are upper-cased and returned without quotes by snowflake, e.g. "OBJECT" == "\"OBJECT\"" (true)
// which is correct (the same ids).
func ContainsIdentifierIgnoringQuotes(ids []string, id string) bool {
	if len(ids) == 0 || len(id) == 0 {
		return false
	}

	idToCompare, err := DecodeSnowflakeParameterID(id)
	if err != nil {
		return false
	}

	for _, stringId := range ids {
		objectIdentifier, err := DecodeSnowflakeParameterID(stringId)
		if err != nil {
			return false
		}
		if idToCompare.FullyQualifiedName() == objectIdentifier.FullyQualifiedName() {
			return true
		}
	}

	return false
}
