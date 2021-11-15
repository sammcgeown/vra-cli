/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sammcgeown/vra-cli/pkg/cmd/orchestrator"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCategory     bool
	categoryType     string
	parentCategoryID string
)

// getCategoryCmd represents the workflows command
var getCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Get Orchestrator Category",
	Long:  `Get Orchestrator Categories by ID, Name`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var categories []*types.WsCategory
		if name != "" {
			categories, err = orchestrator.GetCategoryByName(APIClient, name, categoryType)
		} else if id != "" {
			category, cErr := orchestrator.GetCategoryByID(APIClient, id)
			if cErr != nil {
				log.Fatalln(cErr)
			}
			categories = append(categories, category)
		} else {
			categories, err = orchestrator.GetCategory(APIClient, rootCategory, categoryType)
		}
		if err != nil {
			log.Errorln("Unable to get workflows: ", err)
		}
		var resultCount = len(categories)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		} else {
			if APIClient.Output == "table" {
				// Print result table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Type", "Category Path"})
				for _, c := range categories {
					table.Append([]string{c.ID, c.Name, c.Type, c.Path})
				}
				table.Render()
			} else if APIClient.Output == "export" {
				log.Warnln("Export mode is not supported for this command")
				// Export the Worfklow
				// for _, workflow := range categories {
				// 	err := orchestrator.ExportWorkflow(APIClient, workflow.ID, workflow.Name, category)
				// 	if err != nil {
				// 		log.Warnln("Unable to export workflow: ", err)
				// 	} else {
				// 		log.Infoln("Workflow", workflow.Name, "exported")
				// 	}
				// }

			} else {
				helpers.PrettyPrint(categories)
			}
		}

	},
}

// delCategoryCmd - represents the delete category command
var delCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Delete a Category",
	Long:  `Delete a Category with a specific ID`,
	Run: func(cmd *cobra.Command, args []string) {
		if name != "" {
			categories, err := orchestrator.GetCategoryByName(APIClient, name, categoryType)
			if err != nil {
				log.Fatalln("Unable to find category by name", err)
			}
			id = categories[0].ID
		}

		if id != "" {
			err := orchestrator.DeleteCategory(APIClient, id)
			if err != nil {
				log.Errorln("Unable to delete Category: ", err)
			} else {
				log.Infoln("Category deleted")
			}
		} else {
			log.Fatalln("Unable to delete Category: ID not found")
		}
	},
}

// createCategoryCmd - Create a Category
var createCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Create a Category",
	Long:  `Create a Category`,
	Run: func(cmd *cobra.Command, args []string) {
		newCategory, err := orchestrator.CreateCategory(APIClient, name, categoryType, parentCategoryID)
		if err != nil {
			log.Errorln("Unable to create category: ", err)
		}
		if APIClient.Output == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Type", "Category Path"})
			table.Append([]string{newCategory.ID, newCategory.Name, newCategory.Type, newCategory.Path})
			table.Render()
		} else if APIClient.Output == "json" {
			helpers.PrettyPrint(newCategory)
		} else {
			helpers.PrettyPrint(newCategory)
		}
	},
}

// updateCategoryCmd - Create a Category
var updateCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Update a Category",
	Long:  `Update a Category`,
	Run: func(cmd *cobra.Command, args []string) {
		updatedCategory, err := orchestrator.UpdateCategory(APIClient, id, name, parentCategoryID)
		if err != nil {
			log.Errorln("Unable to update Category: ", err)
		} else {
			if APIClient.Output == "table" {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Type", "Category Path"})
				table.Append([]string{updatedCategory.ID, updatedCategory.Name, updatedCategory.Type, updatedCategory.Path})
				table.Render()
			} else if APIClient.Output == "json" {
				helpers.PrettyPrint(updatedCategory)
			} else {
				helpers.PrettyPrint(updatedCategory)
			}
		}
	},
}

func init() {
	// Get
	getCmd.AddCommand(getCategoryCmd)
	getCategoryCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Category")
	getCategoryCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Category")
	getCategoryCmd.Flags().StringVarP(&categoryType, "type", "t", "", "Type of Category")
	getCategoryCmd.Flags().BoolVarP(&rootCategory, "root", "", false, "List root Categories only")
	// Delete
	deleteCmd.AddCommand(delCategoryCmd)
	delCategoryCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Category to delete")
	delCategoryCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Category to delete")
	// Create
	createCmd.AddCommand(createCategoryCmd)
	createCategoryCmd.Flags().StringVarP(&name, "name", "n", "", "Category Name")
	createCategoryCmd.Flags().StringVarP(&categoryType, "type", "t", "", "Category Type  ['ResourceElementCategory', 'ConfigurationElementCategory', 'WorkflowCategory', 'PolicyTemplateCategory', 'ScriptModuleCategory']")
	createCategoryCmd.Flags().StringVar(&parentCategoryID, "parent", "", "Category Category ID")
	createCategoryCmd.MarkFlagRequired("name")
	createCategoryCmd.MarkFlagRequired("type")
	// Update
	updateCmd.AddCommand(updateCategoryCmd)
	updateCategoryCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Category")
	updateCategoryCmd.Flags().StringVarP(&name, "name", "n", "", "Category Name")
	updateCategoryCmd.Flags().StringVar(&parentCategoryID, "parent", "", "Category Category ID")
}
