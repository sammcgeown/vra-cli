/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"
	"strings"

	"github.com/sammcgeown/vra-cli/pkg/cmd/orchestrator"

	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getWorkflowCmd represents the workflows command
var getWorkflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Get Orchestrator Workflows",
	Long: `Get Orchestrator Workflows by ID, Name, Project and Status

# Get only failed workflows:
vra-cli get workflow --status FAILED
# Get an workflow by ID:
vra-cli get workflow --id bb3f6aff-311a-45fe-8081-5845a529068d
# Get Failed workflows in Project "Field Demo" with the name "Learn Code Stream"
vra-cli get workflow --status FAILED --project "Field Demo" --name "Learn Code Stream"`,
	Run: func(cmd *cobra.Command, args []string) {
		// if err := auth.GetConnection(&targetConfig, debug); err != nil {
		// 	log.Fatalln(err)
		// }

		response, err := orchestrator.GetWorkflow(APIClient, id, category, name)
		if err != nil {
			log.Errorln("Unable to get workflows: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		} else {
			if APIClient.Output == "table" {
				// Print result table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Version", "Description", "Category"})
				for _, c := range response {
					category, _ := orchestrator.GetCategoryByID(APIClient, c.CategoryID)
					table.Append([]string{c.ID, c.Name, c.Version, c.Description, category.Path})
				}
				table.Render()
			} else if APIClient.Output == "export" {
				// Export the Worfklow
				for _, workflow := range response {
					err := orchestrator.ExportWorkflow(APIClient, workflow.ID, workflow.Name, category)
					if err != nil {
						log.Warnln("Unable to export workflow: ", err)
					} else {
						log.Infoln("Workflow", workflow.Name, "exported")
					}
				}

			} else {
				helpers.PrettyPrint(response)
			}
		}

	},
}

// delWorkflowCmd represents the delete workflows command
var delWorkflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Delete an Workflow",
	Long:  `Delete an Workflow with a specific Workflow ID`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := orchestrator.DeleteWorkflow(APIClient, id)
		if err != nil {
			log.Errorln("Unable to delete workflow: ", err)
		} else {
			log.Infoln("Workflow with ID " + id + " deleted")
		}
	},
}

// createWorkflowCmd represents the workflows command
var createWorkflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Create a Workflow",
	Long:  `Create a Workflow`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the category ID
		var CategoryID string
		categoryName := (strings.Split(category, "/"))[len(strings.Split(category, "/"))-1]
		categories, _ := orchestrator.GetCategoryByName(APIClient, categoryName, "WorkflowCategory")
		if len(categories) == 0 {
			log.Fatalln("Unable to find category:", categoryName)
		} else if len(categories) == 1 {
			// Only one category found
			log.Debugln("Category found:", categories[0].Name, categories[0].ID)
			CategoryID = categories[0].ID
		} else {
			for _, matchedCategory := range categories {
				if matchedCategory.Path == category {
					log.Debugln("Category ID:", matchedCategory.ID)
					CategoryID = matchedCategory.ID
					break
				}
			}
			if CategoryID == "" {
				log.Fatalln("Multiple categories found, try using a more specific category - e.g.: path/to/category")
			}
		}
		for _, path := range helpers.GetFilePaths(importPath, ".zip") {
			log.Infoln("Importing workflow:", path)
			err := orchestrator.ImportWorkflow(APIClient, path, CategoryID)
			if err != nil {
				log.Errorln("Unable to import workflow: ", err)
			} else {
				workflow, err := orchestrator.GetWorkflow(APIClient, "", categoryName, name)
				if err != nil {
					log.Errorln("Workflow imported OK, but I'm unable to get imported workflow details: ", err)
				}
				log.Infoln("Workflow imported:", workflow[0].Name, "with ID:", workflow[0].ID)
			}
		}

	},
}

func init() {
	// Get
	getCmd.AddCommand(getWorkflowCmd)
	getWorkflowCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Workflow")
	getWorkflowCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Workflows to list")
	getWorkflowCmd.Flags().StringVarP(&category, "category", "c", "", "Filter Workflows by Category")
	// Delete
	deleteCmd.AddCommand(delWorkflowCmd)
	delWorkflowCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Workflow to delete")
	delWorkflowCmd.MarkFlagRequired("id")

	// Create
	createCmd.AddCommand(createWorkflowCmd)
	createWorkflowCmd.Flags().StringVarP(&category, "category", "c", "", "Category to import")
	createWorkflowCmd.Flags().StringVar(&importPath, "importPath", "", "Path to the zip file, or folder containing zip files, to import")
}
