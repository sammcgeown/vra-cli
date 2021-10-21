/*
Package cloudassembly Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cloudassembly

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetProject - Get Projects
func GetProject(apiClient *client.MulticloudIaaS, apiVersion string, name string, id string) ([]*models.IaaSProject, error) {
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

	ProjectParams := project.NewGetProjectsParams()
	ProjectParams.DollarFilter = &filter
	ProjectParams.APIVersion = &apiVersion

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

// DeleteProject - Delete Project
func DeleteProject(apiClient *client.MulticloudIaaS, apiVersion string, id string) error {

	// Workaround an issue where the cloud regions need to be removed before the project can be deleted.
	_, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithAPIVersion(&apiVersion).WithID(id).WithBody(&models.IaaSProjectSpecification{
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

// CreateProject - Create Project
func CreateProject(apiClient *client.MulticloudIaaS, apiVersion string, name string, description string, administrators []*models.User, members []*models.User, viewers []*models.User, zoneAssignment []*models.ZoneAssignmentSpecification, constraints map[string][]models.Constraint, operationTimeout int64, machineNamingTemplate string, sharedResources *bool) (*models.IaaSProject, error) {
	createdProject, err := apiClient.Project.CreateProject(project.NewCreateProjectParams().WithAPIVersion(&apiVersion).WithBody(&models.IaaSProjectSpecification{
		Administrators:               administrators,
		Constraints:                  constraints,
		Description:                  description,
		MachineNamingTemplate:        machineNamingTemplate,
		Members:                      members,
		Name:                         &name,
		OperationTimeout:             &operationTimeout,
		SharedResources:              *sharedResources,
		Viewers:                      viewers,
		ZoneAssignmentConfigurations: zoneAssignment,
	}))
	if err != nil {
		return nil, err
	}
	return createdProject.Payload, nil
}

// UpdateProject - Update Project
func UpdateProject(apiClient *client.MulticloudIaaS, apiVersion string, id string, name string, description string, administrators []*models.User, members []*models.User, viewers []*models.User, zoneAssignment []*models.ZoneAssignmentSpecification, constraints map[string][]models.Constraint, operationTimeout int64, machineNamingTemplate string, sharedResources *bool) (*models.IaaSProject, error) {
	ProjectSpecification := models.IaaSProjectSpecification{}

	if len(administrators) > 0 {
		ProjectSpecification.Administrators = administrators
	}
	if len(members) > 0 {
		ProjectSpecification.Members = members
	}
	if len(viewers) > 0 {
		ProjectSpecification.Viewers = viewers
	}
	if len(zoneAssignment) > 0 {
		ProjectSpecification.ZoneAssignmentConfigurations = zoneAssignment
	}
	if len(constraints) > 0 {
		ProjectSpecification.Constraints = constraints
	}
	if name != "" {
		ProjectSpecification.Name = &name
	}
	if description != "" {
		ProjectSpecification.Description = description
	}
	if operationTimeout != 0 {
		ProjectSpecification.OperationTimeout = &operationTimeout
	}
	if machineNamingTemplate != "" {
		ProjectSpecification.MachineNamingTemplate = machineNamingTemplate
	}
	if *sharedResources == bool(true) {
		ProjectSpecification.SharedResources = *sharedResources
	}

	updatedProject, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithAPIVersion(&apiVersion).WithID(id).WithBody(&ProjectSpecification))
	if err != nil {
		return nil, err
	}
	return updatedProject.Payload, nil
}
