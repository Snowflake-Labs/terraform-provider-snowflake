package resources

import (
	"testing"
)

func TestSplitID(t *testing.T) {
	// Vanilla
	id := "database-name_privilege"
	dbName, priv, err := splitID(id)
	if err != nil {
		t.Error(err)
	}
	if dbName != "database-name" {
		t.Errorf("Expected dbName to be database-name, got %v", dbName)
	}
	if priv != "privilege" {
		t.Errorf("Expected priv to be privilege, got %v", priv)
	}

	// DB with underscore
	id = "database_name_privilege"
	dbName, priv, err = splitID(id)
	if err != nil {
		t.Error(err)
	}
	if dbName != "database_name" {
		t.Errorf("Expected dbName to be database_name, got %v", dbName)
	}
	if priv != "privilege" {
		t.Errorf("Expected priv to be privilege, got %v", priv)
	}

	// Bad ID
	id = "database-name-privilege"
	dbName, priv, err = splitID(id)
	if err == nil {
		t.Error("Expected an error, got none")
	}
}
