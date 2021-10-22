/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	"github.com/sammcgeown/vra-cli/pkg/cmd/orchestrator"

	"github.com/sammcgeown/vra-cli/pkg/util/auth"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getWorkflowCmd represents the workflows command
var getWorkflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Get Orchestrator Workflows",
	Long: `Get Orchestrator Workflows by ID, Name, Project and Status

# Get only failed workflows:
vra-cli get workflow --status FAILED
# Get an workflow by ID:
vra-cli get workflow --id bb3f6aff-311a-45fe-8081-5845a529068d
# Get Failed workflows in Project "Field Demo" with the name "Learn Code Stream"
vra-cli get workflow --status FAILED --project "Field Demo" --name "Learn Code Stream"`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.GetConnection(&targetConfig, debug); err != nil {
			log.Fatalln(err)
		}

		response, err := orchestrator.GetWorkflow(restClient, id, category, name)
		if err != nil {
			log.Errorln("Unable to get workflows: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		} else {
			if output == "table" {
				// Print result table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Version", "Description", "Category"})
				for _, c := range response {
					category, _ := orchestrator.GetCategoryByID(restClient, c.CategoryID)
					table.Append([]string{c.ID, c.Name, c.Version, c.Description, category.Path})
				}
				table.Render()
			} else {
				helpers.PrettyPrint(response)
			}
		}

	},
}

// // delExecutionCmd represents the workflows command
// var delExecutionCmd = &cobra.Command{
// 	Use:   "workflow",
// 	Short: "Delete an Execution",
// 	Long: `Delete an Execution with a specific Execution ID

// 	`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if err := auth.GetConnection(&targetConfig, debug); err != nil {
// 			log.Fatalln(err)
// 		}
// 		if id != "" {
// 			_, err := codestream.DeleteExecution(restClient, id)
// 			if err != nil {
// 				log.Errorln("Unable to delete workflow: ", err)
// 			} else {
// 				log.Infoln("Execution with id " + id + " deleted")
// 			}
// 		} else if projectName != "" {
// 			response, err := codestream.DeleteExecutions(restClient, confirm, projectName, status, name, nested)
// 			if err != nil {
// 				log.Errorln("Unable to delete workflows: ", err)
// 			} else {
// 				log.Infoln(len(response), "Executions deleted")
// 			}
// 		}
// 	},
// }

// // createExecutionCmd represents the workflows command
// var createExecutionCmd = &cobra.Command{
// 	Use:   "workflow",
// 	Short: "Create an Execution",
// 	Long: `Create an Execution with a specific pipeline ID and form payload.

// 	`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if err := auth.GetConnection(&targetConfig, debug); err != nil {
// 			log.Fatalln(err)
// 		}

// 		response, err := codestream.CreateExecution(restClient, id, inputs, comments)
// 		if err != nil {
// 			log.Errorln("Unable to create workflow: ", err)
// 		}
// 		log.Infoln("Execution " + response.ExecutionLink + " created")

// 	},
// }

func init() {
	// Get
	getCmd.AddCommand(getWorkflowCmd)
	getWorkflowCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the workflow")
	getWorkflowCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the workflows to list")
	getWorkflowCmd.Flags().StringVarP(&category, "category", "c", "", "Filter workflows by Category")
	// // Delete
	// deleteCmd.AddCommand(delExecutionCmd)
	// delExecutionCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to delete workflows for")
	// delExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the workflow to delete")
	// delExecutionCmd.Flags().StringVarP(&status, "status", "s", "", "Delete workflows by status (Completed|Waiting|Pausing|Paused|Resuming|Running)")
	// delExecutionCmd.Flags().StringVarP(&projectName, "project", "p", "", "Delete workflows by Project")
	// delExecutionCmd.Flags().BoolVarP(&nested, "nested", "", false, "Delete nested workflows")
	// // Create
	// createCmd.AddCommand(createExecutionCmd)
	// createExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to execute")
	// createExecutionCmd.Flags().StringVarP(&inputs, "inputs", "", "", "JSON form inputs")
	// createExecutionCmd.Flags().StringVarP(&inputPath, "inputPath", "", "", "JSON input file")
	// createExecutionCmd.Flags().StringVarP(&comments, "comments", "", "", "Execution comments")
	// createExecutionCmd.MarkFlagRequired("id")
}
