package tracking

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

const (
	ProviderVersion string = "v0.98.0" // TODO(SNOW-1814934): Currently hardcoded, make it computed
	MetadataPrefix  string = "terraform_provider_usage_tracking"
)

type key int

var metadataContextKey key

type Operation string

const (
	CreateOperation     Operation = "create"
	ReadOperation       Operation = "read"
	UpdateOperation     Operation = "update"
	DeleteOperation     Operation = "delete"
	ImportOperation     Operation = "import"
	CustomDiffOperation Operation = "custom_diff"
)

type Metadata struct {
	Version   string                 `json:"version,omitempty"`
	Resource  resources.ResourceName `json:"resource,omitempty"`
	Operation Operation              `json:"operation,omitempty"`
}

func NewMetadata(version string, resourceName resources.ResourceName, operation Operation) Metadata {
	return Metadata{
		Version:   version,
		Resource:  resourceName,
		Operation: operation,
	}
}

func NewVersionedMetadata(resourceName resources.ResourceName, operation Operation) Metadata {
	return Metadata{
		Version:   ProviderVersion,
		Resource:  resourceName,
		Operation: operation,
	}
}

func NewContext(ctx context.Context, metadata Metadata) context.Context {
	return context.WithValue(ctx, metadataContextKey, metadata)
}

func FromContext(ctx context.Context) (Metadata, bool) {
	metadata, ok := ctx.Value(metadataContextKey).(Metadata)
	return metadata, ok
}
