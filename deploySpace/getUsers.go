package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Users map[string]*UserConfig

type UserConfig struct {
	Type     string              `json:"userType"`
	Database *RelationalDBConfig `json:"database"`
}

type RelationalDBConfig struct {
	Name string `json:"name"`
	Plan string `json:"plan"`
}

func getNames(u Users) []string {
	var res []string

	for userName := range u {
		res = append(res, userName)
	}

	return res
}

func getSpaces(u Users) []string {
	return getNames(u)
}

func getUserConfigs(u Users) []*UserConfig {
	var res []*UserConfig

	for _, config := range u {
		res = append(res, config)
	}

	return res
}

func GetUsers(input []byte) (Users, error) {
	u := Users{}

	if err := json.Unmarshal(input, &u); err != nil {
		return nil, err
	}

	for userName, userConfig := range u {
		if userConfig.Type == "" {
			return nil, fmt.Errorf("Did not receive a role for user %v", userName)
		}
	}

	return u, nil
}

func GetUsersFromFile(filename string) (Users, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return GetUsers(fileBytes)
}
