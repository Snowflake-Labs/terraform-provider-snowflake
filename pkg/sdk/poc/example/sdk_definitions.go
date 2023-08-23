package example

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// This file's only purpose is to make generated objects compile (or close to compile).
// Later code will be generated inside sdk package, so the objects will be accessible there.
// TODO: check if generating with package name + invoking format removes unnecessary qualifier

type optionsProvider[T any] interface {
	toOpts() *T
}

type validatable interface {
	validate() error
}

type Client struct{}

type ObjectIdentifier interface{}
type AccountObjectIdentifier struct{}
type DatabaseObjectIdentifier struct{}
type ExternalObjectIdentifier struct{}
type SchemaObjectIdentifier struct{}
type TableColumnIdentifier struct{}

func randomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
	_ = t
	return AccountObjectIdentifier{}
}

func randomDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
	_ = t
	return DatabaseObjectIdentifier{}
}

func randomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
	_ = t
	return SchemaObjectIdentifier{}
}

func validObjectidentifier(objectIdentifier ObjectIdentifier) bool {
	_ = objectIdentifier
	return true
}

func valueSet(value interface{}) bool {
	_ = value
	return true
}

func anyValueSet(values ...interface{}) bool {
	_ = values
	return true
}

func everyValueSet(values ...interface{}) bool {
	_ = values
	return true
}

func exactlyOneValueSet(values ...interface{}) bool {
	_ = values
	return true
}

func errOneOf(fieldNames ...string) error {
	return fmt.Errorf("fields %v are incompatible and cannot be set at once", fieldNames)
}

func errExactlyOneOf(fieldNames ...string) error {
	return fmt.Errorf("exactly one of %v must be set", fieldNames)
}

func errAtLeastOneOf(fieldNames ...string) error {
	return fmt.Errorf("at least one of %v must be set", fieldNames)
}

var (
	errNilOptions              = errors.New("options cannot be nil")
	ErrInvalidObjectIdentifier = errors.New("invalid object identifier")
)

func validateAndExec(client *Client, ctx context.Context, opts validatable) error {
	_, _, _ = client, ctx, opts
	return nil
}

func validateAndQuery[T any](client *Client, ctx context.Context, opts validatable) (*[]T, error) {
	_, _, _ = client, ctx, opts
	return nil, nil
}

func validateAndQueryOne[T any](client *Client, ctx context.Context, opts validatable) (*T, error) {
	_, _, _ = client, ctx, opts
	return nil, nil
}
