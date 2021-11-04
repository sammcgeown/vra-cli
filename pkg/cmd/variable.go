/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/sammcgeown/vra-cli/pkg/cmd/codestream"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"

	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// GetVariableCmd represents the variable command
var GetVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Get Variables",
	Long: `Get Code Stream Variables by name, project or by id - e.g:

# Get Variable by ID
vra-cli get variable --id 6b7936d3-a19d-4298-897a-65e9dc6620c8
	
# Get Variable by Name
vra-cli get variable --name my-variable
	
# Get Variable by Project
vra-cli get variable --project production`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := codestream.GetVariable(APIClient, id, name, projectName, exportPath)
		if err != nil {
			log.Fatalln("Unable to get Code Stream Variables: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Fatalln("No results found")
		}
		if APIClient.Output == "json" {
			helpers.PrettyPrint(response[0])
		} else if APIClient.Output == "export" {
			for _, c := range response {
				codestream.ExportVariable(c, exportPath)
			}
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Type", "Description", "Value"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.Project, c.Type, c.Description, c.Value})
			}
			table.Render()
		}
	},
}

// GetVariableCmd represents the variable command
var createVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Create a Variable",
	Long:  `Create a Variable`,
	Run: func(cmd *cobra.Command, args []string) {

		if importPath != "" { // If we are importing a file
			variables := codestream.ImportVariables(importPath)
			for _, value := range variables {
				if projectName != "" { // If the project is specified update the object
					value.Project = projectName
				}
				createResponse, err := codestream.CreateVariable(APIClient, value.Name, value.Description, value.Type, value.Project, value.Value)
				if err != nil {
					log.Warnln("Unable to create Code Stream Variable: ", err)
				} else {
					log.Infoln("Created variable", createResponse.Name, "in", createResponse.Project)
				}
			}
		} else {
			createResponse, err := codestream.CreateVariable(APIClient, name, description, typename, projectName, value)
			if err != nil {
				log.Errorln("Unable to create Code Stream Variable: ", err)
			} else {
				if APIClient.Output == "json" {
					helpers.PrettyPrint(createResponse)
				} else {
					// Print result table
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Id", "Name", "Project", "Type", "Description", "Value"})
					table.Append([]string{createResponse.ID, createResponse.Name, createResponse.Project, createResponse.Type, createResponse.Description, createResponse.Value})
					table.Render()
				}
			}
		}
	},
}

// updateVariableCmd represents the variable command
var updateVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Update a Variable",
	Long:  `Update a Variable`,
	Run: func(cmd *cobra.Command, args []string) {

		if importPath != "" { // If we are importing a file
			variables := codestream.ImportVariables(importPath)
			for _, value := range variables {
				exisitingVariable, err := codestream.GetVariable(APIClient, "", value.Name, value.Project, "")
				if err != nil {
					log.Errorln("Update failed - unable to find existing Code Stream Variable", value.Name, "in", value.Project)
				} else {
					_, err := codestream.UpdateVariable(APIClient, exisitingVariable[0].ID, value.Name, value.Description, value.Type, value.Value)
					if err != nil {
						log.Errorln("Unable to update Code Stream Variable: ", err)
					} else {
						log.Infoln("Updated variable", value.Name)
					}
				}
			}
		} else { // Else we are updating using flags
			updateResponse, err := codestream.UpdateVariable(APIClient, id, name, description, typename, value)
			if err != nil {
				log.Errorln("Unable to update Code Stream Variable: ", err)
			}
			log.Infoln("Variable updated")
			if APIClient.Output == "json" {
				helpers.PrettyPrint(updateResponse)
			} else {
				// Print result table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Project", "Type", "Description", "Value"})
				table.Append([]string{updateResponse.ID, updateResponse.Name, updateResponse.Project, updateResponse.Type, updateResponse.Description, updateResponse.Value})
				table.Render()
			}
		}
	},
}

// deleteVariableCmd represents the executions command
var deleteVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Delete a Variable",
	Long: `Delete a Variable

# Delete Variable by ID
vra-cli delete variable --id "variable ID"

# Delete Variable by Name
vra-cli delete variable --name "My Variable"

# Delete Variable by Name and Project
vra-cli delete variable --name "My Variable" --project "My Project"

# Delete all Variables in Project
vra-cli delete variable --project "My Project"
	`,
	Run: func(cmd *cobra.Command, args []string) {

		if id != "" {
			_, err := codestream.DeleteVariable(APIClient, id)
			if err != nil {
				log.Errorln("Unable to delete variable: ", err)
			} else {
				log.Infoln("Variable " + id + " deleted")
			}
		} else if projectName != "" {
			response, err := codestream.DeleteVariableByProject(APIClient, projectName)
			if err != nil {
				log.Errorln("Delete Variables in "+projectName+" failed:", err)
			} else {
				log.Infoln(len(response), "Variables deleted")
			}
		}
	},
}

func init() {
	// Get Variable
	getCmd.AddCommand(GetVariableCmd)
	GetVariableCmd.Flags().StringVarP(&name, "name", "n", "", "List variable with name")
	GetVariableCmd.Flags().StringVarP(&projectName, "project", "p", "", "List variables in project")
	GetVariableCmd.Flags().StringVarP(&id, "id", "i", "", "List variables by id")
	GetVariableCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	// Create Variable
	createCmd.AddCommand(createVariableCmd)
	createVariableCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the variable to create")
	createVariableCmd.Flags().StringVarP(&typename, "type", "t", "", "The type of the variable to create REGULAR|SECRET|RESTRICTED")
	createVariableCmd.Flags().StringVarP(&projectName, "project", "p", "", "The project in which to create the variable")
	createVariableCmd.Flags().StringVar(&value, "value", "", "The value of the variable to create")
	createVariableCmd.Flags().StringVarP(&description, "description", "d", "", "The description of the variable to create")
	createVariableCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Path to a YAML file with the variables to import")

	// Update Variable
	updateCmd.AddCommand(updateVariableCmd)
	updateVariableCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the variable to update")
	updateVariableCmd.Flags().StringVarP(&name, "name", "n", "", "Update the name of the variable")
	updateVariableCmd.Flags().StringVarP(&typename, "type", "t", "", "Update the type of the variable REGULAR|SECRET|RESTRICTED")
	updateVariableCmd.Flags().StringVar(&value, "value", "", "Update the value of the variable ")
	updateVariableCmd.Flags().StringVarP(&description, "description", "d", "", "Update the description of the variable")
	updateVariableCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Path to a YAML file with the variables to import")
	//updateVariableCmd.MarkFlagRequired("id")

	// Delete Variable
	deleteCmd.AddCommand(deleteVariableCmd)
	deleteVariableCmd.Flags().StringVarP(&id, "id", "i", "", "Delete variable by id")
	deleteVariableCmd.Flags().StringVarP(&projectName, "project", "p", "", "The project in which to delete the variable, or delete all variables in project")
}
