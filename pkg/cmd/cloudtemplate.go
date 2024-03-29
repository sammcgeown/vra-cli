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
	"github.com/sammcgeown/vra-cli/pkg/cmd/cloudassembly"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	content string
	scope   string
	schema  bool
)

// getCloudTemplateCmd represents the Blueprint command
var getCloudTemplateCmd = &cobra.Command{
	Use:   "cloudtemplate",
	Short: "Get Cloud Templates",
	Long:  `Get Cloud Templates by ID, name or status`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := cloudassembly.GetCloudTemplate(APIClient, id, name, projectName)
		if err != nil {
			log.Fatalln("Unable to get Cloud Template(s): ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		} else if resultCount == 1 {
			if schema {
				// if inputSchema, err := getCloudTemplateInputSchema(response[0].ID); err != nil {
				// 	log.Errorln("Unable to retrieve input schema: ", err)
				// } else {
				// 	//inputs := getInputsFromSchema(inputSchema)
				// 	//helpers.PrettyPrint(inputs)
				// 	helpers.PrettyPrint(inputSchema)
				// }
			} else {
				helpers.PrettyPrint(response[0])
			}
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Status", "Valid"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.ProjectName, c.Status, strconv.FormatBool(*c.Valid)})
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
// 		if err := auth.GetConnection(&targetConfig, debug); err != nil {
// 			log.Fatalln(err)
// 		}

// 	},
// }

// createCloudTemplateCmd represents the Blueprint create command
//
// Cloud Template JSON structure:
// {
//     "projectID": "90bb3da1-8e1f-40c0-b431-0838e8ebc28d",
//     "name": "vra-cli Test",
//     "description": "Blueprint to test Packer Image builds",
//     "status": "DRAFT",
//     "content": "formatVersion: 1\ninputs: {}\nresources:\n  Cloud_Machine_CentOS7:\n    type: Cloud.Machine\n    properties:\n      image: '[Packer Test] CentOS7'\n      flavor: small\n      constraints:\n        - tag: 'env:vsphere'\n  Cloud_Machine_CentOS8:\n    type: Cloud.Machine\n    properties:\n      image: '[Packer Test] CentOS8'\n      flavor: small\n      constraints:\n        - tag: 'env:vsphere'\n  Cloud_Machine_Ubuntu1804:\n    type: Cloud.Machine\n    properties:\n      image: '[Packer Test] Ubuntu1804'\n      flavor: small\n      constraints:\n        - tag: 'env:vsphere'"
// }
//
// Cloud Template YAML structure:
//
// formatVersion: 1
// inputs: {}
// resources:
//   Cloud_Machine_CentOS7:
//     type: Cloud.Machine
//     properties:
//       image: '[Packer Test] CentOS7'
//       flavor: small
//       constraints:
//         - tag: 'env:vsphere'
//   Cloud_Machine_CentOS8:
//     type: Cloud.Machine
//     properties:
//       image: '[Packer Test] CentOS8'
//       flavor: small
//       constraints:
//         - tag: 'env:vsphere'
//   Cloud_Machine_Ubuntu1804:
//     type: Cloud.Machine
//     properties:
//       image: '[Packer Test] Ubuntu1804'
//       flavor: small
//       constraints:
//         - tag: 'env:vsphere'
var createCloudTemplateCmd = &cobra.Command{
	Use:   "cloudtemplate",
	Short: "Create a Cloud Template",
	Long: `Create a Cloud Template.

	Create from piped JSON:
	  cat test/cloudtemplate.json | vra-cli create cloud template
	Create from piped JSON, overriding project:
	  cat test/cloudtemplate.json | vra-cli create cloud template --project Test
	Create from flags:
	  vra-cli create cloudtemplate --name Test --project Development --description "My new template" --content "{formatVersion: 1, inputs: {}, resources: {}}" --scope project
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var cloudTemplateReq types.CloudTemplate
		var projectID string

		// Check if input is piped JSON
		if helpers.IsInputFromPipe() {
			if err := json.NewDecoder(os.Stdin).Decode(&cloudTemplateReq); err != nil {
				log.Warnln(err)
			}
		}
		// If project name flag is set, get the project ID and update the request
		if projectName != "" {
			log.Debugln("Project: " + projectName)
			projectObjs, pErr := cloudassembly.GetProject(APIClient, projectName, "")
			if pErr != nil {
				log.Fatalln(pErr)
			} else if len(projectObjs) == 1 {
				projectObj := projectObjs[0]
				projectID = *projectObj.ID
				log.Debugln("Project ID: " + projectID)
				cloudTemplateReq.ProjectID = projectID
			} else {
				log.Fatalln("Unable to find Project \"" + projectName + "\"")
			}
		}
		// If name flag is set, update the request
		if name != "" {
			cloudTemplateReq.Name = name
		}
		// If description flag is set, update the request
		if description != "" {
			cloudTemplateReq.Description = description
		}
		// If content flag is set, update the request
		if content != "" {
			cloudTemplateReq.Content = content
		}
		// If scope flag is set, update the request
		if scope == "org" {
			cloudTemplateReq.RequestScopeOrg = true
		} else if scope == "project" {
			cloudTemplateReq.RequestScopeOrg = false
		}
		// Create the cloud template
		cloudTemplate, err := cloudassembly.CreateCloudTemplate(APIClient, cloudTemplateReq.Name, cloudTemplateReq.Description, cloudTemplateReq.ProjectID, cloudTemplateReq.Content, cloudTemplateReq.RequestScopeOrg)
		if err != nil {
			log.Errorln("Unable to create Cloud Template(s): ", err)
		}
		if APIClient.Output == "json" {
			helpers.PrettyPrint(cloudTemplate)
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Status", "Valid"})
			table.Append([]string{cloudTemplate.ID, cloudTemplate.Name, cloudTemplate.ProjectName, cloudTemplate.Status, strconv.FormatBool(*cloudTemplate.Valid)})
			table.Render()
		}
	},
}

// deleteCloudTemplateCmd represents the delete Blueprint command
var deleteCloudTemplateCmd = &cobra.Command{
	Use:   "cloudtemplate",
	Short: "Delete a Cloud Template",
	Long: `Delete a Blueprint with a specific ID
	`,
	Run: func(cmd *cobra.Command, args []string) {

		if name != "" {
			response, err := cloudassembly.GetCloudTemplate(APIClient, id, name, projectName)
			if err != nil {
				log.Fatalln(err)
			}
			if len(response) > 1 {
				log.Warnln("There are multiple Cloud Templates matching your criteria, please use the Cloud Template ID")
				// Print result table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Project", "Status", "Valid"})
				for _, c := range response {
					table.Append([]string{c.ID, c.Name, c.ProjectName, c.Status, strconv.FormatBool(*c.Valid)})
				}
				table.Render()
			} else {
				id = response[0].ID
			}
		}
		if id != "" {
			if err := cloudassembly.DeleteCloudTemplate(APIClient, id); err != nil {
				log.Errorln("Unable to delete Cloud Template: ", err)
			} else {
				log.Infoln("Cloud Template with id " + id + " deleted")
			}
		}

	},
}

func init() {
	// Get
	getCmd.AddCommand(getCloudTemplateCmd)
	getCloudTemplateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Template to list executions for")
	getCloudTemplateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Cloud Template to list")
	getCloudTemplateCmd.Flags().StringVarP(&projectName, "project", "p", "", "List Cloud Template in project")
	getCloudTemplateCmd.Flags().StringVar(&exportPath, "exportPath", "", "Path to export objects - relative or absolute location")
	getCloudTemplateCmd.Flags().BoolVar(&schema, "schema", false, "Get the Cloud Template Input Schema")

	// Create
	createCmd.AddCommand(createCloudTemplateCmd)
	// createCloudTemplateCmd.Flags().StringVarP(&importPath, "importPath", "", "", "YAML configuration file to import")
	createCloudTemplateCmd.Flags().StringVarP(&projectName, "project", "p", "", "Project in which to create the Cloud Template (overrides piped JSON)")
	createCloudTemplateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Template (overrides piped JSON)")
	createCloudTemplateCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the Cloud Template (overrides piped JSON)")
	createCloudTemplateCmd.Flags().StringVarP(&content, "content", "c", "", "Content of the Cloud Template - YAML as a string (overrides piped JSON)")
	createCloudTemplateCmd.Flags().StringVarP(&scope, "scope", "", "", "Scope of the Cloud Template, false is project, true is any project in the organization (overrides piped JSON)")

	// // Update
	// updateCmd.AddCommand(updateCloudTemplateCmd)
	// updateCloudTemplateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Blueprint to list")
	// updateCloudTemplateCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Configuration file to import")
	// updateCloudTemplateCmd.Flags().StringVarP(&state, "state", "s", "", "Set the state of the Blueprint (ENABLED|DISABLED|RELEASED")

	// Delete
	deleteCmd.AddCommand(deleteCloudTemplateCmd)
	deleteCloudTemplateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Cloud Template to delete")
	deleteCloudTemplateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Template to delete")
	deleteCloudTemplateCmd.Flags().StringVarP(&projectName, "project", "p", "", "Project of the Cloud Template to delete")

}
