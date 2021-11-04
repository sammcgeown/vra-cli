/*
Package codestream Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package codestream

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

// GetEndpoint returns an endpoint
func GetEndpoint(APIClient *types.APIClientOptions, id, name, project, endpointtype string, exportPath string) ([]*types.Endpoint, error) {
	var endpoints []*types.Endpoint

	APIClient.RESTClient.QueryParam.Set("expand", "true")
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
		APIClient.RESTClient.QueryParam.Set("$filter", "("+strings.Join(filters, " and ")+")")
	}

	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/endpoints")

	APIClient.RESTClient.QueryParam.Del("$filter")
	APIClient.RESTClient.QueryParam.Del("expand")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.Endpoint{}
		mapstructure.Decode(value, &c)
		endpoints = append(endpoints, &c)
	}
	return endpoints, err
}

// DeleteEndpoint deletes an endpoint
func DeleteEndpoint(APIClient *types.APIClientOptions, id string) error {
	queryResponse, _ := APIClient.RESTClient.R().
		SetResult(&types.DocumentsList{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/endpoints/" + id)

	if queryResponse.IsError() {
		return queryResponse.Error().(error)
	}
	return nil
}

// DeleteEndpointByProject deletes an endpoint by project
func DeleteEndpointByProject(APIClient *types.APIClientOptions, project string) ([]*types.Endpoint, error) {
	var deletedEndpoints []*types.Endpoint
	Endpoints, err := GetEndpoint(APIClient, "", "", project, "", "")
	if err != nil {
		return nil, err
	}
	confirm := helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(Endpoints)) + " Endpoints in " + project + ", are you sure?")
	if confirm {

		for _, endpoint := range Endpoints {
			err := DeleteEndpoint(APIClient, endpoint.ID)
			if err != nil {
				log.Warnln("Unable to delete "+endpoint.Name, err)
			}
			deletedEndpoints = append(deletedEndpoints, endpoint)
		}
		return deletedEndpoints, nil
	}
	return nil, errors.New("user declined")
}
