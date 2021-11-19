/*
Package auth Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package auth

import (
	"crypto/tls"
	"errors"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-resty/resty/v2"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vmware/vra-sdk-go/pkg/client"
)

// ValidateConfiguration - returns a connection to vRA
func ValidateConfiguration(APIClient *types.APIClientOptions) error {
	// Get a Resty client
	APIClient.RESTClient = GetRESTClient(APIClient.Config, APIClient.Version, APIClient.VerifySSL, APIClient.Debug)
	// Query the API to see if we're authenticated
	queryResponse, err := APIClient.RESTClient.R().
		SetResult(&types.UserPreferences{}).
		SetError(&types.Exception{}).
		Get("/pipeline/api/user-preferences")

	if err != nil {
		return err
	}
	log.Debugln(queryResponse.RawResponse)

	// If we cannot query the API, we're not authenticated
	if queryResponse.StatusCode() == 401 {
		log.Debug("Attempting to authenticate the existing API Refresh Token")
		queryResponse, _ := APIClient.RESTClient.R().
			SetBody(types.Authentication{RefreshToken: APIClient.Config.APIToken}).
			SetResult(&types.AuthenticationResponse{}).
			SetError(&types.AuthenticationError{}).
			Post("/iaas/api/login")

		// If we get an error response, we have an expired token
		if queryResponse.IsError() {
			log.Debugln("Refresh Token failed", queryResponse.Error().(*types.AuthenticationError).Message)
			// If it's vRA Cloud we have no credentials to authenticate
			if APIClient.Config.Server == "api.mgmt.cloud.vmware.com" {
				// Return the token error
				return errors.New(queryResponse.Error().(*types.AuthenticationError).Message)
			}
			// If it's vRA On-premises, we have credentials to authenticate
			log.Debugln("Authenticating vRA with Credentials", APIClient.Config.Username)
			var authPath string
			authBody := &types.AuthenticationRequest{
				Username: APIClient.Config.Username,
				Password: APIClient.Config.Password,
			}
			if APIClient.Config.Domain == "" {
				log.Debugln("Using Basic Authentication")
				authPath = "/csp/gateway/am/api/login?access_token"
			} else {
				log.Debugln("Using Identity Provider Authentication", APIClient.Config.Username, APIClient.Config.Domain)
				authPath = "/csp/gateway/am/idp/auth/login?access_token"
				authBody.Domain = APIClient.Config.Domain
			}

			loginResponse, _ := APIClient.RESTClient.R().
				SetBody(authBody).
				SetResult(&types.AuthenticationResponse{}).
				SetError(&types.AuthenticationError{}).
				Post(authPath)

			if loginResponse.IsError() {
				log.Debugln("Authentication failed")
				return errors.New(loginResponse.Error().(*types.AuthenticationError).ServerMessage)
			}
			APIClient.Config.APIToken = loginResponse.Result().(*types.AuthenticationResponse).RefreshToken

			// Authenticate with the IaaS API
			queryResponse, _ := APIClient.RESTClient.R().
				SetBody(types.Authentication{RefreshToken: APIClient.Config.APIToken}).
				SetResult(&types.AuthenticationResponse{}).
				SetError(&types.AuthenticationError{}).
				Post("/iaas/api/login")
			if queryResponse.IsError() {
				return errors.New(queryResponse.Error().(*types.AuthenticationError).Message)
			}
			APIClient.Config.AccessToken = queryResponse.Result().(*types.AuthenticationResponse).Token
			log.Debugln("Authentication succeeded")
			// Else we have a valid token
		} else {
			log.Debug("Refresh Token succeeded")
			APIClient.Config.AccessToken = queryResponse.Result().(*types.AuthenticationResponse).Token
		}
		// Write tokens to config file
		if viper.ConfigFileUsed() != "" { // If we're using a Config file
			viper.Set("target."+APIClient.Config.Name+".AccessToken", APIClient.Config.AccessToken)
			viper.Set("target."+APIClient.Config.Name+".ApiToken", APIClient.Config.APIToken)
			viper.WriteConfig()
		}
	} else {
		log.Debugln("Access Token OK (Username:", queryResponse.Result().(*types.UserPreferences).UserName, ")")
	}

	APIClient.RESTClient = GetRESTClient(APIClient.Config, APIClient.Version, APIClient.VerifySSL, APIClient.Debug)
	APIClient.SDKClient = GetAPIClient(APIClient.Config, APIClient.Debug)

	return nil
}

// GetAPIClient - returns a vRA API client
func GetAPIClient(config *types.Config, debug bool) *client.MulticloudIaaS {
	transport := httptransport.New(config.Server, "", nil)
	transport.SetDebug(debug)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+config.AccessToken)
	apiclient := client.New(transport, strfmt.Default)
	return apiclient
}

// GetRESTClient - returns a vRA REST client
func GetRESTClient(config *types.Config, apiVersion string, insecure bool, debug bool) *resty.Client {
	// Configure the Resty Client
	client := resty.New().
		SetDebug(debug).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: insecure}).
		SetAuthToken(config.AccessToken).
		SetHostURL("https://"+config.Server).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetError(&types.Exception{}).
		SetQueryParam("apiVersion", apiVersion)
	return client
}
