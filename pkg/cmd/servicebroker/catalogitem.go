/*
Package servicebroker Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package servicebroker

import (
	"github.com/go-openapi/strfmt"
	"github.com/sammcgeown/vra-cli/pkg/cmd/cloudassembly"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_items"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetCatalogItems returns a catalog item by name, id or project
func GetCatalogItems(apiclient *client.MulticloudIaaS, id string, name string, project string) ([]*models.CatalogItem, error) {

	if id != "" {
		CatalogItemParams := catalog_items.
			NewGetCatalogItemUsingGET1Params().
			WithAPIVersion(&apiVersion).
			WithID(strfmt.UUID(id)).
			WithExpandProjects(&expandProjects)
		catalogItems, err := apiclient.CatalogItems.GetCatalogItemUsingGET1(CatalogItemParams)
		if err != nil {
			return nil, err
		}
		return []*models.CatalogItem{catalogItems.Payload}, nil
	}
	CatalogItemParams := catalog_items.
		NewGetCatalogItemsUsingGET1Params().
		WithAPIVersion(&apiVersion).
		WithExpandProjects(&expandProjects)

	if project != "" {
		Projects, err := cloudassembly.GetProject(apiclient, apiVersion, project, "")
		if err != nil {
			return nil, err
		}
		ProjectID := Projects[0].ID
		CatalogItemParams.WithProjects([]string{*ProjectID})
	}
	if name != "" {
		CatalogItemParams.WithSearch(&name)
	}

	catalogItems, err := apiclient.CatalogItems.GetCatalogItemsUsingGET1(CatalogItemParams)
	if err != nil {
		return nil, err
	}
	return catalogItems.Payload.Content, nil
}

// func createCatalogItemRequest(id string, request types.CatalogItemRequest) (*types.CatalogItemRequestResponse, error) {
// 	client := resty.New()
// 	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&types.CatalogItemRequestResponse{}).
// 		SetAuthToken(targetConfig.AccessToken).
// 		SetError(&types.Exception{}).
// 		SetBody(request).
// 		Post("https://" + targetConfig.Server + "/catalog/api/items/" + id + "/request")

// 	log.Debugln(queryResponse.Request.RawRequest.URL)
// 	// log.Debugln(queryResponse.String())

// 	if queryResponse.IsError() {
// 		return nil, errors.New(queryResponse.Error().(*types.Exception).Message)
// 	} else {
// 		response := queryResponse.Result().(*types.CatalogItemRequestResponse)
// 		return response, nil
// 	}

// }
