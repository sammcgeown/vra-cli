/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client/data_collector"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func getDataCollectors(id string) ([]*models.DataCollector, error) {
	var dataCollectors []*models.DataCollector

	apiclient := getApiClient()

	if id != "" || name != "" {
		// Get Data Collector by ID or Name
		log.Debug("Getting Data Collector by ID: ", id)
		ret, err := apiclient.DataCollector.GetDataCollector(data_collector.NewGetDataCollectorParams().WithID(id))
		if err != nil {
			return nil, err
		} else {
			dataCollectors = append(dataCollectors, ret.Payload)
			return dataCollectors, err
		}
	} else {
		log.Debug("Getting Data Collectors")
		ret, err := apiclient.DataCollector.GetDataCollectors(data_collector.NewGetDataCollectorsParams())
		if err != nil {
			return nil, err
		} else {
			return ret.Payload.Content, err
		}
	}
}
