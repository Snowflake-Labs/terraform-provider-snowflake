package sdk

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type unexportedTestHelper struct {
	static bool `ddl:"static" sql:"EXAMPLE_STATIC"`
}

func TestBuilder_parseField(t *testing.T) {
	t.Run("test boolean keyword", func(t *testing.T) {
		s := struct {
			BooleanKeyword *bool `ddl:"keyword" sql:"EXAMPLE_KEYWORD"`
		}{
			BooleanKeyword: Bool(true),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("BooleanKeyword")
		field, ok := typ.FieldByName("BooleanKeyword")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "EXAMPLE_KEYWORD", clause.String())
	})

	t.Run("test boolean keyword with false value", func(t *testing.T) {
		s := struct {
			BooleanKeyword *bool `ddl:"keyword" sql:"EXAMPLE_KEYWORD"`
		}{
			BooleanKeyword: Bool(false),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("BooleanKeyword")
		field, ok := typ.FieldByName("BooleanKeyword")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Nil(t, clause)
	})

	t.Run("test boolean keyword with nil value", func(t *testing.T) {
		s := struct {
			BooleanKeyword *bool `ddl:"keyword" sql:"EXAMPLE_KEYWORD"`
		}{}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("BooleanKeyword")
		field, ok := typ.FieldByName("BooleanKeyword")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Nil(t, clause)
	})

	t.Run("test string keyword with value", func(t *testing.T) {
		s := struct {
			StringKeyword *string `ddl:"keyword"`
		}{
			StringKeyword: String("example"),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("StringKeyword")
		field, ok := typ.FieldByName("StringKeyword")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "example", clause.String())
	})

	t.Run("test string keyword with nil value", func(t *testing.T) {
		s := struct {
			StringKeyword *string `ddl:"keyword" `
		}{}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("StringKeyword")
		field, ok := typ.FieldByName("StringKeyword")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Nil(t, clause)
	})

	t.Run("test string keyword with double quotes", func(t *testing.T) {
		s := struct {
			StringKeyword *string `ddl:"keyword,double_quotes"`
		}{
			StringKeyword: String("example"),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("StringKeyword")
		field, ok := typ.FieldByName("StringKeyword")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, `"example"`, clause.String())
	})

	t.Run("test string keyword with single quotes", func(t *testing.T) {
		s := struct {
			StringKeyword *string `ddl:"keyword,single_quotes"`
		}{
			StringKeyword: String("example"),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("StringKeyword")
		field, ok := typ.FieldByName("StringKeyword")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, `'example'`, clause.String())
	})

	t.Run("test static with value", func(t *testing.T) {
		s := struct {
			Static *bool `ddl:"static" sql:"EXAMPLE_STATIC"`
		}{
			Static: Bool(true),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("Static")
		field, ok := typ.FieldByName("Static")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "EXAMPLE_STATIC", clause.String())
	})

	t.Run("test static with nil value", func(t *testing.T) {
		s := struct {
			Static *bool `ddl:"static" sql:"EXAMPLE_STATIC"`
		}{}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("Static")
		field, ok := typ.FieldByName("Static")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "EXAMPLE_STATIC", clause.String())
	})

	t.Run("test parameter with value", func(t *testing.T) {
		s := struct {
			Parameter *string `ddl:"parameter" sql:"EXAMPLE_PARAMETER"`
		}{
			Parameter: String("example"),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("Parameter")
		field, ok := typ.FieldByName("Parameter")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "EXAMPLE_PARAMETER = example", clause.String())
	})

	t.Run("test parameter with nil value", func(t *testing.T) {
		s := struct {
			Parameter *string `ddl:"parameter" sql:"EXAMPLE_PARAMETER"`
		}{}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("Parameter")
		field, ok := typ.FieldByName("Parameter")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Nil(t, clause)
	})

	t.Run("test parameter with double quotes", func(t *testing.T) {
		s := struct {
			Parameter *string `ddl:"parameter,double_quotes" sql:"EXAMPLE_PARAMETER"`
		}{
			Parameter: String("example"),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("Parameter")
		field, ok := typ.FieldByName("Parameter")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, `EXAMPLE_PARAMETER = "example"`, clause.String())
	})

	t.Run("test parameter with single quotes", func(t *testing.T) {
		s := struct {
			Parameter *string `ddl:"parameter,single_quotes" sql:"EXAMPLE_PARAMETER"`
		}{
			Parameter: String("example"),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("Parameter")
		field, ok := typ.FieldByName("Parameter")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, `EXAMPLE_PARAMETER = 'example'`, clause.String())
	})

	t.Run("test parameter with integer value", func(t *testing.T) {
		s := struct {
			Parameter *int `ddl:"parameter" sql:"EXAMPLE_PARAMETER"`
		}{
			Parameter: Int(1),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("Parameter")
		field, ok := typ.FieldByName("Parameter")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "EXAMPLE_PARAMETER = 1", clause.String())
	})

	t.Run("test parameter with no db", func(t *testing.T) {
		s := struct {
			Parameter *string `ddl:"parameter"`
		}{
			Parameter: String("example"),
		}
		val := reflect.ValueOf(s)
		typ := val.Type()
		value := val.FieldByName("Parameter")
		field, ok := typ.FieldByName("Parameter")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "= example", clause.String())
	})

	t.Run("test unexported static value set", func(t *testing.T) {
		s := &unexportedTestHelper{
			static: true,
		}
		val := reflect.ValueOf(s).Elem()
		typ := val.Type()
		value := val.FieldByName("static")
		field, ok := typ.FieldByName("static")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "EXAMPLE_STATIC", clause.String())
	})

	t.Run("test unexported static value not set", func(t *testing.T) {
		s := &unexportedTestHelper{
			static: false,
		}
		val := reflect.ValueOf(s).Elem()
		typ := val.Type()
		value := val.FieldByName("static")
		field, ok := typ.FieldByName("static")
		require.True(t, ok)
		clause, err := builder.parseField(field, value)
		require.NoError(t, err)
		assert.Equal(t, "EXAMPLE_STATIC", clause.String())
	})
}

func TestReverseModifier(t *testing.T) {
	t.Run("test reverse modifier", func(t *testing.T) {
		result := Reverse.Modify([]string{"example", "DESC"})
		assert.Equal(t, `DESC example`, result)
	})

	t.Run("test no reverse modifier", func(t *testing.T) {
		result := NoReverse.Modify([]string{"example", "DESC"})
		assert.Equal(t, `example DESC`, result)
	})

	t.Run("test unknown reverse modifier", func(t *testing.T) {
		result := reverseModifier("unknown").Modify([]string{"example", "DESC"})
		assert.Equal(t, `example DESC`, result)
	})
}

func TestEqualsModifier(t *testing.T) {
	t.Run("test equals modifier", func(t *testing.T) {
		result := Equals.Modify("example")
		assert.Equal(t, `example = `, result)
	})

	t.Run("test no equals modifier", func(t *testing.T) {
		result := NoEquals.Modify("example")
		assert.Equal(t, `example `, result)
	})

	t.Run("test unknown equals modifier", func(t *testing.T) {
		result := equalsModifier("unknown").Modify("example")
		assert.Equal(t, `example `, result)
	})
}

func TestParenModifier(t *testing.T) {
	t.Run("test paren modifier", func(t *testing.T) {
		result := Parentheses.Modify("example")
		assert.Equal(t, `(example)`, result)
	})

	t.Run("test no paren modifier", func(t *testing.T) {
		result := NoParentheses.Modify("example")
		assert.Equal(t, `example`, result)
	})

	t.Run("test unknown paren modifier", func(t *testing.T) {
		result := parenModifier("unknown").Modify("example")
		assert.Equal(t, `example`, result)
	})
}

func TestQuoteModifier(t *testing.T) {
	t.Run("test quotes modifier", func(t *testing.T) {
		result := DoubleQuotes.Modify("example")
		assert.Equal(t, `"example"`, result)
	})

	t.Run("test no quotes modifier", func(t *testing.T) {
		result := NoQuotes.Modify("example")
		assert.Equal(t, `example`, result)
	})

	t.Run("test single quotes modifier", func(t *testing.T) {
		result := SingleQuotes.Modify("example")
		assert.Equal(t, `'example'`, result)
	})

	t.Run("test unknown modifier", func(t *testing.T) {
		result := quoteModifier("unknown").Modify("example")
		assert.Equal(t, `example`, result)
	})
}

type structTestHelper struct {
	static bool                    `ddl:"static" sql:"EXAMPLE_STATIC"`
	name   AccountObjectIdentifier `ddl:"identifier"`
	Param  *string                 `ddl:"parameter" sql:"EXAMPLE_PARAMETER"`
}

func TestBuilder_parseStruct(t *testing.T) {
	t.Run("test struct with no fields", func(t *testing.T) {
		s := struct{}{}
		clauses, err := builder.parseStruct(s)
		require.NoError(t, err)
		assert.Len(t, clauses, 0)
	})

	t.Run("test struct with all fields", func(t *testing.T) {
		s := &structTestHelper{
			static: true,
			name:   randomAccountObjectIdentifier(t),
			Param:  String("example"),
		}
		clauses, err := builder.parseStruct(s)
		require.NoError(t, err)
		assert.Len(t, clauses, 3)
		assert.Equal(t, "EXAMPLE_STATIC", clauses[0].String())
		assert.Equal(t, s.name.FullyQualifiedName(), clauses[1].String())
		assert.Equal(t, "EXAMPLE_PARAMETER = example", clauses[2].String())
	})

	t.Run("struct with a slice field using ddl: keyword", func(t *testing.T) {
		type testListElement struct {
			K  *string `ddl:"parameter,single_quotes" sql:"KEY"`
			K2 *string `ddl:"parameter,single_quotes" sql:"KEY2"`
		}
		s := &struct {
			List []testListElement `ddl:"keyword,parentheses" sql:"TAG"`
		}{
			List: []testListElement{{K: String("abc"), K2: String("def")}, {K: String("123"), K2: String("456")}},
		}
		clauses, err := builder.parseStruct(s)
		require.NoError(t, err)
		assert.Len(t, clauses, 1)
		assert.Equal(t, "TAG (KEY = 'abc' KEY2 = 'def',KEY = '123' KEY2 = '456')", clauses[0].String())
	})

	t.Run("struct with a slice field using ddl: - (no elements)", func(t *testing.T) {
		type testListElement struct {
			K *string `ddl:"parameter,single_quotes" sql:"KEY"`
		}
		s := &struct {
			List []testListElement `ddl:"-"`
		}{}
		clauses, err := builder.parseStruct(s)
		require.NoError(t, err)
		assert.Len(t, clauses, 0)
	})

	t.Run("struct with a slice field using ddl: - (no_parentheses)", func(t *testing.T) {
		type testListElement struct {
			K *string `ddl:"parameter,single_quotes" sql:"KEY"`
		}
		s := &struct {
			List []testListElement `ddl:"-,no_parentheses"`
		}{
			List: []testListElement{{K: String("abc")}, {K: String("123")}},
		}
		clauses, err := builder.parseStruct(s)
		require.NoError(t, err)
		assert.Len(t, clauses, 1)
		assert.Equal(t, "KEY = 'abc',KEY = '123'", clauses[0].String())
	})

	t.Run("struct with a struct list using ddl: list", func(t *testing.T) {
		type testListElement struct {
			A bool `ddl:"static" sql:"A"`
			B bool `ddl:"static" sql:"B"`
			C bool `ddl:"static" sql:"C"`
		}
		s := &struct {
			List *testListElement `ddl:"list"`
		}{
			List: &testListElement{A: true, B: true, C: true},
		}
		clauses, err := builder.parseStruct(s)
		require.NoError(t, err)
		assert.Len(t, clauses, 1)
		assert.Equal(t, "A,B,C", clauses[0].String())
	})
}

func TestBuilder_sql(t *testing.T) {
	t.Run("test sql with no clauses", func(t *testing.T) {
		s := builder.sql([]sqlClause{}...)
		assert.Equal(t, "", s)
	})

	t.Run("test sql with clauses", func(t *testing.T) {
		clauses := []sqlClause{
			sqlStaticClause("EXAMPLE_STATIC"),
			sqlParameterClause{
				key:   "EXAMPLE_KEYWORD",
				value: "example",
				em:    Equals,
			},
		}
		s := builder.sql(clauses...)
		assert.Equal(t, "EXAMPLE_STATIC EXAMPLE_KEYWORD = example", s)
	})
}
