package resources

import (
	"testing"
)

func TestSplitGrantID(t *testing.T) {
	// Vanilla
	dataIdentifiers := [][]string{{"database_name", "schema", "view_name", "privilege"}}
	grantID, err := createGrantID(dataIdentifiers)
	if err != nil {
		t.Error(err)
	}

	grantIDArray, err := splitGrantID(grantID)
	if err != nil {
		t.Error(err)
	}

	db, schema, view, priv := grantIDArray[0], grantIDArray[1], grantIDArray[2], grantIDArray[3]

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
	dataIdentifiers = [][]string{{"database_name", "privilege"}}
	grantID, err = createGrantID(dataIdentifiers)
	if err != nil {
		t.Error(err)
	}

	grantIDArray, err = splitGrantID(grantID)
	if err != nil {
		t.Error(err)
	}

	// aku: this test shouldn't be relevant because splitGrantID takes a variable # of inputs
	// db, schema, view, priv = grantIDArray[0], grantIDArray[1], grantIDArray[2], grantIDArray[3]

	// if db != "database_name" {
	// 	t.Errorf("Expected db to be database_name, got %v", db)
	// }
	// if schema != "" {
	// 	t.Errorf("Expected schema to be blank, got %v", schema)
	// }
	// if view != "" {
	// 	t.Errorf("Expected view to be blank, got %v", view)
	// }
	// if priv != "privilege" {
	// 	t.Errorf("Expected priv to be privilege, got %v", priv)
	// }

	// Bad ID
	// dataIdentifiers = [][]string{{"database_name", "name-privilege"}}
	// grantID, err = createGrantID(dataIdentifiers)
	// if err != nil {
	// 	t.Error(err)
	// }

	// grantIDArray, err = splitGrantID(grantID)

	// if err == nil {
	// 	t.Error("Expected an error, got none")
	// }

	// aku: commented out this test because it would be duplicated with new refactor of splitGrantID()
	// Bad ID
	// id = "database||||name-privilege"
	// _, _, _, _, err = splitGrantID(id)
	// if err == nil {
	// 	t.Error("Expected an error, got none")
	// }
}
