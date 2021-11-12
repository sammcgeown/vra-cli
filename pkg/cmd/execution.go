/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/sammcgeown/vra-cli/pkg/cmd/codestream"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var nested, rollback bool
var inputs, comments string

// getExecutionCmd represents the executions command
var getExecutionCmd = &cobra.Command{
	Use:   "execution",
	Short: "Get Pipeline Executions",
	Long: `Get Code Stream Pipeline Executions by ID, Name, Project and Status

# Get only failed executions:
vra-cli get execution --status FAILED
# Get an execution by ID:
vra-cli get execution --id bb3f6aff-311a-45fe-8081-5845a529068d
# Get Failed executions in Project "Field Demo" with the name "Learn Code Stream"
vra-cli get execution --status FAILED --project "Field Demo" --name "Learn Code Stream"`,
	Run: func(cmd *cobra.Command, args []string) {

		response, err := codestream.GetExecution(APIClient, id, projectName, status, name, nested, rollback)
		if err != nil {
			log.Errorln("Unable to get executions: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Warnln("No results found")
			return
		}
		if APIClient.Output == "json" {
			helpers.PrettyPrint(response)
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Execution", "Project", "Status", "Message"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name + "#" + fmt.Sprint(c.Index), c.Project, c.Status, c.StatusMessage})
			}
			table.Render()
		}

	},
}

// delExecutionCmd represents the executions command
var delExecutionCmd = &cobra.Command{
	Use:   "execution",
	Short: "Delete an Execution",
	Long: `Delete an Execution with a specific Execution ID
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if id != "" {
			_, err := codestream.DeleteExecution(APIClient, id)
			if err != nil {
				log.Errorln("Unable to delete execution:", err)
			} else {
				log.Infoln("Execution with id", id, "deleted")
			}
		} else {
			response, err := codestream.DeleteExecutions(APIClient, projectName, status, name, nested, rollback)
			if err != nil {
				log.Errorln("Unable to delete executions:", err)
			} else {
				log.Infoln(len(response), "Executions deleted")
			}
		}
	},
}

// createExecutionCmd represents the executions command
var createExecutionCmd = &cobra.Command{
	Use:   "execution",
	Short: "Create an Execution",
	Long: `Create an Execution with a specific pipeline ID and form payload.
	
	`,
	Run: func(cmd *cobra.Command, args []string) {

		response, err := codestream.CreateExecution(APIClient, id, inputs, comments)
		if err != nil {
			log.Errorln("Unable to create execution:", err)
		}
		log.Infoln("Execution", response.ExecutionID, "created")

	},
}

func init() {
	// Get
	getCmd.AddCommand(getExecutionCmd)
	getExecutionCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to list executions for")
	getExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the executions to list")
	getExecutionCmd.Flags().StringVarP(&status, "status", "s", "", "Filter executions by status (Completed|Waiting|Pausing|Paused|Resuming|Running)")
	getExecutionCmd.Flags().StringVarP(&projectName, "project", "p", "", "Filter executions by Project")
	getExecutionCmd.Flags().BoolVarP(&nested, "nested", "", false, "Include nested executions")
	getExecutionCmd.Flags().BoolVarP(&rollback, "rollback", "", false, "Include rollback executions")
	// Delete
	deleteCmd.AddCommand(delExecutionCmd)
	delExecutionCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to delete executions for")
	delExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the execution to delete")
	delExecutionCmd.Flags().StringVarP(&status, "status", "s", "", "Delete executions by status (Completed|Waiting|Pausing|Paused|Resuming|Running)")
	delExecutionCmd.Flags().StringVarP(&projectName, "project", "p", "", "Delete executions by Project")
	delExecutionCmd.Flags().BoolVarP(&nested, "nested", "", false, "Delete nested executions")
	delExecutionCmd.Flags().BoolVarP(&rollback, "rollback", "", false, "Delete rollback executions")
	// Create
	createCmd.AddCommand(createExecutionCmd)
	createExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to execute")
	createExecutionCmd.Flags().StringVarP(&inputs, "inputs", "", "", "JSON form inputs")
	createExecutionCmd.Flags().StringVarP(&importPath, "importPath", "", "", "JSON input file")
	createExecutionCmd.Flags().StringVarP(&comments, "comments", "", "", "Execution comments")
	createExecutionCmd.MarkFlagRequired("id")
}
