/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

func getPipelines(id string, name string, project string, exportPath string) ([]*types.Pipeline, error) {
	var arrResults []*types.Pipeline
	client := resty.New()

	var filters []string
	if id != "" {
		filters = append(filters, "(id eq '"+id+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if project != "" {
		filters = append(filters, "(project eq '"+project+"')")
	}
	if len(filters) > 0 {
		qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
	}
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.DocumentsList{}).
		SetAuthToken(targetConfig.AccessToken).
		SetError(&types.Exception{}).
		Get("https://" + targetConfig.Server + "/pipeline/api/pipelines")

	log.Debugln(queryResponse.Request.RawRequest.URL)
	// log.Debugln(queryResponse.String())

	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)

	}
	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.Pipeline{}
		mapstructure.Decode(value, &c)
		if exportPath != "" {
			if err := exportYaml(c.Name, c.Project, exportPath, "pipelines"); err != nil {
				log.Warnln(err)
			}
			arrResults = append(arrResults, &c)
		} else {
			arrResults = append(arrResults, &c)
		}
	}
	return arrResults, err
}

// patchPipeline - Patch Code Stream Pipeline by ID
func patchPipeline(id string, payload string) (*types.Pipeline, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&types.Pipeline{}).
		SetAuthToken(targetConfig.AccessToken).
		Patch("https://" + targetConfig.Server + "/pipeline/api/pipelines/" + id)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*types.Pipeline), err
}

func deletePipeline(id string) (*types.Pipeline, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.Pipeline{}).
		SetAuthToken(targetConfig.AccessToken).
		SetError(&types.Exception{}).
		Delete("https://" + targetConfig.Server + "/pipeline/api/pipelines/" + id)
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.Pipeline), err
}

func deletePipelineInProject(project string) ([]*types.Pipeline, error) {
	var deletedPipes []*types.Pipeline
	pipelines, err := getPipelines("", "", project, "")
	if err != nil {
		return nil, err
	}
	confirm := helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(pipelines)) + " Pipelines in " + project + ", are you sure?")
	if confirm {
		for _, pipeline := range pipelines {
			deletedPipe, err := deletePipeline(pipeline.ID)
			if err != nil {
				log.Warnln("Unable to delete "+pipeline.Name, err)
			}
			deletedPipes = append(deletedPipes, deletedPipe)
		}
		return deletedPipes, nil
	} else {
		return nil, errors.New("user declined")
	}
}
