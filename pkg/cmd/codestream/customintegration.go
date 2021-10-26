/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
)

// GetCustomIntegration returns a custom integration
func GetCustomIntegration(APIClient *types.APIClientOptions, id, name string) ([]*types.CustomIntegration, error) {
	var arrCustomIntegrations []*types.CustomIntegration

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
		return nil, queryResponse.Error().(error)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.CustomIntegration{}
		mapstructure.Decode(value, &c)
		arrCustomIntegrations = append(arrCustomIntegrations, &c)
	}
	return arrCustomIntegrations, err
}

// // createCustomIntegration - Create a new Code Stream CustomIntegration
// func createCustomIntegration(name string, description string, variableType string, project string, value string) (*types.CustomIntegrationResponse, error) {
// 	client := resty.New()
// 	response, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetBody(
// 			types.CustomIntegrationRequest{
// 				Project:     project,
// 				Kind:        "VARIABLE",
// 				Name:        name,
// 				Description: description,
// 				Type:        variableType,
// 				Value:       value,
// 			}).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&types.CustomIntegrationResponse{}).
// 		SetError(&types.Exception{}).
// 		SetAuthToken(targetConfig.AccessToken).
// 		Post("https://" + targetConfig.Server + "/pipeline/api/variables")
// 	if response.IsError() {
// 		return nil, errors.New(response.Error().(*types.Exception).Message)
// 	}
// 	return response.Result().(*types.CustomIntegrationResponse), err
// }

// // updateCustomIntegration - Create a new Code Stream CustomIntegration
// func updateCustomIntegration(id string, name string, description string, typename string, value string) (*types.CustomIntegrationResponse, error) {
// 	variable, _ := getCustomIntegrationByID(id)
// 	if name != "" {
// 		variable.Name = name
// 	}
// 	if description != "" {
// 		variable.Description = description
// 	}
// 	if typename != "" {
// 		variable.Type = typename
// 	}
// 	if value != "" {
// 		variable.Value = value
// 	}
// 	client := resty.New()
// 	response, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetBody(variable).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&types.CustomIntegrationResponse{}).
// 		SetError(&types.Exception{}).
// 		SetAuthToken(targetConfig.AccessToken).
// 		Put("https://" + targetConfig.Server + "/pipeline/api/variables/" + id)
// 	if response.IsError() {
// 		return nil, errors.New(response.Error().(*types.Exception).Message)
// 	}
// 	return response.Result().(*types.CustomIntegrationResponse), err
// }

// // deleteCustomIntegration - Delete a Code Stream CustomIntegration
// func deleteCustomIntegration(id string) (*types.CustomIntegrationResponse, error) {
// 	client := resty.New()
// 	response, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetHeader("Accept", "application/json").
// 		SetResult(&types.CustomIntegrationResponse{}).
// 		SetAuthToken(targetConfig.AccessToken).
// 		Delete("https://" + targetConfig.Server + "/pipeline/api/variables/" + id)
// 	if response.IsError() {
// 		log.Errorln("Create CustomIntegration failed", err)
// 		os.Exit(1)
// 	}
// 	return response.Result().(*types.CustomIntegrationResponse), err
// }

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
