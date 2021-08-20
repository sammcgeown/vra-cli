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
	log "github.com/sirupsen/logrus"
)

func getCloudTemplate(id string, name string, project string, exportPath string) ([]*CloudAssemblyCloudTemplate, error) {
	var arrResults []*CloudAssemblyCloudTemplate
	var error error
	client := resty.New()

	if id != "" {
		queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetQueryParams(qParams).
			SetHeader("Accept", "application/json").
			SetResult(&CloudAssemblyCloudTemplate{}).
			SetAuthToken(targetConfig.accesstoken).
			SetError(&CloudAssemblyException{}).
			Get("https://" + targetConfig.server + "/blueprint/api/blueprints/" + id)

		log.Debugln(queryResponse.Request.RawRequest.URL)
		log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CloudAssemblyException).Message)
		}

		result := queryResponse.Result().(*CloudAssemblyCloudTemplate)

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
			SetAuthToken(targetConfig.accesstoken).
			SetError(&CloudAssemblyException{}).
			Get("https://" + targetConfig.server + "/blueprint/api/blueprints")

		log.Debugln(queryResponse.Request.RawRequest.URL)
		log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CloudAssemblyException).Message)
		}
		for _, value := range queryResponse.Result().(*contentsList).Content {
			c := CloudAssemblyCloudTemplate{}
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
func getCloudTemplateInputSchema(id string) (*CloudAssemblyCloudTemplateInputSchema, error) {
	client := resty.New()

	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CloudAssemblyCloudTemplateInputSchema{}).
		SetAuthToken(targetConfig.accesstoken).
		SetError(&CloudAssemblyException{}).
		Get("https://" + targetConfig.server + "/blueprint/api/blueprints/" + id + "/inputs-schema")

	log.Debugln(queryResponse.Request.RawRequest.URL)

	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CloudAssemblyException).Message)
	}

	return queryResponse.Result().(*CloudAssemblyCloudTemplateInputSchema), err

}

// // patchBlueprint - Patch Code Stream Blueprint by ID
// func patchBlueprint(id string, payload string) (*CodeStreamBlueprint, error) {
// 	client := resty.New()
// 	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetHeader("Accept", "application/json").
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(payload).
// 		SetResult(&CodeStreamBlueprint{}).
// 		SetAuthToken(targetConfig.accesstoken).
// 		Patch("https://" + targetConfig.server + "/Blueprint/api/Blueprints/" + id)
// 	if queryResponse.IsError() {
// 		return nil, queryResponse.Error().(error)
// 	}
// 	return queryResponse.Result().(*CodeStreamBlueprint), err
// }

func deleteCloudTemplate(id string) error {
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetAuthToken(targetConfig.accesstoken).
		SetError(&CloudAssemblyException{}).
		Delete("https://" + targetConfig.server + "/blueprint/api/blueprints/" + id)
	if queryResponse.IsError() {
		return errors.New(queryResponse.Error().(*CloudAssemblyException).Message)
	}
	return nil
}

// createCloudTemplate - Create a new Cloud Assembly Cloud Template
func createCloudTemplate(name string, description string, projectId string, content string, scope bool) (*CloudAssemblyCloudTemplate, error) {
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetBody(
			CloudAssemblyCloudTemplateRequest{
				ProjectID:       projectId,
				Name:            name,
				Description:     description,
				Content:         content,
				RequestScopeOrg: scope,
			}).
		SetHeader("Accept", "application/json").
		SetResult(&CloudAssemblyCloudTemplate{}).
		SetError(&CloudAssemblyException{}).
		SetAuthToken(targetConfig.accesstoken).
		Post("https://" + targetConfig.server + "/blueprint/api/blueprints")
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CloudAssemblyException).Message)
	}
	newCloudTemplate, cErr := getCloudTemplate(queryResponse.Result().(*CloudAssemblyCloudTemplate).ID, "", "", "")
	return newCloudTemplate[0], cErr
}
