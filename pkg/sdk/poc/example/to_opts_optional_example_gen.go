package example

import "context"

type ToOptsOptionalExamples interface {
	Create(ctx context.Context, request *CreateToOptsOptionalExampleRequest) error
	Alter(ctx context.Context, request *AlterToOptsOptionalExampleRequest) error
}

// CreateToOptsOptionalExampleOptions is based on https://example.com.
type CreateToOptsOptionalExampleOptions struct {
	alter    bool                     `ddl:"static" sql:"ALTER"`
	IfExists *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name     DatabaseObjectIdentifier `ddl:"identifier"`
}

// AlterToOptsOptionalExampleOptions is based on https://example.com.
type AlterToOptsOptionalExampleOptions struct {
	alter         bool                     `ddl:"static" sql:"ALTER"`
	IfExists      *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name          DatabaseObjectIdentifier `ddl:"identifier"`
	OptionalField *OptionalField           `ddl:"keyword"`
	RequiredField RequiredField            `ddl:"keyword"`
}

type OptionalField struct {
	SomeList []DatabaseObjectIdentifier `ddl:"list"`
}

type RequiredField struct {
	SomeRequiredList []DatabaseObjectIdentifier `ddl:"list"`
}
