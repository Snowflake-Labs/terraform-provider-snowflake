package config

type ReplacementPlaceholder string

const (
	SnowflakeProviderConfigNull                      ReplacementPlaceholder = "SF_TF_TEST_NULL_PLACEHOLDER"
	SnowflakeProviderConfigMultilineMarker           ReplacementPlaceholder = "SF_TF_TEST_MULTILINE_MARKER_PLACEHOLDER"
	SnowflakeProviderConfigSingleAttributeWorkaround ReplacementPlaceholder = "SF_TF_TEST_SINGLE_ATTRIBUTE_WORKAROUND"
)
