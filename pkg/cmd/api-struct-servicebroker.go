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

type CatalogItemSchemaProperties struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Default     string `json:"default"`
}

type CatalogItemRequest struct {
	DeploymentName string            `json:"deploymentName"`
	Inputs         map[string]string `json:"inputs"`
	ProjectId      string            `json:"projectId"`
	Reason         string            `json:"reason"`
	Version        string            `json:"version"`
}

type CatalogItemRequestResponse struct {
	DeploymentId   string `json:"deploymentId"`
	DeploymentName string `json:"deploymentName"`
}
