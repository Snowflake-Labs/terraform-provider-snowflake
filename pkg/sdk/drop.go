package sdk

import (
	"context"
	"errors"
)

func SafeDrop[DropRequest any, ID AccountObjectIdentifier | DatabaseObjectIdentifier | SchemaObjectIdentifier | SchemaObjectIdentifierWithArguments](
	client *Client,
	drop func() error,
	ctx context.Context,
	id ID,
) error {
	if err := drop(); err != nil {
		errs := []error{err}

		// ErrObjectNotExistOrAuthorized will be returned
		// when the higher hierarchy object is not accessible for some reason during the "main" showById.
		shouldCheckHigherHierarchies := errors.Is(err, ErrObjectNotExistOrAuthorized)

		switch id := any(id).(type) {
		case AccountObjectIdentifier:
			return err
		case DatabaseObjectIdentifier:
			if shouldCheckHigherHierarchies {
				if err := client.Databases.Drop(ctx, id.DatabaseId(), &DropDatabaseOptions{IfExists: Bool(true)}); err != nil {
					errs = append(errs, err)
				}
			}

			return errors.Join(errs...)
		case SchemaObjectIdentifier:
			if shouldCheckHigherHierarchies {
				if _, err := client.Schemas.ShowByID(ctx, id.SchemaId()); err != nil {
					errs = append(errs, err)

					if errors.Is(err, ErrObjectNotFound) {
						return errors.Join(errs...)
					}
				}

				if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
					errs = append(errs, err)
				}
			}

			return errors.Join(errs...)
		case SchemaObjectIdentifierWithArguments:
			if shouldCheckHigherHierarchies {
				if _, err := client.Schemas.ShowByID(ctx, id.SchemaId()); err != nil {
					errs = append(errs, err)

					if errors.Is(err, ErrObjectNotFound) {
						return errors.Join(errs...)
					}
				}

				if _, err := client.Databases.ShowByID(ctx, id.DatabaseId()); err != nil {
					errs = append(errs, err)
				}
			}

			return errors.Join(errs...)
		}
	}

	return nil
}
