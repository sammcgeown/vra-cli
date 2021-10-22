/*
Package orchestrator Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package orchestrator

import (
	"github.com/go-resty/resty/v2"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
)

var (
	apiVersion          = "2019-10-17"
	expandProjects bool = true
)

// GetCategoryByID returns the category by ID
func GetCategoryByID(client *resty.Client, id string) (*types.WsCategory, error) {
	var Category *types.WsCategory

	queryResponse, err := client.R().
		SetResult(&types.WsCategory{}).
		SetError(&types.Exception{}).
		Get("/vco/api/categories/" + id)

	if err != nil {
		return nil, err
	}

	Category = queryResponse.Result().(*types.WsCategory)

	return Category, nil
}

// GetCategoryByName returns the category by ID
func GetCategoryByName(client *resty.Client, name string) ([]*types.WsCategory, error) {
	var Categories []*types.WsCategory
	client.QueryParam.Set("conditions", "name~"+name)

	queryResponse, err := client.R().
		SetResult(&types.InventoryItemsList{}).
		SetError(&types.Exception{}).
		Get("/vco/api/catalog/System/WorkflowCategory/")

	if err != nil {
		return nil, err
	}

	for _, value := range queryResponse.Result().(*types.InventoryItemsList).Link {
		for _, attribute := range value.Attributes {
			if attribute.Name == "id" {
				Category, _ := GetCategoryByID(client, attribute.Value)
				Categories = append(Categories, Category)
			}

		}
	}

	return Categories, nil
}

// GetCategory returns the category
// func GetCategory(client *resty.Client, root bool) ([]*types.CategoryContext, error) {
// 	var Categories []*types.CategoryContext

// 	if root {
// 		client.SetQueryParam("isRoot", "true")
// 	}

// 	queryResponse, err := client.R().
// 		SetResult(&types.InventoryItemsList{}).
// 		SetError(&types.Exception{}).
// 		Get("/vco/api/categories")

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, category := range queryResponse.Result().(*types.InventoryItemsList).Link {
// 		Categories = append(Categories, category)
// 	}
// 	return Categories, nil
// }
