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

func getCatalogItems(id string, name string, project string) ([]*CatalogItem, error) {
	var arrResults []*CatalogItem
	client := resty.New()
	qParams["expandProjects"] = "true"

	if id != "" {
		queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetQueryParams(qParams).
			SetHeader("Accept", "application/json").
			SetResult(&CatalogItem{}).
			SetAuthToken(targetConfig.AccessToken).
			SetError(&CodeStreamException{}).
			Get("https://" + targetConfig.Server + "/catalog/api/items/" + id)

		log.Debugln(queryResponse.Request.RawRequest.URL)
		// log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
		}

		arrResults = append(arrResults, queryResponse.Result().(*CatalogItem))
	} else {
		queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetQueryParams(qParams).
			SetHeader("Accept", "application/json").
			SetResult(&contentsList{}).
			SetAuthToken(targetConfig.AccessToken).
			SetError(&CodeStreamException{}).
			Get("https://" + targetConfig.Server + "/catalog/api/items")

		log.Debugln(queryResponse.Request.RawRequest.URL)
		// log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
		}

		for _, value := range queryResponse.Result().(*contentsList).Content {
			c := CatalogItem{}
			mapstructure.Decode(value, &c)
			arrResults = append(arrResults, &c)
		}

	}
	return arrResults, nil
}

func createCatalogItemRequest(id string, request CatalogItemRequest) (*CatalogItemRequestResponse, error) {
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CatalogItemRequestResponse{}).
		SetAuthToken(targetConfig.AccessToken).
		SetError(&CodeStreamException{}).
		SetBody(request).
		Post("https://" + targetConfig.Server + "/catalog/api/items/" + id + "/request")

	log.Debugln(queryResponse.Request.RawRequest.URL)
	// log.Debugln(queryResponse.String())

	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
	} else {
		response := queryResponse.Result().(*CatalogItemRequestResponse)
		return response, nil
	}

}
