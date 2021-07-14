/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import "time"

//  *** Cloud Assembly ***

// CloudAssemblyException - Generic exception struct
type CloudAssemblyException struct {
	Message     string      `json:"message"`
	StatusCode  int         `json:"statusCode"`
	ErrorCode   int         `json:"errorCode"`
	ReferenceID interface{} `json:"referenceId"`
}

// CloudAssemblyCloudTemplate - Struct
type CloudAssemblyCloudTemplate struct {
	ID                        string        `json:"id"`
	CreatedAt                 time.Time     `json:"createdAt"`
	CreatedBy                 string        `json:"createdBy"`
	UpdatedAt                 time.Time     `json:"updatedAt"`
	UpdatedBy                 string        `json:"updatedBy"`
	OrgID                     string        `json:"orgId"`
	ProjectID                 string        `json:"projectId"`
	ProjectName               string        `json:"projectName"`
	SelfLink                  string        `json:"selfLink"`
	Name                      string        `json:"name"`
	Status                    string        `json:"status"`
	Content                   string        `json:"content"`
	Valid                     bool          `json:"valid"`
	ValidationMessages        []interface{} `json:"validationMessages"`
	TotalVersions             int           `json:"totalVersions"`
	TotalReleasedVersions     int           `json:"totalReleasedVersions"`
	RequestScopeOrg           bool          `json:"requestScopeOrg"`
	ContentSourceID           string        `json:"contentSourceId"`
	ContentSourcePath         string        `json:"contentSourcePath"`
	ContentSourceType         string        `json:"contentSourceType"`
	ContentSourceSyncStatus   string        `json:"contentSourceSyncStatus"`
	ContentSourceSyncMessages []string      `json:"contentSourceSyncMessages"`
	ContentSourceSyncAt       time.Time     `json:"contentSourceSyncAt"`
}

type CloudAssemblyCloudTemplateRequest struct {
	Content         string `json:"content"`
	Description     string `json:"description"`
	Name            string `json:"name"`
	ProjectID       string `json:"projectId"`
	RequestScopeOrg bool   `json:"requestScopeOrg"`
}

type CloudAssemblyCloudTemplateInputSchema struct {
	Type       string                 `json:"type"`
	Encrypted  bool                   `json:"encrypted"`
	Required   []string               `json:"required"`
	Properties map[string]interface{} `json:"properties"`
}

type CloudAssemblyCloudTemplateInputProperty struct {
	Type      string `json:"type"`
	Encrypted bool   `json:"encrypted"`
	OneOf     []struct {
		Encrypted        bool   `json:"encrypted"`
		Computed         bool   `json:"computed"`
		RecreateOnUpdate bool   `json:"recreateOnUpdate"`
		IgnoreOnUpdate   bool   `json:"ignoreOnUpdate"`
		IgnoreCaseOnDiff bool   `json:"ignoreCaseOnDiff"`
		Title            string `json:"title"`
		Const            string `json:"const"`
	} `json:"oneOf"`
	Enum        []string `json:"enum"`
	Title       string   `json:"title"`
	Default     string   `json:"default"`
	Description string   `json:"description"`
	Pattern     string   `json:"pattern"`
	MaxLength   int      `json:"maxLength"`
	MinLength   int      `json:"minLength"`
}

type DeploymentInput struct {
	Inputs map[string]interface{} `json:"inputs"`
}
