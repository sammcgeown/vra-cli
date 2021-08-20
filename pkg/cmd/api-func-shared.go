/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"gopkg.in/yaml.v2"
)

func ensureTargetConnection() error {
	if testAccessToken() { // If the Access Token is OK
		log.Debugln("Access Token is valid")
	} else {
		var refreshTokenError, credentialError error
		targetConfig.AccessToken, refreshTokenError = authenticateApiToken(targetConfig.Server, targetConfig.ApiToken) // Test the API Token (refresh_token)
		if refreshTokenError != nil {                                                                                  // We could not get an access token from the API Token
			log.Debugln("Refresh Token is invalid")
			if targetConfig.Server == "api.mgmt.cloud.vmware.com" { // If it's vRA Cloud we have no credentials to authenticate
				return refreshTokenError // Return the token error
			}
			targetConfig.ApiToken, credentialError = authenticateCredentials(targetConfig.Server, targetConfig.Username, targetConfig.Password, targetConfig.Domain)
			if credentialError != nil {
				return credentialError // Return the credential error
			}
			// Try again, now we have a new access token
			targetConfig.AccessToken, refreshTokenError = authenticateApiToken(targetConfig.Server, targetConfig.ApiToken) // Test the API Token (refresh_token)
			if refreshTokenError != nil {
				return refreshTokenError
			}
		}

		if viper.ConfigFileUsed() != "" { // If we're using a Config file
			viper.Set("target."+currentTargetName+".AccessToken", targetConfig.AccessToken)
			viper.Set("target."+currentTargetName+".ApiToken", targetConfig.ApiToken)
			viper.WriteConfig()
		}

	}
	return nil
}

// authenticateCredentials - returns the API Refresh Token for vRA On-premises (8.0.1+)
func authenticateCredentials(server string, username string, password string, domain string) (string, error) {
	log.Debugln("Authenticating vRA with Credentials")
	var authPath string
	var authBody AuthenticationRequest
	authBody.Username = username
	authBody.Password = password
	client := resty.New()

	if domain == "" {
		log.Debugln("Basic Auth")
		// Use Basic Authentication
		authPath = "/csp/gateway/am/api/login?access_token"
	} else {
		log.Debugln("Enhanced Auth")
		// Use Enhanced Login (e.g. domain users)
		authPath = "/csp/gateway/am/idp/auth/login?access_token"
		authBody.Domain = domain
	}

	loginResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetBody(authBody).
		SetResult(&AuthenticationResponse{}).
		SetError(&AuthenticationError{}).
		Post("https://" + server + authPath)
	if loginResponse.IsError() {
		log.Debugln("Authentication failed")
		return "", errors.New(loginResponse.Error().(*AuthenticationError).ServerMessage)
	}
	log.Debugln("Authentication succeeded")
	return loginResponse.Result().(*AuthenticationResponse).RefreshToken, err
}

// authenticateApiToken - get vRA Access token (valid for 8h)
func authenticateApiToken(server string, token string) (string, error) {
	log.Debug("Attempting to authenticate the API Refresh Token")
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetBody(ApiAuthentication{token}).
		SetResult(&ApiAuthenticationResponse{}).
		SetError(&ApiAuthenticationError{}).
		Post("https://" + server + "/iaas/api/login")
	if queryResponse.IsError() {
		log.Debug("Refresh Token failed")
		return "", errors.New(queryResponse.Error().(*ApiAuthenticationError).Message)
	}
	log.Debug("Refresh Token succeeded")
	return queryResponse.Result().(*ApiAuthenticationResponse).Token, err
}

func testAccessToken() bool {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetHeader("Accept", "application/json").
		SetAuthToken(targetConfig.AccessToken).
		SetResult(&UserPreferences{}).
		SetError(&CodeStreamException{}).
		Get("https://" + targetConfig.Server + "/pipeline/api/user-preferences")
	if err != nil {
		log.Warnln(err)
		return false
	}
	// log.Debugln(queryResponse.RawResponse)
	if queryResponse.StatusCode() == 401 {
		log.Debugln("Access Token Expired")
		return false
	}
	log.Debugln("Access Token OK (Username:", queryResponse.Result().(*UserPreferences).UserName, ")")
	return true
}

func getApiClient() *client.MulticloudIaaS {
	transport := httptransport.New(targetConfig.Server, "", nil)
	transport.SetDebug(debug)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+targetConfig.AccessToken)
	apiclient := client.New(transport, strfmt.Default)
	return apiclient
}

func exportYaml(name, project, path, object string) error {
	var exportPath string
	var qParams = make(map[string]string)
	qParams[object] = name
	qParams["project"] = project
	if path != "" {
		exportPath = path
	} else {
		exportPath, _ = os.Getwd()
	}
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/x-yaml;charset=UTF-8").
		SetAuthToken(targetConfig.AccessToken).
		SetOutput(filepath.Join(exportPath, name+".yaml")).
		SetError(&CodeStreamException{}).
		Get("https://" + targetConfig.Server + "/pipeline/api/export")
	log.Debugln(queryResponse.Request.RawRequest.URL)

	if queryResponse.IsError() {
		return errors.New(queryResponse.Status())
	}
	return nil
}

// importYaml import a yaml pipeline or endpoint
func importYaml(yamlPath, action, project, importType string) error {
	var pipeline CodeStreamPipelineYaml
	var endpoint CodeStreamEndpointYaml

	qParams["action"] = action
	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return err
	}

	if project != "" { // If the project flag is set we need to update the project value
		if importType == "pipeline" {
			yamlErr := yaml.Unmarshal(yamlBytes, &pipeline)
			if yamlErr != nil {
				return yamlErr
			}
			pipeline.Project = project
			yamlBytes, _ = yaml.Marshal(pipeline)
		} else {
			yamlErr := yaml.Unmarshal(yamlBytes, &endpoint)
			if yamlErr != nil {
				return yamlErr
			}
			endpoint.Project = project
			yamlBytes, _ = yaml.Marshal(endpoint)
		}
	}

	yamlPayload := string(yamlBytes)
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Content-Type", "application/x-yaml").
		SetBody(yamlPayload).
		SetAuthToken(targetConfig.AccessToken).
		SetError(&CodeStreamException{}).
		Post("https://" + targetConfig.Server + "/pipeline/api/import")
	log.Debugln(queryResponse.Request.RawRequest.URL)
	if queryResponse.IsError() {
		return queryResponse.Error().(error)
	}
	var importResponse CodeStreamPipelineImportResponse
	if err = yaml.Unmarshal(queryResponse.Body(), &importResponse); err != nil {
		return err
	}

	if importResponse.Status != "CREATED" && action == "create" {
		return errors.New(importResponse.Status + " - " + importResponse.StatusMessage)
	}
	if importResponse.Status != "UPDATED" && action == "apply" {
		return errors.New(importResponse.Status + " - " + importResponse.StatusMessage)
	}
	return nil
}

func exportCloudTemplate(name, project, content, path string) error {
	var exportPath string
	if path != "" {
		exportPath = path
		_, folderError := os.Stat(exportPath) // Get file system info
		if os.IsNotExist(folderError) {       // If it doesn't exist
			log.Debugln("Folder doesn't exist - creating")
			mkdirErr := os.MkdirAll(exportPath, os.FileMode(0755)) // Attempt to make it
			if mkdirErr != nil {
				return mkdirErr
			}
		}
	} else {
		// If path is not specified, use the current path
		exportPath, _ = os.Getwd()
	}
	exportPath = filepath.Join(exportPath, project+" - "+name+".yaml")
	f, cerr := os.Create(exportPath) // Open the file for writing
	if cerr != nil {
		return cerr
	}
	defer f.Close() // Defer closing until just before this function returns
	_, werr := f.WriteString(content)
	if werr != nil {
		return werr
	}
	return nil
}
