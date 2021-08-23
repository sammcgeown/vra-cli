/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sammcgeown/vra-cli/pkg/util/auth"
	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func getCloudAccounts(id string, name string, cloudaccounttype string) ([]*models.CloudAccount, error) {

	// transport := httptransport.New(targetConfig.Server, "", nil)
	// // transport.SetDebug(debug)
	// transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+targetConfig.AccessToken)
	// apiclient := client.New(transport, strfmt.Default)

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

	apiclient := auth.GetApiClient(targetConfig, debug)

	ret, err := apiclient.CloudAccount.GetCloudAccounts(CloudAccountParams)
	if err != nil {
		return nil, err
	} else {
		return ret.Payload.Content, err
	}

}

func createCloudAccountAws(name, accesskey, secretkey, regions, tags string) (*models.CloudAccountAws, error) {
	AwsSpec := models.CloudAccountAwsSpecification{}
	AwsSpec.Name = &name
	AwsSpec.AccessKeyID = &accesskey
	AwsSpec.SecretAccessKey = &secretkey
	AwsSpec.RegionIds = strings.Split(regions, ",")
	AwsSpec.Tags = helpers.StringToTags(tags)

	apiclient := auth.GetApiClient(targetConfig, debug)
	createResp, err := apiclient.CloudAccount.CreateAwsCloudAccount(cloud_account.NewCreateAwsCloudAccountParams().WithBody(&AwsSpec))
	if err != nil {
		return nil, err
	} else {
		return createResp.Payload, nil
	}
}

func createCloudAccountvSphere(name, description, fqdn, username, password, nsxcloudaccount, cloudproxy, tags string, insecure, createcloudzone bool) (*models.CloudAccountVsphere, error) {
	apiclient := auth.GetApiClient(targetConfig, debug)

	DatacenterIds, _ := getvSphereRegions(fqdn, username, password, cloudproxy, insecure)

	vSphereSpec := models.CloudAccountVsphereSpecification{
		Name:                        &name,
		Description:                 description,
		HostName:                    &fqdn,
		Username:                    &username,
		Password:                    &password,
		CreateDefaultZones:          createcloudzone,
		AcceptSelfSignedCertificate: insecure,
		RegionIds:                   DatacenterIds.ExternalRegionIds,
		Tags:                        helpers.StringToTags(tags),
	}
	if nsxcloudaccount != "" {
		vSphereSpec.AssociatedCloudAccountIds = []string{nsxcloudaccount}
	}

	createResp, err := apiclient.CloudAccount.CreateVSphereCloudAccount(cloud_account.NewCreateVSphereCloudAccountParams().WithBody(&vSphereSpec))
	if err != nil {
		return nil, err
	} else {
		return createResp.Payload, nil
	}
}

func getvSphereRegions(fqdn, username, password, cloudproxy string, insecure bool) (*models.CloudAccountRegions, error) {
	apiclient := auth.GetApiClient(targetConfig, debug)

	vSphereSpec := models.CloudAccountVsphereSpecification{
		AcceptSelfSignedCertificate: insecure,
		HostName:                    &fqdn,
		Password:                    &password,
		Username:                    &username,
	}
	if cloudproxy != "" {
		vSphereSpec.Dcid = cloudproxy
	}
	// Get Regions
	if getResp, err := apiclient.CloudAccount.EnumerateVSphereRegions(cloud_account.NewEnumerateVSphereRegionsParams().WithBody(&vSphereSpec)); err != nil {
		return nil, err
	} else {
		return getResp.Payload, nil
	}

}

func createCloudAccountNsxT(name, description, fqdn, username, password, vccloudaccount, cloudproxy, tags string, global, manager, insecure bool) (*models.CloudAccountNsxT, error) {
	apiclient := auth.GetApiClient(targetConfig, debug)

	if vccloudaccount != "" {
		if vCenter, err := getCloudAccounts("", vccloudaccount, "vsphere"); err != nil {
			log.Warnln("Unable to find a vSphere Cloud Account named "+vccloudaccount+" to associate with NSXT Cloud Account", err)
			vccloudaccount = ""
		} else {
			vccloudaccount = *vCenter[0].ID
		}
	}

	NsxTSpec := models.CloudAccountNsxTSpecification{
		AcceptSelfSignedCertificate: insecure,
		Name:                        &name,
		Description:                 description,
		HostName:                    &fqdn,
		Username:                    &username,
		Password:                    &password,
		ManagerMode:                 manager,
		IsGlobalManager:             global,
		Tags:                        helpers.StringToTags(tags),
	}
	if vccloudaccount != "" {
		NsxTSpec.AssociatedCloudAccountIds = []string{vccloudaccount}
	}
	if cloudproxy != "" {
		NsxTSpec.Dcid = &cloudproxy
	}

	createResp, err := apiclient.CloudAccount.CreateNsxTCloudAccount(cloud_account.NewCreateNsxTCloudAccountParams().WithBody(&NsxTSpec))
	if err != nil {
		return nil, err
	} else {
		return createResp.Payload, nil
	}
}

func deleteCloudAccount(id string) error {

	apiclient := auth.GetApiClient(targetConfig, debug)

	_, err := apiclient.CloudAccount.DeleteAwsCloudAccount(cloud_account.NewDeleteAwsCloudAccountParams().WithID(id))
	if err != nil {
		return err
	} else {
		return nil
	}
}