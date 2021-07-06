/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"errors"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func getCloudTemplates(id string, name string, project string, exportPath string) ([]*CloudAssemblyCloudTemplate, error) {
	var arrResults []*CloudAssemblyCloudTemplate
	var idUrl string
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

		arrResults = append(arrResults, queryResponse.Result().(*CloudAssemblyCloudTemplate))
		return arrResults, err

	} else {
		var filters []string
		if name != "" {
			filters = append(filters, "(name eq '"+name+"')")
		}
		if project != "" {
			filters = append(filters, "(project eq '"+project+"')")
		}
		if len(filters) > 0 {
			qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
		}
		queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetQueryParams(qParams).
			SetHeader("Accept", "application/json").
			SetResult(&contentsList{}).
			SetAuthToken(targetConfig.accesstoken).
			SetError(&CloudAssemblyException{}).
			Get("https://" + targetConfig.server + "/blueprint/api/blueprints" + idUrl)

		log.Debugln(queryResponse.Request.RawRequest.URL)
		log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CloudAssemblyException).Message)
		}
		for _, value := range queryResponse.Result().(*contentsList).Content {
			c := CloudAssemblyCloudTemplate{}
			mapstructure.Decode(value, &c)
			if exportPath != "" {
				if err := exportYaml(c.Name, c.ProjectName, exportPath, "Blueprints"); err != nil {
					log.Warnln(err)
				}
				arrResults = append(arrResults, &c)
			} else {
				arrResults = append(arrResults, &c)
			}
		}
		return arrResults, err
	}
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

// func deleteBlueprint(id string) (*CodeStreamBlueprint, error) {
// 	client := resty.New()
// 	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&CodeStreamBlueprint{}).
// 		SetAuthToken(targetConfig.accesstoken).
// 		Delete("https://" + targetConfig.server + "/Blueprint/api/Blueprints/" + id)
// 	if queryResponse.IsError() {
// 		return nil, queryResponse.Error().(error)
// 	}
// 	return queryResponse.Result().(*CodeStreamBlueprint), err
// }
