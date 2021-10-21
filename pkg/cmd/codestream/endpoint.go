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

// GetEndpoint returns an endpoint
func GetEndpoint(client *resty.Client, id, name, project, endpointtype string, exportPath string) ([]*types.Endpoint, error) {
	var endpoints []*types.Endpoint

	client.QueryParam.Set("expand", "true")
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
	if endpointtype != "" {
		filters = append(filters, "(type eq '"+endpointtype+"')")
	}
	if len(filters) > 0 {
		client.QueryParam.Set("$filter", "("+strings.Join(filters, " and ")+")")
	}

	queryResponse, err := client.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/endpoints")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.Endpoint{}
		mapstructure.Decode(value, &c)
		if exportPath != "" {
			//helpers.ExportYaml(c.Name, c.Project, exportPath, "endpoints")
			endpoints = append(endpoints, &c)
		} else {
			endpoints = append(endpoints, &c)
		}

	}
	return endpoints, err
}

// DeleteEndpoint deletes an endpoint
func DeleteEndpoint(client *resty.Client, id string) error {
	queryResponse, _ := client.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/endpoints/" + id)

	if queryResponse.IsError() {
		return queryResponse.Error().(error)
	}
	return nil
}

// DeleteEndpointByProject deletes an endpoint by project
func DeleteEndpointByProject(client *resty.Client, project string) ([]*types.Endpoint, error) {
	var deletedEndpoints []*types.Endpoint
	Endpoints, err := GetEndpoint(client, "", "", project, "", "")
	if err != nil {
		return nil, err
	}
	confirm := helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(Endpoints)) + " Endpoints in " + project + ", are you sure?")
	if confirm {

		for _, endpoint := range Endpoints {
			err := DeleteEndpoint(client, endpoint.ID)
			if err != nil {
				log.Warnln("Unable to delete "+endpoint.Name, err)
			}
			deletedEndpoints = append(deletedEndpoints, endpoint)
		}
		return deletedEndpoints, nil
	}
	return nil, errors.New("user declined")
}