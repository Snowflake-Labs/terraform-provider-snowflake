package example

import (
	"context"
)

var _ ToOptsOptionalExamples = (*toOptsOptionalExamples)(nil)

type toOptsOptionalExamples struct {
	client *Client
}

func (v *toOptsOptionalExamples) Create(ctx context.Context, request *CreateToOptsOptionalExampleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *toOptsOptionalExamples) Alter(ctx context.Context, request *AlterToOptsOptionalExampleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateToOptsOptionalExampleRequest) toOpts() *CreateToOptsOptionalExampleOptions {
	opts := &CreateToOptsOptionalExampleOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *AlterToOptsOptionalExampleRequest) toOpts() *AlterToOptsOptionalExampleOptions {
	opts := &AlterToOptsOptionalExampleOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}

	if r.OptionalField != nil {
		opts.OptionalField = &OptionalField{
			SomeList: r.OptionalField.SomeList,
		}
	}
	opts.RequiredField = RequiredField{
		SomeRequiredList: r.RequiredField.SomeRequiredList,
	}

	return opts
}
