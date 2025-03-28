package resources

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FullyQualifiedNameAttributeName = "fully_qualified_name"
	AtAttributeName                 = "at"
	BeforeAttributeName             = "before"
)

// DiffSuppressStatement will suppress diffs between statements if they differ in only case or in
// runs of whitespace (\s+ = \s). This is needed because the snowflake api does not faithfully
// round-trip queries, so we cannot do a simple character-wise comparison to detect changes.
//
// Warnings: We will have false positives in cases where a change in case or run of whitespace is
// semantically significant.
//
// If we can find a sql parser that can handle the snowflake dialect then we should switch to parsing
// queries and either comparing ASTs or emitting a canonical serialization for comparison. I couldn't
// find such a library.
func DiffSuppressStatement(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(normalizeQuery(old), normalizeQuery(new))
}

var space = regexp.MustCompile(`\s+`)

func normalizeQuery(str string) string {
	return strings.TrimSpace(space.ReplaceAllString(str, " "))
}

// TODO [SNOW-1325214]: address during stage resource rework
func suppressQuoting(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return false
	} else {
		oldWithoutQuotes := strings.ReplaceAll(oldValue, "'", "")
		newWithoutQuotes := strings.ReplaceAll(newValue, "'", "")
		return oldWithoutQuotes == newWithoutQuotes
	}
}

func ctyValToSliceString(valueElems []cty.Value) []string {
	elems := make([]string, len(valueElems))
	for i, v := range valueElems {
		elems[i] = v.AsString()
	}
	return elems
}

func ImportName[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.AccountIdentifier](ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	switch any(new(T)).(type) {
	case *sdk.AccountObjectIdentifier:
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}
	case *sdk.DatabaseObjectIdentifier:
		id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}

		if err := d.Set("database", id.DatabaseName()); err != nil {
			return nil, err
		}
	case *sdk.SchemaObjectIdentifier:
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}

		if err := d.Set("database", id.DatabaseName()); err != nil {
			return nil, err
		}

		if err := d.Set("schema", id.SchemaName()); err != nil {
			return nil, err
		}
	case *sdk.AccountIdentifier:
		id, err := sdk.ParseAccountIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.AccountName()); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func TrackingImportWrapper(resourceName resources.Resource, importImplementation schema.StateContextFunc) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.ImportOperation))
		return importImplementation(ctx, d, meta)
	}
}

func TrackingCreateWrapper(resourceName resources.Resource, createImplementation schema.CreateContextFunc) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.CreateOperation))
		return createImplementation(ctx, d, meta)
	}
}

func TrackingReadWrapper(resourceName resources.Resource, readImplementation schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.ReadOperation))
		return readImplementation(ctx, d, meta)
	}
}

func TrackingUpdateWrapper(resourceName resources.Resource, updateImplementation schema.UpdateContextFunc) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.UpdateOperation))
		return updateImplementation(ctx, d, meta)
	}
}

func TrackingDeleteWrapper(resourceName resources.Resource, deleteImplementation schema.DeleteContextFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.DeleteOperation))
		return deleteImplementation(ctx, d, meta)
	}
}

func TrackingCustomDiffWrapper(resourceName resources.Resource, customdiffImplementation schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.CustomDiffOperation))
		return customdiffImplementation(ctx, diff, meta)
	}
}

func ensureResourceIsEnabled(featureRaw string, meta any) error {
	enabled := meta.(*provider.Context).EnabledFeatures
	feature, err := previewfeatures.StringToFeature(featureRaw)
	if err != nil {
		return err
	}
	if err := previewfeatures.EnsurePreviewFeatureEnabled(feature, enabled); err != nil {
		return err
	}
	return nil
}

func PreviewFeatureCreateContextWrapper(featureRaw string, createFunc schema.CreateContextFunc) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return createFunc(ctx, d, meta)
	}
}

func PreviewFeatureReadContextWrapper(featureRaw string, readFunc schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return readFunc(ctx, d, meta)
	}
}

func PreviewFeatureUpdateContextWrapper(featureRaw string, updateFunc schema.UpdateContextFunc) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return updateFunc(ctx, d, meta)
	}
}

func PreviewFeatureDeleteContextWrapper(featureRaw string, deleteFunc schema.DeleteContextFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return deleteFunc(ctx, d, meta)
	}
}

func SafeShowById[T any, ID sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.SchemaObjectIdentifierWithArguments](
	resourceName resources.Resource,
	client *sdk.Client,
	showById func(context.Context, ID) (T, error),
	ctx context.Context,
	id ID,
) (T, bool, diag.Diagnostics) {
	result, err := showById(ctx, id)

	if err != nil {
		var zeroValue T

		buildMainDiagnostic := func(objectLevel string, fullyQualifiedName string, err error) diag.Diagnostic {
			return diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Failed to query %s object. Marking the resource as removed.", objectLevel),
				Detail:   fmt.Sprintf("%s identifier: %s, Err: %s", resourceName.String(), fullyQualifiedName, err),
			}
		}

		buildHierarchyErrDiagnostic := func(objectLevel string, fullyQualifiedName string, err error) diag.Diagnostic {
			return diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Failed to query %s for %s.", objectLevel, resourceName.String()),
				Detail:   fmt.Sprintf("%s identifier: %s, Err: %s", resourceName.String(), fullyQualifiedName, err),
			}
		}

		shouldRemoveFromState := false
		if errors.Is(err, sdk.ErrObjectNotFound) {
			shouldRemoveFromState = true
		}

		// ErrObjectNotExistOrAuthorized or ErrDoesNotExistOrOperationCannotBePerformed can only happen
		// when the higher hierarchy object is not accessible for some reason during the "main" showById.
		shouldCheckHigherHierarchies := errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) || errors.Is(err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)

		switch id := any(id).(type) {
		case sdk.AccountObjectIdentifier:
			return zeroValue, shouldRemoveFromState, diag.Diagnostics{
				buildMainDiagnostic("account", id.FullyQualifiedName(), err),
			}
		case sdk.DatabaseObjectIdentifier:
			diags := diag.Diagnostics{
				buildMainDiagnostic("database", id.FullyQualifiedName(), err),
			}

			if shouldCheckHigherHierarchies {
				if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
					if errors.Is(err, sdk.ErrObjectNotFound) {
						shouldRemoveFromState = true
					}

					diags = append(diags, buildHierarchyErrDiagnostic("database", id.FullyQualifiedName(), err))
				}
			}

			return zeroValue, shouldRemoveFromState, diags
		case sdk.SchemaObjectIdentifier:
			diags := diag.Diagnostics{
				buildMainDiagnostic("schema", id.FullyQualifiedName(), err),
			}

			if shouldCheckHigherHierarchies {
				if _, err := client.Schemas.ShowByID(ctx, id.SchemaId()); err != nil {
					// If the underlying database is missing, this can also either throw ErrObjectNotExistOrAuthorized or ErrDoesNotExistOrOperationCannotBePerformed
					if errors.Is(err, sdk.ErrObjectNotFound) || errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) || errors.Is(err, sdk.ErrDoesNotExistOrOperationCannotBePerformed) {
						shouldRemoveFromState = true
					}

					diags = append(diags, buildHierarchyErrDiagnostic("schema", id.FullyQualifiedName(), err))
				}

				if !shouldRemoveFromState {
					if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
						if errors.Is(err, sdk.ErrObjectNotFound) {
							shouldRemoveFromState = true
						}

						diags = append(diags, buildHierarchyErrDiagnostic("database", id.FullyQualifiedName(), err))
					}
				}
			}

			return zeroValue, shouldRemoveFromState, diags
		case sdk.SchemaObjectIdentifierWithArguments:
			diags := diag.Diagnostics{
				buildMainDiagnostic("schema", id.FullyQualifiedName(), err),
			}

			if shouldCheckHigherHierarchies {
				if _, err := client.Schemas.ShowByID(ctx, id.SchemaId()); err != nil {
					// If the underlying database is missing, this can also either throw ErrObjectNotExistOrAuthorized or ErrDoesNotExistOrOperationCannotBePerformed
					if errors.Is(err, sdk.ErrObjectNotFound) || errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) || errors.Is(err, sdk.ErrDoesNotExistOrOperationCannotBePerformed) {
						shouldRemoveFromState = true
					}

					diags = append(diags, buildHierarchyErrDiagnostic("schema", id.FullyQualifiedName(), err))
				}

				if !shouldRemoveFromState {
					if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
						if errors.Is(err, sdk.ErrObjectNotFound) {
							shouldRemoveFromState = true
						}

						diags = append(diags, buildHierarchyErrDiagnostic("database", id.FullyQualifiedName(), err))
					}
				}
			}

			return zeroValue, shouldRemoveFromState, diags
		}
	}

	return result, false, nil
}
