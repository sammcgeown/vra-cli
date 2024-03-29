/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"errors"
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
	exportOptions types.ExportPackageOptions
	importOptions types.ImportPackageOptions
	deleteOption  string
)

// getPackageCmd represents the Packages command
var getPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Get Orchestrator Packages",
	Long: `Get Orchestrator Packages

# Get all Packages
vra-cli get package

# Get Package by Name
vra-cli get package --name <name>

# Export a Package by Name
vra-cli get package --name <name> --out export --exportPath <path> --exportConfigurationAttributeValues <true/false> \
	--exportConfigSecureStringAttributeValues <true/false> --exportGlobalTags <true/false> --viewContents <true/false> \
	--addToPackage <true/false> --editContents <true/false>`,
	Args: func(cmd *cobra.Command, args []string) error {
		if APIClient.Output != "export" && exportPath != "" {
			return errors.New("--exportPath is not required when not exporting")
		}
		return nil
	},
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

					err := orchestrator.ExportPackage(APIClient, Package.Name, exportOptions, exportPath)
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

// delPackageCmd represents the delete Packages command
var delPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Delete a Package",
	Long: `Delete a Package with a specific Package name

# Delete a Package (package only):
vra-cli delete package --name <name>

# Delete a Package and content
vra-cli delete package --name <name> --deleteOption deletePackageWithContent

# Delete a Package but keep shared content
vra-cli delete package --name <name> --deleteOption deletePackageKeepingShared`,
	Args: func(cmd *cobra.Command, args []string) error {
		if deleteOption != "deletePackage" && deleteOption != "deletePackageWithContent" && deleteOption != "deletePackageKeepingShared" {
			return errors.New("Invalid delete option. Available values: deletePackage, deletePackageWithContent, deletePackageKeepingShared")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		err := orchestrator.DeletePackage(APIClient, name, deleteOption)
		if err != nil {
			log.Errorln("Unable to delete Package: ", err)
		} else {
			log.Infoln("Package deleted")
		}
	},
}

// createPackageCmd represents the Packages command
var createPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Create a Package",
	Long:  `Create a Package`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		for _, path := range helpers.GetFilePaths(importPath, ".package") {
			log.Debugln("Importing Package:", path)
			packageDetails, packageErr := orchestrator.GetPackageDetails(APIClient, path, importOptions)
			if packageErr != nil {
				log.Fatalln("Unable to get Package details: ", packageErr)
			}

			log.Infoln("Importing", packageDetails.PackageName)
			if packageDetails.PackageAlreadyExists && !APIClient.Force {
				log.Warnln("Package already exists, only new or newer content will be imported, use --force to override.")
			}
			if !packageDetails.CertificateValid {
				helpers.PrettyPrint(packageDetails.CertificateInfo)
				if !helpers.AskForConfirmation("Certificate is not valid") {
					log.Fatalln("Certificate is not valid, user declined")
				}

			}
			if !packageDetails.CertificateTrusted {
				helpers.PrettyPrint(packageDetails.CertificateInfo)
				if !helpers.AskForConfirmation("Certificate is not trusted, continue?") {
					log.Fatalln("Certificate is not trusted, user declined")
				}
			}

			if packageDetails.PackageAlreadyExists {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Type", "Version", "Imported"})
				for _, value := range packageDetails.ImportElementDetails {
					table.Append([]string{value.FileObjectName, value.Type, value.FileObjectVersion, strconv.FormatBool(value.ImportIt)})
				}
				table.Render()
				if !helpers.AskForConfirmation("Package already exists, continue?") {
					log.Fatalln("Package already exists, user declined")
				}
			}

			importError := orchestrator.CreatePackage(APIClient, path, importOptions)
			if importError != nil {
				log.Errorln("Unable to import Package: ", importError)
			} else {
				Package, err := orchestrator.GetPackage(APIClient, packageDetails.PackageName)
				if err != nil {
					log.Errorln("Package imported OK, but I'm unable to get imported Package details: ", err)
				}
				log.Infoln("Package imported:", Package[0].Name, "with ID:", Package[0].ID)
			}
		}

	},
}

func init() {
	// Get
	getCmd.AddCommand(getPackageCmd)
	getPackageCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Package")
	getPackageCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	getPackageCmd.Flags().BoolVar(&exportOptions.ExportConfigurationAttributeValues, "exportConfigurationAttributeValues", false, "(Export) Add configuration attribute values to package")
	getPackageCmd.Flags().BoolVar(&exportOptions.ExportConfigSecureStringAttributeValues, "exportConfigSecureStringAttributeValues", false, "(Export) Add configuration SecureString attribute values to package")
	getPackageCmd.Flags().BoolVar(&exportOptions.ExportGlobalTags, "exportGlobalTags", false, "(Export) Add global tags to package")
	getPackageCmd.Flags().BoolVar(&exportOptions.ViewContents, "viewContents", true, "(Export) Set `View Contents` permission. Default: true")
	getPackageCmd.Flags().BoolVar(&exportOptions.AddToPackage, "addToPackage", true, "(Export) Set `Add to package` permission. Default: true")
	getPackageCmd.Flags().BoolVar(&exportOptions.EditContents, "editContents", true, "(Export) Set `Edit contents` permission. Default: true")
	// Delete
	deleteCmd.AddCommand(delPackageCmd)
	delPackageCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Package")
	delPackageCmd.Flags().StringVar(&deleteOption, "deleteOption", "deletePackage", "Package deletion options. Available values: deletePackage, deletePackageWithContent, deletePackageKeepingShared. Default: deletePackage")
	delPackageCmd.MarkFlagRequired("name")
	// Create
	createCmd.AddCommand(createPackageCmd)
	createPackageCmd.Flags().BoolVar(&importOptions.ImportConfigurationAttributeValues, "importConfigurationAttributeValues", true, "Import configuration attribute values with the package. Default: true")
	createPackageCmd.Flags().BoolVar(&importOptions.ImportConfigSecureStringAttributeValues, "importConfigSecureStringAttributeValues", true, "Import configuration SecureString attribute values with the package. Default: true")
	createPackageCmd.Flags().StringVar(&importOptions.TagImportMode, "tagImportMode", "ImportButPreserveExistingValue", "Tag import mode. Available values : DoNotImport, ImportAndOverwriteExistingValue, ImportButPreserveExistingValue. Default: ImportButPreserveExistingValue")
	createPackageCmd.Flags().StringVar(&importPath, "importPath", "", "Path to the zip file, or folder containing zip files, to import")
	createPackageCmd.MarkFlagRequired("importPath")
	// Update (alias of create for package import)
	updateCmd.AddCommand(createPackageCmd)
}
