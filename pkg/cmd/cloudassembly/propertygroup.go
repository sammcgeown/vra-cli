/*
Package cloudassembly Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cloudassembly

import (
	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	"github.com/vmware/vra-sdk-go/pkg/client/property_groups"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetPropertyGroups returns a list of Property Groups
func GetPropertyGroups(APIClient *types.APIClientOptions, id string, name string, project string) ([]*models.PropertyGroup, error) {
	// Get deployment by ID
	if id != "" {
		PropertyGroupsParams := property_groups.NewGetPropertyGroupUsingGETParams().
			WithAPIVersion(&APIClient.Version).
			WithPropertyGroupID(strfmt.UUID(id))
		PropertyGroups, err := APIClient.SDKClient.PropertyGroups.GetPropertyGroupUsingGET(PropertyGroupsParams)
		if err != nil {
			return nil, err
		}

		return []*models.PropertyGroup{PropertyGroups.Payload}, nil
	}

	PropertyGroupsParams := property_groups.NewListPropertyGroupsUsingGETParams().
		WithAPIVersion(&APIClient.Version)

	if name != "" {
		PropertyGroupsParams.SetName(&name)
	}
	if project != "" {
		p, perr := GetProject(APIClient, project, "")
		if perr != nil {
			return nil, perr
		} else if len(p) == 0 {
			return nil, errors.New("Project not found")
		}
		PropertyGroupsParams.SetProjects([]string{*(p[0]).ID})
	}

	PropertyGroups, err := APIClient.SDKClient.PropertyGroups.ListPropertyGroupsUsingGET(PropertyGroupsParams)
	if err != nil {
		return nil, err
	}
	return PropertyGroups.Payload.Content, nil
}
