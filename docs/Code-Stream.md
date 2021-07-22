# Code Stream

Currently implemented commands:
* [Pipelines](#Working-with-Pipelines)
* [Variables](#Working-with-Variables)
* [Executions](#Working-with-Executions)
* [Endpoints](#Working-with-Endpoints)
* [Custom Integrations](#Working-with-Custom-Integrations)

## Working with Pipelines

Getting and filtering pipelines
```bash
# List all pipelines
vra-cli get pipeline
# List all pipelines in a project
vra-cli get pipeline --project "Field Demo"
# Get a pipeline by ID
vra-cli get pipeline --id 7a3b41af-0e49-4e3d-999b-6c4c5ec55956
# Get a pipeline by name
vra-cli get pipeline --name "vra-CreateVariable"
```

Exporting pipelines:
```bash
# Export a specific pipeline to current location
vra-cli get pipeline --name "vra-CreateVariable"
# Export a specific pipeline to a specific location
vra-cli get pipeline --name "vra-CreateVariable" --exportPath /path/to/my/folder
# Export all pipelines
vra-cli get pipeline
# Export all pipelines in a project
vra-cli get pipeline --project "Field Demo"
```

Importing pipelines:
```bash
# Import a yaml definition
vra-cli create pipeline --importPath /my/yaml-pipeline.yaml
# Import a folder of YAML files (will attempt to import all YAML files in the folder - .yml/.yaml)
vra-cli create pipeline --importPath /Users/sammcgeown/Desktop/vra-cli/pipelines
# Update an existing pipeline
# Note: You cannot change the pipeline name - this
#       will result in a new Pipeline being created
vra-cli update pipeline --importPath /my/updated-pipe.yaml
# Update existing pipelines from folder
vra-cli update pipeline --importPath /Users/sammcgeown/Desktop/vra-cli/pipelines
# Import a pipeline to a specific Project (overriding the YAML definition)
vra-cli create pipeline --importPath export/pipelines/Field-Demo-Chat-App.yaml --project "Field Demo"
```


Delete a pipeline:
```bash
# Delete pipeline by ID
vra-cli delete pipeline --id 7a3b41af-0e49-4e3d-999b-6c4c5ec55956
```

## Working with Variables

```bash
# Get all variables
vra-cli get variable
# Get a variable by ID
vra-cli get variable --id 50613ab6-6f25-4976-8b3e-5be7a4bc60eb
# Get a variable by name
vra-cli get variable --name vra-cli
# Create a new variable manually
vra-cli create variable --name cli-demo --project "Field Demo"  --type REGULAR --value "New variable..." --description "Now from the CLI\!"

# Export all variables to variables.yaml
vra-cli get variable
# Export all variables to /your/own/filename.yaml
vra-cli get variable --exportPath /your/own/filename.yaml

# Create new variables from file
vra-cli create variable --importfile variables.yaml
# Create new variables from file, overwrite the Project
vra-cli create variable --importfile variables.yaml --project TestProject

# Update existing variables from file
vra-cli update variable --importfile variables.yaml
```
*Note that SECRET variables will not export, so if you export your secrets, be sure to add the value data before re-importing them!*

## Working with Executions

```bash
# List all executions
vra-cli get execution
# View an execution by ID
vra-cli get execution --id 9cc5aedc-db48-4c02-a5e4-086de3160dc0
# View executions of a specific pipeline
get execution --name vra-authenticateUser
# View executions by status
vra-cli get execution --status Failed
```

Create a new execution of a pipeline:
```bash
# Get the input form of the pipeline to execute
vra-cli get pipeline --id 7a3b41af-0e49-4e3d-999b-6c4c5ec55956 --form
# Outputs:
# {
#   "vraFQDN": "",
#   "vraUserName": "",
#   "vraUserPassword": ""
# }

# Create a new execution with the input form from above
vra-cli create execution --id 7a3b41af-0e49-4e3d-999b-6c4c5ec55956 --inputs '{
  "vraFQDN": "vra8-test-ga.cmbu.local",
  "vraUserName": "fakeuser",
  "vraUserPassword": "fakeuser"
}' --comments "Executed from vra-cli"
# Outputs
# Execution /codestream/api/executions/9cc5aedc-db48-4c02-a5e4-086de3160dc0 created

# Inspect the new execution
vra-cli get execution --id 9cc5aedc-db48-4c02-a5e4-086de3160dc0
```



## Working with Endpoints
Getting Endpoints
```bash
# Get all endpoints
vra-cli get endpoint
# Get endpoints by project
vra-cli get endpoint --project "Field Demo"
# Get endpoint by Name
vra-cli get endpoint --name "My-Git-Endpoint"
# Get endpoint by Project and Type
vra-cli get endpoint --type "git" --project "Field Demo"
```

Exporting endpoints:
```bash
# Export all endpoints
vra-cli get endpoint --exportPath my-endpoints/
# Export endpoint by Name
vra-cli get endpoint --name "My-Git-Endpoint"
```

Importing endpoints
```bash
# Create a new endpoint
vra-cli create endpoint --importPath /path/to/my/endpoint.yaml
# Update an existing endpoint
# Note: You cannot change the endpoint name - this
#       will result in a new endpoint being created
vra-cli update endpoint --importPath updated-endpoint.yaml
# Import an endpoint to a specific Project (overriding the YAML)
vra-cli create endpoint --importPath /path/to/my/endpoint.yaml --project "Field Demo"
```

Delete an endpoint
```bash
# Delete endpoint by ID
vra-cli delete endpoint --id 8c36f59a-2fcf-4039-8b48-1026f601a4b0
```
## Working with Custom Integrations

```bash
# Get all custom integrations
vra-cli get customintegration
# Get custom integration by id
vra-cli get customintegration --id c145b52e-c797-49d1-88a5-1d70e7788d03
# Get custom integration by name
vra-cli get customintegration --name base64Encode
```
