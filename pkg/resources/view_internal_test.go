package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitViewID(t *testing.T) {
	r := require.New(t)

	id := "great_db|great_schema|great_view"
	db, schema, view, err := splitViewID(id)
	r.NoError(err)

	r.Equal("great_db", db)
	r.Equal("great_schema", schema)
	r.Equal("great_view", view)

	id = "bad_id"
	_, _, _, err = splitViewID(id)
	r.Error(err)
}
