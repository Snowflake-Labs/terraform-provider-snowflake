package testenvs

import (
	"fmt"
	"os"
	"testing"
)

type env string

const (
	BusinessCriticalAccount env = "SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"

	TestAccountCreate          env = "TEST_SF_TF_TEST_ACCOUNT_CREATE"
	TestFailoverGroups         env = "TEST_SF_TF_TEST_FAILOVER_GROUPS"
	ResourceMonitorNotifyUsers env = "TEST_SF_TF_RESOURCE_MONITOR_NOTIFY_USERS"

	AwsExternalBucketUrl   env = "TEST_SF_TF_AWS_EXTERNAL_BUCKET_URL"
	AwsExternalKeyId       env = "TEST_SF_TF_AWS_EXTERNAL_KEY_ID"
	AwsExternalSecretKey   env = "TEST_SF_TF_AWS_EXTERNAL_SECRET_KEY" // #nosec G101
	AwsExternalRoleArn     env = "TEST_SF_TF_AWS_EXTERNAL_ROLE_ARN"
	AzureExternalBucketUrl env = "TEST_SF_TF_AZURE_EXTERNAL_BUCKET_URL"
	AzureExternalTenantId  env = "TEST_SF_TF_AZURE_EXTERNAL_TENANT_ID"
	AzureExternalSasToken  env = "TEST_SF_TF_AZURE_EXTERNAL_SAS_TOKEN" // #nosec G101
	GcsExternalBuckerUrl   env = "TEST_SF_TF_GCS_EXTERNAL_BUCKET_URL"

	SkipManagedAccountTest  env = "TEST_SF_TF_SKIP_MANAGED_ACCOUNT_TEST"
	SkipSamlIntegrationTest env = "TEST_SF_TF_SKIP_SAML_INTEGRATION_TEST"

	EnableSweep         env = "TEST_SF_TF_ENABLE_SWEEP"
	ConfigureClientOnce env = "SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE"
)

func GetOrSkipTest(t *testing.T, envName Env) string {
	t.Helper()
	env := os.Getenv(fmt.Sprintf("%v", envName))
	if env == "" {
		t.Skipf("Skipping %s, env %v missing", t.Name(), envName)
	}
	return env
}

func SkipTestIfSet(t *testing.T, envName Env, reason string) {
	t.Helper()
	env := os.Getenv(fmt.Sprintf("%v", envName))
	if env != "" {
		t.Skipf("Skipping %s, because env %v is set. Reason: \"%s\"", t.Name(), envName, reason)
	}
}

type Env interface {
	xxxProtected()
}

func (e env) xxxProtected() {}
