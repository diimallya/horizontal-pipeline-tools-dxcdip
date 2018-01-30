package main

import (
	"testing"
)

const exampleUsers = `{
  "Aon": {
    "userType": "Broker",
    "database": {
      "name": "compose-for-mysql",
      "plan": "Standard"
    }
  },
  "Beazley": {
    "userType": "Insurer",
    "database": {
      "name": "compose-for-mysql",
      "plan": "Standard"
    }
  },
  "Chaucer": {
    "userType": "Insurer",
    "database": {
      "name": "compose-for-mysql",
      "plan": "Standard"
    }
  },
  "Zurich": {
    "userType": "Insurer",
    "database": {
      "name": "compose-for-mysql",
      "plan": "Standard"
    }
  },
  "InnovatorIan": {
    "userType": "Insurer"
  }
}`

func TestGetUsers(t *testing.T) {
	users, err := GetUsers([]byte(exampleUsers))
	if err != nil {
		t.Fatal(err)
	}

	if val, exists := users["Aon"]; !exists {
		t.Fatal("Expected to have configuration about user Aon but did not")
	} else {
		if val.Type != "Broker" {
			t.Fatalf("Expected Aon to have user type 'Broker' but got %v", val.Type)
		}
		if val.Database == nil {
			t.Fatal("Expected Aon to have an associated relational database but did not")
		}
		if val.Database.Name != MYSQLDB {
			t.Fatalf("Expected Aon to have DB %v but had %v", MYSQLDB, val.Database.Name)
		}
	}

	if len(users) != 5 {
		t.Fatalf("There were %v users, expected 5", len(users))
	}
}

func TestGetUsersWithUserWithoutRole(t *testing.T) {
	_, err := GetUsers([]byte(`{"Rufus":{}}`))
	if err == nil {
		t.Fatal("Expected error but none was thrown")
	}
}
