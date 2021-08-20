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

func getDeployments(id string) ([]*Deployment, error) {
	var arrResults []*Deployment
	client := resty.New()

	if id != "" {
		queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetQueryParams(qParams).
			SetHeader("Accept", "application/json").
			SetResult(&Deployment{}).
			SetAuthToken(targetConfig.AccessToken).
			SetError(&CodeStreamException{}).
			Get("https://" + targetConfig.Server + "/deployment/api/deployments/" + id)

		log.Debugln(queryResponse.Request.RawRequest.URL)
		// log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
		}

		arrResults = append(arrResults, queryResponse.Result().(*Deployment))
	} else {
		queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetQueryParams(qParams).
			SetHeader("Accept", "application/json").
			SetResult(&contentsList{}).
			SetAuthToken(targetConfig.AccessToken).
			SetError(&CodeStreamException{}).
			Get("https://" + targetConfig.Server + "/deployment/api/deployments")

		log.Debugln(queryResponse.Request.RawRequest.URL)
		// log.Debugln(queryResponse.String())

		if queryResponse.IsError() {
			return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
		}

		for _, value := range queryResponse.Result().(*contentsList).Content {
			c := Deployment{}
			mapstructure.Decode(value, &c)
			arrResults = append(arrResults, &c)
		}

	}
	return arrResults, nil
}
