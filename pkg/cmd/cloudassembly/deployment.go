/*
Package cloudassembly Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cloudassembly

import (
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetDeployments returns a list of deployments
func GetDeployments(apiclient *client.MulticloudIaaS, id string, name string, project string, status string) ([]*models.Deployment, error) {
	// Get deployment by ID
	if id != "" {
		DeploymentsParams := deployments.NewGetDeploymentByIDUsingGETParams().
			WithExpand([]string{"project"}).
			WithDeploymentID(strfmt.UUID(id))
		Deployments, err := apiclient.Deployments.GetDeploymentByIDUsingGET(DeploymentsParams)
		if err != nil {
			return nil, err

		}
		return []*models.Deployment{Deployments.Payload}, nil
	}
	// Else get deployments by name, project, or status
	DeploymentsParams := deployments.NewGetDeploymentsUsingGETParams().
		WithExpand([]string{"project"})

	if name != "" {
		DeploymentsParams.Name = &name
	}

	if project != "" {
		Project, err := GetProject(apiclient, apiVersion, project, "")
		p := Project[0]
		if err != nil {
			return nil, err
		}
		DeploymentsParams.Projects = []string{*(p.ID)}
	}

	if status != "" {
		DeploymentsParams.Status = []string{status}
	}

	log.Debug("GetDeployments: ", DeploymentsParams)

	Deployments, err := apiclient.Deployments.GetDeploymentsUsingGET(DeploymentsParams)
	if err != nil {
		return nil, err

	}
	return Deployments.Payload.Content, nil

	// var arrResults []*types.Deployment
	// client := resty.New()

	// if id != "" {
	// 	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
	// 		SetQueryParams(qParams).
	// 		SetHeader("Accept", "application/json").
	// 		SetResult(&types.Deployment{}).
	// 		SetAuthToken(targetConfig.AccessToken).
	// 		SetError(&types.Exception{}).
	// 		Get("https://" + targetConfig.Server + "/deployment/api/deployments/" + id)

	// 	log.Debugln(queryResponse.Request.RawRequest.URL)
	// 	// log.Debugln(queryResponse.String())

	// 	if queryResponse.IsError() {
	// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	// 	}

	// 	arrResults = append(arrResults, queryResponse.Result().(*types.Deployment))
	// } else {
	// 	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
	// 		SetQueryParams(qParams).
	// 		SetHeader("Accept", "application/json").
	// 		SetResult(&types.ContentsList{}).
	// 		SetAuthToken(targetConfig.AccessToken).
	// 		SetError(&types.Exception{}).
	// 		Get("https://" + targetConfig.Server + "/deployment/api/deployments")

	// 	log.Debugln(queryResponse.Request.RawRequest.URL)
	// 	// log.Debugln(queryResponse.String())

	// 	if queryResponse.IsError() {
	// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	// 	}

	// 	for _, value := range queryResponse.Result().(*types.ContentsList).Content {
	// 		c := types.Deployment{}
	// 		mapstructure.Decode(value, &c)
	// 		arrResults = append(arrResults, &c)
	// 	}

	// }
	// return arrResults, nil
}
