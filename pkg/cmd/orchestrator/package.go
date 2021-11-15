/*
Package orchestrator Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package orchestrator

import (
	"errors"
	"path/filepath"
	"strconv"

	"github.com/sammcgeown/vra-cli/pkg/util/types"
)

// // GetCategoryByID returns the category by ID
// func GetCategoryByID(APIClient *types.APIClientOptions, id string) (*types.WsCategory, error) {
// 	var Category *types.WsCategory

// 	queryResponse, err := APIClient.RESTClient.R().
// 		SetResult(&types.WsCategory{}).
// 		SetError(&types.Exception{}).
// 		Get("/vco/api/categories/" + id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Debugln(string(queryResponse.Body()))

// 	Category = queryResponse.Result().(*types.WsCategory)

// 	return Category, nil
// }

// // GetCategoryByName returns the category by name
// func GetCategoryByName(APIClient *types.APIClientOptions, categoryName string, categoryType string) ([]*types.WsCategory, error) {
// 	var Categories []*types.WsCategory
// 	APIClient.RESTClient.QueryParam.Set("conditions", "name~"+categoryName)

// 	queryResponse, err := APIClient.RESTClient.R().
// 		SetResult(&types.InventoryItemsList{}).
// 		SetError(&types.Exception{}).
// 		Get("/vco/api/catalog/System/" + categoryType + "/")

// 	if err != nil {
// 		return nil, err
// 	}
// 	APIClient.RESTClient.QueryParam.Del("conditions")

// 	for _, value := range queryResponse.Result().(*types.InventoryItemsList).Link {
// 		for _, attribute := range value.Attributes {
// 			if attribute.Name == "id" {
// 				Category, _ := GetCategoryByID(APIClient, attribute.Value)
// 				Categories = append(Categories, Category)
// 			}

// 		}
// 	}

// 	return Categories, nil
// }

// GetPackage returns all packages
func GetPackage(APIClient *types.APIClientOptions, name string) ([]*types.WsPackage, error) {
	var Categories []*types.WsPackage

	if name != "" {
		queryResponse, err := APIClient.RESTClient.R().
			SetResult(&types.WsPackage{}).
			SetError(&types.Exception{}).
			Get("/vco/api/packages/" + name)
		if err != nil {
			return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
		}
		Categories = append(Categories, queryResponse.Result().(*types.WsPackage))
		return Categories, nil
	}

	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.WsPackages{}).
		SetError(&types.Exception{}).
		Get("/vco/api/packages")

	if err != nil {
		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
	}

	for _, value := range queryResponse.Result().(*types.WsPackages).Link {
		for _, attribute := range value.Attribute {
			if attribute.Name == "name" {
				Category, _ := GetPackage(APIClient, attribute.Value)
				Categories = append(Categories, Category...)
			}

		}
	}
	return Categories, nil
}

// ExportPackage exports a package
func ExportPackage(APIClient *types.APIClientOptions, name string, options types.ExportPackageOptions, exportPath string) error {

	APIClient.RESTClient.QueryParam.Set("exportConfigurationAttributeValues", strconv.FormatBool(options.ExportConfigurationAttributeValues))
	APIClient.RESTClient.QueryParam.Set("exportConfigSecureStringAttributeValues", strconv.FormatBool(options.ExportConfigSecureStringAttributeValues))
	APIClient.RESTClient.QueryParam.Set("exportGlobalTags", strconv.FormatBool(options.ExportGlobalTags))
	var allowedOperations string
	if options.ViewContents {
		allowedOperations = allowedOperations + "v"
	}
	if options.AddToPackage {
		allowedOperations = allowedOperations + "f"
	}
	if options.EditContents {
		allowedOperations = allowedOperations + "e"
	}
	APIClient.RESTClient.QueryParam.Set("allowedOperations", allowedOperations)

	queryResponse, err := APIClient.RESTClient.R().
		SetHeader("accept", "application/zip").
		SetOutput(filepath.Join(exportPath, name+".package")).
		SetResult(&types.WsPackage{}).
		SetError(&types.Exception{}).
		Get("/vco/api/packages/" + name)
	if err != nil {
		return errors.New(queryResponse.Error().(*types.Exception).Message)
	}
	return nil
}

// // CreateCategory creates a category
// func CreateCategory(APIClient *types.APIClientOptions, categoryName string, categoryType string, parentCategoryID string) (*types.WsCategory, error) {
// 	var categoryURL string

// 	if parentCategoryID != "" {
// 		categoryURL = "/" + parentCategoryID
// 	}

// 	queryResponse, err := APIClient.RESTClient.R().
// 		SetBody(types.WsCategoryRequest{
// 			Name: categoryName,
// 			Type: categoryType,
// 		}).
// 		SetResult(&types.WsCategory{}).
// 		SetError(&types.Exception{}).
// 		Post("/vco/api/categories" + categoryURL)

// 	if err != nil {
// 		return nil, err
// 	}

// 	category := queryResponse.Result().(*types.WsCategory)

// 	return category, nil
// }

// // UpdateCategory updates a category
// func UpdateCategory(APIClient *types.APIClientOptions, categoryID string, categoryName string, parentCategoryID string) (*types.WsCategory, error) {
// 	category, err := GetCategoryByID(APIClient, categoryID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if parentCategoryID == "" {
// 		for _, link := range category.Relations.Link {
// 			if link.Rel == "up" {
// 				for _, attribute := range link.Attributes {
// 					if attribute.Name == "id" {
// 						parentCategoryID = attribute.Value
// 					}
// 					break
// 				}
// 				break
// 			}
// 		}
// 	}

// 	queryResponse, err := APIClient.RESTClient.R().
// 		SetBody(types.WsCategoryRequest{
// 			Name:             categoryName,
// 			Type:             category.Type,
// 			ParentCategoryID: parentCategoryID,
// 		}).
// 		SetResult(&types.WsCategory{}).
// 		SetError(&types.Exception{}).
// 		Put("/vco/api/categories/" + categoryID)

// 	if err != nil {
// 		return nil, err
// 	}
// 	if queryResponse.IsError() {
// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
// 	}

// 	updatedCategory, err := GetCategoryByID(APIClient, categoryID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return updatedCategory, nil
// }

// // DeleteCategory - deletes a category
// func DeleteCategory(APIClient *types.APIClientOptions, categoryID string) error {
// 	if APIClient.Force {
// 		APIClient.RESTClient.QueryParam.Set("deleteNonEmptyContent", "true")
// 	}
// 	queryResponse, err := APIClient.RESTClient.R().
// 		SetError(&types.Exception{}).
// 		Delete("/vco/api/categories/" + categoryID)

// 	if err != nil {
// 		return err
// 	}
// 	if queryResponse.IsError() {
// 		return errors.New(queryResponse.Error().(*types.Exception).Message)
// 	}
// 	return nil
// }
