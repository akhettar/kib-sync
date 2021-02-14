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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func Test_createKibanaConfig(t *testing.T) {

	// createTest dir
	createDir("config/dashboard")

	// 1. Create insntace of the crate command
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create all kiban objects (monitors, dashbaor, etc) present in the config folder",
		Long:  `create monitors, destinations, dashboards, email accoutns, group emails, search in the given kibana cluster`,
		Run:   createKibanaConfig(testCreateHandler()),
	}
	// 2. Execute comand
	if err := cmd.Execute(); err != nil {
		t.Errorf("failed to run create command: %v", err)
	}

	os.RemoveAll("config")

}

func testCreateHandler() CreateHandler {
	return func(url string, body []byte) (*http.Response, error) {

		bytes, err := ioutil.ReadFile("../data/dashboard_create_res.json")
		if err != nil {
			log.Fatal(err)
		}
		r := ioutil.NopCloser(strings.NewReader(string(bytes)))
		res := http.Response{StatusCode: 201, Body: r}
		return &res, nil
	}
}
