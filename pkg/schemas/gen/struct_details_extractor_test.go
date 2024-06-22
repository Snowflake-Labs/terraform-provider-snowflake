package gen

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

// TODO: test completely new struct with:
//   - slice (string, enum)
//   - slice (identifier)
//   - sdk.FileFormatTypeOptions (sql builder struct)
//   - sdk.ObjectIdentifier (interface)
//   - (?) slice of pointers
//   - (?) pointer to slice
func Test_ExtractStructDetails(t *testing.T) {

	type testStruct struct {
		unexportedString     string
		unexportedInt        int
		unexportedBool       bool
		unexportedFloat64    float64
		unexportedStringPtr  *string
		unexportedIntPtr     *int
		unexportedBoolPtr    *bool
		unexportedFloat64Ptr *float64

		unexportedTime    time.Time
		unexportedTimePtr *time.Time

		unexportedStringEnum    sdk.WarehouseType
		unexportedStringEnumPtr *sdk.WarehouseType
		unexportedIntEnum       sdk.ResourceMonitorLevel
		unexportedIntEnumPtr    *sdk.ResourceMonitorLevel

		unexportedAccountIdentifier           sdk.AccountIdentifier
		unexportedExternalObjectIdentifier    sdk.ExternalObjectIdentifier
		unexportedAccountObjectIdentifier     sdk.AccountObjectIdentifier
		unexportedDatabaseObjectIdentifier    sdk.DatabaseObjectIdentifier
		unexportedSchemaObjectIdentifier      sdk.SchemaObjectIdentifier
		unexportedTableColumnIdentifier       sdk.TableColumnIdentifier
		unexportedAccountIdentifierPtr        *sdk.AccountIdentifier
		unexportedExternalObjectIdentifierPtr *sdk.ExternalObjectIdentifier
		unexportedAccountObjectIdentifierPtr  *sdk.AccountObjectIdentifier
		unexportedDatabaseObjectIdentifierPtr *sdk.DatabaseObjectIdentifier
		unexportedSchemaObjectIdentifierPtr   *sdk.SchemaObjectIdentifier
		unexportedTableColumnIdentifierPtr    *sdk.TableColumnIdentifier

		ExportedString     string
		ExportedInt        int
		ExportedBool       bool
		ExportedFloat64    float64
		ExportedStringPtr  *string
		ExportedIntPtr     *int
		ExportedBoolPtr    *bool
		ExportedFloat64Ptr *float64
	}

	assertFieldExtracted := func(field Field, expectedName string, expectedConcreteType string, expectedUnderlyingType string) {
		assert.Equal(t, expectedName, field.Name)
		assert.Equal(t, expectedConcreteType, field.ConcreteType)
		assert.Equal(t, expectedUnderlyingType, field.UnderlyingType)
	}

	t.Run("test struct details extraction", func(t *testing.T) {
		structDetails := ExtractStructDetails(testStruct{})

		assert.Equal(t, structDetails.Name, "gen.testStruct")

		assertFieldExtracted(structDetails.Fields[0], "unexportedString", "string", "string")
		assertFieldExtracted(structDetails.Fields[1], "unexportedInt", "int", "int")
		assertFieldExtracted(structDetails.Fields[2], "unexportedBool", "bool", "bool")
		assertFieldExtracted(structDetails.Fields[3], "unexportedFloat64", "float64", "float64")
		assertFieldExtracted(structDetails.Fields[4], "unexportedStringPtr", "*string", "*string")
		assertFieldExtracted(structDetails.Fields[5], "unexportedIntPtr", "*int", "*int")
		assertFieldExtracted(structDetails.Fields[6], "unexportedBoolPtr", "*bool", "*bool")
		assertFieldExtracted(structDetails.Fields[7], "unexportedFloat64Ptr", "*float64", "*float64")

		assertFieldExtracted(structDetails.Fields[8], "unexportedTime", "time.Time", "struct")
		assertFieldExtracted(structDetails.Fields[9], "unexportedTimePtr", "*time.Time", "*struct")

		assertFieldExtracted(structDetails.Fields[10], "unexportedStringEnum", "sdk.WarehouseType", "string")
		assertFieldExtracted(structDetails.Fields[11], "unexportedStringEnumPtr", "*sdk.WarehouseType", "*string")
		assertFieldExtracted(structDetails.Fields[12], "unexportedIntEnum", "sdk.ResourceMonitorLevel", "int")
		assertFieldExtracted(structDetails.Fields[13], "unexportedIntEnumPtr", "*sdk.ResourceMonitorLevel", "*int")

		assertFieldExtracted(structDetails.Fields[14], "unexportedAccountIdentifier", "sdk.AccountIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[15], "unexportedExternalObjectIdentifier", "sdk.ExternalObjectIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[16], "unexportedAccountObjectIdentifier", "sdk.AccountObjectIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[17], "unexportedDatabaseObjectIdentifier", "sdk.DatabaseObjectIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[18], "unexportedSchemaObjectIdentifier", "sdk.SchemaObjectIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[19], "unexportedTableColumnIdentifier", "sdk.TableColumnIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[20], "unexportedAccountIdentifierPtr", "*sdk.AccountIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[21], "unexportedExternalObjectIdentifierPtr", "*sdk.ExternalObjectIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[22], "unexportedAccountObjectIdentifierPtr", "*sdk.AccountObjectIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[23], "unexportedDatabaseObjectIdentifierPtr", "*sdk.DatabaseObjectIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[24], "unexportedSchemaObjectIdentifierPtr", "*sdk.SchemaObjectIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[25], "unexportedTableColumnIdentifierPtr", "*sdk.TableColumnIdentifier", "*struct")
	})
}
