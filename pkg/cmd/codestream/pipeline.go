/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

// GetPipeline - Get Code Stream Pipeline
func GetPipeline(client *resty.Client, id string, name string, project string, exportPath string) ([]*types.Pipeline, error) {
	var arrResults []*types.Pipeline

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
		client.QueryParam.Add("$filter", "("+strings.Join(filters, ") and (")+")")
		log.Debugln(client.QueryParam)
	}

	queryResponse, err := client.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/pipelines")

	log.Debugln(queryResponse.Request.RawRequest.URL)
	// log.Debugln(queryResponse.String())

	if err != nil {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.Pipeline{}
		mapstructure.Decode(value, &c)
		if exportPath != "" {
			// if err := cmd.exportYaml(c.Name, c.Project, exportPath, "pipelines"); err != nil {
			// 	log.Warnln(err)
			// }
			arrResults = append(arrResults, &c)
		} else {
			arrResults = append(arrResults, &c)
		}
	}
	return arrResults, err
}

// PatchPipeline - Patch Code Stream Pipeline by ID
func PatchPipeline(client *resty.Client, id string, payload string) (*types.Pipeline, error) {
	queryResponse, err := client.R().
		SetBody(payload).
		SetResult(&types.Pipeline{}).
		SetError(&types.Exception{}).
		Patch("/pipeline/api/pipelines/" + id)
	if err != nil {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.Pipeline), err
}

// DeletePipeline - Delete Code Stream Pipeline by ID
func DeletePipeline(client *resty.Client, id string) (*types.Pipeline, error) {
	queryResponse, err := client.R().
		SetResult(&types.Pipeline{}).
		SetError(&types.Exception{}).
		Delete("/pipeline/api/pipelines/" + id)
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return queryResponse.Result().(*types.Pipeline), err
}

// DeletePipelineInProject - Delete Code Stream Pipeline by Project
func DeletePipelineInProject(client *resty.Client, project string) ([]*types.Pipeline, error) {
	var deletedPipes []*types.Pipeline
	pipelines, err := GetPipeline(client, "", "", project, "")
	if err != nil {
		return nil, err
	}
	confirm := helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(pipelines)) + " Pipelines in " + project + ", are you sure?")
	if confirm {
		for _, pipeline := range pipelines {
			deletedPipe, err := DeletePipeline(client, pipeline.ID)
			if err != nil {
				log.Warnln("Unable to delete "+pipeline.Name, err)
			}
			deletedPipes = append(deletedPipes, deletedPipe)
		}
		return deletedPipes, nil
	}

	return nil, errors.New("user declined")
}
