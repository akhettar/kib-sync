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
	"log"
	"os"
	"github.com/spf13/cobra"
)

const (
	// UserName the elk user name
	UserName = "username"
	// Password the elk password
	Password = "password"
	// URL the elk url
	URL = "url"
)

var (
	// WarnLog instance
	WarnLog *log.Logger
	// InfoLog instanace
	InfoLog *log.Logger
	// ErrorLog instance
	ErrorLog *log.Logger
)

var rootCmd = &cobra.Command{
	Use:   "kibana-sync",
	Short: "Kibsync is a tool that fetches configured monitors",
	Long: `This tool performs the followings:
		1. Fetches configured monitors for the given kibana cluster and store them locally as json files.
		2. Pushes the changes done to the monitor's config to Kiban cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("running kibana synchronizer")
	},
}

// Execute the main execute function
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {

	// set the logger
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLog = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// declare all the possible flags for this app
	rootCmd.PersistentFlags().String("username", "", "The kibana cluster username. This is required argument to connect to the ELK cluster")
	rootCmd.PersistentFlags().String("password", "", "The kibana cluster password. This is a required argument to connect to the ELK cluster")
	rootCmd.PersistentFlags().String("url", "", "The kibana cluster url. This is required argument to connect to the ELk cluster")

	// add the command
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(createCmd)
}

func getValue(flag string) string {
	value, err := rootCmd.PersistentFlags().GetString(flag)
	if err != nil {
		ErrorLog.Fatalf("failed to get the required argement %s", flag)
	}
	if value == "" {
		ErrorLog.Fatalf("%s argument is required", flag)
	}
	return value
}
