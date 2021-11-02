/*
Package orchestrator Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package orchestrator

import (
	"os"
	"testing"

	"github.com/sammcgeown/vra-cli/pkg/util/auth"
	"github.com/sammcgeown/vra-cli/pkg/util/config"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

var (
	APIClient = &types.APIClientOptions{
		Version: "2019-10-17",
		Debug:   false,
	}
	rootCategory = &types.WsCategory{
		Name: "vra-cli-testing",
		Type: "WorkflowCategory",
	}
	childCategory = []struct{ Name, Type string }{
		{"childCategory1", "WorkflowCategory"},
		{"childCategory2", "WorkflowCategory"},
		{"childCategory3", "WorkflowCategory"},
	}
	createdRootCategory  *types.WsCategory
	createdChildCategory = []*types.WsCategory{}
)

func TestMain(m *testing.M) {
	// Configure Logging
	if APIClient.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Debug logging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// Configure API Client
	APIClient.Config = config.GetConfigFromEnv()
	err := auth.ValidateConfiguration(APIClient)
	if err != nil {
		log.Fatalln(err)
	}
	// Clean environment
	CleanUp()
	// Run tests
	code := m.Run()
	// Clean up after tests
	//CleanUp()
	// Exit
	os.Exit(code)
}

func CleanUp() {
	// Delete the category, force delete if it has content
	categories, err := GetCategoryByName(APIClient, rootCategory.Name, rootCategory.Type)
	if err != nil {
		log.Warn(err)
	}
	for _, category := range categories {
		if category.Path == rootCategory.Name { // If path and name are the same, it's the root category
			err := DeleteCategory(APIClient, category.ID)
			if err != nil {
				log.Warn(err)
			}
			log.Debugln("Deleted root category")
			break
		}
	}
}

func TestCreateCategory(t *testing.T) {
	var err error
	createdRootCategory, err = CreateCategory(APIClient, rootCategory.Name, rootCategory.Type, "")
	assert.NilError(t, err)                                      // Create should not throw an error
	assert.Equal(t, rootCategory.Name, createdRootCategory.Name) // Name should be the same
	assert.Equal(t, rootCategory.Type, createdRootCategory.Type) // Type should be the same

	for _, child := range childCategory {
		createdChild, err := CreateCategory(APIClient, child.Name, child.Type, createdRootCategory.ID)
		assert.NilError(t, err)                                              // Create should not throw an error
		assert.Equal(t, child.Name, createdChild.Name)                       // Name should be the same
		assert.Equal(t, child.Type, createdChild.Type)                       // Type should be the same
		assert.Equal(t, rootCategory.Name+"/"+child.Name, createdChild.Path) // Path should be the same
		createdChildCategory = append(createdChildCategory, createdChild)
	}
}

func TestUpdateCategory(t *testing.T) {
	updatedCategory, err := UpdateCategory(APIClient, createdChildCategory[2].ID, "UpdatedCategoryName", createdChildCategory[1].ID)
	assert.NilError(t, err)
	assert.Equal(t, updatedCategory.Name, "UpdatedCategoryName")
	assert.Equal(t, updatedCategory.Type, createdChildCategory[2].Type)
	assert.Equal(t, rootCategory.Name+"/"+createdChildCategory[1].Name+"/UpdatedCategoryName", updatedCategory.Path) // Path should include the parent category

}

func TestDeleteCategoryWithContent(t *testing.T) {
	delErr := DeleteCategory(APIClient, createdRootCategory.ID)
	assert.Error(t, delErr, "Folder '"+rootCategory.Name+"' is not empty") // Should throw an error
}

func TestDeleteCategory(t *testing.T) {
	delErr := DeleteCategory(APIClient, createdChildCategory[2].ID)
	assert.NilError(t, delErr) // Should not throw an error
}

func TestDeleteRootCategory() {
	// Delete the category, force delete if it has content
	categories, err := GetCategoryByName(APIClient, rootCategory.Name, rootCategory.Type)
	if err != nil {
		log.Warn(err)
	}
	for _, category := range categories {
		if category.Path == rootCategory.Name { // If path and name are the same, it's the root category
			err := DeleteCategory(APIClient, category.ID)
			if err != nil {
				log.Warn(err)
			}
			log.Debugln("Deleted root category")
			break
		}
	}
}
