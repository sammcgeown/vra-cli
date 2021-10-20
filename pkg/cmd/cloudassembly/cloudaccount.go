/*
Package cloudassembly Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cloudassembly

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sammcgeown/vra-cli/pkg/util/helpers"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// GetCloudAccounts returns a list of cloud accounts
func GetCloudAccounts(apiclient *client.MulticloudIaaS, id string, name string, cloudaccounttype string) ([]*models.CloudAccount, error) {

	// transport := httptransport.New(&targetConfig.Server, "", nil)
	// // transport.SetDebug(debug)
	// transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+&targetConfig.AccessToken)
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

	ret, err := apiclient.CloudAccount.GetCloudAccounts(CloudAccountParams)
	if err != nil {
		return nil, err
	}
	return ret.Payload.Content, err

}

// CreateCloudAccountAWS creates a new AWS cloud account
func CreateCloudAccountAWS(apiclient *client.MulticloudIaaS, name, accesskey, secretkey, regions, tags string) (*models.CloudAccountAws, error) {
	AwsSpec := models.CloudAccountAwsSpecification{}
	AwsSpec.Name = &name
	AwsSpec.AccessKeyID = &accesskey
	AwsSpec.SecretAccessKey = &secretkey
	AwsSpec.RegionIds = strings.Split(regions, ",")
	AwsSpec.Tags = helpers.StringToTags(tags)

	createResp, err := apiclient.CloudAccount.CreateAwsCloudAccount(cloud_account.NewCreateAwsCloudAccountParams().WithBody(&AwsSpec))
	if err != nil {
		return nil, err
	}
	return createResp.Payload, nil

}

// CreateCloudAccountvSphere creates a new vSphere cloud account
func CreateCloudAccountvSphere(apiclient *client.MulticloudIaaS, name, description, fqdn, username, password, nsxcloudaccount, cloudproxy, tags string, insecure, createcloudzone bool) (*models.CloudAccountVsphere, error) {

	DatacenterIds, _ := GetvSphereRegions(apiclient, fqdn, username, password, cloudproxy, insecure)

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
	}
	return createResp.Payload, nil

}

// GetvSphereRegions returns a list of vSphere regions
func GetvSphereRegions(apiclient *client.MulticloudIaaS, fqdn, username, password, cloudproxy string, insecure bool) (*models.CloudAccountRegions, error) {

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
	getResp, err := apiclient.CloudAccount.EnumerateVSphereRegions(cloud_account.NewEnumerateVSphereRegionsParams().WithBody(&vSphereSpec))
	if err != nil {
		return nil, err
	}
	return getResp.Payload, nil

}

// CreateCloudAccountNsxT creates a new NSX-T cloud account
func CreateCloudAccountNsxT(apiclient *client.MulticloudIaaS, name, description, fqdn, username, password, vccloudaccount, cloudproxy, tags string, global, manager, insecure bool) (*models.CloudAccountNsxT, error) {

	if vccloudaccount != "" {
		if vCenter, err := GetCloudAccounts(apiclient, "", vccloudaccount, "vsphere"); err != nil {
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
	}
	return createResp.Payload, nil

}

// DeleteCloudAccount deletes a cloud account
func DeleteCloudAccount(apiclient *client.MulticloudIaaS, id string) error {

	_, err := apiclient.CloudAccount.DeleteAwsCloudAccount(cloud_account.NewDeleteAwsCloudAccountParams().WithID(id))
	if err != nil {
		return err
	}
	return nil

}
