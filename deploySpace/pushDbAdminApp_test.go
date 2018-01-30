package main

import (
	"fmt"
	"strings"
	"testing"
)

type PushProvisioner struct {
	calls []string
}

func (l *PushProvisioner) CallBluemixWithArgs(args ...string) error {

	stringsToCallOS := parseArgs(args...)
	l.calls = append(l.calls, strings.Join(stringsToCallOS, " "))

	argString := strings.Join(stringsToCallOS, " ")

	if strings.Contains(argString, "cf app Hyperion-admin") {
		return fmt.Errorf(EXPECTEDERROR)
	}

	return nil
}

func (l *PushProvisioner) allCalls() string {
	return strings.Join(l.calls, "\n")
}

func TestPushDbAdminAppBasic(t *testing.T) {
	provisioner := &TestProvisioner{}

	users := make(map[string]*UserConfig)
	users["InnovatorIan"] = &UserConfig{
		Type: "Innovator",
	}
	users["Broker1"] = &UserConfig{
		Type: "Broker",
	}
	users["Insurer1"] = &UserConfig{
		Type: "Insurer",
		Database: &RelationalDBConfig{
			Name: MYSQLDB,
			Plan: "Standard",
		},
	}
	users["Insurer2"] = &UserConfig{
		Type: "Insurer",
	}

	err := PushDbAdminApp(ProvisionInfo{
		users:       users,
		provisioner: provisioner,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(provisioner.calls) != 8 {
		t.Fatalf("Expected %v calls (target and push for each user) but got %v", 8, len(provisioner.calls))
	}
}

func TestOnlyPushIfAppDoesNotExist(t *testing.T) {
	provisioner := &PushProvisioner{}

	users := make(map[string]*UserConfig)
	users["InnovatorIan"] = &UserConfig{
		Type: "Innovator",
	}
	users["Hyperion"] = &UserConfig{
		Type: "Broker",
	}

	err := PushDbAdminApp(ProvisionInfo{
		users:       users,
		provisioner: provisioner,
	})
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(provisioner.allCalls(), "push -f manifest-InnovatorIan") {
		t.Fatal("Did not expect to push for InnovatorIan as 'cf app InnovatorIan' returned no error, indicating the app exists")
	}

	if !strings.Contains(provisioner.allCalls(), "cf app Hyperion-admin") {
		t.Fatal("Expected to check whether app existed but did not")
	}

	if !strings.Contains(provisioner.allCalls(), "push -f manifest-Hyperion") {
		t.Fatal("Expected to push an app for Hyperion as our mock dictates that it does not already exists")
	}
}
