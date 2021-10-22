/*
Package orchestrator Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package orchestrator

import (
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
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

// GetCategory returns the category
func GetCategory(client *resty.Client) ([]*types.CategoryContext, error) {
	var Categories []*types.CategoryContext

	queryResponse, err := client.R().
		SetResult(&types.InventoryItemsList{}).
		SetError(&types.Exception{}).
		Get("/vco/api/categories")

	if err != nil {
		return nil, err
	}

	for _, category := range queryResponse.Result().(*types.InventoryItemsList).Link {
		c := types.CategoryContext{}
		mapstructure.Decode(category, &c)
		Categories = append(Categories, &c)
	}
	return Categories, nil
}
