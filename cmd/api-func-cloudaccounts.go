/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"strings"

	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"

	httptransport "github.com/go-openapi/runtime/client"
	client "github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func getCloudAccounts(id string, name string, cloudaccounttype string) ([]*models.CloudAccount, error) {

	transport := httptransport.New(targetConfig.server, "", nil)
	// transport.SetDebug(debug)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+targetConfig.accesstoken)
	apiclient := client.New(transport, strfmt.Default)

	var filters []string
	var filter string

	if id != "" {
		filters = append(filters, "(id eq '"+id+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if cloudaccounttype != "" {
		filters = append(filters, "(cloudAccountType eq '"+cloudaccounttype+"')")
	}
	if len(filters) > 0 {
		filter = "(" + strings.Join(filters, " and ") + ")"
	}

	log.Debugln("Filter:", filter)
	CloudAccountParams := cloud_account.NewGetCloudAccountsParams()
	CloudAccountParams.DollarFilter = &filter

	ret, err := apiclient.CloudAccount.GetCloudAccounts(CloudAccountParams)

	return ret.Payload.Content, err

}
