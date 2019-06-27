package resources

import (
	"testing"
)

func TestSplitViewID(t *testing.T) {
	id := "great_db|great_schema|great_view"
	db, schema, view, err := splitViewID(id)
	if err != nil {
		t.Error(err)
	}
	if db != "great_db" {
		t.Errorf("Expecting great_db, got %v", db)
	}
	if schema != "great_schema" {
		t.Errorf("Expecting great_schema, got %v", schema)
	}
	if view != "great_view" {
		t.Errorf("Expecting great_view, got %v", view)
	}

	id = "bad_id"
	_, _, _, err = splitViewID(id)
	if err == nil {
		t.Errorf("Expecting an error, got none")
	}
}
