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

// // GetConnection - returns a connection to vRA
// func GetConnection(config *types.Config, insecure bool) error {
// 	if TestAccessToken(config, insecure) { // If the Access Token is OK
// 		log.Debugln("Access Token is valid")
// 	} else {
// 		var refreshTokenError, credentialError error
// 		config.AccessToken, refreshTokenError = AuthenticateAPIToken(config, insecure) // Test the API Token (refresh_token)
// 		if refreshTokenError != nil {                                                  // We could not get an access token from the API Token
// 			log.Debugln("Refresh Token is invalid")
// 			if config.Server == "api.mgmt.cloud.vmware.com" { // If it's vRA Cloud we have no credentials to authenticate
// 				return refreshTokenError // Return the token error
// 			}
// 			config.ApiToken, credentialError = AuthenticateCredentials(*config, insecure)
// 			if credentialError != nil {
// 				return credentialError // Return the credential error
// 			}
// 			// Try again, now we have a new access token
// 			config.AccessToken, refreshTokenError = AuthenticateAPIToken(config, insecure) // Test the API Token (refresh_token)
// 			if refreshTokenError != nil {
// 				return refreshTokenError
// 			}
// 		}

// 		if viper.ConfigFileUsed() != "" { // If we're using a Config file
// 			viper.Set("target."+config.Name+".AccessToken", config.AccessToken)
// 			viper.Set("target."+config.Name+".ApiToken", config.ApiToken)
// 			viper.WriteConfig()
// 		}

// 	}
// 	return nil
// }

// ValidateConfiguration - returns a connection to vRA
func ValidateConfiguration(APIClient *types.APIClientOptions) error {
	// Get a Resty client
	APIClient.RESTClient = GetRESTClient(APIClient.Config, APIClient.VerifySSL, APIClient.Debug)
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
				SetError(&types.Authentication{}).
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

	APIClient.RESTClient = GetRESTClient(APIClient.Config, APIClient.VerifySSL, APIClient.Debug)
	APIClient.SDKClient = GetAPIClient(APIClient.Config, APIClient.Debug)

	return nil
}

// // AuthenticateCredentials - returns the API Refresh Token for vRA On-premises (8.0.1+)
// func AuthenticateCredentials(config types.Config, ignoreCert bool) (string, error) {
// 	log.Debugln("Authenticating vRA with Credentials")
// 	var authPath string
// 	authBody := &types.AuthenticationRequest{
// 		Username: config.Username,
// 		Password: config.Password,
// 	}

// 	client := resty.New()

// 	if config.Domain == "" {
// 		log.Debugln("Basic Auth")
// 		// Use Basic Authentication
// 		authPath = "/csp/gateway/am/api/login?access_token"
// 	} else {
// 		log.Debugln("Enhanced Auth")
// 		// Use Enhanced Login (e.g. domain users)
// 		authPath = "/csp/gateway/am/idp/auth/login?access_token"
// 		authBody.Domain = config.Domain
// 	}

// 	loginResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetBody(authBody).
// 		SetResult(&types.AuthenticationResponse{}).
// 		SetError(&types.Authentication{}).
// 		Post("https://" + config.Server + authPath)
// 	if loginResponse.IsError() {
// 		log.Debugln("Authentication failed")
// 		return "", errors.New(loginResponse.Error().(*types.AuthenticationError).ServerMessage)
// 	}
// 	log.Debugln("Authentication succeeded")
// 	return loginResponse.Result().(*types.AuthenticationResponse).RefreshToken, err
// }

// // AuthenticateAPIToken - get vRA Access token (valid for 8h)
// func AuthenticateAPIToken(config *types.Config, ignoreCert bool) (string, error) {
// 	log.Debug("Attempting to authenticate the API Refresh Token")
// 	client := resty.New()
// 	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
// 		SetBody(types.Authentication{RefreshToken: config.ApiToken}).
// 		SetResult(&types.AuthenticationResponse{}).
// 		SetError(&types.AuthenticationError{}).
// 		Post("https://" + config.Server + "/iaas/api/login")
// 	if queryResponse.IsError() {
// 		log.Debug("Refresh Token failed")
// 		return "", errors.New(queryResponse.Error().(*types.AuthenticationError).Message)
// 	}
// 	log.Debug("Refresh Token succeeded")
// 	return queryResponse.Result().(*types.AuthenticationResponse).Token, err
// }

// TestAccessToken - returns true if the Access Token is valid
// func TestAccessToken(config *types.Config, ignoreCert bool) bool {
// 	client := GetRestClient(config, false, debug)
// 	queryResponse, err := APIClient.RESTClient.R().
// 		SetResult(&types.UserPreferences{}).
// 		SetError(&types.Exception{}).
// 		Get("https://" + config.Server + "/pipeline/api/user-preferences")
// 	if err != nil {
// 		log.Warnln(err)
// 		return false
// 	}
// 	// log.Debugln(queryResponse.RawResponse)
// 	if queryResponse.StatusCode() == 401 {
// 		log.Debugln("Access Token Expired")
// 		return false
// 	}
// 	log.Debugln("Access Token OK (Username:", queryResponse.Result().(*types.UserPreferences).UserName, ")")
// 	return true
// }

// GetAPIClient - returns a vRA API client
func GetAPIClient(config *types.Config, debug bool) *client.MulticloudIaaS {
	transport := httptransport.New(config.Server, "", nil)
	transport.SetDebug(debug)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+config.AccessToken)
	apiclient := client.New(transport, strfmt.Default)
	return apiclient
}

// GetRESTClient - returns a vRA REST client
func GetRESTClient(config *types.Config, insecure bool, debug bool) *resty.Client {
	// Configure the Resty Client
	client := resty.New().
		SetDebug(debug).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: insecure}).
		SetAuthToken(config.AccessToken).
		SetHostURL("https://"+config.Server).
		SetHeader("Accept", "application/json").
		SetError(&types.Exception{})
	return client
}
