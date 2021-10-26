/*
Package cloudassembly Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cloudassembly

import (
	"os"
	"path/filepath"

	"github.com/go-openapi/strfmt"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetCloudTemplate - Get a Cloud Assembly Cloud Template
func GetCloudTemplate(APIClient *types.APIClientOptions, id string, name string, project string) ([]*models.Blueprint, error) {
	var result []*models.Blueprint

	if id == "" {
		CloudTemplateParams := blueprint.NewListBlueprintsUsingGET1Params()
		CloudTemplateParams.DollarSelect = []string{"*"}
		if name != "" {
			CloudTemplateParams.Name = &name
		}
		if project != "" {
			p, perr := GetProject(APIClient, project, "")
			if perr != nil {
				return nil, perr
			}
			CloudTemplateParams.Projects = []string{*(p[0]).ID}
		}

		log.Debug(CloudTemplateParams)

		ret, err := APIClient.SDKClient.Blueprint.ListBlueprintsUsingGET1(CloudTemplateParams)
		if err != nil {
			return nil, err
		}
		result = ret.Payload.Content

	} else {
		CloudTemplateParams := blueprint.NewGetBlueprintUsingGET1Params()
		CloudTemplateParams.BlueprintID = strfmt.UUID(id)

		ret, err := APIClient.SDKClient.Blueprint.GetBlueprintUsingGET1(CloudTemplateParams)
		if err != nil {
			return nil, err
		}
		result = append(result, ret.Payload)

	}
	return result, nil
}

// GetCloudTemplateInputSchema - Get a Cloud Assembly Cloud Template Schema
// func GetCloudTemplateInputSchema(APIClient *types.APIClientOptions, id string) (*models., error) {
// 	SchemaParams := blueprint.NewGetBlueprintInputsSchemaUsingGET1Params()
// 	SchemaParams.BlueprintID = id

// 	ret, err := APIClient.SDKClient.Blueprint.GetBlueprintInputsSchemaUsingGET1(SchemaParams)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return ret.Payload, err

// }

// // patchBlueprint - Patch Code Stream Blueprint by ID
// func patchBlueprint(id string, payload string) (*types.CodeStreamBlueprint, error) {
// 	client := resty.New()
// 	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetHeader("Accept", "application/json").
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(payload).
// 		SetResult(&CodeStreamBlueprint{}).
// 		SetAuthToken(targetConfig.AccessToken).
// 		Patch("https://" + targetConfig.Server + "/Blueprint/api/Blueprints/" + id)
// 	if queryResponse.IsError() {
// 		return nil, queryResponse.Error().(error)
// 	}
// 	return queryResponse.Result().(*CodeStreamBlueprint), err
// }

// DeleteCloudTemplate - Delete a Cloud Assembly Cloud Template
func DeleteCloudTemplate(APIClient *types.APIClientOptions, id string) error {
	DeleteParams := blueprint.NewDeleteBlueprintUsingDELETE1Params()
	DeleteParams.BlueprintID = strfmt.UUID(id)
	_, err := APIClient.SDKClient.Blueprint.DeleteBlueprintUsingDELETE1(DeleteParams)
	if err != nil {
		return err
	}
	return nil
}

// CreateCloudTemplate - Create a new Cloud Assembly Cloud Template
func CreateCloudTemplate(APIClient *types.APIClientOptions, name string, description string, projectID string, content string, scope bool) (*models.Blueprint, error) {
	CreateParams := blueprint.NewCreateBlueprintUsingPOST1Params()
	CreateParams.Blueprint = &models.Blueprint{
		Name:            name,
		Description:     description,
		ProjectID:       projectID,
		Content:         content,
		RequestScopeOrg: scope,
	}
	ret, err := APIClient.SDKClient.Blueprint.CreateBlueprintUsingPOST1(CreateParams)
	if err != nil {
		return nil, err
	}
	return ret.Payload, err
}

// ExportCloudTemplate - Export a Cloud Assembly Cloud Template
func ExportCloudTemplate(name, project, content, path string) error {
	var exportPath string
	if path != "" {
		exportPath = path
		_, folderError := os.Stat(exportPath) // Get file system info
		if os.IsNotExist(folderError) {       // If it doesn't exist
			log.Debugln("Folder doesn't exist - creating")
			mkdirErr := os.MkdirAll(exportPath, os.FileMode(0755)) // Attempt to make it
			if mkdirErr != nil {
				return mkdirErr
			}
		}
	} else {
		// If path is not specified, use the current path
		exportPath, _ = os.Getwd()
	}
	exportPath = filepath.Join(exportPath, project+" - "+name+".yaml")
	f, cerr := os.Create(exportPath) // Open the file for writing
	if cerr != nil {
		return cerr
	}
	defer f.Close() // Defer closing until just before this function returns
	_, werr := f.WriteString(content)
	if werr != nil {
		return werr
	}
	return nil
}
