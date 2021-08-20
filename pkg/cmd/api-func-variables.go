/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func getVariable(id, name, project, exportPath string) ([]*types.VariableResponse, error) {
	var arrVariables []*types.VariableResponse
	//var qParams = make(map[string]string)
	client := resty.New()

	// Get by ID
	if id != "" {
		v, e := getVariableByID(id)
		arrVariables = append(arrVariables, v)
		return arrVariables, e
	}
	if name != "" && project != "" {
		qParams["$filter"] = "((name eq '" + name + "') and (project eq '" + project + "'))"
	} else {
		// Get by name
		if name != "" {
			qParams["$filter"] = "(name eq '" + name + "')"
		}
		// Get by project
		if project != "" {
			qParams["$filter"] = "(project eq '" + project + "')"
		}
	}
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.DocumentsList{}).
		SetAuthToken(targetConfig.AccessToken).
		Get("https://" + targetConfig.Server + "/pipeline/api/variables")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}

	log.Debugln(queryResponse.Request.URL)

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.VariableResponse{}
		mapstructure.Decode(value, &c)
		arrVariables = append(arrVariables, &c)
		if exportPath != "" {
			exportVariable(c, exportPath)
		}
	}
	return arrVariables, err
}

// getVariableByID - get Code Stream Variable by ID
func getVariableByID(id string) (*types.VariableResponse, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.VariableResponse{}).
		SetAuthToken(targetConfig.AccessToken).
		Get("https://" + targetConfig.Server + "/pipeline/api/variables/" + id)
	if queryResponse.IsError() {
		log.Errorln("GET Variable failed", err)
	}
	return queryResponse.Result().(*types.VariableResponse), err
}

// createVariable - Create a new Code Stream Variable
func createVariable(name string, description string, variableType string, project string, value string) (*types.VariableResponse, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetBody(
			types.VariableRequest{
				Project:     project,
				Kind:        "VARIABLE",
				Name:        name,
				Description: description,
				Type:        variableType,
				Value:       value,
			}).
		SetHeader("Accept", "application/json").
		SetResult(&types.VariableResponse{}).
		SetError(&types.Exception{}).
		SetAuthToken(targetConfig.AccessToken).
		Post("https://" + targetConfig.Server + "/pipeline/api/variables")
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.VariableResponse), err
}

// updateVariable - Create a new Code Stream Variable
func updateVariable(id string, name string, description string, typename string, value string) (*types.VariableResponse, error) {
	variable, _ := getVariableByID(id)
	if name != "" {
		variable.Name = name
	}
	if description != "" {
		variable.Description = description
	}
	if typename != "" {
		variable.Type = typename
	}
	if value != "" {
		variable.Value = value
	}
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetBody(variable).
		SetHeader("Accept", "application/json").
		SetResult(&types.VariableResponse{}).
		SetError(&types.Exception{}).
		SetAuthToken(targetConfig.AccessToken).
		Put("https://" + targetConfig.Server + "/pipeline/api/variables/" + id)
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.VariableResponse), err
}

// deleteVariable - Delete a Code Stream Variable
func deleteVariable(id string) (*types.VariableResponse, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.VariableResponse{}).
		SetAuthToken(targetConfig.AccessToken).
		Delete("https://" + targetConfig.Server + "/pipeline/api/variables/" + id)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*types.VariableResponse), err
}

func deleteVariableByProject(project string) ([]*types.VariableResponse, error) {
	var deletedVariables []*types.VariableResponse
	Variables, err := getVariable("", "", project, "")
	if err != nil {
		return nil, err
	}
	confirm := helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(Variables)) + " variables in " + project + ", are you sure?")
	if confirm {

		for _, Variable := range Variables {
			deletedVariable, err := deleteVariable(Variable.ID)
			if err != nil {
				log.Warnln("Unable to delete "+Variable.Name, err)
			}
			deletedVariables = append(deletedVariables, deletedVariable)
		}
		return deletedVariables, nil
	} else {
		return nil, errors.New("user declined")
	}
}

// exportVariable - Export a variable to YAML
func exportVariable(variable interface{}, exportPath string) {
	var exportFile string
	// variable will be a types.VariableResponse, so lets remap to types.VariableRequest
	c := types.VariableRequest{}
	mapstructure.Decode(variable, &c)
	yaml, err := yaml.Marshal(c)
	if err != nil {
		log.Errorln("Unable to export variable ", c.Name)
	}

	if filepath.Ext(exportPath) != ".yaml" {
		exportFile = filepath.Join(exportPath, "variables.yaml")
	} else {
		exportFile = exportPath
	}

	file, err := os.OpenFile(exportFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	file.WriteString("---\n" + string(yaml))
}

// importVariables - Import variables from the filePath
func importVariables(filePath string) []types.VariableRequest {
	var returnVariables []types.VariableRequest
	filename, _ := filepath.Abs(filePath)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(yamlFile)
	decoder := yaml.NewDecoder(reader)
	var request types.VariableRequest
	for decoder.Decode(&request) == nil {
		returnVariables = append(returnVariables, request)
	}
	return returnVariables
}
