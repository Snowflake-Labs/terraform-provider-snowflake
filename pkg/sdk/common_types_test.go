package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStringProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toStringProperty()
		assert.Empty(t, prop.Value)
		assert.Empty(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property row containing values", func(t *testing.T) {
		row := &propertyRow{
			Value:        "value",
			DefaultValue: "default value",
			Description:  "desc",
		}
		prop := row.toStringProperty()
		assert.Equal(t, prop.Value, "value")
		assert.Equal(t, prop.DefaultValue, "default value")
		assert.Equal(t, prop.Description, row.Description)
	})
}

func TestToIntProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Empty(t, prop.Value)
		assert.Empty(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property row not containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "value",
			DefaultValue: "default value",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Empty(t, prop.Value)
		assert.Empty(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property not containing default value", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10",
			DefaultValue: "null",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Equal(t, *prop.Value, 10)
		assert.Empty(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property row containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Equal(t, *prop.Value, 10)
		assert.Equal(t, *prop.DefaultValue, 0)
		assert.Equal(t, prop.Description, row.Description)
	})
}

func TestToBoolProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toBoolProperty()
		assert.Empty(t, prop.Value)
		assert.Empty(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property row containing values", func(t *testing.T) {
		row := &propertyRow{
			Value:        "true",
			DefaultValue: "false",
			Description:  "desc",
		}
		prop := row.toBoolProperty()
		assert.Equal(t, prop.Value, true)
		assert.Equal(t, prop.DefaultValue, false)
		assert.Equal(t, prop.Description, row.Description)
	})
}
