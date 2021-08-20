/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func exportYaml(name, project, path, object string) error {
	var exportPath string
	var qParams = make(map[string]string)
	qParams[object] = name
	qParams["project"] = project
	if path != "" {
		exportPath = path
	} else {
		exportPath, _ = os.Getwd()
	}
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/x-yaml;charset=UTF-8").
		SetAuthToken(targetConfig.AccessToken).
		SetOutput(filepath.Join(exportPath, name+".yaml")).
		SetError(&types.Exception{}).
		Get("https://" + targetConfig.Server + "/pipeline/api/export")
	log.Debugln(queryResponse.Request.RawRequest.URL)

	if queryResponse.IsError() {
		return errors.New(queryResponse.Status())
	}
	return nil
}

// importYaml import a yaml pipeline or endpoint
func importYaml(yamlPath, action, project, importType string) error {
	var pipeline types.PipelineYaml
	var endpoint types.EndpointYaml

	qParams["action"] = action
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
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Content-Type", "application/x-yaml").
		SetBody(yamlPayload).
		SetAuthToken(targetConfig.AccessToken).
		SetError(&types.Exception{}).
		Post("https://" + targetConfig.Server + "/pipeline/api/import")
	log.Debugln(queryResponse.Request.RawRequest.URL)
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

func exportCloudTemplate(name, project, content, path string) error {
	var exportPath string
	if path != "" {
		exportPath = path
		_, folderError := os.Stat(exportPath) // Get file system info
		if os.IsNotExist(folderError) {       // If it doesn't exist
			log.Debugln("Folder doesn't exist - creating")
			mkdirErr := os.MkdirAll(exportPath, os.FileMode(0755)) // Attempt to make it
			if mkdirErr != nil {
				return mkdirErr
			}
		}
	} else {
		// If path is not specified, use the current path
		exportPath, _ = os.Getwd()
	}
	exportPath = filepath.Join(exportPath, project+" - "+name+".yaml")
	f, cerr := os.Create(exportPath) // Open the file for writing
	if cerr != nil {
		return cerr
	}
	defer f.Close() // Defer closing until just before this function returns
	_, werr := f.WriteString(content)
	if werr != nil {
		return werr
	}
	return nil
}
