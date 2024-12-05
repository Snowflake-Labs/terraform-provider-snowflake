package example

import (
	"context"
)

var _ FeaturesExample = (*featuresExample)(nil)

type featuresExample struct {
	client *Client
}

func (v *featuresExample) Alter(ctx context.Context, request *AlterFeaturesExamplesRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *AlterFeaturesExamplesRequest) toOpts() *AlterFeaturesExamplesOptions {
	opts := &AlterFeaturesExamplesOptions{
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
