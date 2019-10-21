package resources

import (
	"testing"
)

func TestSplitGrantID(t *testing.T) {
	// Vanilla
	id := "database_name|schema|view_name|privilege"
	db, schema, view, priv, err := splitGrantID(id)
	if err != nil {
		t.Error(err)
	}
	if db != "database_name" {
		t.Errorf("Expected db to be database_name, got %v", db)
	}
	if schema != "schema" {
		t.Errorf("Expected schema to be schema, got %v", schema)
	}
	if view != "view_name" {
		t.Errorf("Expected view to be view_name, got %v", view)
	}
	if priv != "privilege" {
		t.Errorf("Expected priv to be privilege, got %v", priv)
	}

	// No view
	id = "database_name|||privilege"
	db, schema, view, priv, err = splitGrantID(id)
	if err != nil {
		t.Error(err)
	}
	if db != "database_name" {
		t.Errorf("Expected db to be database_name, got %v", db)
	}
	if schema != "" {
		t.Errorf("Expected schema to be blank, got %v", schema)
	}
	if view != "" {
		t.Errorf("Expected view to be blank, got %v", view)
	}
	if priv != "privilege" {
		t.Errorf("Expected priv to be privilege, got %v", priv)
	}

	// Bad ID
	id = "database|name-privilege"
	_, _, _, _, err = splitGrantID(id)
	if err == nil {
		t.Error("Expected an error, got none")
	}

	// Bad ID
	id = "database||||name-privilege"
	_, _, _, _, err = splitGrantID(id)
	if err == nil {
		t.Error("Expected an error, got none")
	}
}
