/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sammcgeown/vra-cli/pkg/util/types"
	"gopkg.in/yaml.v2"
)

// ExportYaml exports the Pipeline or Endpoint to a YAML file
func ExportYaml(APIClient *types.APIClientOptions, id, name, project, path, object string) error {
	var exportPath string
	if path != "" {
		exportPath = path
	} else {
		exportPath, _ = os.Getwd()
	}
	APIClient.RESTClient.QueryParam.Set(object, name)
	APIClient.RESTClient.QueryParam.Set("project", project)

	queryResponse, err := APIClient.RESTClient.R().
		SetError(&types.Exception{}).
		SetOutput(filepath.Join(exportPath, name+".yaml")).
		SetHeader("Accept", "application/x-yaml;charset=UTF-8").
		Get("/pipeline/api/export")

	if err != nil {
		return err
	}

	if queryResponse.IsError() {
		return errors.New(queryResponse.Status())
	}
	return nil
}

// ImportYaml import a yaml pipeline or endpoint
func ImportYaml(APIClient *types.APIClientOptions, yamlPath, action, project, importType string) error {
	var pipeline types.PipelineYaml
	var endpoint types.EndpointYaml

	APIClient.RESTClient.QueryParam.Set("action", action)

	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return err
	}

	if project != "" { // If the project flag is set we need to update the project value
		if importType == "pipeline" {
			yamlErr := yaml.Unmarshal(yamlBytes, &pipeline)
			if yamlErr != nil {
				return yamlErr
			}
			pipeline.Project = project
			yamlBytes, _ = yaml.Marshal(pipeline)
		} else {
			yamlErr := yaml.Unmarshal(yamlBytes, &endpoint)
			if yamlErr != nil {
				return yamlErr
			}
			endpoint.Project = project
			yamlBytes, _ = yaml.Marshal(endpoint)
		}
	}
	yamlPayload := string(yamlBytes)

	queryResponse, err := APIClient.RESTClient.R().
		SetError(&types.Exception{}).
		SetBody(yamlPayload).
		SetHeader("Content-Type", "application/x-yaml").
		Post("/pipeline/api/import")

	if err != nil {
		return err
	}

	if queryResponse.IsError() {
		return queryResponse.Error().(error)
	}

	var importResponse types.PipelineImportResponse
	if err = yaml.Unmarshal(queryResponse.Body(), &importResponse); err != nil {
		return err
	}

	if importResponse.Status != "CREATED" && action == "create" {
		return errors.New(importResponse.Status + " - " + importResponse.StatusMessage)
	}
	if importResponse.Status != "UPDATED" && action == "apply" {
		return errors.New(importResponse.Status + " - " + importResponse.StatusMessage)
	}
	return nil
}
