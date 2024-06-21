package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

var ShowDatabaseParametersSchema = map[string]*schema.Schema{
	strings.ToLower(string(sdk.ObjectParameterDataRetentionTimeInDays)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterMaxDataExtensionTimeInDays)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterSuspendTaskAfterNumFailures)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterTaskAutoRetryAttempts)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterUserTaskTimeoutMs)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterReplaceInvalidCharacters)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterQuotedIdentifiersIgnoreCase)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterEnableConsoleOutput)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterExternalVolume)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterCatalog)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterDefaultDDLCollation)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterStorageSerializationPolicy)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterLogLevel)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterTraceLevel)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	strings.ToLower(string(sdk.ObjectParameterUserTaskManagedInitialWarehouseSize)): {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
}

// TODO: New parameters had to be added in 4 places (schema, database validations, and here in 2 places)

func DatabaseParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	databaseParameters := make(map[string]any)
	for _, param := range parameters {
		parameterSchema := ParameterToSchema(param)
		switch strings.ToUpper(param.Key) {
		case string(sdk.ObjectParameterDataRetentionTimeInDays):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterDataRetentionTimeInDays))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterMaxDataExtensionTimeInDays):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterMaxDataExtensionTimeInDays))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterSuspendTaskAfterNumFailures):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterSuspendTaskAfterNumFailures))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterTaskAutoRetryAttempts):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterTaskAutoRetryAttempts))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterUserTaskTimeoutMs):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterUserTaskTimeoutMs))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterReplaceInvalidCharacters):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterReplaceInvalidCharacters))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterQuotedIdentifiersIgnoreCase):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterQuotedIdentifiersIgnoreCase))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterEnableConsoleOutput):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterEnableConsoleOutput))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterExternalVolume):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterExternalVolume))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterCatalog):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterCatalog))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterDefaultDDLCollation):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterDefaultDDLCollation))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterStorageSerializationPolicy):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterStorageSerializationPolicy))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterLogLevel):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterLogLevel))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterTraceLevel):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterTraceLevel))] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterUserTaskManagedInitialWarehouseSize):
			databaseParameters[strings.ToLower(string(sdk.ObjectParameterUserTaskManagedInitialWarehouseSize))] = []map[string]any{parameterSchema}
		}
	}
	return databaseParameters
}
