/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func getProject(id, name string) ([]*models.Project, error) {
	var filters []string
	var filter string
	if id != "" {
		filters = append(filters, "(id eq '"+id+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if len(filters) > 0 {
		filter = "(" + strings.Join(filters, " and ") + ")"
	}

	log.Debugln("Filter:", filter)

	apiClient := getApiClient()

	ProjectParams := project.NewGetProjectsParams()
	ProjectParams.DollarFilter = &filter

	ret, err := apiClient.Project.GetProjects(ProjectParams)
	if err != nil {
		switch err.(type) {
		case *project.GetProjectNotFound:
			return nil, errors.New("Project with ID " + id + " not found")
		}
		return nil, err
	}
	return ret.Payload.Content, nil
}

// func getProject(id, name string) ([]*CodeStreamProject, error) {
// 	var projects []*CodeStreamProject
// 	client := resty.New()

// 	var filters []string
// 	if id != "" {
// 		filters = append(filters, "(id eq '"+id+"')")
// 	}
// 	if name != "" {
// 		filters = append(filters, "(name eq '"+name+"')")
// 	}
// 	if len(filters) > 0 {
// 		qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
// 	}

// 	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetQueryParams(qParams).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&CodeStreamProjectList{}).
// 		SetAuthToken(targetConfig.accesstoken).
// 		Get("https://" + targetConfig.server + "/project-service/api/projects")

// 	if queryResponse.IsError() {
// 		return nil, queryResponse.Error().(error)
// 	}

// 	log.Debugln(queryResponse.Request.URL)

// 	for _, value := range queryResponse.Result().(*CodeStreamProjectList).Content {
// 		c := CodeStreamProject{}
// 		mapstructure.Decode(value, &c)
// 		projects = append(projects, &c)
// 	}
// 	if len(projects) == 0 {
// 		return nil, errors.New(fmt.Sprint("Unable to find Project \"", name, id, "\""))
// 	} else {
// 		return projects, err
// 	}
// }
