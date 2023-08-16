package sdk

import (
	"context"
)

// validatableOpts is just a proposal how we can remove some of the boilerplate.
type validatableOpts interface {
	// validateProp will be renamed if accepted.
	// Having the name `validate` now would make all other `...Opts` structs validatableOpts because of Go nature.
	validateProp() error
}

// validateAndExec is just a proposal how we can remove some of the boilerplate.
func validateAndExec(client *Client, ctx context.Context, opts validatableOpts) error {
	if err := opts.validateProp(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = client.exec(ctx, sql)
	return err
}

// validateAndQuery is just a proposal how we can remove some of the boilerplate.
func validateAndQuery[T any](client *Client, ctx context.Context, opts validatableOpts) (*[]T, error) {
	if err := opts.validateProp(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}

	var dest []T
	err = client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

// validateAndQueryOne is just a proposal how we can remove some of the boilerplate.
func validateAndQueryOne[T any](client *Client, ctx context.Context, opts validatableOpts) (*T, error) {
	if err := opts.validateProp(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}

	var dest T
	err = client.queryOne(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func createIfNil[T any](t *T) *T {
	if t == nil {
		return new(T)
	}
	return t
}

type convertibleRow[T any] interface {
	convert() *T
}

func convertRows[T convertibleRow[U], U any](dbRows *[]T) []U {
	resultList := make([]U, len(*dbRows))
	for i, row := range *dbRows {
		resultList[i] = *(row.convert())
	}
	return resultList
}

type optionsProvider[T any] interface {
	toOpts() *T
}
