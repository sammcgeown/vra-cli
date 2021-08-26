/*
Package config Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package config

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/mrz1836/go-sanitize"
	"github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetConfigFromEnv returns a config object from environment variables
func GetConfigFromEnv() *types.Config {
	log.Debugln("Using config: ENV")
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
	log.Debugln("Config:", config)
	return &config
}

// GetConfigFromFile returns a config object from a file
func GetConfigFromFile(configFile string) *types.Config {
	if configFile != "" { // If the user has specified a config file
		if _, err := os.Stat(configFile); err == nil { // Check if it exists
			viper.SetConfigFile(configFile)
		} else {
			log.Fatalln("File specified with --config does not exist (" + configFile + ")")
		}
	} else {
		// Home directory
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
		}

		viper.SetConfigName(".vra-cli")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
	}

	// Attempt to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	} else {
		log.Debugln("Using config:", viper.ConfigFileUsed())
	}

	config := types.Config{}

	currentTargetName := viper.GetString("currentTargetName")
	if currentTargetName != "" {
		log.Infoln("Context:", currentTargetName)
		configuration := viper.Sub("target." + currentTargetName)
		if configuration == nil { // Sub returns nil if the key cannot be found
			log.Fatalln("Target configuration not found")
		}
		config = types.Config{
			Name:        currentTargetName,
			Domain:      configuration.GetString("domain"),
			Server:      sanitize.URL(configuration.GetString("server")),
			Username:    configuration.GetString("username"),
			Password:    configuration.GetString("password"),
			ApiToken:    configuration.GetString("apitoken"),
			AccessToken: configuration.GetString("accesstoken"),
		}
	} else {
		log.Fatalln("No target specified, use `vra-cli config use-target --name <target name>` to specify a name")
	}

	return &config
}
