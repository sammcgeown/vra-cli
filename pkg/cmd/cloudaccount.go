/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"errors"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sammcgeown/vra-cli/pkg/cmd/cloudassembly"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cloudaccounttype string
	tags             string
	regions          string
	// AWS
	awsaccesskeyid     string
	awssecretaccesskey string
	// Azure
	subscriptionID string
	tenantID       string
	clientID       string
	clientSecret   string
	// vSphere
	fqdn            string
	username        string
	password        string
	nsxaccount      string
	cloudproxy      string
	insecure        bool
	createcloudzone bool
	// NSX
	vccloudaccount string
	nsxtglobal     bool
	nsxtmanager    bool
)

// getCloudAccountCmd represents the Blueprint command
var getCloudAccountCmd = &cobra.Command{
	Use:   "cloudaccount",
	Short: "Get Cloud Accounts",
	Long: `Get Cloud Accounts by ID, name or type

Get Cloud Account by ID:
  vra-cli get cloudaccount --id <cloudaccount-id>

Get Cloud Account by Name:
  vra-cli get cloudaccount --name <cloudaccount-name>

Get Cloud Accounts by Type:
  vra-cli get cloudaccount --type <cloudaccount-type>`,
	Run: func(cmd *cobra.Command, args []string) {
		cloudAccounts, err := cloudassembly.GetCloudAccounts(APIClient, id, name, cloudaccounttype)
		if err != nil {
			log.Fatalln(err)
		}

		if len(cloudAccounts) == 0 {
			log.Warnln("No Cloud Accounts found")
		} else if len(cloudAccounts) == 1 {
			helpers.PrettyPrint(cloudAccounts[0])
		} else {
			if APIClient.Output == "json" {
				helpers.PrettyPrint(cloudAccounts)
			} else if APIClient.Output == "export" {
				log.Warnln("Exporting Cloud Accounts is not supported yet")
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Description", "Type"})
				for _, c := range cloudAccounts {
					table.Append([]string{*c.ID, c.Name, c.Description, *c.CloudAccountType})
				}
				table.Render()
			}
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
// 		if err := auth.GetConnection(&targetConfig, debug); err != nil {
// 			log.Fatalln(err)
// 		}

// 	},
// }

// createCloudAccountCmd represents the Blueprint create command
var createCloudAccountCmd = &cobra.Command{
	Use:   "cloudaccount",
	Short: "Create a Cloud Account",
	Long: `Create a Cloud Account.

Create a new AWS Cloud Account:
  vra-cli create cloudaccount --name <name> --description <description> --type aws --awsaccesskeyid <AWS Access Key ID> \
    --awssecretaccesskey <AWS Secret Access Key> --tags "<tag:value>,<tag:value>" \
	--regions "<region>,<region>"
	
Create a new Azure Cloud Account:
  vra-cli create cloudaccount --name <name> --description <description> --type azure \
    --subscriptionid <subscription ID> --tenantid <tenant ID> --clientid <client ID> \
	 --clientsecret <client secret> --tags "<tag:value>,<tag:value>" \
	 --regions "<region>,<region>"`,
	Args: func(cmd *cobra.Command, args []string) error {
		switch cloudaccounttype {
		case "aws":
			if awsaccesskeyid == "" ||
				awssecretaccesskey == "" {
				return errors.New("--awsaccesskeyid and --awssecretaccesskey are required for AWS Cloud Accounts")
			}
		case "azure":
			if subscriptionID == "" ||
				tenantID == "" ||
				clientID == "" ||
				clientSecret == "" {
				return errors.New("--subscriptionID, --tenantID, --clientID, and --clientSecret are required for Azure Cloud Accounts")
			}
		case "vsphere":
			if fqdn == "" ||
				username == "" ||
				password == "" {
				return errors.New("--fqdn, --username, and --password are required for vSphere Cloud Accounts")
			}
		case "nsxt":
			if fqdn == "" ||
				username == "" ||
				password == "" {
				return errors.New("--fqdn, --username, and --password are required for NSX-T Accounts")
			}
		default:
			return errors.New("--type is required (aws/azure/vsphere/nsxt)")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// if helpers.IsInputFromPipe() { // If it's a pipe, then read from stdin
		// 	// Decode JSON to struct
		// 	if err := json.NewDecoder(os.Stdin).Decode(&CloudAccount); err != nil {
		// 		log.Warnln(err)
		// 	}
		// }
		if cloudaccounttype == "aws" {
			newAccount, err := cloudassembly.CreateCloudAccountAWS(
				APIClient,
				name,
				awsaccesskeyid,
				awssecretaccesskey,
				regions,
				tags,
			)
			if err != nil {
				log.Fatalln(err)
			}
			helpers.PrettyPrint(newAccount)
		} else if cloudaccounttype == "azure" {
			newAccount, err := cloudassembly.CreateCloudAccountAzure(
				APIClient,
				name,
				description,
				subscriptionID,
				tenantID,
				clientID,
				clientSecret,
				regions,
				tags,
			)
			if err != nil {
				log.Fatalln(err)
			}
			helpers.PrettyPrint(newAccount)
		} else if cloudaccounttype == "vsphere" {
			newAccount, err := cloudassembly.CreateCloudAccountvSphere(
				APIClient,
				name,
				description,
				fqdn,
				username,
				password,
				nsxaccount,
				cloudproxy,
				tags,
				insecure,
				createcloudzone,
			)
			if err != nil {
				log.Fatalln(err)
			}
			helpers.PrettyPrint(newAccount)
		} else if cloudaccounttype == "nsxt" {
			newAccount, err := cloudassembly.CreateCloudAccountNsxT(
				APIClient,
				name,
				description,
				fqdn,
				username,
				password,
				vccloudaccount,
				cloudproxy,
				tags,
				nsxtglobal,
				nsxtmanager,
				insecure,
			)
			if err != nil {
				log.Fatalln(err)
			}
			helpers.PrettyPrint(newAccount)
		}
	},
}

// deleteCloudAccountCmd represents the delete Blueprint command
var deleteCloudAccountCmd = &cobra.Command{
	Use:   "cloudaccount",
	Short: "Delete a Cloud Account",
	Long: `Delete a Cloud Account with a specific ID or Name

Delete a Cloud Account by Name:
  vra-cli delete cloudaccount --name <Cloud Account Name>

Delete a Cloud Account by ID:
  vra-cli delete cloudaccount --id <Cloud Account ID>`,
	Run: func(cmd *cobra.Command, args []string) {
		if account, err := cloudassembly.GetCloudAccounts(APIClient, id, name, ""); err != nil {
			log.Fatalln(err) // There was an error getting the cloud account
		} else {
			if len(account) == 0 {
				// No error was throw, but there was no cloud account
				log.Fatalln("No Cloud Account matching the request was found")
			} else if len(account) > 1 {
				// There was more than one cloud account
				log.Fatalln("More than one Cloud Account matching the request was found")
			} else {
				// There was only one cloud account
				if err := cloudassembly.DeleteCloudAccount(APIClient, *account[0].ID); err != nil {
					log.Fatalln(err) // There was an error deleting the cloud account
				} else {
					log.Infoln("Cloud Account deleted successfully")
				}
			}
		}

	},
}

func init() {
	// Get
	getCmd.AddCommand(getCloudAccountCmd)
	getCloudAccountCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Account")
	getCloudAccountCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Cloud Account")
	getCloudAccountCmd.Flags().StringVarP(&cloudaccounttype, "type", "t", "", "List by Type of the Cloud Account")

	// Create
	createCmd.AddCommand(createCloudAccountCmd)
	createCloudAccountCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Account")
	createCloudAccountCmd.MarkFlagRequired("name")
	createCloudAccountCmd.Flags().StringVarP(&description, "description", "d", "", "Decscription of the Cloud Account")
	createCloudAccountCmd.Flags().StringVarP(&cloudaccounttype, "type", "t", "", "Type of the Cloud Account")
	createCloudAccountCmd.MarkFlagRequired("type")
	createCloudAccountCmd.Flags().StringVar(&tags, "tags", "", "List of Tags (comma separated e.g. \"name1:value2,name2:value2\") to apply to the Cloud Account")
	createCloudAccountCmd.Flags().StringVar(&cloudproxy, "cloudproxy", "", "vRA Cloud only - ID of the Data Collector (Cloud Proxy) (use: vra-cli get datacollector)")
	createCloudAccountCmd.Flags().BoolVar(&insecure, "insecure", false, "Ignore Self-Signed Certificates")
	createCloudAccountCmd.Flags().StringVar(&regions, "regions", "", "List of Regions (comma separated) of the Cloud Account")
	// Create Azure Cloud Account
	createCloudAccountCmd.Flags().StringVar(&subscriptionID, "subscriptionid", "", "Azure Subscription ID")
	createCloudAccountCmd.Flags().StringVar(&tenantID, "tenantid", "", "Azure Tenant ID")
	createCloudAccountCmd.Flags().StringVar(&clientID, "clientid", "", "Azure Client Application ID")
	createCloudAccountCmd.Flags().StringVar(&clientSecret, "clientsecret", "", "Azure Client Application Secret ID")
	// Create AWS Cloud Account
	createCloudAccountCmd.Flags().StringVar(&awsaccesskeyid, "awsaccesskeyid", "", "AWS Access Key ID of the Cloud Account")
	createCloudAccountCmd.Flags().StringVar(&awssecretaccesskey, "awssecretaccesskey", "", "AWS Secret Access Key of the Cloud Account")
	// Create vSphere Cloud Account
	createCloudAccountCmd.Flags().StringVar(&fqdn, "fqdn", "", "vCenter Server FQDN")
	createCloudAccountCmd.Flags().StringVar(&username, "username", "", "User Name")
	createCloudAccountCmd.Flags().StringVar(&password, "password", "", "Password")
	createCloudAccountCmd.Flags().StringVar(&nsxaccount, "nsxaccount", "", "ID of the NSX-T or NSX-v Cloud Account to link (use: vra-cli get cloudaccount --type nsxt/nsxv)")
	createCloudAccountCmd.Flags().BoolVar(&createcloudzone, "createcloudzone", false, "Automatically create a Cloud Zone for this Account")
	// Create NSX-T Cloud Account
	createCloudAccountCmd.Flags().StringVar(&vccloudaccount, "vccloudaccount", "", "Name of the vCenter Cloud Account to associate with NSX")
	createCloudAccountCmd.Flags().BoolVar(&nsxtglobal, "nsxtglobal", false, "NSX-T is Global")
	createCloudAccountCmd.Flags().BoolVar(&nsxtmanager, "nsxtmanager", false, "NSX-T Manager mode (true: manager, false: policy)")

	// // Update
	// updateCmd.AddCommand(updateCloudAccountCmd)
	// updateCloudAccountCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Blueprint to list")
	// updateCloudAccountCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Configuration file to import")
	// updateCloudAccountCmd.Flags().StringVarP(&state, "state", "s", "", "Set the state of the Blueprint (ENABLED|DISABLED|RELEASED")

	// Delete
	deleteCmd.AddCommand(deleteCloudAccountCmd)
	deleteCloudAccountCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Cloud Account to delete")
	deleteCloudAccountCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Cloud Account to delete")
}
