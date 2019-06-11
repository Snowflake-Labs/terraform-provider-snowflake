package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	var args []interface{}
	a := assert.New(t)
	v := View("test", "SELECT * FROM DUMMY LIMIT 1", args)
	a.NotNil(v)
	a.False(v.secure)

	v.WithSecure()
	a.True(v.secure)

	v.WithComment("great comment")
	a.Equal("great comment", v.comment)

	q, qArgs := v.Create()
	a.Equal("CREATE SECURE VIEW ? COMMENT = ? AS SELECT * FROM DUMMY LIMIT 1", q)
	a.Len(qArgs, 2)
	a.Equal(qArgs[0], "test")
	a.Equal(qArgs[1], "great comment")

	q, qArgs = v.Rename("test2")
	a.Equal("ALTER VIEW ? RENAME TO ?", q)
	a.Len(qArgs, 2)
	a.Equal(qArgs[0], "test")
	a.Equal(qArgs[1], "test2")
}
