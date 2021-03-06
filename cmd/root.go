/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "puppet-changes",
	Short: "A small cli to look for continious changes or get a history of changes",
	Long:  `Scans changes for hosts and displays if they are recurring. Can also show a history of changes for all or single nodes.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		node, _ := cmd.Flags().GetString("node")
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		cert, _ := cmd.Flags().GetString("cert")
		ca, _ := cmd.Flags().GetString("ca")
		port, _ := cmd.Flags().GetInt("port")
		history, _ := cmd.Flags().GetBool("history")
		sWarns, _ := cmd.Flags().GetBool("show-warnings")
		sErrors, _ := cmd.Flags().GetBool("show-errors")

		master := Master{
			Name:     "default",
			Host:     host,
			Port:     port,
			SSL:      false,
			Key:      key,
			Ca:       ca,
			Insecure: false,
			Cert:     cert,
		}

		if key != "" && ca != "" && cert != "" {
			master.SSL = true
		}

		if !history {
			if node == "" {
				GetContiniousChangesForAll(master, sWarns, sErrors)
			} else {
				GetContiniousChangesForNode(node, master, sWarns, sErrors)
			}
		} else {
			if node == "" {
				GetHistoryForAll(master, sWarns, sErrors)
			} else {
				GetHistoryForNode(node, master, sWarns, sErrors)
			}
		}

	},
}

func GetContiniousChangesForAll(master Master, sWarns bool, sErrors bool) {
	names := GetCertNames(master)
	for _, name := range names {
		GetContiniousChangesForNode(name, master, sWarns, sErrors)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.puppet-changes.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringP("node", "n", "", "If you only want changes for a specific node/certname.")
	rootCmd.Flags().StringP("host", "H", "localhost", "The puppetdb host.")
	rootCmd.Flags().IntP("port", "p", 8080, "The puppetdb port.")
	rootCmd.Flags().StringP("key", "k", "", "The private key.")
	rootCmd.Flags().StringP("cert", "c", "", "The certificate.")
	rootCmd.Flags().StringP("ca", "C", "", "The ca certificate.")
	rootCmd.Flags().BoolP("history", "r", false, "Show all changes by time.")
	rootCmd.Flags().BoolP("show-warnings", "W", false, "Show the warnings as well.")
	rootCmd.Flags().BoolP("show-errors", "E", false, "Show the errors as well.")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".puppet-changes" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".puppet-changes")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
