/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sammcgeown/vra-cli/pkg/cmd/codestream"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getEndpointCmd represents the endpoint command
var getEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Get Endpoint Configurations",
	Long:  `Get Code Stream Endpoint Configurations`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := codestream.GetEndpoint(APIClient, id, name, projectName, typename, exportPath)
		if err != nil {
			log.Infoln("Unable to get endpoints: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Warnln("No results found")
			return
		}
		if APIClient.Output == "json" {
			helpers.PrettyPrint(response)
		} else if APIClient.Output == "export" {
			for _, c := range response {
				err := codestream.ExportYaml(APIClient, c.ID, c.Name, c.Project, exportPath, "endpoints")
				if err != nil {
					log.Infoln("Endpoint", c.Name, "export failed: ", err)
				} else {
					log.Infoln("Endpoint", c.Name, "exported successfully")
				}
			}
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Project", "Type", "Description"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.Project, c.Type, c.Description})
			}
			table.Render()
		}
	},
}

// createEndpointCmd represents the endpoint create command
var createEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Create an Endpoint",
	Long: `Create an Endpoint by importing a YAML specification.
	
	Create from YAML
	  vra-cli create endpoint --importPath "/Users/sammcgeown/Desktop/endpoint.yaml"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		if importPath != "" {
			yamlFilePaths := helpers.GetFilePaths(importPath, "yaml")
			if len(yamlFilePaths) == 0 {
				log.Warnln("No YAML files were found in", importPath)
			}
			for _, yamlFilePath := range yamlFilePaths {
				yamlFileName := filepath.Base(yamlFilePath)
				err := codestream.ImportYaml(APIClient, yamlFilePath, "create", projectName, "endpoint")
				if err != nil {
					log.Warnln("Failed to import", yamlFilePath, "as Endpoint", err)
				} else {
					fmt.Println("Imported", yamlFileName, "successfully - Endpoint created.")
				}
			}
		}
	},
}

// updateEndpointCmd represents the endpoint update command
var updateEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Update an Endpoint",
	Long: `Update an Endpoint by importing the YAML specification

	Update from a YAML file
	vra-cli update endpoint --importPath "/Users/sammcgeown/vra-cli/endpoints/updated-endpoint.yaml"
	Update from a folder of YAML files
	vra-cli update endpoint --importPath "/Users/sammcgeown/vra-cli/endpoints"
	`,
	Run: func(cmd *cobra.Command, args []string) {

		if importPath != "" {
			yamlFilePaths := helpers.GetFilePaths(importPath, ".yaml")
			if len(yamlFilePaths) == 0 {
				log.Warnln("No YAML files were found in", importPath)
			}
			for _, yamlFilePath := range yamlFilePaths {
				yamlFileName := filepath.Base(yamlFilePath)
				err := codestream.ImportYaml(APIClient, yamlFilePath, "apply", "", "endpoint")
				if err != nil {
					log.Warnln("Failed to import", yamlFilePath, "as Endpoint", err)
				} else {
					fmt.Println("Imported", yamlFileName, "successfully - Endpoint updated.")
				}
			}
		}
	},
}

// deleteEndpointCmd represents the executions command
var deleteEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Delete an Endpoint",
	Long: `Delete an Endpoint with a specific Endpoint ID or Name

# Delete Endpoint by ID:
vra-cli delete endpoint --id "Endpoint ID"

# Delete Endpoint by Name:
vra-cli delete endpoint --name "Endpoint Name"

# Delete Endpoint by Project and Name:
vra-cli delete endpoint --project "My Project" --name "Endpoint Name"

# Delete all Endpoints in Project (prompts for confirmation):
vra-cli delete endpoint --project "My Project"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		// if id != "" && name != "" {
		// 	return errors.New("please specify either endpoint name or endpoint id")
		// }
		// if id == "" && name == "" {
		// 	return errors.New("please specify endpoint name or endpoint id")
		// }

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if name != "" {
			response, err := codestream.GetEndpoint(APIClient, id, name, projectName, typename, exportPath)
			if err != nil {
				log.Fatalln(err)
			}
			// return first element of map[string]
			for _, c := range response {
				id = c.ID
				break
			}
		}

		if id != "" {

			err := codestream.DeleteEndpoint(APIClient, id)
			if err != nil {
				log.Errorln("Unable to delete Endpoint: ", err)
			}
			log.Infoln("Endpoint with id " + id + " deleted")
		} else if projectName != "" {
			response, err := codestream.DeleteEndpointByProject(APIClient, projectName)
			if err != nil {
				log.Errorln("Unable to delete Endpoint: ", err)
			}
			log.Infoln(len(response), "Endpoints deleted")

		}

	},
}

func init() {
	getCmd.AddCommand(getEndpointCmd)
	getEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Get Endpoint by Name")
	getEndpointCmd.Flags().StringVarP(&id, "id", "i", "", "Get Endpoint by ID")
	getEndpointCmd.Flags().StringVarP(&projectName, "project", "p", "", "Filter Endpoint by Project")
	getEndpointCmd.Flags().StringVarP(&typename, "type", "t", "", "Filter Endpoint by Type")
	getEndpointCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	// Create
	createCmd.AddCommand(createEndpointCmd)
	createEndpointCmd.Flags().StringVarP(&importPath, "importPath", "c", "", "YAML configuration file to import")
	createEndpointCmd.Flags().StringVarP(&projectName, "project", "p", "", "Manually specify the Project in which to create the Endpoint (overrides YAML)")
	createEndpointCmd.MarkFlagRequired("importPath")
	// Update
	updateCmd.AddCommand(updateEndpointCmd)
	updateEndpointCmd.Flags().StringVarP(&importPath, "importPath", "c", "", "YAML configuration file to import")
	updateEndpointCmd.MarkFlagRequired("importPath")
	// Delete
	deleteCmd.AddCommand(deleteEndpointCmd)
	deleteEndpointCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Endpoint to delete")
	deleteEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Endpoint to delete")
	deleteEndpointCmd.Flags().StringVarP(&projectName, "project", "p", "", "Delete Endpoints by Project")

}
