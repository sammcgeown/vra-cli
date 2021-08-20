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

func getEndpoint(id, name, project, endpointtype string, exportPath string) ([]*types.Endpoint, error) {
	var endpoints []*types.Endpoint
	var qParams = make(map[string]string)
	qParams["expand"] = "true"
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
	if endpointtype != "" {
		filters = append(filters, "(type eq '"+endpointtype+"')")
	}
	if len(filters) > 0 {
		qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
	}

	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.DocumentsList{}).
		SetAuthToken(targetConfig.AccessToken).
		Get("https://" + targetConfig.Server + "/pipeline/api/endpoints")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}

	for _, value := range queryResponse.Result().(*types.DocumentsList).Documents {
		c := types.Endpoint{}
		mapstructure.Decode(value, &c)
		if exportPath != "" {
			exportYaml(c.Name, c.Project, exportPath, "endpoints")
			endpoints = append(endpoints, &c)
		} else {
			endpoints = append(endpoints, &c)
		}

	}
	return endpoints, err
}

func deleteEndpoint(id string) (*types.Endpoint, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&types.Endpoint{}).
		SetAuthToken(targetConfig.AccessToken).
		Delete("https://" + targetConfig.Server + "/pipeline/api/endpoints/" + id)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*types.Endpoint), err
}

func deleteEndpointByProject(project string) ([]*types.Endpoint, error) {
	var deletedEndpoints []*types.Endpoint
	Endpoints, err := getEndpoint("", "", project, "", "")
	if err != nil {
		return nil, err
	}
	confirm := helpers.AskForConfirmation("This will attempt to delete " + fmt.Sprint(len(Endpoints)) + " Endpoints in " + project + ", are you sure?")
	if confirm {

		for _, endpoint := range Endpoints {
			deletedEndpoint, err := deleteEndpoint(endpoint.ID)
			if err != nil {
				log.Warnln("Unable to delete "+endpoint.Name, err)
			}
			deletedEndpoints = append(deletedEndpoints, deletedEndpoint)
		}
		return deletedEndpoints, nil
	} else {
		return nil, errors.New("user declined")
	}
}
