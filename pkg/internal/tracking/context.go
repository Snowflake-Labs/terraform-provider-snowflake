package tracking

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

const (
	ProviderVersion string = "v0.99.0" // TODO(SNOW-1814934): Currently hardcoded, make it computed
	MetadataPrefix  string = "terraform_provider_usage_tracking"
)

type key struct{}

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
	Version   string    `json:"version,omitempty"`
	Resource  string    `json:"resource,omitempty"`
	Operation Operation `json:"operation,omitempty"`
}

func (m Metadata) validate() error {
	errs := make([]error, 0)
	if m.Version == "" {
		errs = append(errs, errors.New("version for metadata should not be empty"))
	}
	if m.Resource == "" {
		errs = append(errs, errors.New("resource name for metadata should not be empty"))
	}
	if m.Operation == "" {
		errs = append(errs, errors.New("operation for metadata should not be empty"))
	}
	return errors.Join(errs...)
}

func NewMetadata(version string, resource resources.Resource, operation Operation) Metadata {
	return Metadata{
		Version:   version,
		Resource:  resource.String(),
		Operation: operation,
	}
}

func NewVersionedMetadata(resource resources.Resource, operation Operation) Metadata {
	return Metadata{
		Version:   ProviderVersion,
		Resource:  resource.String(),
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
