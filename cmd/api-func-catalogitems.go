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

func getCatalogItems(name string, project string) ([]*CatalogItem, error) {
	var arrResults []*CatalogItem
	client := resty.New()
	qParams["expandProjects"] = "true"

	// var filters []string
	// if id != "" {
	// 	filters = append(filters, "(id eq '"+id+"')")
	// }
	// if name != "" {
	// 	filters = append(filters, "(name eq '"+name+"')")
	// }
	// if project != "" {
	// 	filters = append(filters, "(project eq '"+project+"')")
	// }
	// if len(filters) > 0 {
	// 	qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
	// }
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&contentsList{}).
		SetAuthToken(targetConfig.accesstoken).
		SetError(&CodeStreamException{}).
		Get("https://" + targetConfig.server + "/catalog/api/items")

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
	return arrResults, err
}
