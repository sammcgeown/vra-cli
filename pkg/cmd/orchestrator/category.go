/*
Package orchestrator Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package orchestrator

import (
	"errors"

	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
)

// GetCategoryByID returns the category by ID
func GetCategoryByID(APIClient *types.APIClientOptions, id string) (*types.WsCategory, error) {
	var Category *types.WsCategory

	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.WsCategory{}).
		SetError(&types.Exception{}).
		Get("/vco/api/categories/" + id)

	if err != nil {
		return nil, err
	}

	log.Debugln(string(queryResponse.Body()))

	Category = queryResponse.Result().(*types.WsCategory)

	return Category, nil
}

// GetCategoryByName returns the category by name
func GetCategoryByName(APIClient *types.APIClientOptions, categoryName string, categoryType string) ([]*types.WsCategory, error) {
	var Categories []*types.WsCategory
	APIClient.RESTClient.QueryParam.Set("conditions", "name~"+categoryName)

	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.InventoryItemsList{}).
		SetError(&types.Exception{}).
		Get("/vco/api/catalog/System/" + categoryType + "/")

	if err != nil {
		return nil, err
	}

	for _, value := range queryResponse.Result().(*types.InventoryItemsList).Link {
		for _, attribute := range value.Attributes {
			if attribute.Name == "id" {
				Category, _ := GetCategoryByID(APIClient, attribute.Value)
				Categories = append(Categories, Category)
			}

		}
	}
	APIClient.RESTClient.QueryParam.Del("conditions")

	return Categories, nil
}

// GetCategory returns the categories
func GetCategory(APIClient *types.APIClientOptions, root bool, categoryType string) ([]*types.WsCategory, error) {
	var Categories []*types.WsCategory

	if categoryType != "" {
		APIClient.RESTClient.QueryParam.Set("categoryType", categoryType)
	}
	if root {
		APIClient.RESTClient.QueryParam.Set("isRoot", "true")
	}

	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.InventoryItemsList{}).
		SetError(&types.Exception{}).
		Get("/vco/api/categories")

	if err != nil {
		return nil, err
	}

	for _, value := range queryResponse.Result().(*types.InventoryItemsList).Link {
		for _, attribute := range value.Attributes {
			if attribute.Name == "id" {
				Category, _ := GetCategoryByID(APIClient, attribute.Value)
				Categories = append(Categories, Category)
			}

		}
	}
	return Categories, nil
}

// CreateCategory creates a category
func CreateCategory(APIClient *types.APIClientOptions, categoryName string, categoryType string, parentCategoryID string) (*types.WsCategory, error) {
	var categoryURL string

	if parentCategoryID != "" {
		categoryURL = "/" + parentCategoryID
	}

	queryResponse, err := APIClient.RESTClient.R().
		SetBody(types.WsCategoryRequest{
			Name: categoryName,
			Type: categoryType,
		}).
		SetResult(&types.WsCategory{}).
		SetError(&types.Exception{}).
		Post("/vco/api/categories" + categoryURL)

	if err != nil {
		return nil, err
	}

	category := queryResponse.Result().(*types.WsCategory)

	return category, nil
}

// UpdateCategory updates a category
func UpdateCategory(APIClient *types.APIClientOptions, categoryID string, categoryName string, parentCategoryID string) (*types.WsCategory, error) {
	category, err := GetCategoryByID(APIClient, categoryID)
	if err != nil {
		return nil, err
	}
	if parentCategoryID == "" {
		for _, link := range category.Relations.Link {
			if link.Rel == "up" {
				for _, attribute := range link.Attributes {
					if attribute.Name == "id" {
						parentCategoryID = attribute.Value
					}
					break
				}
				break
			}
		}
	}

	queryResponse, err := APIClient.RESTClient.R().
		SetBody(types.WsCategoryRequest{
			Name:             categoryName,
			Type:             category.Type,
			ParentCategoryID: parentCategoryID,
		}).
		SetResult(&types.WsCategory{}).
		SetError(&types.Exception{}).
		Put("/vco/api/categories/" + categoryID)

	if err != nil {
		return nil, err
	}
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	updatedCategory, err := GetCategoryByID(APIClient, categoryID)
	if err != nil {
		return nil, err
	}
	return updatedCategory, nil
}

// DeleteCategory - deletes a category
func DeleteCategory(APIClient *types.APIClientOptions, categoryID string, force bool) error {
	if force {
		APIClient.RESTClient.QueryParam.Set("deleteNonEmptyContent", "true")
	}
	queryResponse, err := APIClient.RESTClient.R().
		SetError(&types.Exception{}).
		Delete("/vco/api/categories/" + categoryID)

	if err != nil {
		return err
	}
	if queryResponse.IsError() {
		return errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return nil
}
