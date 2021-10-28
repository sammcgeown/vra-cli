/*
Package orchestrator Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package orchestrator

import (
	"os"
	"testing"

	"github.com/sammcgeown/vra-cli/pkg/util/auth"
	"github.com/sammcgeown/vra-cli/pkg/util/config"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

var (
	APIClient = &types.APIClientOptions{
		Version: "2019-10-17",
		Debug:   false,
	}
)

func TestMain(m *testing.M) {
	// Configure Logging
	if APIClient.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Debug logging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// Configure API Client
	APIClient.Config = config.GetConfigFromEnv()
	err := auth.ValidateConfiguration(APIClient)
	if err != nil {
		log.Fatalln(err)
	}
	// Clean environment
	CleanUp()
	// Run tests
	code := m.Run()
	// Clean up after tests
	CleanUp()
	// Exit
	os.Exit(code)
}

func CleanUp() {

}

// func TestCreateProject(t *testing.T) {
// 	newProject, err := cloudassembly.CreateProject(APIClient, project.Name, project.Description, project.Administrators, project.Members, project.Viewers, nil, nil, 60, project.MachineNamingTemplate, &project.SharedResources)
// 	if err != nil {
// 		log.Warnln("Unable to create Project", err)
// 	}
// 	project.ID = newProject.ID
// 	assert.Equal(t, newProject.Name, project.Name)
// 	assert.Equal(t, newProject.Description, project.Description)
// 	assert.Equal(t, newProject.MachineNamingTemplate, project.MachineNamingTemplate)
// }

// func TestCreateVariable(t *testing.T) {
// 	//helpers.PrettyPrint(APIClient.Config)
// 	for _, c := range cases {
// 		log.Debugln("Creating variable:", c.name)
// 		variable, err := CreateVariable(APIClient, c.name, c.description, c.variableType, project.Name, c.value)
// 		if err != nil {
// 			log.Warnln(err)
// 		}
// 		assert.Equal(t, variable.Name, c.name)
// 		assert.Equal(t, variable.Type, c.variableType)
// 		assert.Equal(t, variable.Description, c.description)
// 		if c.variableType == "REGULAR" { // Regular variables should have a non-hidden value
// 			assert.Equal(t, variable.Value, c.value)
// 		}
// 		assert.Equal(t, variable.Project, project.Name)
// 	}

// }

// func TestGetVariable(t *testing.T) {
// 	for _, c := range cases { // Test each case has been created
// 		log.Debugln("Getting variable: ", c.name)

// 		variable, err := GetVariable(APIClient, "", c.name, project.Name, "")

// 		if err != nil {
// 			log.Warnln(err)
// 		}
// 		if len(variable) == 0 {
// 			t.Errorf("No variables returned")
// 		} else if len(variable) > 1 {
// 			t.Errorf("More than one variable returned")
// 		} else {
// 			assert.Equal(t, variable[0].Name, c.name)
// 			assert.Equal(t, variable[0].Type, c.variableType)
// 			assert.Equal(t, variable[0].Description, c.description)
// 			if c.variableType == "REGULAR" { // Regular variables should have a non-hidden value
// 				assert.Equal(t, variable[0].Value, c.value)
// 			}
// 			assert.Equal(t, variable[0].Project, project.Name)
// 		}
// 	}
// }

// func TestUpdateVariable(t *testing.T) {
// 	variable, err := GetVariable(APIClient, "", "Test1", project.Name, "")
// 	if err != nil {
// 		log.Warnln(err)
// 	}
// 	updateVariable := variable[0]

// 	updatedVariable, err := UpdateVariable(APIClient, updateVariable.ID, "Test1-Updated", "Test 1 Updated Description", "REGULAR", "UpdatedValue")
// 	if err != nil {
// 		log.Warnln(err)
// 	}
// 	assert.Equal(t, updatedVariable.Name, "Test1-Updated")
// 	assert.Equal(t, updatedVariable.Description, "Test 1 Updated Description")
// 	assert.Equal(t, updatedVariable.Type, "REGULAR")
// 	assert.Equal(t, updatedVariable.Value, "UpdatedValue")
// }

// func TestDeleteVariable(t *testing.T) {
// 	c := cases[2] // Delete one test case
// 	log.Debugln("Deleting variable: ", c.name)
// 	variable, err := GetVariable(APIClient, "", c.name, project.Name, "")
// 	if err != nil {
// 		log.Warnln(err)
// 	}
// 	deleted, err := DeleteVariable(APIClient, variable[0].ID)
// 	assert.Assert(t, len(variable) == 1) // Getter should return exactly one variable
// 	assert.NilError(t, err)              // Deleter should not return an error
// 	assert.Equal(t, deleted, true)       // Deleter should return true
// }

// func TestDeleteVariableByProject(t *testing.T) {
// 	log.Debugln("Deleting Variables in", project.Name)
// 	deletedVariables, vErr := DeleteVariableByProject(APIClient, true, project.Name)
// 	if vErr != nil {
// 		log.Warnln(vErr)
// 	}
// 	assert.Assert(t, len(deletedVariables) == len(cases)-1) // Should delete one less than the number of cases
// 	assert.NilError(t, vErr)                                // Deleter should not return an error
// }
