/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import "time"

type Deployment struct {
	Id                 string            `json:"id"`
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	OrgId              string            `json:"orgId"`
	BlueprintId        string            `json:"blueprintId"`
	CreatedAt          time.Time         `json:"createdAt"`
	CreatedBy          string            `json:"createdBy"`
	LastUpdatedAt      time.Time         `json:"lastUpdatedAt"`
	LastUpdatedBy      string            `json:"lastUpdatedBy"`
	Inputs             map[string]string `json:"inputs"`
	ProjectId          string            `json:"projectId"`
	Status             string            `json:"status"`
	BlueprintVersion   string            `json:"blueprintVersion"`
	CatalogItemId      string            `json:"catalogItemId"`
	IconId             string            `json:"iconId"`
	CatalogItemVersion string            `json:"catalogItemVersion"`
}
