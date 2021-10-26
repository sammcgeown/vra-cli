/*
Package types Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package types

import (
	"time"
)

// *** Service Broker ***

// CatalogItem - A service broker catalog item
type CatalogItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type struct {
		ID   string `json:"id"`
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"type"`
	Projects []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"projects"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     string    `json:"createdBy"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy string    `json:"lastUpdatedBy"`
	Schema        struct {
		Type       string                                 `json:"type"`
		Properties map[string]CatalogItemSchemaProperties `json:"properties"`
	} `json:"schema"`
	BulkRequestLimit int `json:"bulkRequestLimit"`
}

// CatalogItemSchemaProperties - A service broker catalog item schema properties
type CatalogItemSchemaProperties struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Default     string `json:"default"`
}

// CatalogItemRequest - A service broker catalog item request
type CatalogItemRequest struct {
	DeploymentName string            `json:"deploymentName"`
	Inputs         map[string]string `json:"inputs"`
	ProjectID      string            `json:"projectId"`
	Reason         string            `json:"reason"`
	Version        string            `json:"version"`
}

// CatalogItemRequestResponse - A service broker catalog item request response
type CatalogItemRequestResponse struct {
	DeploymentID   string `json:"deploymentId"`
	DeploymentName string `json:"deploymentName"`
}
