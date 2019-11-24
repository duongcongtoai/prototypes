/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var configFile string

var config struct {
	ElasticSearch struct {
		Port     int
		HostName string
		Scheme   string
	}
	Indexing struct {
		ItemNumbs int
		IndexName string
	}
}
var mapping = `
{
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties": {
			"title": {
				"type":"text",
				"fields": {
					"keyword": {
						"type":"keyword"
					}
				}
			},
			"desc":{
				"type": "text",
				"fields":{
					"keyword":{
						"type":"keyword"
					}
				}
			},
			"number": {
				"type": "integer"
			}
		}
	}
}
`
var configDefault = []byte(`
elasticsearch:
  port: 9200
  hostname: localhost
  scheme: http://
indexing:
  itemnumbs: 100
  indexname: testindex
`)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "indexer",
	Short: "A brief description of your application",
	Long:  `Index to elaticsearch`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE:          index,
	SilenceErrors: true,
	SilenceUsage:  true,
}
var generateConfigCmd = &cobra.Command{
	Use:   "gen",
	Short: "Something",
	Long:  `Hello world`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s", bytes.TrimSpace(configDefault))
		return nil
	},
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping",
	RunE:  ping,
}

var showConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "nothing",
	Long:  `Display configuration`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("config is %v", config)
		return nil
	},
}
var fakeDataCmd = &cobra.Command{
	Use:   "fake",
	Short: "nothing",
	RunE:  fakeData,
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file")
	rootCmd.AddCommand(generateConfigCmd)
	rootCmd.AddCommand(showConfigCmd)
	rootCmd.AddCommand(pingCmd)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return errors.New("SilentErr")
	})
	rootCmd.AddCommand(fakeDataCmd)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("yaml")

	viper.ReadConfig(bytes.NewBuffer(configDefault))
	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.ReadInConfig()
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.Unmarshal(&config)
}
