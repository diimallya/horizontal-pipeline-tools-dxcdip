package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	CONFIGFILENAME = ".actionconfig.json"
	ALL            = "ALL"
	NONE           = "NONE"
	BROKERS        = "BROKERS"
	INSURERS       = "INSURERS"
)

type ActionFileConfig struct {
	Innovator  string   `json:"innovator"`
	PublishTo  []string `json:"publishTo"`
	ActionName string   `json:"name"`
	ActionFile string   `json:"file"`
}

func (a *ActionFileConfig) PublishToAll() bool {
	for _, item := range a.PublishTo {
		if strings.ToUpper(item) == ALL {
			return true
		}
	}
	return false
}

func (a *ActionFileConfig) PublishToType(InsuranceOrgType string) bool {
	return strings.Contains(strings.Join(a.PublishTo, " "), strings.ToUpper(InsuranceOrgType))
}

func (a *ActionFileConfig) PublishToUser(user string) bool {
	if user == a.Innovator {
		return true
	}
	for _, name := range a.PublishTo {
		if user == name {
			return true
		}
	}
	return false

}

func readInActionConfigs(path string, info os.FileInfo, err error) error {
	if err != nil {
		return nil
	}

	if info.Name() != CONFIGFILENAME {
		return nil
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	a, err := getActionFileConfig(fileBytes)
	if err != nil {
		return err
	}

	// Check whether the function file specified exists
	actionFile := filepath.Dir(path) + "/" + a.ActionFile
	if _, err := os.Stat(actionFile); os.IsNotExist(err) {
		return err
	}

	functionsToDeploy[filepath.Dir(path)] = a

	return nil
}

func getActionFileConfig(input []byte) (*ActionFileConfig, error) {
	a := &ActionFileConfig{}

	if err := json.Unmarshal(input, &a); err != nil {
		return nil, err
	}

	if strings.Contains(a.ActionName, " ") {
		return nil, fmt.Errorf("Action names cannot contain spaces - see %v by %v", a.ActionName, a.Innovator)
	}

	return a, nil
}
