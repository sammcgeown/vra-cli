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
}

// DeleteDeployment - Delete a deployment
func DeleteDeployment(apiclient *client.MulticloudIaaS, id string) error {
	DeleteParams := deployments.NewDeleteDeploymentUsingDELETEParams().WithDeploymentID(strfmt.UUID(id))
	_, err := apiclient.Deployments.DeleteDeploymentUsingDELETE(DeleteParams)
	if err != nil {
		return err
	}
	return nil
}
