/*
Package types Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package types

import (
	"github.com/go-resty/resty/v2"
	"github.com/vmware/vra-sdk-go/pkg/client"
)

// Config - main configuration struct
type Config struct {
	Name        string
	Domain      string
	Password    string
	Server      string
	Username    string
	APIToken    string
	AccessToken string
}

// APIClientOptions - options for the API client
type APIClientOptions struct {
	Version    string
	Debug      bool
	VerifySSL  bool
	Confirm    bool
	Force      bool
	RESTClient *resty.Client
	SDKClient  *client.MulticloudIaaS
	Pagination struct {
		PageSize int
		Page     int
		Skip     int
	}
	Config *Config
	Output string
}

// Exception - Generic exception struct
type Exception struct {
	Timestamp   int64       `json:"timestamp"`
	Path        string      `json:"path"`
	Status      int         `json:"status"`
	Error       string      `json:"error"`
	Message     string      `json:"message"`
	RequestID   string      `json:"requestId"`
	Type        string      `json:"@type"`
	StatusCode  int         `json:"statusCode"`
	ErrorCode   int         `json:"errorCode"`
	ReferenceID interface{} `json:"referenceId"`
}

// DocumentsList - Code Stream Documents List structure
type DocumentsList struct {
	Count      int                    `json:"count"`
	TotalCount int                    `json:"totalCount"`
	Links      []string               `json:"links"`
	Documents  map[string]interface{} `json:"documents"`
}

// ContentsList - Generic Contents List Structure
type ContentsList struct {
	Content  []interface{} `json:"content"`
	Pageable struct {
		Sort struct {
			Sorted   bool `json:"sorted"`
			Unsorted bool `json:"unsorted"`
			Empty    bool `json:"empty"`
		} `json:"sort"`
		PageNumber int  `json:"pageNumber"`
		PageSize   int  `json:"pageSize"`
		Offset     int  `json:"offset"`
		Paged      bool `json:"paged"`
		Unpaged    bool `json:"unpaged"`
	} `json:"pageable"`
	TotalElements int  `json:"totalElements"`
	TotalPages    int  `json:"totalPages"`
	Last          bool `json:"last"`
	Sort          struct {
		Sorted   bool `json:"sorted"`
		Unsorted bool `json:"unsorted"`
		Empty    bool `json:"empty"`
	} `json:"sort"`
	First            bool `json:"first"`
	Number           int  `json:"number"`
	NumberOfElements int  `json:"numberOfElements"`
	Size             int  `json:"size"`
	Empty            bool `json:"empty"`
}

// UserPreferences -
type UserPreferences struct {
	Link               string      `json:"_link"`
	UpdateTimeInMicros int         `json:"_updateTimeInMicros"`
	CreateTimeInMicros int         `json:"_createTimeInMicros"`
	Preferences        interface{} `json:"preferences"`
	UserName           string      `json:"userName"`
}

// AuthenticationRequest - vRA Authentication request structure
type AuthenticationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

// TokenRequest - vRA Authentication request structure
type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthenticationResponse - Authentication response structure
type AuthenticationResponse struct {
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Token        string `json:"token"`
}

// Authentication - vRA Authentication request structure for API login with a refresh token
type Authentication struct {
	RefreshToken string `json:"refreshToken"`
}

// AuthenticationError - Authentication error structure
type AuthenticationError struct {
	Timestamp     int64  `json:"timestamp"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Error         string `json:"error"`
	ServerMessage string `json:"serverMessage"`
	Message       string `json:"message"`
	StatusCode    int64  `json:"statusCode"`
	ErrorCode     int64  `json:"errorCode"`
	ServerErrorID string `json:"serverErrorId"`
	DocumentKind  string `json:"documentKind"`
}
