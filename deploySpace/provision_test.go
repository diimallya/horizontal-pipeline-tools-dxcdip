package main

import (
	"strings"
	"testing"
)

type TestProvisioner struct {
	calls []string
}

func (l *TestProvisioner) CallBluemixWithArgs(args ...string) error {

	stringsToCallOS := parseArgs(args...)

	l.calls = append(l.calls, strings.Join(stringsToCallOS, " "))

	return nil
}

func (l *TestProvisioner) allCalls() string {
	return strings.Join(l.calls, "\n")
}

func TestProvisionEndToEnd(t *testing.T) {
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

	org := "testOrg"

	actionFilePath := "my/dev/path/to/action"
	a := &ActionFileConfig{
		ActionFile: "testDevActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{"NONE"},
	}
	actionFilePath2 := "my/prod/path/to/action"
	a2 := &ActionFileConfig{
		ActionFile: "testProdActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{"ALL"},
	}

	functions := make(map[string]*ActionFileConfig)

	functions[actionFilePath] = a
	functions[actionFilePath2] = a2

	if err := Provision(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	if len(provisioner.calls) == 0 {
		t.Fatal("Provisioner was never called")
	}

	if provisioner.calls[0] != "target -o "+org {
		t.Fatalf("First call to provision should target org, got %v", provisioner.calls[0])
	}

	if !strings.Contains(provisioner.allCalls(), "plugin install Cloud-Functions -r Bluemix -f") {
		t.Fatal("Expected to install Cloud-Functions plugin, but was not called")
	}

	if len(provisioner.calls) != 24 {
		t.Fatalf("Expected 24 bluemix calls, got %v: \n%v", len(provisioner.calls), provisioner.allCalls())
	}
}

func TestProvisionVertical(t *testing.T) {

	provisioner := &TestProvisioner{}

	userName := "InnovatorIan"

	users := make(map[string]*UserConfig)
	users[userName] = &UserConfig{
		Type: "Innovator",
	}

	org := "testOrg"

	functions := make(map[string]*ActionFileConfig)

	if err := ProvisionVertical(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	if len(provisioner.calls) == 0 {
		t.Fatal("Provisioner was never called")
	}

	if provisioner.calls[0] != "cf create-space "+userName {
		t.Fatalf("Second call to provision should create space for InnovatorIan, got %v", provisioner.calls[0])
	}
	if provisioner.calls[1] != "target -s "+userName {
		t.Fatalf("Second call to provision should create space for InnovatorIan, got %v", provisioner.calls[1])
	}

	if !strings.Contains(provisioner.calls[2], "cloudant") {
		t.Fatalf("Expected 3rd call to provision cloudant, got [%v]", provisioner.calls[2])
	}
}
func TestProvisionVerticalWithRelationalDB(t *testing.T) {

	provisioner := &TestProvisioner{}

	userName := "Beazley"

	users := make(map[string]*UserConfig)
	users[userName] = &UserConfig{
		Type: "Insurer",
		Database: &RelationalDBConfig{
			Name: MYSQLDB,
			Plan: "Standard",
		},
	}

	org := "testOrg"

	if err := ProvisionVertical(ProvisionInfo{
		users:       users,
		org:         org,
		provisioner: provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	if len(provisioner.calls) == 0 {
		t.Fatal("Provisioner was never called")
	}

	if strings.Count(provisioner.allCalls(), "service create compose-for-mysql Standard Beazley-compose-for-mysql") != 1 {
		t.Fatal("Expected to provision a single compose-for-mysql but did not", provisioner.allCalls())
	}

	if strings.Count(provisioner.allCalls(), "Aon") != 0 {
		t.Fatal("Did not expect to see user Aon but did")
	}
}
func TestProvisionVerticalWithDb2DB(t *testing.T) {

	provisioner := &TestProvisioner{}

	userName := "Aon"

	users := make(map[string]*UserConfig)
	users[userName] = &UserConfig{
		Type: "Insurer",
		Database: &RelationalDBConfig{
			Name: DB2DB,
			Plan: "Small",
		},
	}

	org := "testOrg"

	if err := ProvisionVertical(ProvisionInfo{
		users:       users,
		org:         org,
		provisioner: provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	if len(provisioner.calls) == 0 {
		t.Fatal("Provisioner was never called")
	}

	if strings.Count(provisioner.allCalls(), "service create db2oncloud Small Aon-db2") != 1 {
		t.Fatal("Expected to provision a single db2oncloud but did not")
	}
}
func TestDoNotProvisionVerticalDb2DBForInnovator(t *testing.T) {

	provisioner := &TestProvisioner{}

	userName := "InnovatorIan"

	users := make(map[string]*UserConfig)
	users[userName] = &UserConfig{
		Type: "Innovator",
		Database: &RelationalDBConfig{
			Name: DB2DB,
			Plan: "Small",
		},
	}

	org := "testOrg"

	if err := ProvisionVertical(ProvisionInfo{
		users:       users,
		org:         org,
		provisioner: provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	if len(provisioner.calls) == 0 {
		t.Fatal("Provisioner was never called")
	}

	if strings.Count(provisioner.allCalls(), "service create db2oncloud Small InnovatorIan-db2") != 0 {
		t.Fatalf("Expected not to provision a single db2oncloud but did: \n%v", provisioner.allCalls())
	}
}
func TestProvisionHorizontalBasic(t *testing.T) {

	provisioner := &TestProvisioner{}

	userName := "InnovatorIan"

	users := make(map[string]*UserConfig)
	users[userName] = &UserConfig{
		Type: "Innovator",
	}

	org := "testOrg"

	actionFilePath := "my/path/to/action"
	a := &ActionFileConfig{
		ActionFile: "testActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{"ALL"},
	}

	functions := make(map[string]*ActionFileConfig)

	functions[actionFilePath] = a

	if err := ProvisionHorizontal(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	if len(provisioner.calls) == 0 {
		t.Fatal("Provisioner was never called")
	}

	if provisioner.calls[0] != "target -s "+userName {
		t.Fatalf("Second call to provision should create space for InnovatorIan, got %v", provisioner.calls[0])
	}

	expectedCall := "wsk action update " + a.ActionName + " " + actionFilePath + "/" + a.ActionFile
	if provisioner.calls[1] != expectedCall {
		t.Fatalf("Expected [%v], got [%v]", expectedCall, provisioner.calls[1])
	}
}
func TestProvisionHorizontalProdAction(t *testing.T) {

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
	}
	users["Insurer2"] = &UserConfig{
		Type: "Insurer",
	}

	org := "testOrg"

	actionFilePath := "my/path/to/action"
	a := &ActionFileConfig{
		ActionFile: "testActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{"ALL"},
	}

	functions := make(map[string]*ActionFileConfig)

	functions[actionFilePath] = a

	if err := ProvisionHorizontal(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	calls := provisioner.allCalls()

	if strings.Count(calls, "wsk action update") != 4 {
		t.Fatalf("Expected wsk action update to occur 4 times, got %v", strings.Count(calls, "wsk action update"))
	}

	if strings.Count(calls, "Broker1") != 1 {
		t.Fatalf("Expected to deploy once to Broker1")
	}
	if strings.Count(calls, "Insurer1") != 1 {
		t.Fatalf("Expected to deploy once to Insurer1")
	}
	if strings.Count(calls, "Insurer2") != 1 {
		t.Fatalf("Expected to deploy once to Insurer2")
	}
	if strings.Count(calls, "InnovatorIan") != 1 {
		t.Fatalf("Expected to deploy once to InnovatorIan")
	}
}
func TestProvisionHorizontalDevAction(t *testing.T) {

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
	}
	users["Insurer2"] = &UserConfig{
		Type: "Insurer",
	}

	org := "testOrg"

	actionFilePath := "my/path/to/action"
	a := &ActionFileConfig{
		ActionFile: "testActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{"NONE"},
	}

	functions := make(map[string]*ActionFileConfig)

	functions[actionFilePath] = a

	if err := ProvisionHorizontal(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	calls := provisioner.allCalls()

	if strings.Count(calls, "wsk action update") != 1 {
		t.Fatalf("Expected wsk action update to occur 1 time, got %v", strings.Count(calls, "wsk action update"))
	}

	if strings.Count(calls, "InnovatorIan") != 1 {
		t.Fatalf("Expected to deploy once to InnovatorIan")
	}
}

func TestProvisionHorizontalDevActionTargetsOnce(t *testing.T) {
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
	}
	users["Insurer2"] = &UserConfig{
		Type: "Insurer",
	}

	org := "testOrg"

	actionFilePath := "my/path/to/action"
	a := &ActionFileConfig{
		ActionFile: "testActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{"NONE"},
	}
	actionFilePath2 := "my/prod/path/to/action"
	a2 := &ActionFileConfig{
		ActionFile: "testProdActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{"ALL"},
	}

	functions := make(map[string]*ActionFileConfig)

	functions[actionFilePath] = a
	functions[actionFilePath2] = a2

	if err := ProvisionHorizontal(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	calls := provisioner.allCalls()

	if strings.Count(calls, "target -s InnovatorIan") != 1 {
		t.Fatalf("Expected to only target InnovatorIan's space once during horizontal pipeline, got %v: \n%v", strings.Count(calls, "target -s InnovatorIan"), calls)
	}
}

func TestProvisionHorizontalDeployActionToCustomList(t *testing.T) {
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
	}

	org := "testOrg"

	actionFilePath := "my/path/to/action"
	a := &ActionFileConfig{
		ActionFile: "testActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{"Broker1"},
	}

	functions := make(map[string]*ActionFileConfig)

	functions[actionFilePath] = a

	if err := ProvisionHorizontal(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	calls := provisioner.allCalls()

	if strings.Count(calls, "wsk action update") != 2 {
		t.Fatalf("Expected to deploy the action to InnovatorIan and Broker1, but only received wsk action update %v times: \n%v", strings.Count(calls, "wsk action update"), calls)
	}
	if strings.Count(calls, "Insurer1") != 0 {
		t.Fatal("Expected not to deploy the action to Insurer1")
	}
}
func TestProvisionHorizontalDeployActionToInsurers(t *testing.T) {
	provisioner := &TestProvisioner{}

	users := make(map[string]*UserConfig)
	users["InnovatorIan"] = &UserConfig{
		Type: "Innovator",
	}
	users["Broker1"] = &UserConfig{
		Type: "Broker",
	}
	users["Broker2"] = &UserConfig{
		Type: "Broker",
	}
	users["Insurer1"] = &UserConfig{
		Type: "Insurer",
	}
	users["Insurer2"] = &UserConfig{
		Type: "Insurer",
	}

	org := "testOrg"

	actionFilePath := "my/path/to/action"
	a := &ActionFileConfig{
		ActionFile: "testActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{INSURERS},
	}

	functions := make(map[string]*ActionFileConfig)

	functions[actionFilePath] = a

	if err := ProvisionHorizontal(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	calls := provisioner.allCalls()

	if strings.Count(calls, "wsk action update") != 3 {
		t.Fatalf("Expected to deploy the action to InnovatorIan and Broker1 and Broker2, but only received wsk action update %v times: \n%v", strings.Count(calls, "wsk action update"), calls)
	}
	if strings.Count(calls, "Broker1") != 0 {
		t.Fatal("Expected not to deploy the action to Broker1")
	}
}

func TestProvisionHorizontalDeployActionToBrokers(t *testing.T) {
	provisioner := &TestProvisioner{}

	users := make(map[string]*UserConfig)
	users["InnovatorIan"] = &UserConfig{
		Type: "Innovator",
	}
	users["Broker1"] = &UserConfig{
		Type: "Broker",
	}
	users["Broker2"] = &UserConfig{
		Type: "Broker",
	}
	users["Insurer1"] = &UserConfig{
		Type: "Insurer",
	}
	users["Insurer2"] = &UserConfig{
		Type: "Insurer",
	}

	org := "testOrg"

	actionFilePath := "my/path/to/action"
	a := &ActionFileConfig{
		ActionFile: "testActionFile.js",
		ActionName: "testAction",
		Innovator:  "InnovatorIan",
		PublishTo:  []string{BROKERS},
	}

	functions := make(map[string]*ActionFileConfig)

	functions[actionFilePath] = a

	if err := ProvisionHorizontal(ProvisionInfo{
		users:               users,
		deploymentFunctions: functions,
		org:                 org,
		provisioner:         provisioner,
	}); err != nil {
		t.Fatal(err)
	}

	calls := provisioner.allCalls()

	if strings.Count(calls, "wsk action update") != 3 {
		t.Fatalf("Expected to deploy the action to InnovatorIan and Broker1 and Broker2, but only received wsk action update %v times: \n%v", strings.Count(calls, "wsk action update"), calls)
	}
	if strings.Count(calls, "Insurer1") != 0 {
		t.Fatal("Expected not to deploy the action to Insurer1")
	}
}
