package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal"
)

func RandomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
	t.Helper()
	return NewSchemaObjectIdentifier(internal.RandomStringN(t, 12), internal.RandomStringN(t, 12), internal.RandomStringN(t, 12))
}

func RandomDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
	t.Helper()
	return NewDatabaseObjectIdentifier(internal.RandomStringN(t, 12), internal.RandomStringN(t, 12))
}

func RandomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
	t.Helper()
	return NewAccountObjectIdentifier(internal.RandomStringN(t, 12))
}
