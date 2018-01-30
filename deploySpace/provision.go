package main

import (
	"fmt"
	"strings"
)

const (
	MYSQLDB = "compose-for-mysql"
	DB2DB   = "db2oncloud"
)

type Provisioner interface {
	CallBluemixWithArgs(args ...string) error
}

type ProvisionInfo struct {
	users               Users
	deploymentFunctions map[string]*ActionFileConfig
	org                 string
	provisioner         Provisioner
}

func addQuotes(name string) string {
	return `"` + name + `"`
}

func removeSpaces(name string) string {
	return strings.Replace(name, " ", "", -1)
}

func Provision(provisionInfo ProvisionInfo) error {

	// // Target specified org
	err := provisionInfo.provisioner.CallBluemixWithArgs("target -o", provisionInfo.org)
	if err != nil {
		return err
	}

	// Vertical pipeline step - create space, create cloudant
	if err = ProvisionVertical(provisionInfo); err != nil {
		return err
	}

	// Install openwhisk CLI
	if err = provisionInfo.provisioner.CallBluemixWithArgs("plugin install Cloud-Functions -r Bluemix -f"); err != nil {
		return err
	}

	// Horizontal pipeline - Deploy cloud functions
	return ProvisionHorizontal(provisionInfo)
}

func ProvisionVertical(provisionInfo ProvisionInfo) error {

	// Vertical pipeline
	for userName, userConfig := range provisionInfo.users {

		err := provisionInfo.provisioner.CallBluemixWithArgs("cf create-space", userName)
		if err != nil {
			return err
		}

		if err = provisionInfo.provisioner.CallBluemixWithArgs("target -s", userName); err != nil {
			return err
		}

		if err = provisionInfo.provisioner.CallBluemixWithArgs("service create cloudantNoSQLDB Lite", removeSpaces(userName)+"-cloudant"); err != nil {
			return err
		}

		if userConfig.Database != nil && userConfig.Type != "Innovator" {
			serviceCall := fmt.Sprintf("service create %v %v %v-%v", userConfig.Database.Name, userConfig.Database.Plan, removeSpaces(userName), userConfig.Database.Name)
			if err = provisionInfo.provisioner.CallBluemixWithArgs(serviceCall); err != nil {
				return err
			}
		}

	}
	return nil
}

func ProvisionHorizontal(provisionInfo ProvisionInfo) error {

	var prevSpace string

	for userName, userConfig := range provisionInfo.users {

		for path, actionConfig := range provisionInfo.deploymentFunctions {

			if actionConfig.PublishToAll() || actionConfig.PublishToUser(userName) || actionConfig.PublishToType(userConfig.Type) {

				if userName != prevSpace {
					err := provisionInfo.provisioner.CallBluemixWithArgs("target -s", userName)
					if err != nil {
						return err
					}
					prevSpace = userName
				}

				call := fmt.Sprintf("wsk action update %v %v/%v", actionConfig.ActionName, path, actionConfig.ActionFile)

				if err := provisionInfo.provisioner.CallBluemixWithArgs(call); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
