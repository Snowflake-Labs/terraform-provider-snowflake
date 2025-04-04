package sdk

import (
	"context"
	"errors"
)

// SafeDrop is a helper function that wraps a drop function and handles common error cases that
// relate to missing high hierarchy objects when dropping lower ones like schemas, tables, views, etc.
// Whenever an object is missing or the higher hierarchy object is not accessible, it will return ErrObjectNotFound error,
// which can be leveraged with [errors.Is] to handle the logic in case of missing objects.
func SafeDrop[ID AccountObjectIdentifier | DatabaseObjectIdentifier | SchemaObjectIdentifier | SchemaObjectIdentifierWithArguments](
	client *Client,
	drop func() error,
	ctx context.Context,
	id ID,
) error {
	err := drop()

	// ErrObjectNotExistOrAuthorized can only happen
	// when the higher hierarchy object is not accessible for some reason during the "main" drop operation.
	shouldCheckHigherHierarchies := errors.Is(err, ErrObjectNotExistOrAuthorized)

	if !shouldCheckHigherHierarchies {
		return err
	}

	if err != nil {
		errs := []error{err}

		switch id := any(id).(type) {
		case AccountObjectIdentifier:
			return err
		case DatabaseObjectIdentifier:
			if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
				errs = append(errs, err)

				if errors.Is(err, ErrObjectNotFound) {
					return errors.Join(append(errs, ErrSkippable)...)
				}
			}

			return errors.Join(errs...)
		case SchemaObjectIdentifier, SchemaObjectIdentifierWithArguments:
			schemaObjectId := id.(interface {
				SchemaId() DatabaseObjectIdentifier
				DatabaseId() AccountObjectIdentifier
			})

			if _, err := client.Schemas.ShowByID(ctx, schemaObjectId.SchemaId()); err != nil {
				errs = append(errs, err)

				if errors.Is(err, ErrObjectNotFound) {
					return errors.Join(append(errs, ErrSkippable)...)
				}
			}

			if _, err := client.Databases.ShowByID(ctx, schemaObjectId.DatabaseId()); err != nil {
				errs = append(errs, err)

				if errors.Is(err, ErrObjectNotFound) {
					return errors.Join(append(errs, ErrSkippable)...)
				}
			}

			return errors.Join(errs...)
		}
	}

	return nil
}
