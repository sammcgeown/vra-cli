/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/sammcgeown/vra-cli/pkg/util/auth"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var nested bool
var inputs string
var comments string
var inputPath string

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
		if err := auth.GetConnection(targetConfig, debug); err != nil {
			log.Fatalln(err)
		}

		response, err := getExecutions(id, projectName, status, name, nested)
		if err != nil {
			log.Errorln("Unable to get executions: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		} else if resultCount == 1 {
			helpers.PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Status", "Message"})
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
		if err := auth.GetConnection(targetConfig, debug); err != nil {
			log.Fatalln(err)
		}
		if id != "" {
			response, err := deleteExecution(id)
			if err != nil {
				log.Errorln("Unable to delete execution: ", err)
			} else {
				log.Infoln("Execution with id " + response.ID + " deleted")
			}
		} else if projectName != "" {
			response, err := deleteExecutions(projectName, status, name, nested)
			if err != nil {
				log.Errorln("Unable to delete executions: ", err)
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
		if err := auth.GetConnection(targetConfig, debug); err != nil {
			log.Fatalln(err)
		}

		response, err := createExecution(id, inputs, comments)
		if err != nil {
			log.Errorln("Unable to create execution: ", err)
		}
		log.Infoln("Execution " + response.ExecutionLink + " created")

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
	// Delete
	deleteCmd.AddCommand(delExecutionCmd)
	delExecutionCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to delete executions for")
	delExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the execution to delete")
	delExecutionCmd.Flags().StringVarP(&status, "status", "s", "", "Delete executions by status (Completed|Waiting|Pausing|Paused|Resuming|Running)")
	delExecutionCmd.Flags().StringVarP(&projectName, "project", "p", "", "Delete executions by Project")
	delExecutionCmd.Flags().BoolVarP(&nested, "nested", "", false, "Delete nested executions")
	// Create
	createCmd.AddCommand(createExecutionCmd)
	createExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to execute")
	createExecutionCmd.Flags().StringVarP(&inputs, "inputs", "", "", "JSON form inputs")
	createExecutionCmd.Flags().StringVarP(&inputPath, "inputPath", "", "", "JSON input file")
	createExecutionCmd.Flags().StringVarP(&comments, "comments", "", "", "Execution comments")
	createExecutionCmd.MarkFlagRequired("id")
}
