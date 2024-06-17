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

func TestToStorageSerializationPolicy(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected StorageSerializationPolicy
		Error    string
	}{
		{
			Name:     "optimized storage serialization policy",
			Input:    string(StorageSerializationPolicyOptimized),
			Expected: StorageSerializationPolicyOptimized,
		},
		{
			Name:     "compatible storage serialization policy",
			Input:    string(StorageSerializationPolicyCompatible),
			Expected: StorageSerializationPolicyCompatible,
		},
		{
			Name:  "validation: incorrect storage serialization policy",
			Input: "incorrect",
			Error: "unknown storage serialization policy: incorrect",
		},
		{
			Name:  "validation: empty input",
			Input: "",
			Error: "unknown storage serialization policy: ",
		},
		{
			Name:     "validation: lower case input",
			Input:    "optimized",
			Expected: StorageSerializationPolicyOptimized,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
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
		{
			Name:     "trace log level",
			Input:    string(LogLevelTrace),
			Expected: LogLevelTrace,
		},
		{
			Name:     "debug log level",
			Input:    string(LogLevelDebug),
			Expected: LogLevelDebug,
		},
		{
			Name:     "info log level",
			Input:    string(LogLevelInfo),
			Expected: LogLevelInfo,
		},
		{
			Name:     "warn log level",
			Input:    string(LogLevelWarn),
			Expected: LogLevelWarn,
		},
		{
			Name:     "error log level",
			Input:    string(LogLevelError),
			Expected: LogLevelError,
		},
		{
			Name:     "fatal log level",
			Input:    string(LogLevelFatal),
			Expected: LogLevelFatal,
		},
		{
			Name:     "off log level",
			Input:    string(LogLevelOff),
			Expected: LogLevelOff,
		},
		{
			Name:  "validation: incorrect log level",
			Input: "incorrect",
			Error: "unknown log level: incorrect",
		},
		{
			Name:  "validation: empty input",
			Input: "",
			Error: "unknown log level: ",
		},
		{
			Name:     "validation: lower case input",
			Input:    "info",
			Expected: LogLevelInfo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
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
		{
			Name:     "always trace level",
			Input:    string(TraceLevelAlways),
			Expected: TraceLevelAlways,
		},
		{
			Name:     "on event trace level",
			Input:    string(TraceLevelOnEvent),
			Expected: TraceLevelOnEvent,
		},
		{
			Name:     "off trace level",
			Input:    string(TraceLevelOff),
			Expected: TraceLevelOff,
		},
		{
			Name:  "validation: incorrect trace level",
			Input: "incorrect",
			Error: "unknown trace level: incorrect",
		},
		{
			Name:  "validation: empty input",
			Input: "",
			Error: "unknown trace level: ",
		},
		{
			Name:     "validation: lower case input",
			Input:    "always",
			Expected: TraceLevelAlways,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
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
