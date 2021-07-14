/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"time"
)

// *** Service Broker ***

type CatalogItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type struct {
		Id   string `json:"id"`
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"type"`
	Projects []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"projects"`
	CreatedAt        time.Time `json:"createdAt"`
	CreatedBy        string    `json:"createdBy"`
	LastUpdatedAt    time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy    string    `json:"lastUpdatedBy"`
	BulkRequestLimit int       `json:"bulkRequestLimit"`
}
