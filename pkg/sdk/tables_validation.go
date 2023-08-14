package sdk

import "errors"

var (
	_ validatableOpts = &CreateTableOptions{}
	_ validatableOpts = &AlterTableOptions{}
)

func (opts *CreateTableOptions) validateProp() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if len(opts.Columns) == 0 {
		return errTableNeedsAtLeastOneColumn
	}
	return nil
}
func (opts *AlterTableOptions) validateProp() error {
	return nil
}

var (
	errTableNeedsAtLeastOneColumn = errors.New("table create statement needs at least one column")
)
