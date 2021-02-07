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
	"encoding/json"
	"fmt"
	"io/ioutil"
	 
	"log"
	"net/http"
	"strings"
	"kib-sync/client"

	"github.com/spf13/cobra"
)

const CrateMonitorPath = "_opendistro/_alerting/monitors"

// creates cobra command that represents the push command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create monitors in the given kiban cluster",
	Long:  `create monitors in the given kiban cluster`,
	Run:   CreateMonitorsConfig(NewCreateHandler()),
}

type CreateHandler func(url string, body []byte) (*http.Response, error)

func NewCreateHandler() CreateHandler {
	return func(url string, body []byte) (*http.Response, error) {
		// create http client
		client := es.NewClient( getValue(URL), getValue(UserName), getValue(Password))

		// invoke create
		return client.Do(CrateMonitorPath, http.MethodPost, body)
	}
}

// PushMonitorsConfig pushes the monitor config to kibana cluster
func CreateMonitorsConfig(handler CreateHandler) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		log.Println("invoking create command")

		fileInfos, err := ioutil.ReadDir("monitors")
		if err != nil {
			log.Fatal(err)
		}

		// Process all the kibana monitor config and push the update to Kibana
		for _, fileInfo := range fileInfos {
			log.Printf("creating kibana monitor: %s to kiban cluster", strings.SplitAfter(fileInfo.Name(), ".")[0])
			bytes, _ := ioutil.ReadFile(fmt.Sprintf("monitors/%s", fileInfo.Name()))

			var m map[string]interface{}
			if err := json.NewDecoder(strings.NewReader(string(bytes))).Decode(&m); err != nil {
				ErrorLog.Printf("failed to decode the monitor config: %s", err.Error())
				continue
			}
			// index := m["_index"].(string)
			documentID := m["_id"].(string)
			monitorDoc := m["_source"].(map[string]interface{})
			docBytes, err := json.Marshal(&monitorDoc)
			if err != nil {
				log.Fatal("failed to marshall the monitor doc")
			}
			log.Printf("pushing monitor config for document id: %s", documentID)

			url := fmt.Sprintf("%s/%s", getValue(URL), "_opendistro/_alerting/monitors")
			res, err := handler(url, docBytes)

			if err != nil {
				ErrorLog.Printf("Error getting response: %s", err.Error())
				continue
			}

			defer res.Body.Close()

			if res.StatusCode > 201 {
				InfoLog.Printf("[%d] Error indexing document", res.StatusCode)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					ErrorLog.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					InfoLog.Printf("[%d] %s; version=%d; monitor=%s", res.StatusCode, r["result"], int(r["_version"].(float64)), strings.SplitAfter(fileInfo.Name(), ".")[0])
				}
			}
		}

	}

}

