/*
Package orchestrator Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package orchestrator

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

// GetWorkflow - returns a list of executions
func GetWorkflow(APIClient *types.APIClientOptions, id string, category string, name string) ([]*types.WsWorkflow, error) {

	var Workflows []*types.WsWorkflow
	if id != "" {
		queryResponse, err := APIClient.RESTClient.R().
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
	APIClient.RESTClient.QueryParam.Set("conditions", strings.Join(conditions, ","))
	log.Debugln("query params:", APIClient.RESTClient.QueryParam)

	queryResponse, err := APIClient.RESTClient.R().
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

				Workflow, _ := GetWorkflow(APIClient, attribute.Value, "", "")
				Workflows = append(Workflows, Workflow...)
			}

		}
	}
	return Workflows, err
}

// ExportWorkflow - exports a workflow
func ExportWorkflow(APIClient *types.APIClientOptions, id string, name string, path string) error {
	log.Debugln("ID:", id, "Name:", name, "Path:", path)
	var exportPath string
	if path != "" {
		exportPath = path
	} else {
		exportPath, _ = os.Getwd()
	}

	queryResponse, err := APIClient.RESTClient.R().
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
func ImportWorkflow(APIClient *types.APIClientOptions, path string, categoryID string) error {
	log.Debugln("Path:", path, "CategoryID:", categoryID, "Overwrite:", APIClient.Force)
	zipFileBytes, _ := ioutil.ReadFile(path)
	APIClient.RESTClient.QueryParam.Set("categoryId", categoryID)
	APIClient.RESTClient.QueryParam.Set("overwrite", strconv.FormatBool(APIClient.Force))
	queryResponse, err := APIClient.RESTClient.R().
		SetFileReader("file", "upload.zip", bytes.NewReader(zipFileBytes)).
		Post("/vco/api/workflows")
	log.Debugln(queryResponse.Request.URL)
	log.Debugln(queryResponse.StatusCode())
	if err != nil {
		return err
	}
	return nil
}

// DeleteWorkflow - deletes an Workflow by ID
func DeleteWorkflow(APIClient *types.APIClientOptions, id string, force bool) (bool, error) {
	APIClient.RESTClient.QueryParam.Set("force", strconv.FormatBool(force))
	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.Executions{}).
		SetError(&types.Exception{}).
		Delete("/vco/api/workflows" + id)

	if err != nil {
		return false, err
	}
	if queryResponse.IsError() {
		return false, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	return true, err
}

// // DeleteExecutions - deletes an execution by project, status, or pipeline name
// func DeleteExecutions(APIClient *types.APIClientOptions, confirm bool, project string, status string, name string, nested bool) ([]*types.Executions, error) {
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
// func CreateExecution(APIClient *types.APIClientOptions, id string, inputs string, comment string) (*types.CreateExecutionResponse, error) {
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
// 	queryResponse, _ := APIClient.RESTClient.R().
// 		SetBody(executionBytes).
// 		SetResult(&types.CreateExecutionResponse{}).
// 		SetError(&types.Exception{}).
// 		Post("/pipeline/api/pipelines/" + id + "/executions")

// 	if err != nil {
// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
// 	}
// 	return queryResponse.Result().(*types.CreateExecutionResponse), nil
// }
