package previewfeatures

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ensureResourceIsEnabled(featureRaw string, meta any) error {
	enabled := meta.(*provider.Context).EnabledFeatures
	feature, err := StringToFeature(featureRaw)
	if err != nil {
		return err
	}
	if err := EnsurePreviewFeatureEnabled(feature, enabled); err != nil {
		return err
	}
	return nil
}

func PreviewFeatureReadContextWrapper(featureRaw string, readFunc schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return readFunc(ctx, d, meta)
	}
}

func PreviewFeatureCreateContextWrapper(featureRaw string, createFunc schema.CreateContextFunc) schema.CreateContextFunc { //nolint
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return createFunc(ctx, d, meta)
	}
}

func PreviewFeatureUpdateContextWrapper(featureRaw string, updateFunc schema.UpdateContextFunc) schema.UpdateContextFunc { //nolint
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return updateFunc(ctx, d, meta)
	}
}

func PreviewFeatureDeleteContextWrapper(featureRaw string, deleteFunc schema.DeleteContextFunc) schema.DeleteContextFunc { //nolint
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return deleteFunc(ctx, d, meta)
	}
}
