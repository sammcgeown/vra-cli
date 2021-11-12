/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

// GetExecution - returns a list of executions
func GetExecution(APIClient *types.APIClientOptions, id string, project string, status string, name string, nested bool, rollback bool) ([]*types.Executions, error) {
	var arrExecutions []*types.Executions
	if id != "" {
		queryResponse, err := APIClient.RESTClient.R().
			SetResult(&types.Executions{}).
			SetError(&types.Exception{}).
			Get("/pipeline/api/executions/" + id)
		if err != nil {
			return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
		}
		arrExecutions = append(arrExecutions, queryResponse.Result().(*types.Executions))
		return arrExecutions, nil
	}

	APIClient.RESTClient.QueryParam.Set("$top", strconv.Itoa(APIClient.Pagination.PageSize))
	APIClient.RESTClient.QueryParam.Set("$page", strconv.Itoa(APIClient.Pagination.Page))
	APIClient.RESTClient.QueryParam.Set("$skip", strconv.Itoa(APIClient.Pagination.Skip))

	var filters []string
	if status != "" {
		filters = append(filters, "(status eq '"+strings.ToUpper(status)+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if nested {
		filters = append(filters, "(_nested eq '"+strconv.FormatBool(nested)+"')")
	}
	if rollback {
		filters = append(filters, "(_rollback eq '"+strconv.FormatBool(rollback)+"')")
	}
	if project != "" {
		filters = append(filters, "(project eq '"+project+"')")
	}
	if len(filters) > 0 {
		APIClient.RESTClient.QueryParam.Set("$filter", "("+strings.Join(filters, ") and (")+")")
		log.Debugln(APIClient.RESTClient.QueryParam)
	}

	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/executions")

	if err != nil {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.Executions{}
		mapstructure.Decode(value, &c)
		arrExecutions = append(arrExecutions, &c)
	}
	return arrExecutions, err
}

// DeleteExecution - deletes an execution by ID
func DeleteExecution(APIClient *types.APIClientOptions, id string) (bool, error) {
	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.Executions{}).
		SetError(&types.Exception{}).
		Delete("/pipeline/api/executions/" + id)

	if err != nil {
		return false, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	return true, err
}

// DeleteExecutions - deletes an execution by project, status, or pipeline name
func DeleteExecutions(APIClient *types.APIClientOptions, project string, status string, name string, nested bool, rollback bool) ([]*types.Executions, error) {
	var deletedExecutions []*types.Executions
	Executions, err := GetExecution(APIClient, "", project, status, name, nested, rollback)
	if err != nil {
		return nil, err
	}
	if !APIClient.Confirm {
		APIClient.Confirm = helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(Executions)) + ", are you sure?")
	}
	if APIClient.Confirm {
		for _, Execution := range Executions {
			_, err := DeleteExecution(APIClient, Execution.ID)
			if err != nil {
				log.Warnln("Unable to delete "+Execution.ID, err)
			}
			deletedExecutions = append(deletedExecutions, Execution)
		}
		return deletedExecutions, nil
	}
	return nil, errors.New("user declined")

}

// CreateExecution - creates an execution
func CreateExecution(APIClient *types.APIClientOptions, id string, inputs string, comment string) (*types.CreateExecutionResponse, error) {
	// Convert JSON string to byte array
	var inputBytes = []byte(inputs)
	// Unmarshal inputs using a generic interface
	var inputsInterface interface{}
	err := json.Unmarshal(inputBytes, &inputsInterface)
	if err != nil {
		return nil, err
	}
	// Create types.CreateExecutionRequest struct
	var execution types.CreateExecutionRequest
	execution.Comments = comment
	execution.Input = inputsInterface
	//Marshal struct to JSON []byte
	executionBytes, err := json.Marshal(execution)
	if err != nil {
		return nil, err
	}
	queryResponse, _ := APIClient.RESTClient.R().
		SetBody(executionBytes).
		SetResult(&types.CreateExecutionResponse{}).
		SetError(&types.Exception{}).
		Post("/pipeline/api/pipelines/" + id + "/executions")

	if err != nil {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.CreateExecutionResponse), nil
}
