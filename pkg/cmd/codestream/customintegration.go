/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
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
func CreateCustomIntegration(APIClient *types.APIClientOptions, name string, description string, yaml string) (*types.CustomIntegration, error) {
	response, err := APIClient.RESTClient.R().
		SetBody(
			types.CustomIntegration{
				Name:        name,
				Description: description,
				Yaml:        yaml,
			}).
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
func DeleteCustomIntegration(APIClient *types.APIClientOptions, id string) error {
	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.CustomIntegration{}).
		SetError(&types.Exception{}).
		Delete("/pipeline/api/custom-integrations/" + id)
	if queryResponse.IsError() {
		return errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return err
}

// // exportCustomIntegration - Export a variable to YAML
// func exportCustomIntegration(variable interface{}, exportFile string) {
// 	// variable will be a types.CustomIntegrationResponse, so lets remap to types.CustomIntegrationRequest
// 	c := types.CustomIntegrationRequest{}
// 	mapstructure.Decode(variable, &c)
// 	yaml, err := yaml.Marshal(c)
// 	if err != nil {
// 		log.Errorln("Unable to export variable ", c.Name)
// 	}
// 	if exportFile == "" {
// 		exportFile = "variables.yaml"
// 	}
// 	file, err := os.OpenFile(exportFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer file.Close()
// 	file.WriteString("---\n" + string(yaml))
// }

// // importCustomIntegrations - Import variables from the filePath
// func importCustomIntegrations(filePath string) []types.CustomIntegrationRequest {
// 	var returnCustomIntegrations []types.CustomIntegrationRequest
// 	filename, _ := filepath.Abs(filePath)
// 	yamlFile, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	reader := bytes.NewReader(yamlFile)
// 	decoder := yaml.NewDecoder(reader)
// 	var request types.CustomIntegrationRequest
// 	for decoder.Decode(&request) == nil {
// 		returnCustomIntegrations = append(returnCustomIntegrations, request)
// 	}
// 	return returnCustomIntegrations
// }
