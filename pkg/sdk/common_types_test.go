package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ToStringProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toStringProperty()
		assert.Empty(t, prop.Value)
		assert.Empty(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property row containing values", func(t *testing.T) {
		row := &propertyRow{
			Value:        "value",
			DefaultValue: "default value",
			Description:  "desc",
		}
		prop := row.toStringProperty()
		assert.Equal(t, "value", prop.Value)
		assert.Equal(t, "default value", prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})
}

func Test_ToIntProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Nil(t, prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property row not containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "value",
			DefaultValue: "default value",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Nil(t, prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property not containing default value", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10",
			DefaultValue: "null",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Equal(t, 10, *prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property row containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Equal(t, 10, *prop.Value)
		assert.Equal(t, 0, *prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with negative value row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "-1",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		require.NotNil(t, prop.Value)
		assert.Equal(t, -1, *prop.Value)
	})

	t.Run("with decimal part value - not parsed correctly", func(t *testing.T) {
		row := &propertyRow{
			Value:        "0.85",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		require.Nil(t, prop.Value)
	})
}

func Test_ToBoolProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toBoolProperty()
		assert.False(t, prop.Value)
		assert.False(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property row containing values", func(t *testing.T) {
		row := &propertyRow{
			Value:        "true",
			DefaultValue: "false",
			Description:  "desc",
		}
		prop := row.toBoolProperty()
		assert.Equal(t, true, prop.Value)
		assert.Equal(t, false, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})
}

func Test_ToFloatProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "null",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.Nil(t, prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property row not containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "value",
			DefaultValue: "default value",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.Nil(t, prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property not containing default value", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10.5",
			DefaultValue: "null",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.Equal(t, 10.5, *prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property row containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10.1",
			DefaultValue: "10.5",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.Equal(t, 10.1, *prop.Value)
		assert.Equal(t, 10.5, *prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with negative value row and zero", func(t *testing.T) {
		row := &propertyRow{
			Value:        "-1.0",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.Equal(t, float64(-1), *prop.Value)
		assert.Equal(t, float64(0), *prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})
}

func TestToStorageSerializationPolicy(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected StorageSerializationPolicy
		Error    string
	}{
		{Input: string(StorageSerializationPolicyOptimized), Expected: StorageSerializationPolicyOptimized},
		{Input: string(StorageSerializationPolicyCompatible), Expected: StorageSerializationPolicyCompatible},
		{Name: "validation: incorrect storage serialization policy", Input: "incorrect", Error: "unknown storage serialization policy: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown storage serialization policy: "},
		{Name: "validation: lower case input", Input: "optimized", Expected: StorageSerializationPolicyOptimized},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v storage serialization policy", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToStorageSerializationPolicy(testCase.Input)
			if testCase.Error != "" {
				assert.Empty(t, value)
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.Expected, value)
			}
		})
	}
}

func TestToLogLevel(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected LogLevel
		Error    string
	}{
		{Input: string(LogLevelTrace), Expected: LogLevelTrace},
		{Input: string(LogLevelDebug), Expected: LogLevelDebug},
		{Input: string(LogLevelInfo), Expected: LogLevelInfo},
		{Input: string(LogLevelWarn), Expected: LogLevelWarn},
		{Input: string(LogLevelError), Expected: LogLevelError},
		{Input: string(LogLevelFatal), Expected: LogLevelFatal},
		{Input: string(LogLevelOff), Expected: LogLevelOff},
		{Name: "validation: incorrect log level", Input: "incorrect", Error: "unknown log level: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown log level: "},
		{Name: "validation: lower case input", Input: "info", Expected: LogLevelInfo},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v log level", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToLogLevel(testCase.Input)
			if testCase.Error != "" {
				assert.Empty(t, value)
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.Expected, value)
			}
		})
	}
}

func TestToTraceLevel(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected TraceLevel
		Error    string
	}{
		{Input: string(TraceLevelAlways), Expected: TraceLevelAlways},
		{Input: string(TraceLevelOnEvent), Expected: TraceLevelOnEvent},
		{Input: string(TraceLevelOff), Expected: TraceLevelOff},
		{Name: "validation: incorrect trace level", Input: "incorrect", Error: "unknown trace level: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown trace level: "},
		{Name: "validation: lower case input", Input: "always", Expected: TraceLevelAlways},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v trace level", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToTraceLevel(testCase.Input)
			if testCase.Error != "" {
				assert.Empty(t, value)
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.Expected, value)
			}
		})
	}
}
