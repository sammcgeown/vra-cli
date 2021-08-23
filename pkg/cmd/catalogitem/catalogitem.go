/*
Package catalogitem Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package catalogitem

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
)

func GetCatalogItems(client *resty.Client, id string, name string, project string) ([]*types.CatalogItem, error) {
	var arrResults []*types.CatalogItem

	client.QueryParam.Add("expandProjects", "true")

	if id != "" {
		queryResponse, _ := client.R().
			SetResult(&types.CatalogItem{}).
			SetError(&types.Exception{}).
			Get("/catalog/api/items/" + id)

		log.Debugln(queryResponse.Request.RawRequest.URL)
		// log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
		}

		arrResults = append(arrResults, queryResponse.Result().(*types.CatalogItem))
	} else {
		queryResponse, _ := client.R().
			SetResult(&types.ContentsList{}).
			SetError(&types.Exception{}).
			Get("/catalog/api/items")

		log.Debugln(queryResponse.Request.RawRequest.URL)
		// log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
		}

		for _, value := range queryResponse.Result().(*types.ContentsList).Content {
			c := types.CatalogItem{}
			mapstructure.Decode(value, &c)
			arrResults = append(arrResults, &c)
		}

	}
	return arrResults, nil
}

// func createCatalogItemRequest(id string, request types.CatalogItemRequest) (*types.CatalogItemRequestResponse, error) {
// 	client := resty.New()
// 	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&types.CatalogItemRequestResponse{}).
// 		SetAuthToken(targetConfig.AccessToken).
// 		SetError(&types.Exception{}).
// 		SetBody(request).
// 		Post("https://" + targetConfig.Server + "/catalog/api/items/" + id + "/request")

// 	log.Debugln(queryResponse.Request.RawRequest.URL)
// 	// log.Debugln(queryResponse.String())

// 	if queryResponse.IsError() {
// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
// 	} else {
// 		response := queryResponse.Result().(*types.CatalogItemRequestResponse)
// 		return response, nil
// 	}

// }
