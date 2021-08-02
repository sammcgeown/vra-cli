/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func getProject(id, name string) ([]*models.Project, error) {
	var filters []string
	var filter string
	if id != "" {
		filters = append(filters, "(id eq '"+id+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if len(filters) > 0 {
		filter = "(" + strings.Join(filters, " and ") + ")"
	}

	log.Debugln("Filter:", filter)

	apiClient := getApiClient()

	ProjectParams := project.NewGetProjectsParams()
	ProjectParams.DollarFilter = &filter

	ret, err := apiClient.Project.GetProjects(ProjectParams)
	if err != nil {
		switch err.(type) {
		case *project.GetProjectNotFound:
			return nil, errors.New("Project with ID " + id + " not found")
		}
		return nil, err
	}
	return ret.Payload.Content, nil
}

func deleteProject(id string) error {
	apiClient := getApiClient()

	// Workaround an issue where the cloud regions need to be removed before the project can be deleted.
	_, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithID(id).WithBody(&models.ProjectSpecification{
		ZoneAssignmentConfigurations: []*models.ZoneAssignmentSpecification{},
	}))
	if err != nil {
		return err
	}

	_, err = apiClient.Project.DeleteProject(project.NewDeleteProjectParams().WithID(id))
	if err != nil {
		return err
	}
	return nil
}
