/*
Package config Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package config

// func GetConfigFromViper() (string, cmd.Config) {

// }

// func GetConfigFromEnv() cmd.Config {
// 	var envConfig cmd.Config
// 	// Bind ENV variables
// 	viper.SetEnvPrefix("vra")
// 	viper.AutomaticEnv()

// 	// If we're using ENV variables
// 	if viper.Get("server") != nil { // VRA_SERVER environment variable is set
// 		log.Debugln("Using ENV variables")
// 		envConfig = cmd.Config{
// 			Domain:      viper.GetString("domain"),
// 			Server:      sanitize.URL(viper.GetString("server")),
// 			Username:    viper.GetString("username"),
// 			Password:    viper.GetString("password"),
// 			ApiToken:    viper.GetString("apitoken"),
// 			AccessToken: viper.GetString("accesstoken"),
// 		}
// 	}

// 	return envConfig
// }

// func GetConfigFromFile(configFile string) (error, cmd.Config) {
// 	var fileConfig cmd.Config

// 	if configFile != "" {
// 		if file, err := os.Stat(configFile); err == nil { // Check if it exists
// 			viper.SetConfigFile(file.Name())
// 		} else {
// 			return errors.New("File specified with --config does not exist"), cmd.Config{}
// 		}
// 	}
// 	return nil, fileConfig
// }
