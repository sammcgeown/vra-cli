/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sammcgeown/vra-cli/pkg/util/auth"
	"github.com/sammcgeown/vra-cli/pkg/util/config"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

var (
	targetConfig *types.Config
	RESTClient   *resty.Client
	// Config
	insecure = false
	cases    = []struct{ name, description, variableType, project, value string }{
		{"Test1", "Test 1 Description", "REGULAR", "vra-cli-testing", "Test1"},
		{"Test2", "Test 2 Description", "SECRET", "vra-cli-testing", "Test2"},
		{"Test3", "Test 3 Description", "RESTRICTED", "vra-cli-testing", "Test3"},
	}
)

func init() {
	// log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	// log.SetLevel(log.DebugLevel)

	log.Infoln("Initializing tests")
	targetConfig = config.GetConfigFromEnv()
	if err := auth.GetConnection(targetConfig, insecure); err != nil {
		log.Fatalln(err)
	}

}

func TestCreateVariable(t *testing.T) {
	RESTClient = auth.GetRestClient(targetConfig, insecure)

	for _, c := range cases {
		log.Infoln("Creating variable: ", c.name)
		variable, err := CreateVariable(RESTClient, c.name, c.description, c.variableType, c.project, c.value)
		if err != nil {
			log.Warnln(err)
		}
		assert.Equal(t, variable.Name, c.name)
		assert.Equal(t, variable.Type, c.variableType)
		assert.Equal(t, variable.Description, c.description)
		if c.variableType == "REGULAR" { // Regular variables should have a non-hidden value
			assert.Equal(t, variable.Value, c.value)
		}
		assert.Equal(t, variable.Project, c.project)
	}

}

func TestGetVariable(t *testing.T) {
	for _, c := range cases { // Test each case has been created
		RESTClient = auth.GetRestClient(targetConfig, insecure) // Re-initialize client
		log.Infoln("Getting variable: ", c.name)

		variable, err := GetVariable(RESTClient, "", c.name, c.project, "")

		if err != nil {
			log.Warnln(err)
		}
		if len(variable) == 0 {
			t.Errorf("No variables returned")
		} else if len(variable) > 1 {
			t.Errorf("More than one variable returned")
		} else {
			assert.Equal(t, variable[0].Name, c.name)
			assert.Equal(t, variable[0].Type, c.variableType)
			assert.Equal(t, variable[0].Description, c.description)
			if c.variableType == "REGULAR" { // Regular variables should have a non-hidden value
				assert.Equal(t, variable[0].Value, c.value)
			}
			assert.Equal(t, variable[0].Project, c.project)
		}
	}
}

func TestDeleteVariable(t *testing.T) {
	for _, c := range cases { // Delete each test case
		log.Infoln("Deleting variable: ", c.name)
		variable, err := GetVariable(RESTClient, "", c.name, c.project, "")
		if err != nil {
			log.Warnln(err)
		}
		assert.Assert(t, len(variable) == 1) // Getter should return exactly one variable
		deleted, err := DeleteVariable(RESTClient, variable[0].ID)
		assert.NilError(t, err) // Deleter should not return an error

		assert.Equal(t, deleted, true) // Deleter should return true
	}

}
