/*
Package orchestrator Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package orchestrator

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

// GetWorkflow - returns a list of executions
func GetWorkflow(client *resty.Client, id string, category string, name string) ([]*types.WsWorkflow, error) {

	var Workflows []*types.WsWorkflow
	if id != "" {
		queryResponse, err := client.R().
			SetResult(&types.WsWorkflow{}).
			SetError(&types.Exception{}).
			Get("/vco/api/workflows/" + id)
		if err != nil {
			return nil, err
		}
		Workflows = append(Workflows, queryResponse.Result().(*types.WsWorkflow))
		return Workflows, nil
	}

	// Configure query string
	var conditions []string
	if name != "" {
		conditions = append(conditions, "name~"+name)
	}
	if category != "" {
		conditions = append(conditions, "categoryName~"+url.QueryEscape(category))
	}
	client.QueryParam.Set("conditions", strings.Join(conditions, ","))
	log.Debugln("query params:", client.QueryParam)

	queryResponse, err := client.R().
		SetResult(&types.InventoryItemsList{}).
		SetError(&types.Exception{}).
		Get("/vco/api/workflows")

	log.Debugln("Query", queryResponse.Request.URL)

	if err != nil {
		return nil, err
	}

	for _, value := range queryResponse.Result().(*types.InventoryItemsList).Link {
		for _, attribute := range value.Attributes {
			if attribute.Name == "id" {
				Workflow, _ := GetWorkflow(client, attribute.Value, "", "")
				Workflows = append(Workflows, Workflow...)
			}

		}
	}
	return Workflows, err
}

// ExportWorkflow - exports a workflow
func ExportWorkflow(client *resty.Client, id string, name string, path string) error {
	var exportPath string
	if path != "" {
		exportPath = path
	} else {
		exportPath, _ = os.Getwd()
	}

	queryResponse, err := client.R().
		SetError(&types.Exception{}).
		SetOutput(filepath.Join(exportPath, name+".zip")).
		SetHeader("Accept", "application/zip").
		Get("/vco/api/workflows/" + id)

	if err != nil {
		return err
	}

	if queryResponse.IsError() {
		return errors.New(queryResponse.Status())
	}
	return nil
}

// ImportWorkflow - imports a workflow
func ImportWorkflow(client *resty.Client, path string, categoryID string, overwrite bool) (*types.WsWorkflow, error) {
	var formData = map[string]string{}
	formData["categoryId"] = categoryID
	formData["overwrite"] = strconv.FormatBool(overwrite)

	queryResponse, err := client.R().
		SetError(&types.Exception{}).
		SetFile("file", path).
		SetFormData(formData).
		SetHeader("Accept", "application/zip").
		Post("/vco/api/workflows/")

	if err != nil {
		return nil, err
	}

	helpers.PrettyPrint(queryResponse.Result())

	return nil, nil
}

// // DeleteExecution - deletes an execution by ID
// func DeleteExecution(client *resty.Client, id string) (bool, error) {
// 	queryResponse, err := client.R().
// 		SetResult(&types.Executions{}).
// 		SetError(&types.Exception{}).
// 		Delete("/pipeline/api/executions/" + id)

// 	if err != nil {
// 		return false, errors.New(queryResponse.Error().(*types.Exception).Message)
// 	}

// 	return true, err
// }

// // DeleteExecutions - deletes an execution by project, status, or pipeline name
// func DeleteExecutions(client *resty.Client, confirm bool, project string, status string, name string, nested bool) ([]*types.Executions, error) {
// 	var deletedExecutions []*types.Executions
// 	Executions, err := GetExecution(client, "", project, status, name, nested)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !confirm {
// 		confirm = helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(Executions)) + " Executions in " + project + ", are you sure?")
// 	}
// 	if confirm {
// 		for _, Execution := range Executions {
// 			_, err := DeleteExecution(client, Execution.ID)
// 			if err != nil {
// 				log.Warnln("Unable to delete "+Execution.ID, err)
// 			}
// 			deletedExecutions = append(deletedExecutions, Execution)
// 		}
// 		return deletedExecutions, nil
// 	}
// 	return nil, errors.New("user declined")

// }

// // CreateExecution - creates an execution
// func CreateExecution(client *resty.Client, id string, inputs string, comment string) (*types.CreateExecutionResponse, error) {
// 	// Convert JSON string to byte array
// 	var inputBytes = []byte(inputs)
// 	// Unmarshal inputs using a generic interface
// 	var inputsInterface interface{}
// 	err := json.Unmarshal(inputBytes, &inputsInterface)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Create types.CreateExecutionRequest struct
// 	var execution types.CreateExecutionRequest
// 	execution.Comments = comment
// 	execution.Input = inputsInterface
// 	//Marshal struct to JSON []byte
// 	executionBytes, err := json.Marshal(execution)
// 	if err != nil {
// 		return nil, err
// 	}
// 	queryResponse, _ := client.R().
// 		SetBody(executionBytes).
// 		SetResult(&types.CreateExecutionResponse{}).
// 		SetError(&types.Exception{}).
// 		Post("/pipeline/api/pipelines/" + id + "/executions")

// 	if err != nil {
// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
// 	}
// 	return queryResponse.Result().(*types.CreateExecutionResponse), nil
// }
