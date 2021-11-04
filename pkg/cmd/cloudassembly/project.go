/*
Package cloudassembly Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cloudassembly

import (
	"errors"
	"strings"

	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetProject - Get Projects
func GetProject(APIClient *types.APIClientOptions, name string, id string) ([]*models.IaaSProject, error) {
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
	ProjectParams.APIVersion = &APIClient.Version

	ret, err := APIClient.SDKClient.Project.GetProjects(ProjectParams)
	if err != nil {
		switch err.(type) {
		case *project.GetProjectNotFound:
			return nil, errors.New("Project with ID " + id + " not found")
		}
		return nil, err
	}
	log.Debugln(ret.Payload.NumberOfElements, "Projects found")
	return ret.Payload.Content, nil
}

// DeleteProject - Delete Project
func DeleteProject(APIClient *types.APIClientOptions, id string) error {

	// Workaround an issue where the cloud regions need to be removed before the project can be deleted.
	_, err := APIClient.SDKClient.Project.UpdateProject(project.NewUpdateProjectParams().WithAPIVersion(&APIClient.Version).WithID(id).WithBody(&models.IaaSProjectSpecification{
		ZoneAssignmentConfigurations: []*models.ZoneAssignmentSpecification{},
	}))
	if err != nil {
		return err
	}

	_, err = APIClient.SDKClient.Project.DeleteProject(project.NewDeleteProjectParams().WithID(id))
	if err != nil {
		return err
	}
	return nil
}

// CreateProject - Create Project
func CreateProject(APIClient *types.APIClientOptions, name string, description string, administrators []*models.User, members []*models.User, viewers []*models.User, zoneAssignment []*models.ZoneAssignmentSpecification, constraints map[string][]models.Constraint, operationTimeout int64, machineNamingTemplate string, sharedResources *bool) (*models.IaaSProject, error) {
	createdProject, err := APIClient.SDKClient.Project.CreateProject(project.NewCreateProjectParams().WithAPIVersion(&APIClient.Version).WithBody(&models.IaaSProjectSpecification{
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
func UpdateProject(APIClient *types.APIClientOptions, id string, name string, description string, administrators []*models.User, members []*models.User, viewers []*models.User, zoneAssignment []*models.ZoneAssignmentSpecification, constraints map[string][]models.Constraint, operationTimeout int64, machineNamingTemplate string, sharedResources *bool) (*models.IaaSProject, error) {
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

	updatedProject, err := APIClient.SDKClient.Project.UpdateProject(project.NewUpdateProjectParams().WithAPIVersion(&APIClient.Version).WithID(id).WithBody(&ProjectSpecification))
	if err != nil {
		return nil, err
	}
	return updatedProject.Payload, nil
}
