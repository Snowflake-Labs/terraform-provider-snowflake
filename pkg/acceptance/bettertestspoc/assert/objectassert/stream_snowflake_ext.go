package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StreamAssert) HasTableId(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.TableName == nil {
			return fmt.Errorf("expected table name to have value; got: nil")
		}
		gotTableId, err := sdk.ParseSchemaObjectIdentifier(*o.TableName)
		if err != nil {
			return err
		}
		if gotTableId.FullyQualifiedName() != expected {
			return fmt.Errorf("expected table name: %v; got: %v", expected, *o.TableName)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasStageName(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.TableName == nil {
			return fmt.Errorf("expected table name to have value; got: nil")
		}
		if *o.TableName != expected {
			return fmt.Errorf("expected table name: %v; got: %v", expected, *o.TableName)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasSourceType(expected sdk.StreamSourceType) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.SourceType == nil {
			return fmt.Errorf("expected source type to have value; got: nil")
		}
		if *o.SourceType != expected {
			return fmt.Errorf("expected source type: %v; got: %v", expected, *o.SourceType)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasBaseTables(expected []sdk.SchemaObjectIdentifier) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.BaseTables == nil {
			return fmt.Errorf("expected base tables to have value; got: nil")
		}
		if len(o.BaseTables) != len(expected) {
			return fmt.Errorf("expected base tables length: %v; got: %v", len(expected), len(o.BaseTables))
		}
		for i := range o.BaseTables {
			if o.BaseTables[i].FullyQualifiedName() != expected[i].FullyQualifiedName() {
				return fmt.Errorf("expected base table id: %v; got: %v", expected[i], o.BaseTables[i])
			}
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasMode(expected sdk.StreamMode) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.Mode == nil {
			return fmt.Errorf("expected mode to have value; got: nil")
		}
		if *o.Mode != expected {
			return fmt.Errorf("expected mode: %v; got: %v", expected, *o.Mode)
		}
		return nil
	})
	return s
}
