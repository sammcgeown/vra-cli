/*
Package config Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package config

import (
	"github.com/mrz1836/go-sanitize"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	"github.com/spf13/viper"
)

func GetConfigFromEnv() *types.Config {
	viper.SetEnvPrefix("vra")
	viper.AutomaticEnv()
	config := types.Config{
		Domain:      viper.GetString("domain"),
		Server:      sanitize.URL(viper.GetString("server")),
		Username:    viper.GetString("username"),
		Password:    viper.GetString("password"),
		ApiToken:    viper.GetString("apitoken"),
		AccessToken: viper.GetString("accesstoken"),
	}
	return &config
}
