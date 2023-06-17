package sdk

import (
	"database/sql"
	"testing"
)

func TestExtractTriggerInts(t *testing.T) {
	// TODO rewrite to use testify/assert
	resp := sql.NullString{String: "51%,63%", Valid: true}
	out, err := extractTriggers(resp, Suspend)
	if err != nil {
		t.Error(err)
	}
	if l := len(out); l != 2 {
		t.Errorf("Expected 2 values, got %d", l)
	}

	first := TriggerDefinition{Threshold: 51, TriggerAction: Suspend}
	if out[0] != first {
		t.Errorf("Expected first value to be 51, got %d", out[0].Threshold)
	}

	second := TriggerDefinition{Threshold: 63, TriggerAction: Suspend}
	if out[1] != second {
		t.Errorf("Expected second value to be 63, got %d", out[1].Threshold)
	}
}
