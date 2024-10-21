package resources_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetPropertyAsPointer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"second_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"third_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"second_string": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"third_string": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"second_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"third_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}, map[string]interface{}{
		"integer":        123,
		"second_integer": 0,
		"string":         "some string",
		"second_string":  "",
		"boolean":        true,
		"second_boolean": false,
		"invalid":        true,
	})

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Equal(t, true, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[int](d, "second_integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[int](d, "third_integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Nil(t, resources.GetPropertyAsPointer[string](d, "second_integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[string](d, "third_string"))
	assert.Equal(t, true, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "second_boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "third_boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))
}

// TODO [SNOW-1511594]: provide TestResourceDataRaw with working GetRawConfig()
func Test_GetConfigPropertyAsPointerAllowingZeroValue(t *testing.T) {
	t.Skip("TestResourceDataRaw does not set up the ResourceData correctly - GetRawConfig is nil")
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"second_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"third_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"second_string": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"third_string": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"second_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"third_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}, map[string]interface{}{
		"integer":        123,
		"second_integer": 0,
		"string":         "some string",
		"second_string":  "",
		"boolean":        true,
		"second_boolean": false,
		"invalid":        true,
	})

	assert.Equal(t, 123, *resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "integer"))
	assert.Equal(t, 0, *resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "second_integer"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "third_integer"))
	assert.Equal(t, "some string", *resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "string"))
	assert.Equal(t, "", *resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "second_integer"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "third_string"))
	assert.Equal(t, true, *resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "boolean"))
	assert.Equal(t, false, *resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "second_boolean"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "third_boolean"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "invalid"))
}

// queriedAccountRolePrivilegesEqualTo will check if all the privileges specified in the argument are granted in Snowflake.
func queriedPrivilegesEqualTo(query func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error), privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		grants, err := query(client, ctx)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			if (grant.GrantTo == sdk.ObjectTypeDatabaseRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole) && grant.Privilege == "USAGE" {
				continue
			}
			if !slices.Contains(privileges, grant.Privilege) {
				return fmt.Errorf("grant not expected, grant: %v, not in %v", grants, privileges)
			}
		}

		return nil
	}
}

// queriedAccountRolePrivilegesContainAtLeast will check if all the privileges specified in the argument are granted in Snowflake.
// Any additional grants will be ignored.
func queriedPrivilegesContainAtLeast(query func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error), roleName sdk.ObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		grants, err := query(client, ctx)
		if err != nil {
			return err
		}
		var grantedPrivileges []string
		for _, grant := range grants {
			grantedPrivileges = append(grantedPrivileges, grant.Privilege)
		}
		notAllPrivilegesInGrantedPrivileges := slices.ContainsFunc(privileges, func(privilege string) bool {
			return !slices.Contains(grantedPrivileges, privilege)
		})
		if notAllPrivilegesInGrantedPrivileges {
			return fmt.Errorf("not every privilege from the list: %v was found in grant privileges: %v, for role name: %s", privileges, grantedPrivileges, roleName.FullyQualifiedName())
		}

		return nil
	}
}

func TestListDiff(t *testing.T) {
	testCases := []struct {
		Name    string
		Before  []any
		After   []any
		Added   []any
		Removed []any
	}{
		{
			Name:    "no changes",
			Before:  []any{1, 2, 3, 4},
			After:   []any{1, 2, 3, 4},
			Removed: []any{},
			Added:   []any{},
		},
		{
			Name:    "only removed",
			Before:  []any{1, 2, 3, 4},
			After:   []any{},
			Removed: []any{1, 2, 3, 4},
			Added:   []any{},
		},
		{
			Name:    "only added",
			Before:  []any{},
			After:   []any{1, 2, 3, 4},
			Removed: []any{},
			Added:   []any{1, 2, 3, 4},
		},
		{
			Name:    "added repeated items",
			Before:  []any{2},
			After:   []any{1, 2, 1},
			Removed: []any{},
			Added:   []any{1, 1},
		},
		{
			Name:    "removed repeated items",
			Before:  []any{1, 2, 1},
			After:   []any{2},
			Removed: []any{1, 1},
			Added:   []any{},
		},
		{
			Name:    "simple diff: ints",
			Before:  []any{1, 2, 3, 4, 5, 6, 7, 8, 9},
			After:   []any{1, 3, 5, 7, 9, 12, 13, 14},
			Removed: []any{2, 4, 6, 8},
			Added:   []any{12, 13, 14},
		},
		{
			Name:    "simple diff: strings",
			Before:  []any{"one", "two", "three", "four"},
			After:   []any{"five", "two", "four", "six"},
			Removed: []any{"one", "three"},
			Added:   []any{"five", "six"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			added, removed := resources.ListDiff(tc.Before, tc.After)
			assert.Equal(t, tc.Added, added)
			assert.Equal(t, tc.Removed, removed)
		})
	}
}

func Test_DataTypeIssue3007DiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{
			name:     "different data type",
			old:      string(sdk.DataTypeVARCHAR),
			new:      string(sdk.DataTypeNumber),
			expected: false,
		},
		{
			name:     "same number data type without arguments",
			old:      string(sdk.DataTypeNumber),
			new:      string(sdk.DataTypeNumber),
			expected: true,
		},
		{
			name:     "same number data type different casing",
			old:      string(sdk.DataTypeNumber),
			new:      "number",
			expected: true,
		},
		{
			name:     "same text data type without arguments",
			old:      string(sdk.DataTypeVARCHAR),
			new:      string(sdk.DataTypeVARCHAR),
			expected: true,
		},
		{
			name:     "same other data type",
			old:      string(sdk.DataTypeFloat),
			new:      string(sdk.DataTypeFloat),
			expected: true,
		},
		{
			name:     "synonym number data type without arguments",
			old:      string(sdk.DataTypeNumber),
			new:      "DECIMAL",
			expected: true,
		},
		{
			name:     "synonym text data type without arguments",
			old:      string(sdk.DataTypeVARCHAR),
			new:      "TEXT",
			expected: true,
		},
		{
			name:     "synonym other data type without arguments",
			old:      string(sdk.DataTypeFloat),
			new:      "DOUBLE",
			expected: true,
		},
		{
			name:     "synonym number data type same precision, no scale",
			old:      "NUMBER(30)",
			new:      "DECIMAL(30)",
			expected: true,
		},
		{
			name:     "synonym number data type precision implicit and same",
			old:      "NUMBER",
			new:      fmt.Sprintf("DECIMAL(%d)", sdk.DefaultNumberPrecision),
			expected: true,
		},
		{
			name:     "synonym number data type precision implicit and different",
			old:      "NUMBER",
			new:      "DECIMAL(30)",
			expected: false,
		},
		{
			name:     "number data type different precisions, no scale",
			old:      "NUMBER(35)",
			new:      "NUMBER(30)",
			expected: false,
		},
		{
			name:     "synonym number data type same precision, different scale",
			old:      "NUMBER(30, 2)",
			new:      "DECIMAL(30, 1)",
			expected: false,
		},
		{
			name:     "synonym number data type default scale implicit and explicit",
			old:      "NUMBER(30)",
			new:      fmt.Sprintf("DECIMAL(30, %d)", sdk.DefaultNumberScale),
			expected: true,
		},
		{
			name:     "synonym number data type default scale implicit and different",
			old:      "NUMBER(30)",
			new:      "DECIMAL(30, 3)",
			expected: false,
		},
		{
			name:     "synonym number data type both precision and scale implicit and explicit",
			old:      "NUMBER",
			new:      fmt.Sprintf("DECIMAL(%d, %d)", sdk.DefaultNumberPrecision, sdk.DefaultNumberScale),
			expected: true,
		},
		{
			name:     "synonym number data type both precision and scale implicit and scale different",
			old:      "NUMBER",
			new:      fmt.Sprintf("DECIMAL(%d, 2)", sdk.DefaultNumberPrecision),
			expected: false,
		},
		{
			name:     "synonym text data type same length",
			old:      "VARCHAR(30)",
			new:      "TEXT(30)",
			expected: true,
		},
		{
			name:     "synonym text data type different length",
			old:      "VARCHAR(30)",
			new:      "TEXT(40)",
			expected: false,
		},
		{
			name:     "synonym text data type length implicit and same",
			old:      "VARCHAR",
			new:      fmt.Sprintf("TEXT(%d)", sdk.DefaultVarcharLength),
			expected: true,
		},
		{
			name:     "synonym text data type length implicit and different",
			old:      "VARCHAR",
			new:      "TEXT(40)",
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := resources.DataTypeIssue3007DiffSuppressFunc("", tc.old, tc.new, nil)
			require.Equal(t, tc.expected, result)
		})
	}
}

// External volume helper tests

func Test_GetStorageLocationName(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := sdk.S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      sdk.S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &sdk.ExternalVolumeS3Encryption{
			Type:     sdk.S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := sdk.AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := sdk.GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &sdk.ExternalVolumeGCSEncryption{
			Type:     sdk.GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := sdk.S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   sdk.S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}

	testCases := []struct {
		Name            string
		StorageLocation sdk.ExternalVolumeStorageLocation
		ExpectedName    string
	}{
		{
			Name:            "S3 storage location name succesfully read",
			StorageLocation: sdk.ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA},
			ExpectedName:    s3StorageLocationA.Name,
		},
		{
			Name:            "S3GOV storage location name succesfully read",
			StorageLocation: sdk.ExternalVolumeStorageLocation{S3StorageLocationParams: &s3GovStorageLocationA},
			ExpectedName:    s3GovStorageLocationA.Name,
		},
		{
			Name:            "GCS storage location name succesfully read",
			StorageLocation: sdk.ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA},
			ExpectedName:    gcsStorageLocationA.Name,
		},
		{
			Name:            "Azure storage location name succesfully read",
			StorageLocation: sdk.ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA},
			ExpectedName:    azureStorageLocationA.Name,
		},
	}

	invalidTestCases := []struct {
		Name            string
		StorageLocation sdk.ExternalVolumeStorageLocation
	}{
		{
			Name:            "Empty S3 storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{S3StorageLocationParams: &sdk.S3StorageLocationParams{}},
		},
		{
			Name:            "Empty GCS storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{GCSStorageLocationParams: &sdk.GCSStorageLocationParams{}},
		},
		{
			Name:            "Empty Azure storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{AzureStorageLocationParams: &sdk.AzureStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			name, err := resources.GetStorageLocationName(tc.StorageLocation)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedName, name)
		})
	}
	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := resources.GetStorageLocationName(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

func Test_GetStorageLocationStorageProvider(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := sdk.S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      sdk.S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &sdk.ExternalVolumeS3Encryption{
			Type:     sdk.S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := sdk.AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := sdk.GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &sdk.ExternalVolumeGCSEncryption{
			Type:     sdk.GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := sdk.S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   sdk.S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}
	testCases := []struct {
		Name                    string
		StorageLocation         sdk.ExternalVolumeStorageLocation
		ExpectedStorageProvider sdk.StorageProvider
	}{
		{
			Name:                    "S3 storage provider",
			StorageLocation:         sdk.ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA},
			ExpectedStorageProvider: sdk.StorageProviderS3,
		},
		{
			Name:                    "S3GOV storage provider",
			StorageLocation:         sdk.ExternalVolumeStorageLocation{S3StorageLocationParams: &s3GovStorageLocationA},
			ExpectedStorageProvider: sdk.StorageProviderS3GOV,
		},
		{
			Name:                    "GCS storage provider",
			StorageLocation:         sdk.ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA},
			ExpectedStorageProvider: sdk.StorageProviderGCS,
		},
		{
			Name:                    "Azure storage provider",
			StorageLocation:         sdk.ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA},
			ExpectedStorageProvider: sdk.StorageProviderAzure,
		},
	}

	invalidTestCases := []struct {
		Name            string
		StorageLocation sdk.ExternalVolumeStorageLocation
	}{
		{
			Name:            "Empty S3 storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{S3StorageLocationParams: &sdk.S3StorageLocationParams{}},
		},
		{
			Name:            "Empty GCS storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{GCSStorageLocationParams: &sdk.GCSStorageLocationParams{}},
		},
		{
			Name:            "Empty Azure storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{AzureStorageLocationParams: &sdk.AzureStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			storageProvider, err := resources.GetStorageLocationStorageProvider(tc.StorageLocation)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedStorageProvider, storageProvider)
		})
	}
	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := resources.GetStorageLocationName(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

var s3StorageAwsExternalId = "1234567890"

func Test_CopyStorageLocationWithTempName(t *testing.T) {
	tempStorageLocationName := "terraform_provider_sentinel_storage_location"
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := sdk.S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      sdk.S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &sdk.ExternalVolumeS3Encryption{
			Type:     sdk.S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := sdk.AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := sdk.GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &sdk.ExternalVolumeGCSEncryption{
			Type:     sdk.GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	t.Run("S3 storage location", func(t *testing.T) {
		storageLocationInput := sdk.ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA}
		copiedStorageLocation, err := resources.CopyStorageLocationWithTempName(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageProvider, s3StorageLocationA.StorageProvider)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageBaseUrl, s3StorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsRoleArn, s3StorageLocationA.StorageAwsRoleArn)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsExternalId, s3StorageLocationA.StorageAwsExternalId)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.Encryption.Type, s3StorageLocationA.Encryption.Type)
		assert.Equal(t, *copiedStorageLocation.S3StorageLocationParams.Encryption.KmsKeyId, *s3StorageLocationA.Encryption.KmsKeyId)
	})

	t.Run("GCS storage location", func(t *testing.T) {
		storageLocationInput := sdk.ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA}
		copiedStorageLocation, err := resources.CopyStorageLocationWithTempName(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.StorageBaseUrl, gcsStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.Encryption.Type, gcsStorageLocationA.Encryption.Type)
		assert.Equal(t, *copiedStorageLocation.GCSStorageLocationParams.Encryption.KmsKeyId, *gcsStorageLocationA.Encryption.KmsKeyId)
	})

	t.Run("Azure storage location", func(t *testing.T) {
		storageLocationInput := sdk.ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA}
		copiedStorageLocation, err := resources.CopyStorageLocationWithTempName(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.StorageBaseUrl, azureStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.AzureTenantId, azureStorageLocationA.AzureTenantId)
	})

	invalidTestCases := []struct {
		Name            string
		StorageLocation sdk.ExternalVolumeStorageLocation
	}{
		{
			Name:            "Empty S3 storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{S3StorageLocationParams: &sdk.S3StorageLocationParams{}},
		},
		{
			Name:            "Empty GCS storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{GCSStorageLocationParams: &sdk.GCSStorageLocationParams{}},
		},
		{
			Name:            "Empty Azure storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{AzureStorageLocationParams: &sdk.AzureStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: sdk.ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := resources.CopyStorageLocationWithTempName(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

func Test_CommonPrefixLastIndex(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageLocationName2 := "s3Test2"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageLocationName2 := "gcsTest2"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageLocationName2 := "azureTest2"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := sdk.S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      sdk.S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &sdk.ExternalVolumeS3Encryption{
			Type:     sdk.S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	s3StorageLocationB := sdk.S3StorageLocationParams{
		Name:                 s3StorageLocationName2,
		StorageProvider:      sdk.S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &sdk.ExternalVolumeS3Encryption{
			Type:     sdk.S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := sdk.AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	azureStorageLocationB := sdk.AzureStorageLocationParams{
		Name:           azureStorageLocationName2,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := sdk.GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &sdk.ExternalVolumeGCSEncryption{
			Type:     sdk.GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	gcsStorageLocationB := sdk.GCSStorageLocationParams{
		Name:           gcsStorageLocationName2,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &sdk.ExternalVolumeGCSEncryption{
			Type:     sdk.GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	gcsStorageLocationC := sdk.GCSStorageLocationParams{
		Name:           "test",
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &sdk.ExternalVolumeGCSEncryption{
			Type:     sdk.GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := sdk.S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   sdk.S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}

	testCases := []struct {
		Name           string
		ListA          []sdk.ExternalVolumeStorageLocation
		ListB          []sdk.ExternalVolumeStorageLocation
		ExpectedOutput int
	}{
		{
			Name:           "Two empty lists",
			ListA:          []sdk.ExternalVolumeStorageLocation{},
			ListB:          []sdk.ExternalVolumeStorageLocation{},
			ExpectedOutput: -1,
		},
		{
			Name:           "First list empty",
			ListA:          []sdk.ExternalVolumeStorageLocation{},
			ListB:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Second list empty",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{},
			ExpectedOutput: -1,
		},
		{
			Name:           "Lists with no common prefix - length 1",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationB}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Lists with no common prefix - length 2",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationB}, {AzureStorageLocationParams: &azureStorageLocationB}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Identical lists - length 1",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ExpectedOutput: 0,
		},
		{
			Name:           "Identical lists - length 2",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ExpectedOutput: 1,
		},
		{
			Name: "Identical lists - length 3",
			ListA: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3GovStorageLocationA},
			},
			ListB: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3GovStorageLocationA},
			},
			ExpectedOutput: 2,
		},
		{
			Name: "Lists with a common prefix - length 3, matching up to and including index 1",
			ListA: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
			},
			ListB: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - length 4, matching up to and including index 2",
			ListA: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
			},
			ListB: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ExpectedOutput: 2,
		},
		{
			Name: "Lists with a common prefix - length 4, matching up to and including index 1",
			ListA: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ListB: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - different lengths, matching up to and including index 1 (last index of shorter list)",
			ListA: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
			},
			ListB: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - different lengths, matching up to and including index 2",
			ListA: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3StorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{AzureStorageLocationParams: &azureStorageLocationB},
			},
			ListB: []sdk.ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3StorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{AzureStorageLocationParams: &azureStorageLocationB},
			},
			ExpectedOutput: 2,
		},
		{
			Name:           "Empty S3 storage location",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &sdk.S3StorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty GCS storage location",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{{GCSStorageLocationParams: &sdk.GCSStorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty Azure storage location",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{{AzureStorageLocationParams: &sdk.AzureStorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty storage location",
			ListA:          []sdk.ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []sdk.ExternalVolumeStorageLocation{{}},
			ExpectedOutput: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			commonPrefixLastIndex, err := resources.CommonPrefixLastIndex(tc.ListA, tc.ListB)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedOutput, commonPrefixLastIndex)
		})
	}
}
