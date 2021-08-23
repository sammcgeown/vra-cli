/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sammcgeown/vra-cli/pkg/cmd/variable"
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

		response, err := getDeployments(id)
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
				variable.ExportVariable(response[0], exportPath)
			}
			helpers.PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Description"})
			for _, c := range response {
				table.Append([]string{c.Id, c.Name, c.ProjectId, c.Description})
			}
			table.Render()
		}
	},
}

func init() {
	// Get Variable
	getCmd.AddCommand(getDeploymentCmd)
	getDeploymentCmd.Flags().StringVarP(&name, "name", "n", "", "List variable with name")
	getDeploymentCmd.Flags().StringVarP(&projectName, "project", "p", "", "List variables in project")
	getDeploymentCmd.Flags().StringVarP(&id, "id", "i", "", "List variables by id")
	getDeploymentCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")

}
