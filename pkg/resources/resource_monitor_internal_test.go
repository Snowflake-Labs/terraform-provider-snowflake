package resources

import (
	"database/sql"
	"testing"
)

func TestExtractTriggerInts(t *testing.T) {
	// TODO rewrite to use testify/assert
	resp := sql.NullString{String: "51%,63%", Valid: true}
	out, err := extractTriggerInts(resp)
	if err != nil {
		t.Error(err)
	}
	if l := len(out); l != 2 {
		t.Errorf("Expected 2 values, got %d", l)
	}

	if out[0] != 51 {
		t.Errorf("Expected first value to be 51, got %d", out[0])
	}

	if out[1] != 63 {
		t.Errorf("Expected second value to be 63, got %d", out[1])
	}
}
