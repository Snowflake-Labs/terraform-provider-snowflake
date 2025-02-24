package objectassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StreamAssert) HasTableId(expected sdk.SchemaObjectIdentifier) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.TableName == nil {
			return fmt.Errorf("expected table name to have value; got: nil")
		}
		gotTableId, err := sdk.ParseSchemaObjectIdentifier(*o.TableName)
		if err != nil {
			return err
		}
		if gotTableId.FullyQualifiedName() != expected.FullyQualifiedName() {
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

func (s *StreamAssert) HasBaseTablesPartiallyQualified(expected ...string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if len(o.BaseTables) != len(expected) {
			return fmt.Errorf("expected base tables length: %v; got: %v", len(expected), len(o.BaseTables))
		}
		var errs []error
		for _, wantName := range expected {
			if !slices.Contains(o.BaseTables, wantName) {
				errs = append(errs, fmt.Errorf("expected name: %s, to be in the list ids: %v", wantName, o.BaseTables))
			}
		}
		return errors.Join(errs...)
	})
	return s
}

// StreamWithTestClient is temporary. It's here to show the changes proposed to the assertions setup.
// It will be generated later.
func StreamWithTestClient(t *testing.T, id sdk.SchemaObjectIdentifier) *StreamAssert {
	t.Helper()
	return &StreamAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeStream, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.Stream, sdk.SchemaObjectIdentifier] {
			return testClient.Stream.Show
		}),
	}
}
