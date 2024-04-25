package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

// Helper methods in this file are used both in SDK tests and also in integration tests.
// These methods are not supposed to be used in production code, so it would be better to not export them.
// Simply moving them to internal/random along other random helper methods is not an easy option, because it would create import cycle:
// sdk needed here -> sdk imports this to use the methods for tests.
//
// To move it there, we have to break a cycle. We have a few, slightly different options:
// 1. SDK reorganization
//    We can extract separate package for identifiers only, it could be imported both by sdk and by internal/random packages.
//    While this is an elegant solution, we are keeping all sdk files in one package now, and this would be exception.
// 2. Changing the package for unit tests in sdk package
//    This would be different from extracting integration tests because it would require only changing package to sdk_test without moving the file.
//    The reason for that is that go allows keeping x_test package together with x package inside one directory.
//    Then one downside (or not a downside?) would be, that we would be running our unit tests as a black box;
//    we would be restricted only to exported objects.
// 3. Using go:linkname compiler directive
//    This is a hacky solution which allows linking without importing a package.
//    We could e.g. replace the XxxIdentifier structs and NewXxxIdentifier methods with our types with such directive.
//    That way we would be able to use them without importing them directly.
//    As stated in an official documentation this is not recommended, unsafe method.
// 4. Do nothing
//    Just leave these helpers exported.

// UPDATE: with introduction of ids generation for integration and acceptance tests, this file will contain only methods used in sdk unit tests

func randomSchemaObjectIdentifier() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(random.StringN(12), random.StringN(12), random.StringN(12))
}

func randomExternalObjectIdentifier() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator(random.StringN(12)), RandomAccountObjectIdentifier())
}

func randomDatabaseObjectIdentifier() DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(random.StringN(12), random.StringN(12))
}

func RandomAccountObjectIdentifier() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(random.StringN(12))
}

func RandomAlphanumericAccountObjectIdentifier() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(random.AlphanumericN(12))
}
