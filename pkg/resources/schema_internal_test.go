package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id := "database_name|schema"
	schema, err := schemaIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", schema.DatabaseName)
	r.Equal("schema", schema.SchemaName)

	// Bad ID -- not enough fields
	id = "database"
	_, err = schemaIDFromString(id)
	r.Equal(fmt.Errorf("2 fields allowed"), err)

	// Bad ID
	id = "|"
	_, err = schemaIDFromString(id)
	r.NoError(err)

	// 0 lines
	id = ""
	_, err = schemaIDFromString(id)
	r.Equal(fmt.Errorf("1 line per schema"), err)

	// 2 lines
	id = `database_name|schema
	database_name|schema`
	_, err = schemaIDFromString(id)
	r.Equal(fmt.Errorf("1 line per schema"), err)
}

func TestSchemaStruct(t *testing.T) {
	r := require.New(t)

	// Vanilla
	schema := &schemaID{
		DatabaseName: "database_name",
		SchemaName:   "schema",
	}
	sID, err := schema.String()
	r.NoError(err)
	r.Equal("database_name|schema", sID)

	// Empty grant
	schema = &schemaID{}
	sID, err = schema.String()
	r.NoError(err)
	r.Equal("|", sID)

	// Grant with extra delimiters
	schema = &schemaID{
		DatabaseName: "database|name",
		SchemaName:   "schema|name",
	}
	sID, err = schema.String()
	r.NoError(err)
	newSchema, err := schemaIDFromString(sID)
	r.NoError(err)
	r.Equal("database|name", newSchema.DatabaseName)
	r.Equal("schema|name", newSchema.SchemaName)
}
