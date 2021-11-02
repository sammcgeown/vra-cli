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

// GetAction - returns a list of executions
func GetAction(APIClient *types.APIClientOptions, id string, category string, name string) ([]*types.WsAction, error) {

	var Actions []*types.WsAction
	if id != "" {
		queryResponse, err := APIClient.RESTClient.R().
			SetResult(&types.WsAction{}).
			SetError(&types.Exception{}).
			Get("/vco/api/actions/" + id)
		if err != nil {
			return nil, err
		}
		Actions = append(Actions, queryResponse.Result().(*types.WsAction))
		return Actions, nil
	}

	// If no ID is specified, use the category and name to find the action
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
		Get("/vco/api/catalog/System/Action")

	if err != nil {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	APIClient.RESTClient.QueryParam.Del("conditions")

	for _, value := range queryResponse.Result().(*types.InventoryItemsList).Link {
		for _, attribute := range value.Attributes {
			if attribute.Name == "id" {
				Action, _ := GetAction(APIClient, attribute.Value, "", "")
				Actions = append(Actions, Action...)
			}

		}
	}

	// Configure query string
	// var conditions []string
	// if name != "" {
	// 	conditions = append(conditions, "name~"+name)
	// }
	// if category != "" {
	// 	conditions = append(conditions, "categoryName~"+url.QueryEscape(category))
	// }
	// APIClient.RESTClient.QueryParam.Set("conditions", strings.Join(conditions, ","))
	// log.Debugln("query params:", APIClient.RESTClient.QueryParam)

	// queryResponse, err := APIClient.RESTClient.R().
	// 	SetResult(&types.InventoryItemsList{}).
	// 	SetError(&types.Exception{}).
	// 	Get("/vco/api/actions")

	// log.Debugln("Query", queryResponse.Request.URL)

	// if err != nil {
	// 	return nil, err
	// }

	// for _, value := range queryResponse.Result().(*types.InventoryItemsList).Link {
	// 	for _, attribute := range value.Attributes {
	// 		if attribute.Name == "id" {

	// 			Action, _ := GetAction(APIClient, attribute.Value, "", "")
	// 			Actions = append(Actions, Action...)
	// 		}

	// 	}
	// }
	return Actions, err
}

// ExportAction - exports a workflow
func ExportAction(APIClient *types.APIClientOptions, id string, name string, path string) error {
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

// ImportAction - imports a workflow
func ImportAction(APIClient *types.APIClientOptions, path string, categoryID string) error {
	log.Debugln("Path:", path, "CategoryID:", categoryID, "Overwrite:", APIClient.Force)
	zipFileBytes, _ := ioutil.ReadFile(path)
	APIClient.RESTClient.QueryParam.Set("categoryId", categoryID)
	APIClient.RESTClient.QueryParam.Set("overwrite", strconv.FormatBool(APIClient.Force))
	queryResponse, err := APIClient.RESTClient.R().
		SetError(&types.Exception{}).
		SetFileReader("file", "upload.zip", bytes.NewReader(zipFileBytes)).
		Post("/vco/api/workflows")
	APIClient.RESTClient.QueryParam.Del("categoryId")
	APIClient.RESTClient.QueryParam.Del("overwrite")
	if err != nil {
		return errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	if queryResponse.IsError() {
		return errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return nil
}

// DeleteAction - deletes an Action by ID
func DeleteAction(APIClient *types.APIClientOptions, id string) (bool, error) {
	APIClient.RESTClient.QueryParam.Set("force", strconv.FormatBool(APIClient.Force))
	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.Executions{}).
		SetError(&types.Exception{}).
		Delete("/vco/api/workflows/" + id)

	if queryResponse.IsError() {
		return false, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	if err != nil {
		return false, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	return true, err
}
