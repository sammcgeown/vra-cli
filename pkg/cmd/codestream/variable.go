/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// GetVariable - Get a Code Stream Variable
func GetVariable(client *resty.Client, id, name, project, exportPath string) ([]*types.VariableResponse, error) {
	var arrVariables []*types.VariableResponse

	// Get by ID
	if id != "" {
		queryResponse, err := client.R().
			SetResult(&types.VariableResponse{}).
			Get("/pipeline/api/variables/" + id)

		if queryResponse.IsError() {
			log.Errorln("GET Variable failed", err)
		}
		arrVariables = append(arrVariables, queryResponse.Result().(*types.VariableResponse))
		return arrVariables, err
	}

	var filters []string
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if project != "" {
		filters = append(filters, "(project eq '"+project+"')")
	}
	if len(filters) > 0 {
		client.QueryParam.Set("$filter", "("+strings.Join(filters, ") and (")+")")
		log.Debugln(client.QueryParam)
	}

	queryResponse, err := client.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/variables")

	if err != nil {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	log.Debugln(queryResponse.Request.URL)
	log.Debugln(queryResponse.RawResponse)

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.VariableResponse{}
		mapstructure.Decode(value, &c)
		arrVariables = append(arrVariables, &c)
		if exportPath != "" {
			ExportVariable(c, exportPath)
		}
	}
	return arrVariables, err
}

// CreateVariable - Create a new Code Stream Variable
func CreateVariable(client *resty.Client, name string, description string, variableType string, project string, value string) (*types.VariableResponse, error) {
	queryResponse, err := client.R().
		SetBody(
			types.VariableRequest{
				Project:     project,
				Kind:        "VARIABLE",
				Name:        name,
				Description: description,
				Type:        variableType,
				Value:       value,
			}).
		SetResult(&types.VariableResponse{}).
		SetError(&types.Exception{}).
		Post("/pipeline/api/variables")

	log.Debugln(queryResponse.RawResponse)
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.VariableResponse), err
}

// UpdateVariable - Update an existing Code Stream Variable
func UpdateVariable(client *resty.Client, id string, name string, description string, typename string, value string) (*types.VariableResponse, error) {
	variables, _ := GetVariable(client, id, "", "", "")
	variable := variables[0]
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

	queryResponse, err := client.R().
		SetBody(variable).
		SetResult(&types.VariableResponse{}).
		SetError(&types.Exception{}).
		Put("/pipeline/api/variables/" + id)

	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.VariableResponse), err
}

// DeleteVariable - Delete a Code Stream Variable
func DeleteVariable(client *resty.Client, id string) (bool, error) {
	queryResponse, err := client.R().
		SetResult(&types.VariableResponse{}).
		Delete("/pipeline/api/variables/" + id)
	if queryResponse.IsError() {
		return false, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return true, err
}

// DeleteVariableByProject - Delete all Variables in a Project
func DeleteVariableByProject(client *resty.Client, confirm bool, project string) ([]*types.VariableResponse, error) {
	var deletedVariables []*types.VariableResponse
	Variables, err := GetVariable(client, "", "", project, "")
	if err != nil {
		return nil, err
	}
	if len(Variables) == 0 {
		log.Infoln("No variables found for project:", project)
		return deletedVariables, nil
	}
	if !confirm {
		confirm = helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(Variables)) + " variables in " + project + ", are you sure?")
	}
	if confirm {
		for _, Variable := range Variables {
			_, err := DeleteVariable(client, Variable.ID)
			if err != nil {
				log.Warnln("Unable to delete "+Variable.Name, err)
			}
			deletedVariables = append(deletedVariables, Variable)
		}
		return deletedVariables, nil
	}
	return nil, errors.New("user declined")
}

// ExportVariable - Export a variable to YAML
func ExportVariable(variable interface{}, exportPath string) {
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

// ImportVariables - Import variables from the filePath
func ImportVariables(filePath string) []types.VariableRequest {
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
