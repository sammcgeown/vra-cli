/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mrz1836/go-sanitize"
	"github.com/sammcgeown/vra-cli/pkg/util/auth"
	"github.com/sammcgeown/vra-cli/pkg/util/config"
	types "github.com/sammcgeown/vra-cli/pkg/util/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd *cobra.Command
	// Configuration
	cfgFile           string
	currentTargetName string
	targetConfig      types.Config
	version           = "dev"
	commit            = "none"
	date              = "unknown"
	builtBy           = "unknown"
	apiVersion        = "2019-10-17"
	client            *resty.Client
	// Global Flags
	debug      bool
	ignoreCert bool
	// API Paging
	count int
	skip  int
	// Command Flags
	id          string
	name        string
	projectName string
	typename    string
	value       string
	description string
	status      string
	printJson   bool
	exportPath  string
	importPath  string
)

var qParams = map[string]string{
	"apiVersion": "2019-10-17",
}

// rootCmd represents the base command when called without any subcommands
// var rootCmd = &cobra.Command{
// 	Use:   "vra-cli",
// 	Short: "CLI Interface for VMware vRealize Automation Code Stream",
// 	Long:  `Command line interface for VMware vRealize Automation Code Stream`,
// }

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "vra-cli",
		Short: "CLI Interface for VMware vRealize Automation Code Stream",
		Long:  `Command line interface for VMware vRealize Automation Code Stream`,
	}
}

// Execute is the main process
func Execute() {
	rootCmd = NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Warnln(err)
	}
}

func init() {
	cobra.OnInitialize(InitConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vra-cli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&ignoreCert, "ignoreCertificateWarnings", false, "Disable HTTPS Certificate Validation")
	// API Paging
	rootCmd.PersistentFlags().IntVar(&count, "count", 100, "API Paging - Count")
	rootCmd.PersistentFlags().IntVar(&skip, "skip", 0, "API Paging - Skip")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
}

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	// Debug logging
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Debug logging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// Home directory
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err)
	}

	viper.SetConfigName(".vra-cli")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)

	// Bind ENV variables
	viper.SetEnvPrefix("vra")
	viper.AutomaticEnv()

	// If we're using ENV variables
	if viper.Get("server") != nil { // VRA_SERVER environment variable is set
		log.Debugln("Using ENV variables")
		targetConfig = *config.GetConfigFromEnv()
	} else {
		if cfgFile != "" { // If the user has specified a config file
			if file, err := os.Stat(cfgFile); err == nil { // Check if it exists
				viper.SetConfigFile(file.Name())
			} else {
				log.Fatalln("File specified with --config does not exist")
			}
		}
		// Attempt to read the configuration file
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				viper.SetConfigType("yaml")
				viper.WriteConfigAs(filepath.Join(home, ".vra-cli"))
				viper.ReadInConfig()
			} else {
				log.Fatalln(err)
			}
		}
		currentTargetName = viper.GetString("currentTargetName")
		if currentTargetName != "" {
			log.Debugln("Using config:", viper.ConfigFileUsed())
			log.Infoln("Context:", currentTargetName)
			configuration := viper.Sub("target." + currentTargetName)
			if configuration == nil { // Sub returns nil if the key cannot be found
				log.Fatalln("Target configuration not found")
			}
			targetConfig = types.Config{
				Name:        currentTargetName,
				Domain:      configuration.GetString("domain"),
				Server:      sanitize.URL(configuration.GetString("server")),
				Username:    configuration.GetString("username"),
				Password:    configuration.GetString("password"),
				ApiToken:    configuration.GetString("apitoken"),
				AccessToken: configuration.GetString("accesstoken"),
			}
		}

		// Validate the configuration and credentials
		if err := auth.GetConnection(&targetConfig, debug); err != nil {
			log.Fatalln(err)
		}

		// Configure the REST client defaults
		client = resty.New().
			SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).
			SetAuthToken(targetConfig.AccessToken).
			SetHostURL("https://"+targetConfig.Server).
			SetHeader("Accept", "application/json").
			SetError(&types.Exception{})
	}
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run:  func(cmd *cobra.Command, args []string) {},
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run:  func(cmd *cobra.Command, args []string) {},
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run:  func(cmd *cobra.Command, args []string) {},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run:  func(cmd *cobra.Command, args []string) {},
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify the configuration of vra-cli",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run:  func(cmd *cobra.Command, args []string) {},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version information",
	Long:  `Print the current version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("*** vra-cli ***")
		fmt.Println("Build version :", version)
		fmt.Println("Build date    :", date)
		fmt.Println("Build commit  :", commit)
		fmt.Println("Built by      :", builtBy)
	},
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(vra-cli completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ vra-cli completion bash > /etc/bash_completion.d/vra-cli
  # macOS:
  $ vra-cli completion bash > /usr/local/etc/bash_completion.d/vra-cli

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ vra-cli completion zsh > "${fpath[1]}/_vra-cli"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ vra-cli completion fish | source

  # To load completions for each session, execute once:
  $ vra-cli completion fish > ~/.config/fish/completions/vra-cli.fish

PowerShell:

  PS> vra-cli completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> vra-cli completion powershell > vra-cli.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}
