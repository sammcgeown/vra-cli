/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var cloudTemplate CloudAssemblyCloudTemplate

// getCloudTemplateCmd represents the Blueprint command
var getCloudTemplateCmd = &cobra.Command{
	Use:   "cloudtemplate",
	Short: "Get Cloud Templates",
	Long:  `Get Cloud Templates by ID, name or status`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}
		response, err := getCloudTemplates(id, name, project, exportPath)
		if err != nil {
			log.Warnln("Unable to get Cloud Template(s): ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		} else if resultCount == 1 {
			PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Status", "Valid"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.ProjectName, c.Status, strconv.FormatBool(c.Valid)})
			}
			table.Render()
		}
	},
}

// // updateCloudTemplateCmd represents the Blueprint update command
// var updateCloudTemplateCmd = &cobra.Command{
// 	Use:   "Blueprint",
// 	Short: "Update a Blueprint",
// 	Long: `Update a Blueprint
// 	Enable/Disable/Release:
// 	vra-cli update Blueprint --id d0185f04-2e87-4f3c-b6d7-ee58abba3e92 --state enabled/disabled/released
// 	Update from YAML
// 	vra-cli update Blueprint --importPath "/Users/sammcgeown/Desktop/Blueprints/SSH Exports.yaml"
// 	`,
// 	Args: func(cmd *cobra.Command, args []string) error {
// 		if state != "" {
// 			switch strings.ToUpper(state) {
// 			case "ENABLED", "DISABLED", "RELEASED":
// 				// Valid states
// 				return nil
// 			}
// 			return errors.New("--state is not valid, must be ENABLED, DISABLED or RELEASED")
// 		}
// 		return nil
// 	},
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if err := ensureTargetConnection(); err != nil {
// 			log.Fatalln(err)
// 		}

// 	},
// }

// createCloudTemplateCmd represents the Blueprint create command
var createCloudTemplateCmd = &cobra.Command{
	Use:   "cloudtemplate",
	Short: "Create a Cloud Template",
	Long: `Create a Cloud Template by importing a JSON specification.

	Create from YAML
	  vra-cli create Blueprint --importPath "/Users/sammcgeown/Desktop/Blueprints/SSH Exports.yaml"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		if isInputFromPipe() {
			if err := json.NewDecoder(os.Stdin).Decode(&cloudTemplate); err != nil {
				log.Warnln(err)
			}
		}

		log.Infoln(cloudTemplate.Name)

		// yamlFilePaths := getYamlFilePaths(importPath)
		// if len(yamlFilePaths) == 0 {
		// 	log.Warnln("No YAML files were found in", importPath)
		// }
		// for _, yamlFilePath := range yamlFilePaths {
		// 	yamlFileName := filepath.Base(yamlFilePath)
		// 	err := importYaml(yamlFilePath, "create", project, "Blueprint")
		// 	if err != nil {
		// 		log.Warnln("Failed to import", yamlFilePath, "as Blueprint", err)
		// 	} else {
		// 		fmt.Println("Imported", yamlFileName, "successfully - Blueprint created.")
		// 	}
		// }
	},
}

// deleteCloudTemplateCmd represents the delete Blueprint command
var deleteCloudTemplateCmd = &cobra.Command{
	Use:   "cloudtemplate",
	Short: "Delete a Cloud Template",
	Long: `Delete a Blueprint with a specific ID
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		if name != "" {
			response, err := getEndpoint(id, name, project, typename, exportPath)
			if err != nil {
				log.Fatalln(err)
			}
			id = response[0].ID
		}

		if err := deleteCloudTemplate(id); err != nil {
			log.Errorln("Unable to delete Cloud Template: ", err)
		} else {
			log.Infoln("Cloud Template with id " + id + " deleted")
		}

	},
}

func init() {
	// Get
	getCmd.AddCommand(getCloudTemplateCmd)
	getCloudTemplateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Template to list executions for")
	getCloudTemplateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Cloud Template to list")
	getCloudTemplateCmd.Flags().StringVarP(&project, "project", "p", "", "List Cloud Template in project")
	getCloudTemplateCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")

	// Create
	createCmd.AddCommand(createCloudTemplateCmd)
	createCloudTemplateCmd.Flags().StringVarP(&importPath, "importPath", "", "", "YAML configuration file to import")
	createCloudTemplateCmd.Flags().StringVarP(&project, "project", "p", "", "Manually specify the Project in which to create the Blueprint (overrides YAML)")

	// // Update
	// updateCmd.AddCommand(updateCloudTemplateCmd)
	// updateCloudTemplateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Blueprint to list")
	// updateCloudTemplateCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Configuration file to import")
	// updateCloudTemplateCmd.Flags().StringVarP(&state, "state", "s", "", "Set the state of the Blueprint (ENABLED|DISABLED|RELEASED")

	// Delete
	deleteCmd.AddCommand(deleteCloudTemplateCmd)
	deleteCloudTemplateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Cloud Template to delete")

}
