/*
Package cloudassembly Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cloudassembly

import (
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client/data_collector"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetDataCollector gets the data collector
func GetDataCollector(APIClient *types.APIClientOptions, id string) ([]*models.DataCollector, error) {
	var dataCollectors []*models.DataCollector

	if id != "" {
		// Get Data Collector by ID
		log.Debug("Getting Data Collector by ID: ", id)
		ret, err := APIClient.SDKClient.DataCollector.GetDataCollector(data_collector.NewGetDataCollectorParams().WithID(id))
		if err != nil {
			return nil, err
		}
		dataCollectors = append(dataCollectors, ret.Payload)
		return dataCollectors, err

	}
	log.Debug("Getting Data Collectors")
	ret, err := APIClient.SDKClient.DataCollector.GetDataCollectors(data_collector.NewGetDataCollectorsParams())
	if err != nil {
		return nil, err
	}
	return ret.Payload.Content, err

}
