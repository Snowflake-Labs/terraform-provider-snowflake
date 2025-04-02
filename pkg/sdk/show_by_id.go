package sdk

import (
	"context"
	"errors"
)

// SafeShowById is a helper function that wraps a showById function and handles common error cases that
// relate to missing high hierarchy objects when querying lower ones like schemas, tables, views, etc.
// Whenever an object is missing or the higher hierarchy object is not accessible, it will return ErrObjectNotFound error,
// which can be leveraged with [errors.Is] to handle the logic in case of missing objects.
func SafeShowById[T any, ID AccountObjectIdentifier | DatabaseObjectIdentifier | SchemaObjectIdentifier | SchemaObjectIdentifierWithArguments](
	client *Client,
	showById func(context.Context, ID) (T, error),
	ctx context.Context,
	id ID,
) (T, error) {
	result, err := showById(ctx, id)
	if errors.Is(err, ErrObjectNotFound) {
		return result, err
	}

	if err != nil {
		var zeroValue T
		errs := []error{err}

		// ErrObjectNotExistOrAuthorized or ErrDoesNotExistOrOperationCannotBePerformed can only happen
		// when the higher hierarchy object is not accessible for some reason during the "main" showById.
		shouldCheckHigherHierarchies := errors.Is(err, ErrObjectNotExistOrAuthorized) || errors.Is(err, ErrDoesNotExistOrOperationCannotBePerformed)

		switch id := any(id).(type) {
		case AccountObjectIdentifier:
			return result, err
		case DatabaseObjectIdentifier:
			if shouldCheckHigherHierarchies {
				if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
					errs = append(errs, err)
				}
			}

			return zeroValue, errors.Join(errs...)
		case SchemaObjectIdentifier:
			if shouldCheckHigherHierarchies {
				if _, err := client.Schemas.ShowByID(ctx, id.SchemaId()); err != nil {
					errs = append(errs, err)

					if errors.Is(err, ErrObjectNotFound) {
						return zeroValue, errors.Join(errs...)
					}
				}

				if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
					errs = append(errs, err)
				}
			}

			return zeroValue, errors.Join(errs...)
		case SchemaObjectIdentifierWithArguments:
			if shouldCheckHigherHierarchies {
				if _, err := client.Schemas.ShowByID(ctx, id.SchemaId()); err != nil {
					errs = append(errs, err)

					if errors.Is(err, ErrObjectNotFound) {
						return zeroValue, errors.Join(errs...)
					}
				}

				if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
					errs = append(errs, err)
				}
			}

			return zeroValue, errors.Join(errs...)
		}
	}

	return result, nil
}
