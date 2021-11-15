/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"
	"strconv"

	"github.com/sammcgeown/vra-cli/pkg/cmd/orchestrator"

	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	exportConfigurationAttributeValues      bool
	exportConfigSecureStringAttributeValues bool
	exportGlobalTags                        bool
	viewContents                            bool
	addToPackage                            bool
	editContents                            bool
)

// getPackageCmd represents the Packages command
var getPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Get Orchestrator Packages",
	Long:  `Get Orchestrator Packages by Name`,
	Run: func(cmd *cobra.Command, args []string) {

		response, err := orchestrator.GetPackage(APIClient, name)
		if err != nil {
			log.Errorln("Unable to get Packages: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No Packages found")
		} else {
			if APIClient.Output == "table" {
				// Print result table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Workflows", "Actions", "Configurations", "Resources"})
				for _, c := range response {
					table.Append([]string{c.ID, c.Name, strconv.Itoa(len(c.Workflows)), strconv.Itoa(len(c.Actions)), strconv.Itoa(len(c.Configurations)), strconv.Itoa(len(c.Resources))})
				}
				table.Render()
			} else if APIClient.Output == "export" {
				// Export the Package
				for _, Package := range response {
					var options = types.ExportPackageOptions{
						ExportConfigurationAttributeValues:      exportConfigurationAttributeValues,
						ExportConfigSecureStringAttributeValues: exportConfigSecureStringAttributeValues,
						ExportGlobalTags:                        exportGlobalTags,
						ViewContents:                            viewContents,
						AddToPackage:                            addToPackage,
						EditContents:                            editContents,
					}
					err := orchestrator.ExportPackage(APIClient, Package.Name, options, exportPath)
					if err != nil {
						log.Warnln("Unable to export Package: ", err)
					} else {
						log.Infoln("Package", Package.Name, "exported")
					}
				}

			} else {
				helpers.PrettyPrint(response)
			}
		}

	},
}

// // delPackageCmd represents the delete Packages command
// var delPackageCmd = &cobra.Command{
// 	Use:   "package",
// 	Short: "Delete an Package",
// 	Long:  `Delete an Package with a specific Package ID`,
// 	Run: func(cmd *cobra.Command, args []string) {

// 		_, err := orchestrator.DeletePackage(APIClient, id)
// 		if err != nil {
// 			log.Errorln("Unable to delete Package: ", err)
// 		} else {
// 			log.Infoln("Package with ID " + id + " deleted")
// 		}
// 	},
// }

// // createPackageCmd represents the Packages command
// var createPackageCmd = &cobra.Command{
// 	Use:   "package",
// 	Short: "Create a Package",
// 	Long:  `Create a Package`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// Get the category ID
// 		var CategoryID string
// 		categoryName := (strings.Split(category, "/"))[len(strings.Split(category, "/"))-1]
// 		categories, _ := orchestrator.GetCategoryByName(APIClient, categoryName, "PackageCategory")
// 		if len(categories) == 0 {
// 			log.Fatalln("Unable to find category:", categoryName)
// 		} else if len(categories) == 1 {
// 			// Only one category found
// 			log.Debugln("Category found:", categories[0].Name, categories[0].ID)
// 			CategoryID = categories[0].ID
// 		} else {
// 			for _, matchedCategory := range categories {
// 				if matchedCategory.Path == category {
// 					log.Debugln("Category ID:", matchedCategory.ID)
// 					CategoryID = matchedCategory.ID
// 					break
// 				}
// 			}
// 			if CategoryID == "" {
// 				log.Fatalln("Multiple categories found, try using a more specific category - e.g.: path/to/category")
// 			}
// 		}
// 		for _, path := range helpers.GetFilePaths(importPath, ".zip") {
// 			log.Infoln("Importing Package:", path)
// 			err := orchestrator.ImportPackage(APIClient, path, CategoryID)
// 			if err != nil {
// 				log.Errorln("Unable to import Package: ", err)
// 			} else {
// 				Package, err := orchestrator.GetPackage(APIClient, "", categoryName, name)
// 				if err != nil {
// 					log.Errorln("Package imported OK, but I'm unable to get imported Package details: ", err)
// 				}
// 				log.Infoln("Package imported:", Package[0].Name, "with ID:", Package[0].ID)
// 			}
// 		}

// 	},
// }

func init() {
	// Get
	getCmd.AddCommand(getPackageCmd)
	getPackageCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Package")
	getPackageCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	getPackageCmd.Flags().BoolVar(&exportConfigurationAttributeValues, "exportConfigurationAttributeValues", false, "(Export) Add configuration attribute values to package")
	getPackageCmd.Flags().BoolVar(&exportConfigSecureStringAttributeValues, "exportConfigSecureStringAttributeValues", false, "(Export) Add configuration SecureString attribute values to package")
	getPackageCmd.Flags().BoolVar(&exportGlobalTags, "exportGlobalTags", false, "(Export) Add global tags to package")
	getPackageCmd.Flags().BoolVar(&viewContents, "viewContents", false, "(Export) Set `View Contents` permission")
	getPackageCmd.Flags().BoolVar(&addToPackage, "addToPackage", false, "(Export) Set `Add to package` permission")
	getPackageCmd.Flags().BoolVar(&editContents, "editContents", false, "(Export) Set `Edit contents` permission")
	// // Delete
	// deleteCmd.AddCommand(delPackageCmd)
	// delPackageCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Package")
	// delPackageCmd.MarkFlagRequired("id")
	// // Create
	// createCmd.AddCommand(createPackageCmd)
	// createPackageCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Package")
	// createPackageCmd.Flags().StringVar(&importPath, "importPath", "", "Path to the zip file, or folder containing zip files, to import")
	// createPackageCmd.MarkFlagRequired("importPath")
}
