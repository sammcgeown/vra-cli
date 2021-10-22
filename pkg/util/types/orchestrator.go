/*
Package types Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package types

// *** Orchestrator Types ***

// InventoryItemsList is a list of InventoryItems
type InventoryItemsList struct {
	Link []struct {
		Attributes []struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"attributes"`
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"link"`
	Total int `json:"total"`
}

// WsWorkflow is a workflow
type WsWorkflow struct {
	Href      string `json:"href"`
	Relations struct {
		Link []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"link"`
	} `json:"relations"`
	ID              string `json:"id"`
	CustomizedIcon  bool   `json:"customized-icon"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Description     string `json:"description"`
	CategoryID      string `json:"category-id"`
	InputParameters []struct {
		Description string `json:"description"`
		Type        string `json:"type"`
		Name        string `json:"name"`
	} `json:"input-parameters"`
	OutputParameters []struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"output-parameters"`
}

// CategoryContext is a category context
type CategoryContext struct {
	CategoryType     string `json:"categoryType"`
	Description      string `json:"description"`
	Name             string `json:"name"`
	ParentCategoryID string `json:"parentCategoryId"`
}

// WsCategory is a category
type WsCategory struct {
	Description string   `json:"description"`
	Href        string   `json:"href"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	PathIds     []string `json:"pathIds"`
	Relations   struct {
		Count int `json:"count"`
		Links []struct {
			Attributes []struct {
				DisplayName  string `json:"displayName"`
				DisplayValue string `json:"displayValue"`
				Name         string `json:"name"`
				Value        string `json:"value"`
			} `json:"attributes"`
			Href string `json:"href"`
			Rel  string `json:"rel"`
			Type string `json:"type"`
		} `json:"links"`
		StartIndex int `json:"startIndex"`
	} `json:"relations"`
	Type string `json:"type"`
}
