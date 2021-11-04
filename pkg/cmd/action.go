/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/sammcgeown/vra-cli/pkg/cmd/orchestrator"

	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getActionCmd represents the actions command
var getActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Get Orchestrator Actions",
	Long: `Get Orchestrator Actions by ID, Name, Project and Status

# Get only failed actions:
vra-cli get action
# Get an action by ID:
vra-cli get action --id bb3f6aff-311a-45fe-8081-5845a529068d`,
	Run: func(cmd *cobra.Command, args []string) {

		response, err := orchestrator.GetAction(APIClient, id, category, name)
		if err != nil {
			log.Errorln("Unable to get actions: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		} else {
			if APIClient.Output == "table" {
				// Print result table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Version", "Description", "Module"})
				for _, c := range response {
					table.Append([]string{c.ID, c.Name, c.Version, c.Description, c.Module})
				}
				table.Render()
			} else if APIClient.Output == "export" {
				// Export the Worfklow
				for _, action := range response {
					err := orchestrator.ExportAction(APIClient, action.ID, action.Name, exportPath)
					if err != nil {
						log.Warnln("Unable to export action: ", err)
					} else {
						log.Infoln("Action", action.Name, "exported")
					}
				}

			} else {
				helpers.PrettyPrint(response)
			}
		}

	},
}

// delActionCmd represents the delete actions command
var delActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Delete an Action",
	Long:  `Delete an Action with a specific Action ID`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := orchestrator.DeleteAction(APIClient, id)
		if err != nil {
			log.Errorln("Unable to delete action: ", err)
		} else {
			log.Infoln("Action with ID " + id + " deleted")
		}
	},
}

// createActionCmd represents the actions command
var createActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Create a Action",
	Long:  `Create a Action`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, path := range helpers.GetFilePaths(importPath, ".action") {
			log.Infoln("Importing action:", path)
			err := orchestrator.ImportAction(APIClient, path, category)
			if err != nil {
				log.Errorln("Unable to import action: ", err)
			} else {
				action, err := orchestrator.GetAction(APIClient, "", category, name)
				if err != nil {
					log.Errorln("Action imported OK, but I'm unable to get imported action details: ", err)
				}
				log.Infoln("Action imported:", action[0].Name, "with ID:", action[0].ID)
			}
		}

	},
}

func init() {
	// Get
	getCmd.AddCommand(getActionCmd)
	getActionCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Action")
	getActionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Actions to list")
	getActionCmd.Flags().StringVarP(&category, "category", "c", "", "Filter Actions by Category")
	getActionCmd.Flags().StringVar(&exportPath, "exportPath", "", "Path to export the file")

	// Delete
	deleteCmd.AddCommand(delActionCmd)
	delActionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Action to delete")
	delActionCmd.MarkFlagRequired("id")

	// Create
	createCmd.AddCommand(createActionCmd)
	createActionCmd.Flags().StringVarP(&category, "category", "c", "", "Category to import")
	createActionCmd.Flags().StringVar(&importPath, "importPath", "", "Path to the zip file, or folder containing zip files, to import")
}
