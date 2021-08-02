/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	deploymentName   string
	deploymentReason string
)

// getCatalogItemCmd represents the CatalogItem command
var getCatalogItemCmd = &cobra.Command{
	Use:   "catalogitem",
	Short: "Get Catalog Items",
	Long:  `Get Catalog Items`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		response, err := getCatalogItems(id, name, projectName)
		if err != nil {
			log.Infoln("Unable to get CatalogItems: ", err)
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
			table.SetHeader([]string{"ID", "Name", "Type", "Project"})
			for _, c := range response {
				var projectList []string
				for _, Project := range c.Projects {
					projectList = append(projectList, Project.Name)
				}
				table.Append([]string{c.Id, c.Name, c.Type.Name, strings.Join(projectList, ", ")})
			}
			table.Render()
		}

	},
}

// createCatalogItemCmd represents the CatalogItem create command
var createCatalogItemCmd = &cobra.Command{
	Use:   "catalogitem",
	Short: "Create a Catalog Item request",
	Long: `Create a CatalogItem request

# Create a request using Catalog Item ID (prompts for inputs)
vra-cli create catalogitem --id 69787c80-b5d8-3d03-8ec0-a0fe67edc9e2 --project "Field Demo" --deploymentName "My Deployment"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}
		requestContent := CatalogItemRequest{}

		if isInputFromPipe() {
			if err := json.NewDecoder(os.Stdin).Decode(&requestContent); err != nil {
				log.Warnln(err)
			}
		} else {

			requestContent.DeploymentName = deploymentName
			requestContent.Reason = fmt.Sprint("[vra-cli]", deploymentReason)

			targetProject, pErr := getProject("", projectName)
			if pErr != nil {
				log.Fatalln(pErr)
			} else {
				requestContent.ProjectId = *targetProject[0].ID
				log.Debugln("Found Project ID:", requestContent.ProjectId)
			}

			catalogItems, cErr := getCatalogItems(id, name, projectName)
			if cErr != nil {
				log.Fatalln(cErr)
			} else {
				if len(catalogItems) == 1 {
					log.Debugln("Found Catalog Item ID:", catalogItems[0].Id)
					requestContent.Inputs = getCatalogItemInputs(catalogItems[0].Schema.Properties)
				} else {
					log.Errorln(len(catalogItems), "Catalog Items found")
				}
				log.Debugln(requestContent)
			}
		}

		requestResponse, rErr := createCatalogItemRequest(id, requestContent)
		if rErr != nil {
			log.Fatalln(rErr)
		} else {
			log.Infoln("Catalog Item request created successfully", requestResponse.DeploymentId)
			log.Infoln("Use vra-cli get deployment --id", requestResponse.DeploymentId, "to view the deployment status")
		}

	},
}

// // updateCatalogItemCmd represents the CatalogItem update command
// var updateCatalogItemCmd = &cobra.Command{
// 	Use:   "CatalogItem",
// 	Short: "Update an CatalogItem",
// 	Long: `Update an CatalogItem by importing the YAML specification

// 	Update from a YAML file
// 	vra-cli update CatalogItem --importPath "/Users/sammcgeown/vra-cli/CatalogItems/updated-CatalogItem.yaml"
// 	Update from a folder of YAML files
// 	vra-cli update CatalogItem --importPath "/Users/sammcgeown/vra-cli/CatalogItems"
// 	`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if err := ensureTargetConnection(); err != nil {
// 			log.Fatalln(err)
// 		}

// 		if importPath != "" {
// 			yamlFilePaths := getYamlFilePaths(importPath)
// 			if len(yamlFilePaths) == 0 {
// 				log.Warnln("No YAML files were found in", importPath)
// 			}
// 			for _, yamlFilePath := range yamlFilePaths {
// 				yamlFileName := filepath.Base(yamlFilePath)
// 				err := importYaml(yamlFilePath, "apply", "", "CatalogItem")
// 				if err != nil {
// 					log.Warnln("Failed to import", yamlFilePath, "as CatalogItem", err)
// 				} else {
// 					fmt.Println("Imported", yamlFileName, "successfully - CatalogItem updated.")
// 				}
// 			}
// 		}
// 	},
// }

// // deleteCatalogItemCmd represents the executions command
// var deleteCatalogItemCmd = &cobra.Command{
// 	Use:   "CatalogItem",
// 	Short: "Delete an CatalogItem",
// 	Long: `Delete an CatalogItem with a specific CatalogItem ID or Name

// # Delete CatalogItem by ID:
// vra-cli delete CatalogItem --id "CatalogItem ID"

// # Delete CatalogItem by Name:
// vra-cli delete CatalogItem --name "CatalogItem Name"

// # Delete CatalogItem by Project and Name:
// vra-cli delete CatalogItem --project "My Project" --name "CatalogItem Name"

// # Delete all CatalogItems in Project (prompts for confirmation):
// vra-cli delete CatalogItem --project "My Project"
// 	`,
// 	Args: func(cmd *cobra.Command, args []string) error {
// 		// if id != "" && name != "" {
// 		// 	return errors.New("please specify either CatalogItem name or CatalogItem id")
// 		// }
// 		// if id == "" && name == "" {
// 		// 	return errors.New("please specify CatalogItem name or CatalogItem id")
// 		// }

// 		return nil
// 	},
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if err := ensureTargetConnection(); err != nil {
// 			log.Fatalln(err)
// 		}
// 		if name != "" {
// 			response, err := getCatalogItem(id, name, project, typename, exportPath)
// 			if err != nil {
// 				log.Fatalln(err)
// 			}
// 			id = response[0].ID
// 		}

// 		if id != "" {

// 			response, err := deleteCatalogItem(id)
// 			if err != nil {
// 				log.Errorln("Unable to delete CatalogItem: ", err)
// 			}
// 			log.Infoln("CatalogItem with id " + response.ID + " deleted")
// 		} else if project != "" {
// 			response, err := deleteCatalogItemByProject(project)
// 			if err != nil {
// 				log.Errorln("Unable to delete CatalogItem: ", err)
// 			}
// 			log.Infoln(len(response), "CatalogItems deleted")

// 		}

// 	},
// }

func init() {
	getCmd.AddCommand(getCatalogItemCmd)
	getCatalogItemCmd.Flags().StringVarP(&name, "name", "n", "", "Get CatalogItem by Name")
	getCatalogItemCmd.Flags().StringVarP(&id, "id", "i", "", "Get CatalogItem by ID")
	getCatalogItemCmd.Flags().StringVarP(&projectName, "project", "p", "", "Filter CatalogItem by Project")
	getCatalogItemCmd.Flags().StringVarP(&typename, "type", "t", "", "Filter CatalogItem by Type")
	getCatalogItemCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	// // Create
	createCmd.AddCommand(createCatalogItemCmd)
	createCatalogItemCmd.Flags().StringVar(&deploymentName, "deploymentName", "", "Get CatalogItem by Name")
	createCatalogItemCmd.Flags().StringVar(&deploymentReason, "deploymentReason", "", "Get CatalogItem by ID")
	createCatalogItemCmd.Flags().StringVarP(&id, "id", "i", "", "Get CatalogItem by ID")
	createCatalogItemCmd.Flags().StringVarP(&name, "name", "n", "", "Get CatalogItem by Name")
	createCatalogItemCmd.Flags().StringVarP(&projectName, "project", "p", "", "Manually specify the Project in which to create the CatalogItem (overrides YAML)")
	createCatalogItemCmd.MarkFlagRequired("deploymentName")
	createCatalogItemCmd.MarkFlagRequired("project")
	// // Update
	// updateCmd.AddCommand(updateCatalogItemCmd)
	// updateCatalogItemCmd.Flags().StringVarP(&importPath, "importPath", "c", "", "YAML configuration file to import")
	// updateCatalogItemCmd.MarkFlagRequired("importPath")
	// // Delete
	// deleteCmd.AddCommand(deleteCatalogItemCmd)
	// deleteCatalogItemCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the CatalogItem to delete")
	// deleteCatalogItemCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the CatalogItem to delete")
	// deleteCatalogItemCmd.Flags().StringVarP(&project, "project", "p", "", "Delete CatalogItems by Project")

}
