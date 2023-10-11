package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
)

func RandomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
	t.Helper()
	return NewSchemaObjectIdentifier(random.RandomStringN(t, 12), random.RandomStringN(t, 12), random.RandomStringN(t, 12))
}

func RandomDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
	t.Helper()
	return NewDatabaseObjectIdentifier(random.RandomStringN(t, 12), random.RandomStringN(t, 12))
}

func RandomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
	t.Helper()
	return NewAccountObjectIdentifier(random.RandomStringN(t, 12))
}
