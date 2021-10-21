/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sammcgeown/vra-cli/pkg/cmd/cloudassembly"
	"github.com/sammcgeown/vra-cli/pkg/util/auth"
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
		if err := auth.GetConnection(&targetConfig, debug); err != nil {
			log.Fatalln(err)
		}

		response, err := cloudassembly.GetDeployments(apiClient, id, name, projectName, status)
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

func init() {
	// Get Variable
	getCmd.AddCommand(getDeploymentCmd)
	getDeploymentCmd.Flags().StringVarP(&name, "name", "n", "", "List Deployments with name")
	getDeploymentCmd.Flags().StringVarP(&projectName, "project", "p", "", "List Deployments in Project")
	getDeploymentCmd.Flags().StringVarP(&id, "id", "i", "", "List Deployments by ID")
	getDeploymentCmd.Flags().StringVarP(&status, "status", "s", "", "List Deployments by Status")
	getDeploymentCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")

}
