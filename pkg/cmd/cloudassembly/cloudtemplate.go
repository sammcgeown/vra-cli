/*
Package cloudassembly Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cloudassembly

import (
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetCloudTemplate - Get a Cloud Assembly Cloud Template
func GetCloudTemplate(apiclient *client.MulticloudIaaS, id string, name string, project string, exportPath string) ([]*models.Blueprint, error) {
	var result []*models.Blueprint

	if id == "" {
		CloudTemplateParams := blueprint.NewListBlueprintsUsingGET1Params()
		if name != "" {
			CloudTemplateParams.Name = &name
		}
		if project != "" {
			p, perr := GetProject(apiclient, "2019-10-17", project, "")
			if perr != nil {
				return nil, perr
			}
			CloudTemplateParams.Projects = []string{*p[0].ID}
		}

		log.Debug(CloudTemplateParams)

		ret, err := apiclient.Blueprint.ListBlueprintsUsingGET1(CloudTemplateParams)
		if err != nil {
			return nil, err
		}
		if len(ret.Payload.Content) == 1 {
			result, err = GetCloudTemplate(apiclient, ret.Payload.Content[0].ID, "", "", "")
		} else {
			result = ret.Payload.Content
		}
	} else {
		CloudTemplateParams := blueprint.NewGetBlueprintUsingGET1Params()
		CloudTemplateParams.BlueprintID = strfmt.UUID(id)

		ret, err := apiclient.Blueprint.GetBlueprintUsingGET1(CloudTemplateParams)
		if err != nil {
			return nil, err
		}
		result = append(result, ret.Payload)

	}
	return result, nil

	// if id != "" {
	// 	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
	// 		SetQueryParams(qParams).
	// 		SetHeader("Accept", "application/json").
	// 		SetResult(&types.CloudTemplate{}).
	// 		SetAuthToken(targetConfig.AccessToken).
	// 		SetError(&types.Exception{}).
	// 		Get("https://" + targetConfig.Server + "/blueprint/api/blueprints/" + id)

	// 	log.Debugln(queryResponse.Request.RawRequest.URL)
	// 	log.Debugln(queryResponse.String())

	// 	if queryResponse.IsError() {
	// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	// 	}

	// 	result := queryResponse.Result().(*types.CloudTemplate)

	// 	arrResults = append(arrResults, result)

	// 	if exportPath != "" {
	// 		if err := exportCloudTemplate(result.Name, result.ProjectName, result.Content, exportPath); err != nil {
	// 			log.Warnln(err)
	// 		}
	// 	}
	// 	error = err

	// } else {
	// 	qParams["$select"] = "*" // Expand the blueprint content!
	// 	if name != "" {
	// 		qParams["name"] = name
	// 	}
	// 	// Project expects an array (e.g. project[]=Development) which I can't figure out
	// 	// how to provide through the resty query parameters...so I'm manually filtering
	// 	// results later.
	// 	// if project != "" {
	// 	// 	qParams["projects"] = project
	// 	// }
	// 	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
	// 		SetQueryParams(qParams).
	// 		SetHeader("Accept", "application/json").
	// 		SetResult(&types.ContentsList{}).
	// 		SetAuthToken(targetConfig.AccessToken).
	// 		SetError(&types.Exception{}).
	// 		Get("https://" + targetConfig.Server + "/blueprint/api/blueprints")

	// 	log.Debugln(queryResponse.Request.RawRequest.URL)
	// 	log.Debugln(queryResponse.String())

	// 	if queryResponse.IsError() {
	// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	// 	}
	// 	for _, value := range queryResponse.Result().(*types.ContentsList).Content {
	// 		c := types.CloudTemplate{}
	// 		mapstructure.Decode(value, &c)
	// 		if project != "" {
	// 			if project == c.ProjectName {
	// 				arrResults = append(arrResults, &c)
	// 			}
	// 		} else {
	// 			arrResults = append(arrResults, &c)
	// 		}
	// 		if exportPath != "" {
	// 			for _, result := range arrResults {
	// 				if err := exportCloudTemplate(result.Name, result.ProjectName, result.Content, exportPath); err != nil {
	// 					log.Warnln(err)
	// 				}
	// 			}

	// 		}
	// 	}

	// 	error = err
	// }

	// return arrResults, error

}

// getCloudTemplateSchema
// func getCloudTemplateInputSchema(id string) (*types.CloudTemplateInputSchema, error) {
// 	client := resty.New()

// 	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&types.CloudTemplateInputSchema{}).
// 		SetAuthToken(targetConfig.AccessToken).
// 		SetError(&types.Exception{}).
// 		Get("https://" + targetConfig.Server + "/blueprint/api/blueprints/" + id + "/inputs-schema")

// 	log.Debugln(queryResponse.Request.RawRequest.URL)

// 	if queryResponse.IsError() {
// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
// 	}

// 	return queryResponse.Result().(*types.CloudTemplateInputSchema), err

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

// func deleteCloudTemplate(id string) error {
// 	client := resty.New()
// 	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetAuthToken(targetConfig.AccessToken).
// 		SetError(&types.Exception{}).
// 		Delete("https://" + targetConfig.Server + "/blueprint/api/blueprints/" + id)
// 	if queryResponse.IsError() {
// 		return errors.New(queryResponse.Error().(*types.Exception).Message)
// 	}
// 	return nil
// }

// createCloudTemplate - Create a new Cloud Assembly Cloud Template
// func createCloudTemplate(name string, description string, projectID string, content string, scope bool) (*types.CloudTemplate, error) {
// 	client := resty.New()
// 	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetBody(
// 			types.CloudTemplateRequest{
// 				ProjectID:       projectID,
// 				Name:            name,
// 				Description:     description,
// 				Content:         content,
// 				RequestScopeOrg: scope,
// 			}).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&types.CloudTemplate{}).
// 		SetError(&types.Exception{}).
// 		SetAuthToken(targetConfig.AccessToken).
// 		Post("https://" + targetConfig.Server + "/blueprint/api/blueprints")
// 	if queryResponse.IsError() {
// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
// 	}
// 	newCloudTemplate, cErr := getCloudTemplate(queryResponse.Result().(*types.CloudTemplate).ID, "", "", "")
// 	return newCloudTemplate[0], cErr
// }
