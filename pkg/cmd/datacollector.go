/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sammcgeown/vra-cli/pkg/cmd/cloudassembly"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getCloudAccountCmd represents the Blueprint command
var getDataCollectorCmd = &cobra.Command{
	Use:   "datacollector",
	Short: "Get Data Collectors",
	Long: `Get Data Collectors by ID

Get Data Collector by ID:
  vra-cli get datacollector --id <datacollector-id>

Get all Data Collectors:
  vra-cli get datacollector`,
	Run: func(cmd *cobra.Command, args []string) {
		if APIClient.Config.Server != "api.mgmt.cloud.vmware.com" {
			log.Fatalln("Data Collectors (Cloud Proxies) are only supported on vRealize Automation Cloud")
		}
		dataCollectors, err := cloudassembly.GetDataCollector(APIClient, id)
		if err != nil {
			log.Fatalln(err)
		}

		if len(dataCollectors) == 0 {
			log.Warnln("No Data Collector (Cloud Proxy) found")
		} else if len(dataCollectors) == 1 {
			helpers.PrettyPrint(dataCollectors[0])
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Hostname", "IP Address", "Status"})
			for _, c := range dataCollectors {
				table.Append([]string{*c.DcID, *c.Name, *c.HostName, *c.IPAddress, *c.Status})
			}
			table.Render()
		}
	},
}

func init() {
	// Get
	getCmd.AddCommand(getDataCollectorCmd)
	getDataCollectorCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Data Collector (Cloud Proxy)")

}
