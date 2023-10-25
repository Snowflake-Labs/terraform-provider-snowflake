// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

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

	first := 51
	if out[0] != first {
		t.Errorf("Expected first value to be 51, got %d", out[0])
	}

	second := 63
	if out[1] != second {
		t.Errorf("Expected second value to be 63, got %d", out[1])
	}
}
