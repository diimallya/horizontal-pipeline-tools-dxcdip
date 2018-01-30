package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var pOrg *string
var pUsersFile *string
var pCloudFunctionsDir *string
var pManifestTemplate *string
var pDatabaseAdminDir *string
var functionsToDeploy map[string]*ActionFileConfig

func init() {
	functionsToDeploy = make(map[string]*ActionFileConfig)

	// Set up command line flag parsing parameters.
	pOrg = flag.String("o", "", "Your bluemix organisation")
	pUsersFile = flag.String("usersFile", "", "The file describing brokers, insurers and innovators")
	pCloudFunctionsDir = flag.String("cloudFunctions", "", "The top level folder for all cloud functions")
	pManifestTemplate = flag.String("manifestTemplate", "", "The template Cloud Foundry manifest which will be used to generate cf manifest files for each insurance organisation or innovator")
	pDatabaseAdminDir = flag.String("databaseAdminDir", "", "The directory which contains the code to setup the databases with default schemas, tables and documents")
}

func main() {
	// Parse the flags.
	flag.Parse()

	// Capture any flags which have not been set.
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			flag.PrintDefaults()
			log.Fatalf("Required flag -%v (%v) was not set.", f.Name, f.Usage)
		}
	})

	// Get users
	users, err := GetUsersFromFile(*pUsersFile)
	if err != nil {
		log.Fatal(err)
	}

	// Get functions
	if err := filepath.Walk(*pCloudFunctionsDir, readInActionConfigs); err != nil {
		log.Fatal(err)
	}

	if len(functionsToDeploy) == 0 {
		log.Println("Warning: No functions found to deploy")
	}

	// Generate manifests for each space
	for user, userType := range users {
		m := &ManifestDetails{
			Spacename:        removeSpaces(user),
			UserConfig:       userType,
			ManifestTemplate: *pManifestTemplate,
			OutputDirectory:  *pDatabaseAdminDir,
		}
		err := GenerateManifest(m)
		if err != nil {
			log.Fatal(err)
		}
	}

	provisionInfo := ProvisionInfo{
		users:               users,
		deploymentFunctions: functionsToDeploy,
		org:                 *pOrg,
		provisioner:         &LiveProvisioner{},
	}

	// Provision relevant spaces and functions
	if err := Provision(provisionInfo); err != nil {
		log.Fatal(err)
	}

	// Change directory to manifest one
	if err = os.Chdir(*pDatabaseAdminDir); err != nil {
		log.Fatal(err)
	}
	// Push db admin app in each space
	if err := PushDbAdminApp(provisionInfo); err != nil {
		log.Fatal(err)
	}

}
