# DxC Pipeline Tools

This repository contains the scripts, configurations, and programs written for the DevOps pipeline to deploy spaces, services and functions across Bluemix.

[This is a link to the pipeline](https://console.bluemix.net/devops/pipelines/4f66dc99-d97c-4379-bc3b-303af1e06de6?env_id=ibm%3Ayp%3Aus-south). 
You will only be able to see it if you have the appropriate access rights.

[Similarly, here is the trello board](https://trello.com/b/yBtaZipF/dxc-assistance).

The [parties file](dbsetup/parties.json) is the most important in the repo, as it describes the current spaces deployed, for each insurance organisation (IO), or innovator.

## Pipeline overview

The pipeline takes the party (i.e. insurance organisation or innovator) found in the parties file and creates a space inside the targeted org for each party.
The targeted org is currently `DXC-Digital-Innovation-Platform`. 
This can be changed via the `ORG` environment variable inside the configuration for the `Deploy Spaces` stage of the pipeline.

Inside each space, the following is provisioned:
- A Cloudant service instance named `(spacename)-cloudant`, containing the JSON documents detailed in `dbsetup/cloudant/`
- The relevant Cloud Functions (or OpenWhisk actions) as applicable to the user ([more on this](#cloud-function-deployment-horizontal-pipeline)).
- A mysql database, if specified in the `dbsetup/parties.json` file, containing tables as detailed [here](dbsetup/mysql/create_table_dxcdip.sql).
The `dxcdip.parties` table will contain each party (name, uuid, and type) as described in the [parties file](dbsetup/parties.json).

### IO and Innovator Deployment (Vertical Pipeline)

To add a new IO or Innovator, you must first add them to the [parties file](dbsetup/parties.json). In this file, each party is described as per the following example:
```
"name": {
    "userType": "Broker/Insurer/Innovator",
    "database": {
        "name": "compose-for-mysql/db2forcloud",
        "plan": "Standard"
    },
    "uuid": "A valid UUID"
}
```

This describes the following fields:
- `name` is the user's name e.g. "Hyperion" or "Eagle Star"
- `userType` is "Broker", "Insurer", or "Innovator"
- `database` describes the relational database name and plan required. Currently the pipeline only fully supports the `compose-for-mysql` database. 
The name must be the name as it appears in the bluemix catalogue from the CLI, e.g. `compose-for-mysql`, or `db2forcloud`.
The plan must be a valid plan for the database type, e.g. "Standard" for `compose-for-mysql`.
- `uuid` must be a valid unique ID. We generated the first few from this [Version 4 UUID generator](https://www.uuidgenerator.net/). 
This is important as it is used as the primary key in the `dxcdip.parties` table in the relational database.

Once you have ensured the parties file contains valid JSON, check it in to this repository.
Once checked in, [the pipeline](https://console.bluemix.net/devops/pipelines/4f66dc99-d97c-4379-bc3b-303af1e06de6?env_id=ibm%3Ayp%3Aus-south) should be started.
If it has not, you can manually start it. This will handle provisioning of the new party based on the details provided.

Note: If you removed any party that has already been provisioned in the parties.json, any new run of the pipeline will have no effect -- i.e. it will not remove any configuration. Removing an invidivual space remains a manual step.

### Cloud Function Deployment (Horizontal Pipeline)

The program deploys functions found in subdirectories of the `cloudFunctions/` folder in the repo.
For each function, the function file and a configuration file (`.actionconfig.json`) are required.

The configuration file describes the innovator who created the action, and where to deploy it. 
It also describes the action name, and action file name, as IBM Cloud Functions require these both.
IBM Cloud Functions also require the action names to not have spaces in them.

#### How to specify where to deploy actions

Using the `publishTo` list in `.actionconfig.json`, you can specify where you want the action to be published.

Possible options for this list are:
- `"publishTo": ["NONE"]`: This will just deploy to the Innovator's space.
- `"publishTo": ["ALL"]`: This will deploy to every insurance organisation.
- `"publishTo": ["INSURERS"]`: This will deploy to all insurers as named in `parties.json`.
- `"publishTo": ["BROKERS"]`: This will deploy to all brokers as named in `parties.json`.
- `"publishTo": ["Aon", "Devonshire", "Zurich"]`: This will deploy to each named insurance organisation.

It is possible to combine the fields above e.g. `"publishTo": ["INSURERS", "BROKERS"]`, or `"publishTo": ["INSURERS", "Aon"]`.

### Pipeline Credentials

If the pipeline were to be set up elsewhere, 2 environment variables would need to be set. These are:
- `ORG` - should be the Bluemix organisation which deployments should be done under.
- `API_KEY` - should be an API key for Bluemix with "manager" access rights, as these are sufficient to create spaces.

#### Removing all parties provisioned by the pipeline

To remove ALL parties provisioned by the pipeline (as currently described in the parties file), in a shell run [this script](demo_scripts/deleteSpaces.sh), e.g. `./demo_scripts/deleteSpaces.sh dbsetup/parties.json`.

WARNING: This script is non-discriminatory. Once run there is no means of reversing - however you can provision it all again (losing any changes made outside of the pipeline) by re-running the pipeline.

### Gotchas

Currently, adding a new party will not append the party to the `dxcdip.parties` table in the relational database owned by other parties. 
This is something that we would address in a further engagement pending any contractual agreement.