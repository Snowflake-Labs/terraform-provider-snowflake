package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	a := assert.New(t)
	v := View("test")
	a.NotNil(v)
	a.False(v.secure)

	v.WithSecure()
	a.True(v.secure)

	v.WithComment("great comment")
	a.Equal("great comment", v.comment)

	v.WithStatement("SELECT * FROM DUMMY LIMIT 1")
	a.Equal("SELECT * FROM DUMMY LIMIT 1", v.statement)

	v.WithStatement("SELECT * FROM DUMMY WHERE blah = '?' LIMIT 1")
	v.WithStatementArgs([]interface{}{"blahblah"})
	a.Len(v.statementArgs, 1)

	q, qArgs := v.Create()
	a.Equal("CREATE SECURE VIEW ? COMMENT = ? AS SELECT * FROM DUMMY WHERE blah = '?' LIMIT 1", q)
	a.Len(qArgs, 3)
	a.Equal(qArgs[0], "test")
	a.Equal(qArgs[1], "great comment")
	a.Equal(qArgs[2], "blahblah")

	q, qArgs = v.Rename("test2")
	a.Equal("ALTER VIEW ? RENAME TO ?", q)
	a.Len(qArgs, 2)
	a.Equal(qArgs[0], "test")
	a.Equal(qArgs[1], "test2")

	q, qArgs = v.Secure()
	a.Equal("ALTER VIEW ? SET SECURE", q)
	a.Len(qArgs, 1)
	a.Equal(qArgs[0], "test")

	q, qArgs = v.Unsecure()
	a.Equal("ALTER VIEW ? UNSET SECURE", q)
	a.Len(qArgs, 1)
	a.Equal(qArgs[0], "test")

	q, qArgs = v.ChangeComment("bad comment")
	a.Equal("ALTER VIEW ? SET COMMENT = ?", q)
	a.Len(qArgs, 2)
	a.Equal(qArgs[0], "test")
	a.Equal(qArgs[1], "bad comment")

	q, qArgs = v.RemoveComment()
	a.Equal("ALTER VIEW ? UNSET COMMENT", q)
	a.Len(qArgs, 1)
	a.Equal(qArgs[0], "test")

	q, qArgs = v.Drop()
	a.Equal("DROP VIEW ?", q)
	a.Len(qArgs, 1)
	a.Equal(qArgs[0], "test")

	q, qArgs = v.Show()
	a.Equal("SHOW VIEWS LIKE ?", q)
	a.Len(qArgs, 1)
	a.Equal(qArgs[0], "test")
}
