package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

var (
	invalidAccountObjectIdentifier = NewAccountObjectIdentifier(random.StringN(256))
	invalidSchemaObjectIdentifier  = NewSchemaObjectIdentifier(random.StringN(255), random.StringN(255), random.StringN(255))
)

func randomSchemaObjectIdentifier() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(random.StringN(12), random.StringN(12), random.StringN(12))
}

func randomSchemaObjectIdentifierInSchema(schemaId DatabaseObjectIdentifier) SchemaObjectIdentifier {
	return NewSchemaObjectIdentifierInSchema(schemaId, random.StringN(12))
}

func randomExternalObjectIdentifier() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator(random.StringN(12)), randomAccountObjectIdentifier())
}

func randomDatabaseObjectIdentifier() DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(random.StringN(12), random.StringN(12))
}

func randomDatabaseObjectIdentifierInDatabase(databaseId AccountObjectIdentifier) DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(databaseId.Name(), random.StringN(12))
}

func randomAccountObjectIdentifier() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(random.StringN(12))
}
