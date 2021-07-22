/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cloudaccounttype string
)

// getCloudAccountCmd represents the Blueprint command
var getCloudAccountCmd = &cobra.Command{
	Use:   "cloudaccount",
	Short: "Get Cloud Accounts",
	Long:  `Get Cloud Accounts by ID, name or status`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}
		cloudAccounts, err := getCloudAccounts(id, name, cloudaccounttype)
		if err != nil {
			log.Fatalln(err)
		}

		if len(cloudAccounts) == 0 {
			log.Warnln("No Cloud Accounts found")
		} else if len(cloudAccounts) == 1 {
			PrettyPrint(cloudAccounts[0])
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Description", "Type"})
			for _, c := range cloudAccounts {
				table.Append([]string{*c.ID, c.Name, c.Description, *c.CloudAccountType})
			}
			table.Render()
		}
	},
}

// // updateCloudAccountCmd represents the Blueprint update command
// var updateCloudAccountCmd = &cobra.Command{
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

// createCloudAccountCmd represents the Blueprint create command
// var createCloudAccountCmd = &cobra.Command{
// 	Use:   "cloudtemplate",
// 	Short: "Create a Cloud Template",
// 	Long: `Create a Cloud Template.`,
// 	Run: func(cmd *cobra.Command, args []string) {

// 	},
// }

// deleteCloudAccountCmd represents the delete Blueprint command
// var deleteCloudAccountCmd = &cobra.Command{
// 	Use:   "cloudtemplate",
// 	Short: "Delete a Cloud Template",
// 	Long: `Delete a Blueprint with a specific ID
// 	`,
// 	Run: func(cmd *cobra.Command, args []string) {

// 	},
// }

func init() {
	// Get
	getCmd.AddCommand(getCloudAccountCmd)
	getCloudAccountCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Account")
	getCloudAccountCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Cloud Account")
	getCloudAccountCmd.Flags().StringVarP(&cloudaccounttype, "type", "t", "", "List by Type of the Cloud Account")

	// Create
	// createCmd.AddCommand(createCloudAccountCmd)
	// // createCloudAccountCmd.Flags().StringVarP(&importPath, "importPath", "", "", "YAML configuration file to import")
	// createCloudAccountCmd.Flags().StringVarP(&project, "project", "p", "", "Project in which to create the Cloud Template (overrides piped JSON)")
	// createCloudAccountCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Template (overrides piped JSON)")
	// createCloudAccountCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the Cloud Template (overrides piped JSON)")
	// createCloudAccountCmd.Flags().StringVarP(&content, "content", "c", "", "Content of the Cloud Template - YAML as a string (overrides piped JSON)")
	// createCloudAccountCmd.Flags().StringVarP(&scope, "scope", "", "", "Scope of the Cloud Template, false is project, true is any project in the organization (overrides piped JSON)")

	// // Update
	// updateCmd.AddCommand(updateCloudAccountCmd)
	// updateCloudAccountCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Blueprint to list")
	// updateCloudAccountCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Configuration file to import")
	// updateCloudAccountCmd.Flags().StringVarP(&state, "state", "s", "", "Set the state of the Blueprint (ENABLED|DISABLED|RELEASED")

	// Delete
	// deleteCmd.AddCommand(deleteCloudAccountCmd)
	// deleteCloudAccountCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Cloud Template to delete")
	// deleteCloudAccountCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Template to delete")
	// deleteCloudAccountCmd.Flags().StringVarP(&project, "project", "p", "", "Project of the Cloud Template to delete")

}
