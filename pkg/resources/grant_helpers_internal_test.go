package resources

import (
	"testing"
)

func TestSplitGrantID(t *testing.T) {
	// Vanilla
	id := "database_name|schema|view_name|privilege"
	grant, err := grantIDFromString(id)
	if err != nil {
		t.Error(err)
	}
	if grant.ResourceName != "database_name" {
		t.Errorf("Expected db to be database_name, got %v", grant.ResourceName)
	}
	if grant.SchemaName != "schema" {
		t.Errorf("Expected schema to be schema, got %v", grant.SchemaName)
	}
	if grant.ViewName != "view_name" {
		t.Errorf("Expected view to be view_name, got %v", grant.ViewName)
	}
	if grant.Privilege != "privilege" {
		t.Errorf("Expected priv to be privilege, got %v", grant.Privilege)
	}

	// No view
	id = "database_name|||privilege"
	grant, err = grantIDFromString(id)
	if err != nil {
		t.Error(err)
	}
	if grant.ResourceName != "database_name" {
		t.Errorf("Expected db to be database_name, got %v", grant.ResourceName)
	}
	if grant.SchemaName != "" {
		t.Errorf("Expected schema to be blank, got %v", grant.SchemaName)
	}
	if grant.ViewName != "" {
		t.Errorf("Expected view to be blank, got %v", grant.ViewName)
	}
	if grant.Privilege != "privilege" {
		t.Errorf("Expected priv to be privilege, got %v", grant.Privilege)
	}

	// Bad ID
	id = "database|name-privilege"
	grant, err = grantIDFromString(id)
	if err == nil {
		t.Error("Expected an error, got none")
	}

	// Bad ID
	id = "database||||name-privilege"
	grant, err = grantIDFromString(id)
	if err == nil {
		t.Error("Expected an error, got none")
	}
}
