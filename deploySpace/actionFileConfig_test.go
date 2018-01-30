package main

import (
	"testing"
)

const exampleActionFileConfig = `{
  "innovator": "InnovatorIan",
  "publishTo": ["NONE"],
  "name": "iansNewAction",
  "file": "iansNewAction.js"
}`

func TestGetActionFileConfig(t *testing.T) {
	a, err := getActionFileConfig([]byte(exampleActionFileConfig))
	if err != nil {
		t.Fatal(err)
	}

	if a.Innovator != "InnovatorIan" {
		t.Fatalf("Expected innovator InnovatorIan got %v", a.Innovator)
	}

	if a.PublishToAll() {
		t.Fatal("Expected not to be publishing to all but boolean returns true")
	}

	if a.PublishToUser("NotARealUser") {
		t.Fatal("Did not expect to be publishing to anywhere but InnovatorIan")
	}

	if !a.PublishToUser("InnovatorIan") {
		t.Fatal("Expected to publish to InnovatorIan but PublishToUser does not mark as true")
	}
}

func TestPublishToAllGivenBrokers(t *testing.T) {
	a := &ActionFileConfig{
		PublishTo: []string{BROKERS},
	}

	if a.PublishToAll() {
		t.Fatal("Did not expect BROKERS to pass 'PublishToAll' function")
	}
}
func TestPublishToAllBasic(t *testing.T) {
	a := &ActionFileConfig{
		PublishTo: []string{ALL},
	}

	if !a.PublishToAll() {
		t.Fatal("expected ALL to pass 'PublishToAll' function")
	}
}
func TestPublishToAllAllianz(t *testing.T) {
	a := &ActionFileConfig{
		PublishTo: []string{"Allianz"},
	}

	if a.PublishToAll() {
		t.Fatal("expected 'Allianz' not to pass 'PublishToAll' function")
	}
}

func TestPublishToType(t *testing.T) {
	a := &ActionFileConfig{
		PublishTo: []string{INSURERS},
	}

	if !a.PublishToType("Insurer") {
		t.Fatalf("Expected type 'Insurer' to pass 'PublishToType' check for %+v", a)
	}
}

func TestActionNameWithSpaceErrors(t *testing.T) {
	_, err := getActionFileConfig([]byte(`{
		"innovator": "InnovatorIan",
		"publishTo": ["NONE"],
		"name": "ians New Action",
		"file": "iansNewAction.js"
	}`))
	if err == nil {
		t.Fatalf("Expected an error to be thrown with a space in the action name but none was thrown")
	}
}
