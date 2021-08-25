/*
Package variable Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package variable

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sammcgeown/vra-cli/pkg/util/auth"
	"github.com/sammcgeown/vra-cli/pkg/util/config"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

var (
	targetConfig *types.Config
	client       *resty.Client
	// Config
	insecure = false
)

func TestCreateVariable(t *testing.T) {

	targetConfig = config.GetConfigFromEnv()

	// Validate the configuration and credentials
	if err := auth.GetConnection(targetConfig, insecure); err != nil {
		log.Fatalln(err)
	}

	client = auth.GetRestClient(targetConfig, insecure)

	variable, err := CreateVariable(client, "test", "test", "REGULAR", "vra-cli-testing", "test")

	assert.NilError(t, err)
	assert.Equal(t, variable.Name, "test")
	assert.Equal(t, variable.Type, "REGULAR")
	assert.Equal(t, variable.Description, "test")
	assert.Equal(t, variable.Value, "test")
	assert.Equal(t, variable.Project, "vra-cli-testing")

}
