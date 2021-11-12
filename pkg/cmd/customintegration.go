/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"
	"strings"

	"github.com/sammcgeown/vra-cli/pkg/cmd/codestream"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	yaml         string
	versionState string
	versionName  string
)

// getCustomIntegrationCmd represents the customintegration command
var getCustomIntegrationCmd = &cobra.Command{
	Use:   "customintegration",
	Short: "Get Custom Integrations",
	Long: `Get Code Stream Custom Integrations by name, project or by id - e.g:

Get by ID
	vra-cli get customintegration --id 6b7936d3-a19d-4298-897a-65e9dc6620c8
	
Get by Name
	vra-cli get customintegration --name my-customintegration
	
Get by Project
	vra-cli get customintegration --project production`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := codestream.GetCustomIntegration(APIClient, id, name)
		if err != nil {
			log.Errorln("Unable to get Code Stream CustomIntegrations: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		}

		if APIClient.Output == "json" {
			helpers.PrettyPrint(response)
		} else if APIClient.Output == "export" {
			log.Warnln("Exporting Custom Integrations is not supported yet")
			// for _, c := range response {
			// 	exportCustomIntegration(c, exportFile)
			// }
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Description", "Status", "Current Version", "Versions"})
			for _, c := range response {
				versions, err := codestream.GetCustomIntegrationVersions(APIClient, c.ID)
				if err != nil {
					log.Errorln("Unable to get Code Stream CustomIntegration Versions: ", err)
				}
				table.Append([]string{c.ID, c.Name, c.Description, c.Status, c.Version, strings.Join(versions, ", ")})
			}
			table.Render()
		}
	},
}

// createCustomIntegrationCmd represents the customintegration command
var createCustomIntegrationCmd = &cobra.Command{
	Use:   "customintegration",
	Short: "Create a new Custom Integration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		createResponse, err := codestream.CreateCustomIntegration(APIClient, name, description, yaml)
		if err != nil {
			log.Errorln("Unable to create Custom Integration:", err)
		}
		helpers.PrettyPrint(createResponse)

	},
}

// updateCustomIntegrationCmd represents the customintegration command
var updateCustomIntegrationCmd = &cobra.Command{
	Use:   "customintegration",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := codestream.UpdateCustomIntegration(APIClient, id, description, yaml, versionName, versionState)
		if err != nil {
			log.Infoln("Unable to update Custom Integration: ", err)
		}
		log.Infoln("Updated Custom Integration")

	},
}

// deleteCustomIntegrationCmd represents the executions command
var deleteCustomIntegrationCmd = &cobra.Command{
	Use:   "customintegration",
	Short: "Delete Custom Integration by ID",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := codestream.DeleteCustomIntegration(APIClient, id)
		if err != nil {
			log.Infoln("Unable to delete customintegration: ", err)
		}
		log.Infoln("CustomIntegration deleted")
	},
}

func init() {
	// Get CustomIntegration
	getCmd.AddCommand(getCustomIntegrationCmd)
	getCustomIntegrationCmd.Flags().StringVarP(&name, "name", "n", "", "List customintegration with name")
	getCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "List customintegrations by id")
	getCustomIntegrationCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	// Create CustomIntegration
	createCmd.AddCommand(createCustomIntegrationCmd)
	createCustomIntegrationCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the customintegration to create")
	createCustomIntegrationCmd.Flags().StringVarP(&description, "description", "d", "", "The description of the customintegration to create")
	createCustomIntegrationCmd.Flags().StringVar(&yaml, "yaml", "", "Custom Integration YAML")
	// Update CustomIntegration
	updateCmd.AddCommand(updateCustomIntegrationCmd)
	updateCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the customintegration to update")
	updateCustomIntegrationCmd.Flags().StringVarP(&description, "description", "d", "", "Update the description of the customintegration")
	updateCustomIntegrationCmd.Flags().StringVar(&yaml, "yaml", "", "Custom Integration YAML")
	updateCustomIntegrationCmd.Flags().StringVar(&versionName, "versionName", "", "Create a new version using this name")
	updateCustomIntegrationCmd.Flags().StringVar(&versionState, "versionState", "", "Update the version state (delete|release|deprecate|restore|withdraw)")
	updateCustomIntegrationCmd.MarkFlagRequired("id")
	// Delete CustomIntegration
	deleteCmd.AddCommand(deleteCustomIntegrationCmd)
	deleteCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "Delete customintegration by id")
	deleteCustomIntegrationCmd.MarkFlagRequired("id")
}
