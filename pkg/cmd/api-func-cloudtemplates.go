/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	CloudAssembly "github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

func getCloudTemplate(id string, name string, project string, exportPath string) ([]*CloudAssembly.CloudTemplate, error) {
	var arrResults []*CloudAssembly.CloudTemplate
	var error error
	client := resty.New()

	if id != "" {
		queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetQueryParams(qParams).
			SetHeader("Accept", "application/json").
			SetResult(&CloudAssembly.CloudTemplate{}).
			SetAuthToken(targetConfig.AccessToken).
			SetError(&CloudAssembly.Exception{}).
			Get("https://" + targetConfig.Server + "/blueprint/api/blueprints/" + id)

		log.Debugln(queryResponse.Request.RawRequest.URL)
		log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CloudAssembly.Exception).Message)
		}

		result := queryResponse.Result().(*CloudAssembly.CloudTemplate)

		arrResults = append(arrResults, result)

		if exportPath != "" {
			if err := exportCloudTemplate(result.Name, result.ProjectName, result.Content, exportPath); err != nil {
				log.Warnln(err)
			}
		}
		error = err

	} else {
		qParams["$select"] = "*" // Expand the blueprint content!
		if name != "" {
			qParams["name"] = name
		}
		// Project expects an array (e.g. project[]=Development) which I can't figure out
		// how to provide through the resty query parameters...so I'm manually filtering
		// results later.
		// if project != "" {
		// 	qParams["projects"] = project
		// }
		queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetQueryParams(qParams).
			SetHeader("Accept", "application/json").
			SetResult(&contentsList{}).
			SetAuthToken(targetConfig.AccessToken).
			SetError(&CloudAssembly.Exception{}).
			Get("https://" + targetConfig.Server + "/blueprint/api/blueprints")

		log.Debugln(queryResponse.Request.RawRequest.URL)
		log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CloudAssembly.Exception).Message)
		}
		for _, value := range queryResponse.Result().(*contentsList).Content {
			c := CloudAssembly.CloudTemplate{}
			mapstructure.Decode(value, &c)
			if project != "" {
				if project == c.ProjectName {
					arrResults = append(arrResults, &c)
				}
			} else {
				arrResults = append(arrResults, &c)
			}
			if exportPath != "" {
				for _, result := range arrResults {
					if err := exportCloudTemplate(result.Name, result.ProjectName, result.Content, exportPath); err != nil {
						log.Warnln(err)
					}
				}

			}
		}

		error = err
	}

	return arrResults, error

}

// getCloudTemplateSchema
func getCloudTemplateInputSchema(id string) (*CloudAssembly.CloudTemplateInputSchema, error) {
	client := resty.New()

	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CloudAssembly.CloudTemplateInputSchema{}).
		SetAuthToken(targetConfig.AccessToken).
		SetError(&CloudAssembly.Exception{}).
		Get("https://" + targetConfig.Server + "/blueprint/api/blueprints/" + id + "/inputs-schema")

	log.Debugln(queryResponse.Request.RawRequest.URL)

	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CloudAssembly.Exception).Message)
	}

	return queryResponse.Result().(*CloudAssembly.CloudTemplateInputSchema), err

}

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

func deleteCloudTemplate(id string) error {
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetAuthToken(targetConfig.AccessToken).
		SetError(&CloudAssembly.Exception{}).
		Delete("https://" + targetConfig.Server + "/blueprint/api/blueprints/" + id)
	if queryResponse.IsError() {
		return errors.New(queryResponse.Error().(*CloudAssembly.Exception).Message)
	}
	return nil
}

// createCloudTemplate - Create a new Cloud Assembly Cloud Template
func createCloudTemplate(name string, description string, projectId string, content string, scope bool) (*CloudAssembly.CloudTemplate, error) {
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetBody(
			CloudAssembly.CloudTemplateRequest{
				ProjectID:       projectId,
				Name:            name,
				Description:     description,
				Content:         content,
				RequestScopeOrg: scope,
			}).
		SetHeader("Accept", "application/json").
		SetResult(&CloudAssembly.CloudTemplate{}).
		SetError(&CloudAssembly.Exception{}).
		SetAuthToken(targetConfig.AccessToken).
		Post("https://" + targetConfig.Server + "/blueprint/api/blueprints")
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CloudAssembly.Exception).Message)
	}
	newCloudTemplate, cErr := getCloudTemplate(queryResponse.Result().(*CloudAssembly.CloudTemplate).ID, "", "", "")
	return newCloudTemplate[0], cErr
}
