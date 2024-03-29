/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/sammcgeown/vra-cli/pkg/cmd/cloudassembly"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

var (
	admins                string
	members               string
	viewers               string
	operationTimeout      int64
	machineNamingTemplate string
	sharedResources       bool
)

// getProjectCommand represents the project command
var getProjectCommand = &cobra.Command{
	Use:   "project",
	Short: "Get Projects",
	Long: `Get all Projects, or get Project by Name or ID
Get all projects:
  vra-cli get project
Get Projects by ID:
  vra-cli get project --id <project ID>
Get Project by Name (case sensitive):
  vra-cli get project --name <project name>`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := cloudassembly.GetProject(APIClient, projectName, id)
		if err != nil {
			log.Errorln("Unable to get Code Stream Projects: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Fatalln("No results found")
		}
		if APIClient.Output == "json" {
			helpers.PrettyPrint(response)
		} else if APIClient.Output == "export" {
			// Export the Project
			// if exportPath != "" {
			// 	tmpDir, err := ioutil.TempDir(os.TempDir(), "vra-cli-*")
			// 	if err != nil {
			// 		log.Fatalln(err)
			// 	}
			// 	zipFile := filepath.Join(exportPath, p.Name+".zip")
			// 	var zipFiles []string
			// 	log.Debugln(zipFile)
			// 	pipelines, _ := getPipelines("", "", p.Name, filepath.Join(tmpDir, p.Name, "pipelines"))
			// 	//pipelineTable.SetHeader([]string{"Id", "Name", "Project", "Description"})
			// 	for _, c := range pipelines {
			// 		zipFiles = append(zipFiles, filepath.Join(tmpDir, p.Name, "pipelines", c.Name+".yaml"))
			// 		//pipelineTable.Append([]string{c.ID, c.Name, c.Project, c.Description})
			// 	}
			// 	variables, _ := getVariable("", "", p.Name, filepath.Join(tmpDir, p.Name))
			// 	//variableTable.SetHeader([]string{"Id", "Name", "Project", "Description"})
			// 	if len(variables) > 0 {
			// 		zipFiles = append(zipFiles, filepath.Join(tmpDir, p.Name, "variables.yaml"))
			// 	}
			// 	// for _, c := range variables {
			// 	// 	//variableTable.Append([]string{c.ID, c.Name, c.Project, c.Description})
			// 	// }
			// 	endpoints, _ := getEndpoint("", "", p.Name, "", filepath.Join(tmpDir, p.Name, "endpoints"))
			// 	//endpointTable.SetHeader([]string{"ID", "Name", "Project", "Type", "Description"})
			// 	for _, c := range endpoints {
			// 		zipFiles = append(zipFiles, filepath.Join(tmpDir, p.Name, c.Name+".yaml"))
			// 		//endpointTable.Append([]string{c.ID, c.Name, c.Project, c.Type, c.Description})
			// 	}
			// 	if err := ZipFiles(zipFile, zipFiles, tmpDir); err != nil {
			// 		log.Fatalln(err)
			// 	}
			// }
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Description", "Admins", "Members", "Viewers"})
			for _, p := range response {
				table.Append([]string{*p.ID, p.Name, p.Description, helpers.UserArrayToString(p.Administrators), helpers.UserArrayToString(p.Members), helpers.UserArrayToString(p.Viewers)})
			}
			table.Render()
		}
	},
}

// createProjectCommand creates a project
var createProjectCommand = &cobra.Command{
	Use:   "project",
	Short: "Create a Project",
	Long: `Create a Project:
  vra-cli create project --name "vra-cli-test" --description "Test from vRA CLI" --admins smcgeown@vmware.com --members cas-user1@definit.co.uk,cas-user2@definit.co.uk --viewers cas-user3@definit.co.uk,cas-user1@definit.co.uk --machineNamingTemplate 'bob-${###}'	`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		adminUsers := helpers.CreateUserArray(strings.Split(admins, ","))
		memberUsers := helpers.CreateUserArray(strings.Split(members, ","))
		viewerUsers := helpers.CreateUserArray(strings.Split(viewers, ","))

		newProject, err := cloudassembly.CreateProject(APIClient, projectName, description, adminUsers, memberUsers, viewerUsers, nil, nil, operationTimeout, machineNamingTemplate, &sharedResources)
		if err != nil {
			log.Fatalln("Unable to create Project", err)
		}
		if APIClient.Output == "json" {
			helpers.PrettyPrint(newProject)
		} else if APIClient.Output == "export" {

		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Description", "Admins", "Members", "Viewers"})
			table.Append([]string{*newProject.ID, newProject.Name, newProject.Description, helpers.UserArrayToString(adminUsers), helpers.UserArrayToString(memberUsers), helpers.UserArrayToString(viewerUsers)})
			table.Render()
		}

	},
}

// updateProjectCommand creates a project
var updateProjectCommand = &cobra.Command{
	Use:   "project",
	Short: "Update a Project",
	Long:  `Update a Project`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		currentProject, err := cloudassembly.GetProject(APIClient, "", id)
		if err != nil || len(currentProject) == 0 {
			log.Fatalln("Unable to get Project", err)
		}
		var adminUsers, memberUsers, viewerUsers []*models.User
		if admins != "" {
			adminUsers = helpers.CreateUserArray(strings.Split(admins, ","))
		} else {
			adminUsers = currentProject[0].Administrators
		}
		if members != "" {
			memberUsers = helpers.CreateUserArray(strings.Split(members, ","))
		} else {
			memberUsers = currentProject[0].Members
		}
		if viewers != "" {
			viewerUsers = helpers.CreateUserArray(strings.Split(viewers, ","))
		} else {
			viewerUsers = currentProject[0].Viewers
		}

		newProject, err := cloudassembly.UpdateProject(APIClient, id, projectName, description, adminUsers, memberUsers, viewerUsers, nil, nil, operationTimeout, machineNamingTemplate, &sharedResources)
		if err != nil {
			log.Fatal("Unable to update Project", err)
		}
		if APIClient.Output == "json" {
			helpers.PrettyPrint(newProject)
		} else if APIClient.Output == "export" {

		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Description", "Admins", "Members", "Viewers"})
			table.Append([]string{*newProject.ID, newProject.Name, newProject.Description, helpers.UserArrayToString(adminUsers), helpers.UserArrayToString(memberUsers), helpers.UserArrayToString(viewerUsers)})
			table.Render()
		}

	},
}

// deleteProjectCommand deletes a project
var deleteProjectCommand = &cobra.Command{
	Use:   "project",
	Short: "Delete a Project",
	Long: `Delete a Project by ID

Delete by ID:
  vra-cli delete project --id <project ID>`,
	Run: func(cmd *cobra.Command, args []string) {
		if id != "" {
			err := cloudassembly.DeleteProject(APIClient, id)
			if err != nil {
				log.Errorln("Delete Project failed:", err)
			}
			log.Infoln("Pipeline id " + id + " deleted")
		}
	},
}

func init() {
	// Get
	getCmd.AddCommand(getProjectCommand)
	getProjectCommand.Flags().StringVarP(&projectName, "name", "n", "", "Name of the Project (case sensitive)")
	getProjectCommand.Flags().StringVarP(&id, "id", "i", "", "ID of the Project")

	// Create
	createCmd.AddCommand(createProjectCommand)
	createProjectCommand.Flags().StringVarP(&projectName, "name", "n", "", "Name of the Project")
	createProjectCommand.Flags().StringVarP(&description, "description", "d", "", "Description of the Project")
	createProjectCommand.Flags().StringVar(&admins, "admins", "", "Comma separated list of email addresses to assign administrator role for this project")
	createProjectCommand.Flags().StringVar(&members, "members", "", "Comma separated list of email addresses to assign member role for this project")
	createProjectCommand.Flags().StringVar(&viewers, "viewers", "", "Comma separated list of email addresses to assign viewer role for this project")
	createProjectCommand.Flags().Int64Var(&operationTimeout, "timeout", 0, "Operation Timeout setting for this project")
	createProjectCommand.Flags().StringVar(&machineNamingTemplate, "machineNamingTemplate", "", "Machine naming template for this project")
	createProjectCommand.Flags().BoolVar(&sharedResources, "sharedResources", false, "If true, Deployments are shared between all users in the project")
	createProjectCommand.MarkFlagRequired("name")

	// Update
	updateCmd.AddCommand(updateProjectCommand)
	updateProjectCommand.Flags().StringVarP(&id, "id", "i", "", "ID of the Project")
	updateProjectCommand.MarkFlagRequired("id")
	updateProjectCommand.Flags().StringVarP(&projectName, "name", "n", "", "Name of the Project")
	updateProjectCommand.Flags().StringVarP(&description, "description", "d", "", "Description of the Project")
	updateProjectCommand.Flags().StringVar(&admins, "admins", "", "Comma separated list of email addresses to assign administrator role for this project")
	updateProjectCommand.Flags().StringVar(&members, "members", "", "Comma separated list of email addresses to assign member role for this project")
	updateProjectCommand.Flags().StringVar(&viewers, "viewers", "", "Comma separated list of email addresses to assign viewer role for this project")
	updateProjectCommand.Flags().Int64Var(&operationTimeout, "timeout", 0, "Operation Timeout setting for this project")
	updateProjectCommand.Flags().StringVar(&machineNamingTemplate, "machineNamingTemplate", "", "Machine naming template for this project")
	updateProjectCommand.Flags().BoolVar(&sharedResources, "sharedResources", false, "If true, Deployments are shared between all users in the project")

	// Delete
	deleteCmd.AddCommand(deleteProjectCommand)
	deleteProjectCommand.Flags().StringVarP(&id, "id", "i", "", "ID of the Project to delete")
	deleteProjectCommand.MarkFlagRequired("id")

}
