/*
Package cmd Copyright 2021 VMware, Inc.
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

func GetConnection(config types.Config, insecure bool) error {
	if TestAccessToken(config, insecure) { // If the Access Token is OK
		log.Debugln("Access Token is valid")
	} else {
		var refreshTokenError, credentialError error
		config.AccessToken, refreshTokenError = AuthenticateApiToken(config, insecure) // Test the API Token (refresh_token)
		if refreshTokenError != nil {                                                  // We could not get an access token from the API Token
			log.Debugln("Refresh Token is invalid")
			if config.Server == "api.mgmt.cloud.vmware.com" { // If it's vRA Cloud we have no credentials to authenticate
				return refreshTokenError // Return the token error
			}
			config.ApiToken, credentialError = AuthenticateCredentials(config, insecure)
			if credentialError != nil {
				return credentialError // Return the credential error
			}
			// Try again, now we have a new access token
			config.AccessToken, refreshTokenError = AuthenticateApiToken(config, insecure) // Test the API Token (refresh_token)
			if refreshTokenError != nil {
				return refreshTokenError
			}
		}

		if viper.ConfigFileUsed() != "" { // If we're using a Config file
			viper.Set("target."+config.Name+".AccessToken", config.AccessToken)
			viper.Set("target."+config.Name+".ApiToken", config.ApiToken)
			viper.WriteConfig()
		}

	}
	return nil
}

// authenticateCredentials - returns the API Refresh Token for vRA On-premises (8.0.1+)
func AuthenticateCredentials(config types.Config, ignoreCert bool) (string, error) {
	log.Debugln("Authenticating vRA with Credentials")
	var authPath string
	var authBody types.AuthenticationRequest
	authBody.Username = config.Username
	authBody.Password = config.Password
	client := resty.New()

	if config.Domain == "" {
		log.Debugln("Basic Auth")
		// Use Basic Authentication
		authPath = "/csp/gateway/am/api/login?access_token"
	} else {
		log.Debugln("Enhanced Auth")
		// Use Enhanced Login (e.g. domain users)
		authPath = "/csp/gateway/am/idp/auth/login?access_token"
		authBody.Domain = config.Domain
	}

	loginResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetBody(authBody).
		SetResult(&types.AuthenticationResponse{}).
		SetError(&types.Authentication{}).
		Post("https://" + config.Server + authPath)
	if loginResponse.IsError() {
		log.Debugln("Authentication failed")
		return "", errors.New(loginResponse.Error().(*types.AuthenticationError).ServerMessage)
	}
	log.Debugln("Authentication succeeded")
	return loginResponse.Result().(*types.AuthenticationResponse).RefreshToken, err
}

// authenticateApiToken - get vRA Access token (valid for 8h)
func AuthenticateApiToken(config types.Config, ignoreCert bool) (string, error) {
	log.Debug("Attempting to authenticate the API Refresh Token")
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetBody(types.Authentication{RefreshToken: config.ApiToken}).
		SetResult(&types.AuthenticationResponse{}).
		SetError(&types.Authentication{}).
		Post("https://" + config.Server + "/iaas/api/login")
	if queryResponse.IsError() {
		log.Debug("Refresh Token failed")
		return "", errors.New(queryResponse.Error().(*types.AuthenticationError).Message)
	}
	log.Debug("Refresh Token succeeded")
	return queryResponse.Result().(*types.AuthenticationResponse).Token, err
}

func TestAccessToken(config types.Config, ignoreCert bool) bool {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetHeader("Accept", "application/json").
		SetAuthToken(config.AccessToken).
		SetResult(&types.UserPreferences{}).
		SetError(&types.Exception{}).
		Get("https://" + config.Server + "/pipeline/api/user-preferences")
	if err != nil {
		log.Warnln(err)
		return false
	}
	// log.Debugln(queryResponse.RawResponse)
	if queryResponse.StatusCode() == 401 {
		log.Debugln("Access Token Expired")
		return false
	}
	log.Debugln("Access Token OK (Username:", queryResponse.Result().(*types.UserPreferences).UserName, ")")
	return true
}

func GetApiClient(config types.Config, debug bool) *client.MulticloudIaaS {
	transport := httptransport.New(config.Server, "", nil)
	transport.SetDebug(debug)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+config.AccessToken)
	apiclient := client.New(transport, strfmt.Default)
	return apiclient
}
