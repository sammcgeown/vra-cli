/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sammcgeown/vra-cli/pkg/cmd/cloudassembly"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// getDeploymentCmd represents the variable command
var getDeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Get Deployments",
	Long:  `Get Deployments`,
	Run: func(cmd *cobra.Command, args []string) {
		// if err := auth.GetConnection(&targetConfig, debug); err != nil {
		// 	log.Fatalln(err)
		// }

		response, err := cloudassembly.GetDeployments(APIClient, id, name, projectName, status)
		if err != nil {
			log.Fatalln("Unable to get Deployments: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Warnln("No results found")
		} else if resultCount == 1 {
			// Print the single result
			if exportPath != "" {
				//variable.ExportVariable(response[0], exportPath)
			}
			helpers.PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Description", "Owner", "Status"})
			for _, c := range response {
				table.Append([]string{c.ID.String(), *c.Name, c.Project.Name, c.Description, c.OwnedBy, c.Status})
			}
			table.Render()
		}
	},
}

// deleteDeploymentCmd represents the delete Deployment command
var deleteDeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Delete a Deployment",
	Long: `Delete a Deployment with a specific ID

Delete a Deployment by ID:
  vra-cli delete deployment --id <Deployment ID>`,
	Run: func(cmd *cobra.Command, args []string) {
		// if err := auth.GetConnection(&targetConfig, debug); err != nil {
		// 	log.Fatalln(err)
		// }
		deployment, err := cloudassembly.GetDeployments(APIClient, id, "", "", "")
		if err != nil {
			log.Debug(err) // There was an error getting the Deployment
		}

		if len(deployment) == 0 {
			// No error was throw, but there was no Deployment
			log.Fatalln("No Deployment matching the request was found")
		} else if len(deployment) > 1 {
			// There was more than one Deployment
			log.Fatalln("More than one Deployment matching the request was found")
		} else {
			// There was only one Deployment
			if err := cloudassembly.DeleteDeployment(APIClient, (deployment[0].ID).String()); err != nil {
				log.Fatalln(err) // There was an error deleting the Deployment
			} else {
				log.Infoln("Deployment deleted successfully")
			}
		}

	},
}

func init() {
	// Get Deployment
	getCmd.AddCommand(getDeploymentCmd)
	getDeploymentCmd.Flags().StringVarP(&name, "name", "n", "", "List Deployments with name")
	getDeploymentCmd.Flags().StringVarP(&projectName, "project", "p", "", "List Deployments in Project")
	getDeploymentCmd.Flags().StringVarP(&id, "id", "i", "", "List Deployments by ID")
	getDeploymentCmd.Flags().StringVarP(&status, "status", "s", "", "List Deployments by Status")
	getDeploymentCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")

	// Delete Deployment
	deleteCmd.AddCommand(deleteDeploymentCmd)
	deleteDeploymentCmd.Flags().StringVarP(&id, "id", "i", "", "Delete Deployment by ID")

}
