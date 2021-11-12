/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

// GetCustomIntegration returns a custom integration
func GetCustomIntegration(APIClient *types.APIClientOptions, id, name string) ([]*types.CustomIntegration, error) {
	var arrCustomIntegrations []*types.CustomIntegration

	APIClient.RESTClient.QueryParam.Set("$top", strconv.Itoa(APIClient.Pagination.PageSize))
	APIClient.RESTClient.QueryParam.Set("$page", strconv.Itoa(APIClient.Pagination.Page))
	APIClient.RESTClient.QueryParam.Set("$skip", strconv.Itoa(APIClient.Pagination.Skip))

	var filters []string
	if id != "" {
		filters = append(filters, "(id eq '"+id+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if len(filters) > 0 {
		APIClient.RESTClient.QueryParam.Set("$filter", "("+strings.Join(filters, " and ")+")")
	}
	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/custom-integrations")

	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.CustomIntegration{}
		mapstructure.Decode(value, &c)
		arrCustomIntegrations = append(arrCustomIntegrations, &c)
	}
	APIClient.RESTClient.QueryParam.Del("$top")
	APIClient.RESTClient.QueryParam.Del("$page")
	APIClient.RESTClient.QueryParam.Del("$skip")
	APIClient.RESTClient.QueryParam.Del("$filter")
	return arrCustomIntegrations, err
}

// GetCustomIntegrationVersions returns all versions of a custom integration
func GetCustomIntegrationVersions(APIClient *types.APIClientOptions, id string) ([]string, error) {
	var versions []string

	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/custom-integrations/" + id + "/versions")

	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.CustomIntegration{}
		mapstructure.Decode(value, &c)
		versions = append(versions, c.Version)
	}
	return versions, err
}

// CreateCustomIntegration - Create a new Code Stream CustomIntegration
func CreateCustomIntegration(APIClient *types.APIClientOptions, name string, description string, yaml string, importPath string) (*types.CustomIntegration, error) {
	var customIntegration *types.CustomIntegration
	if importPath != "" {
		var importErr error
		customIntegration, importErr = ImportCustomIntegration(importPath)
		if importErr != nil {
			return nil, importErr
		}
	} else {
		customIntegration.Name = name
		customIntegration.Description = description
		customIntegration.Yaml = yaml
	}
	response, err := APIClient.RESTClient.R().
		SetBody(customIntegration).
		SetResult(&types.CustomIntegration{}).
		SetError(&types.Exception{}).
		Post("/pipeline/api/custom-integrations")
	if response.IsError() {
		return nil, errors.New(response.Error().(*types.Exception).Message)
	}
	return response.Result().(*types.CustomIntegration), err
}

// UpdateCustomIntegration - Create a new Code Stream CustomIntegration
func UpdateCustomIntegration(APIClient *types.APIClientOptions, id string, description string, yaml string, version string, state string) (*types.CustomIntegration, error) {
	var updatedCustomIntegration *types.CustomIntegration
	CustomIntegration, _ := GetCustomIntegration(APIClient, id, "")
	if len(CustomIntegration) == 0 {
		return nil, errors.New("Custom Integration not found")
	} else if len(CustomIntegration) > 1 {
		return nil, errors.New("Multiple Custom Integrations found")
	}

	if description != "" || yaml != "" {
		if description != "" {
			CustomIntegration[0].Description = description
		}
		if yaml != "" {
			CustomIntegration[0].Yaml = yaml
		}

		queryResponse, _ := APIClient.RESTClient.R().
			SetBody(CustomIntegration[0]).
			SetResult(&types.CustomIntegration{}).
			SetError(&types.Exception{}).
			Put("/pipeline/api/custom-integrations/" + CustomIntegration[0].ID)

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
		}
		updatedCustomIntegration = queryResponse.Result().(*types.CustomIntegration)
	}
	if version != "" {
		currentVersions, err := GetCustomIntegrationVersions(APIClient, CustomIntegration[0].ID)
		if err != nil {
			return nil, err
		}
		if !helpers.StringArrayContains(currentVersions, version) {
			// Version doesn't already exist
			queryResponse, _ := APIClient.RESTClient.R().
				SetBody(`{"changeLog":"Updated by vra-cli", "description":"Updated by vRealize Automation CLI", "version":"` + version + `"}`).
				SetResult(&types.CustomIntegration{}).
				SetError(&types.Exception{}).
				Post("/pipeline/api/custom-integrations/" + CustomIntegration[0].ID + "/versions")

			if queryResponse.IsError() {
				return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
			}
			updatedCustomIntegration = queryResponse.Result().(*types.CustomIntegration)
		}
		if state == "delete" {
			// Delete the version
			queryResponse, _ := APIClient.RESTClient.R().
				SetResult(&types.CustomIntegration{}).
				SetError(&types.Exception{}).
				Delete("/pipeline/api/custom-integrations/" + CustomIntegration[0].ID + "/versions/" + version)

			if queryResponse.IsError() {
				return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
			}
			updatedCustomIntegration = queryResponse.Result().(*types.CustomIntegration)
		} else if state != "" {
			// Update the version's state
			queryResponse, _ := APIClient.RESTClient.R().
				SetResult(&types.CustomIntegration{}).
				SetError(&types.Exception{}).
				Post("/pipeline/api/custom-integrations/" + CustomIntegration[0].ID + "/versions/" + version + "/" + state)

			if queryResponse.IsError() {
				return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
			}
			updatedCustomIntegration = queryResponse.Result().(*types.CustomIntegration)
		}

	}
	return updatedCustomIntegration, nil
}

// DeleteCustomIntegration - Delete a Code Stream CustomIntegration
func DeleteCustomIntegration(APIClient *types.APIClientOptions, id, name string) error {
	if name != "" && id == "" {
		customIntegration, err := GetCustomIntegration(APIClient, "", name)
		if err != nil {
			return err
		}
		if len(customIntegration) == 0 {
			return errors.New("Custom Integration not found")
		}
		if len(customIntegration) > 1 {
			return errors.New("Multiple Custom Integrations found")
		}
		id = customIntegration[0].ID
	}
	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.CustomIntegration{}).
		SetError(&types.Exception{}).
		Delete("/pipeline/api/custom-integrations/" + id)
	if queryResponse.IsError() {
		return errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return err
}

// ExportCustomIntegration - Export a custom integration
func ExportCustomIntegration(customintegration types.CustomIntegration, exportPath string, overwrite bool) error {
	_, pathError := os.Stat(exportPath) // Check if folder exists
	if pathError != nil {
		if os.IsNotExist(pathError) { // If it doesn't exist
			os.MkdirAll(exportPath, 0755) // Create the folder
		} else {
			return pathError // Return the error
		}
	}

	// Create the absolute path, with file name
	filePath, _ := filepath.Abs(filepath.Join(exportPath, customintegration.Name+".json"))
	fileStat, fileErr := os.Stat(filePath)
	if fileErr != nil {
		if os.IsNotExist(fileErr) {
			// Create the file
		} else {
			return fileErr
		}
	} else if fileStat.Mode().IsRegular() && !overwrite {
		return errors.New(filePath + " exists, use --force to overwrite")
	}

	ci, _ := json.MarshalIndent(customintegration, "", "  ")
	writeErr := os.WriteFile(filePath, ci, 0644)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

// ImportCustomIntegration - Import Custom Integrations from the importPath
func ImportCustomIntegration(importPath string) (*types.CustomIntegration, error) {
	filename, _ := filepath.Abs(importPath)
	log.Debugln("Importing Custom Integration from", filename)
	customIntegration, readErr := os.ReadFile(filename)
	if readErr != nil {
		return nil, readErr
	}
	ci := &types.CustomIntegration{}
	json.Unmarshal(customIntegration, ci)
	return ci, nil
}
