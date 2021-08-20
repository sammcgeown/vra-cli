/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

func getExecutions(id string, project string, status string, name string, nested bool) ([]*types.Executions, error) {
	var arrExecutions []*types.Executions
	if id != "" {
		x, err := getExecution("/codestream/api/executions/" + id)
		if err != nil {
			return nil, err
		}
		arrExecutions = append(arrExecutions, x)
		return arrExecutions, err
	}
	client := resty.New()
	var qParams = make(map[string]string)

	qParams["$orderby"] = "_requestTimeInMicros desc"
	qParams["$top"] = fmt.Sprint(count)
	qParams["$skip"] = fmt.Sprint(skip)

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
	if project != "" {
		filters = append(filters, "(project eq '"+project+"')")
	}
	if len(filters) > 0 {
		qParams["$filter"] = "(" + strings.Join(filters, ") and (") + ")"
		log.Debugln(qParams["$filter"])
	}

	log.Debug(qParams)

	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		SetAuthToken(targetConfig.AccessToken).
		Get("https://" + targetConfig.Server + "/pipeline/api/executions")
	if queryResponse.IsError() {
		//return nil, queryResponse.Error().(error)
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.Executions{}
		mapstructure.Decode(value, &c)
		arrExecutions = append(arrExecutions, &c)
	}
	return arrExecutions, err
}

func getExecution(executionLink string) (*types.Executions, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.Executions{}).
		SetAuthToken(targetConfig.AccessToken).
		Get("https://" + targetConfig.Server + executionLink)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*types.Executions), err
}

func deleteExecution(id string) (*types.Executions, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.Executions{}).
		SetAuthToken(targetConfig.AccessToken).
		Delete("https://" + targetConfig.Server + "/pipeline/api/executions/" + id)
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.Executions), err
}

func deleteExecutions(project string, status string, name string, nested bool) ([]*types.Executions, error) {
	var deletedExecutions []*types.Executions
	Executions, err := getExecutions("", project, status, name, nested)
	if err != nil {
		return nil, err
	}
	confirm := helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(Executions)) + " Executions in " + project + ", are you sure?")
	if confirm {
		for _, Execution := range Executions {
			deletedExecution, err := deleteExecution(Execution.ID)
			if err != nil {
				log.Warnln("Unable to delete "+Execution.ID, err)
			}
			deletedExecutions = append(deletedExecutions, deletedExecution)
		}
		return deletedExecutions, nil
	} else {
		return nil, errors.New("user declined")
	}
}

func createExecution(id string, inputs string, comment string) (*types.CreateExecutionResponse, error) {
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
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Content-Type", "application/json").
		SetBody(executionBytes).
		SetResult(&types.CreateExecutionResponse{}).
		SetAuthToken(targetConfig.AccessToken).
		Post("https://" + targetConfig.Server + "/pipeline/api/pipelines/" + id + "/executions")
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.CreateExecutionResponse), nil
}
